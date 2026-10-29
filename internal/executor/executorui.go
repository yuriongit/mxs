/*
Package executor provides the execution screen for locating scripts,
streaming their output, and displaying their completion status.

The screens directory contains terminal user interfaces used by the application.
*/
package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/scr"
	"github.com/yuriongit/mxs/internal/ui"
)

// ============================================================================
// Application State Enum
// ============================================================================

type state int

const (
	stateSearching state = iota // 0: Finding script in ~/.mxs/scripts
	stateRunning                // 1: Streaming script output
	stateFinished               // 2: Completed execution or failed
)

// ============================================================================
// Bubble Tea Message Definitions
// ============================================================================

type scriptFoundMsg struct{ path string }
type scriptNotFoundMsg struct{ err error }
type validScriptMsg struct{ scriptDesc string }
type invalidScriptMsg struct{ err error }
type outputLineMsg string
type executionFinishedMsg struct{ err error }
type tickMsg time.Time // Stopwatch update tick message

// ============================================================================
// Bubble Tea Model
// ============================================================================

type model struct {
	targetName      string   // Name provided by user (e.g. "gitupd")
	scriptArgs      []string // Arguments forwarded to script
	scriptPath      string   // Resolved path on disk
	scriptDesc      string
	cfg             *cfg.Cfg
	scr             *scr.Scr
	state           state         // Current lifecycle state
	spinner         spinner.Model // Bubbles spinner component
	ellipsisSpinner spinner.Model
	output          []string      // Streamed stdout/stderr log buffer
	err             error         // Execution error, if any
	startTime       time.Time     // Timer start
	duration        time.Duration // Total elapsed run time
	sub             chan tea.Msg  // Async channel for streaming output
}

// formatScriptLog formats and styles each line of output from the scirpt as a "log".
// Appends the updated slice into [model.output].
func (m *model) formatScriptLog(msg outputLineMsg) {
	m.output = append(m.output, ui.Log.Render(string(msg)))
}

// InitialModel creates the initial TUI model for script execution.
//
//nolint:revive // unexported-return
func InitialModel(
	c *cfg.Cfg,
	scr *scr.Scr,
	targetName string,
	scriptArgs []string,
) model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = ui.Spinner

	ellipsisSpinner := spinner.New()
	ellipsisSpinner.Spinner = spinner.Ellipsis
	ellipsisSpinner.Style = ui.Subtle

	return model{
		targetName:      targetName,
		scriptArgs:      scriptArgs,
		cfg:             c,
		scr:             scr,
		state:           stateSearching,
		spinner:         s,
		ellipsisSpinner: ellipsisSpinner,
		output:          make([]string, 0),
		sub:             make(chan tea.Msg),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.ellipsisSpinner.Tick,
		findScriptCmd(m.targetName),
	)
}

// Run executes the UI with reserved terminal buffer space to prevent screen scrolling
func Run(c *cfg.Cfg, scr *scr.Scr, targetName string, scriptArgs []string) error {
	// Pre-scroll 20 lines so the terminal never scrolls up during rendering
	fmt.Print(strings.Repeat("\n", 20) + "\033[20A")

	p := tea.NewProgram(InitialModel(c, scr, targetName, scriptArgs))
	_, err := p.Run()
	return err
}

// 80ms timer tick for smooth stopwatch updates
func tickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ============================================================================
// Async Commands (tea.Cmd)
// ============================================================================

func findScriptCmd(name string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(150 * time.Millisecond)

		home, err := os.UserHomeDir()
		if err != nil {
			return scriptNotFoundMsg{err: fmt.Errorf("could not locate home directory")}
		}

		scriptsDir := filepath.Join(home, ".mxs", "scripts")
		candidates := []string{
			filepath.Join(scriptsDir, name),
			filepath.Join(scriptsDir, name+".sh"),
		}

		for _, path := range candidates {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				return scriptFoundMsg{path: path}
			}
		}

		return scriptNotFoundMsg{
			err: fmt.Errorf("script '%s' not found in %s", name, scriptsDir),
		}
	}
}

func validateScriptCmd(name string, c *cfg.Cfg, scr *scr.Scr) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join(c.FullScriptsPath, name+".sh")
		// #nosec G304 -- variable safe; uses predetermined path from Cfg and then the script name
		// then the .sh file extension
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		scriptDesc, err := scr.ExtractDesc(string(file))
		if err != nil {
			return invalidScriptMsg{err: err}
		}

		return validScriptMsg{scriptDesc: scriptDesc}
	}
}

func runScriptCmd(scriptPath string, args []string, sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		execArgs := append([]string{scriptPath}, args...)
		// #nosec G204 -- scriptPath from internal execution pipeline
		cmd := exec.Command("bash", execArgs...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return executionFinishedMsg{err: err}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return executionFinishedMsg{err: err}
		}

		if err := cmd.Start(); err != nil {
			return executionFinishedMsg{err: err}
		}

		streamReader := func(r io.Reader) {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				sub <- outputLineMsg(scanner.Text())
			}
		}

		go streamReader(stdout)
		go streamReader(stderr)

		err = cmd.Wait()
		return executionFinishedMsg{err: err}
	}
}

func waitForActivity(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// ============================================================================
// Bubble Tea Update Loop
// ============================================================================

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case scriptFoundMsg:
		m.scriptPath = msg.path
		m.state = stateRunning
		cmds = append(cmds, validateScriptCmd(m.targetName, m.cfg, m.scr))

	case invalidScriptMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit

	case validScriptMsg:
		m.state = stateRunning
		m.startTime = time.Now()
		m.scriptDesc = msg.scriptDesc
		cmds = append(cmds,
			runScriptCmd(m.scriptPath, m.scriptArgs, m.sub),
			waitForActivity(m.sub),
			tickCmd(), // Start the stopwatch
		)

	case tickMsg:
		// Re-trigger tickCmd as long as execution is running
		if m.state == stateRunning {
			cmds = append(cmds, tickCmd())
		}

	case scriptNotFoundMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit

	case outputLineMsg:
		m.formatScriptLog(msg)
		cmds = append(cmds, waitForActivity(m.sub))

	case executionFinishedMsg:
		m.state = stateFinished
		m.duration = time.Since(m.startTime)
		m.err = msg.err
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		m.ellipsisSpinner, cmd = m.ellipsisSpinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ============================================================================
// Bubble Tea View Renderer
// ============================================================================

func (m model) View() string {
	var b strings.Builder

	switch m.state {
	case stateSearching:
		fmt.Fprintf(&b, "locating %s%s\n",
			ui.ScrName.Render(m.targetName),
			m.ellipsisSpinner.View(),
		)

	case stateRunning:
		var logs string
		if len(m.output) == 0 {
			logs = m.ellipsisSpinner.View()
		} else {
			logs = ui.ScrName.Render(strings.Join(m.output, "\n"))
		}
		fmt.Fprintf(
			&b,
			"%s%s\n%s",
			ui.CreateInfoHeader("desc", ui.Primary, true),
			ui.Primary.Faint(true).Render(m.scriptDesc),
			ui.CreateInfoHeader("logs", ui.ScrName, false),
		)
		// TODO: group all of separate b.WriteString writes
		// into one [fmt.Fprintf] statements
		b.WriteString(ui.ExecBox.Render(logs))

		// Live 80ms stopwatch readout
		elapsed := time.Since(m.startTime).Round(10 * time.Millisecond)
		fmt.Fprintf(&b, "\n%s executing %s for %s",
			m.spinner.View(),
			ui.ScrName.Render(m.targetName),
			ui.Primary.Render(elapsed.String()),
		)

	case stateFinished:
		if m.err != nil {
			if len(m.output) > 0 {
				fmt.Fprintf(&b, "%s%s\n%s",
					ui.CreateInfoHeader("desc", ui.Primary, true),
					ui.Primary.Render(m.scriptDesc),
					ui.CreateInfoHeader("logs", ui.Err, false),
				)
				// TODO: group all of separate b.WriteString writes
				// into one [fmt.Fprintf] statements
				b.WriteString(ui.ErrBox.Foreground(lg.Color(ui.Red)).Render(strings.Join(m.output, "\n")))
				b.WriteString("\n")
			}

			// TODO: group all of separate b.WriteString writes
			// into one [fmt.Fprintf] statements
			b.WriteString(ui.Err.Render("✗"))
			b.WriteString(" executing ")
			b.WriteString(ui.Err.Render(m.targetName))
			b.WriteString(" failed, ")
			fmt.Fprintf(&b, "%s\n", m.err.Error())

		} else {
			if len(m.output) > 0 {
				fmt.Fprintf(&b, "%s%s\n%s",
					ui.CreateInfoHeader("desc", ui.Primary, true),
					ui.Primary.Render(m.scriptDesc),
					ui.CreateInfoHeader("logs", ui.Succ, false),
				)
				// TODO: group all of separate b.WriteString writes
				// into one [fmt.Fprintf] statements
				b.WriteString(ui.SuccBox.Render(strings.Join(m.output, "\n")))
			}

			// TODO: group all of separate b.WriteString writes
			// into one [fmt.Fprintf] statements
			b.WriteString("\n")
			b.WriteString(ui.CreateSuccMsg(fmt.Sprintf("executed %s for ", ui.Succ.Render(m.targetName))))

			b.WriteString(
				ui.Primary.Render(m.duration.Round(time.Millisecond).String(), "\n"),
			)
		}
	}

	return b.String()
}
