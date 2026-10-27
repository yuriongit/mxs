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
	"github.com/yuriongit/mxs/internal/cnf"
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
type outputLineMsg string
type executionFinishedMsg struct{ err error }
type tickMsg time.Time // Stopwatch update tick message

// ============================================================================
// Bubble Tea Model
// ============================================================================

type model struct {
	targetName string   // Name provided by user (e.g. "gitupd")
	scriptArgs []string // Arguments forwarded to script
	scriptPath string   // Resolved path on disk
	scriptDesc string
	cnf        *cnf.Cnf
	state      state         // Current lifecycle state
	spinner    spinner.Model // Bubbles spinner component
	output     []string      // Streamed stdout/stderr log buffer
	err        error         // Execution error, if any
	startTime  time.Time     // Timer start
	duration   time.Duration // Total elapsed run time
	sub        chan tea.Msg  // Async channel for streaming output
}

// formatScriptLog formats and styles each line of output from the scirpt as a "log".
// Appends the updated slice into [model.output].
func (m *model) formatScriptLog(msg outputLineMsg) {
	m.output = append(m.output, ui.SubMsgStyle.Render(string(msg)))
}

func InitialModel(c *cnf.Cnf, targetName string, scriptArgs []string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.Spinner

	return model{
		targetName: targetName,
		scriptArgs: scriptArgs,
		cnf:        c,
		state:      stateSearching,
		spinner:    s,
		output:     make([]string, 0),
		sub:        make(chan tea.Msg),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		findScriptCmd(m.targetName),
	)
}

// Run executes the UI with reserved terminal buffer space to prevent screen scrolling
func Run(c *cnf.Cnf, targetName string, scriptArgs []string) error {
	// Pre-scroll 20 lines so the terminal never scrolls up during rendering
	fmt.Print(strings.Repeat("\n", 20) + "\033[20A")

	p := tea.NewProgram(InitialModel(c, targetName, scriptArgs))
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

type validScriptMsg struct{ scriptDesc string }

type InvalidScriptMsg struct{ err error }

func validateScriptCmd(name string, c *cnf.Cnf) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join(c.FullScriptsPath, name+".sh")
		// #nosec G304 -- variable safe; uses predetermined path from Cnf and then the script name
		// then the .sh file extension
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		scriptDesc, err := c.ExtractScriptDescription(string(file))
		if err != nil {
			return InvalidScriptMsg{err: err}
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

func infoHeader(header string, s lg.Style, wall bool) string {
	pipe := ""
	if wall {
		pipe = s.Render("│ ")
	}
	return fmt.Sprintf(
		"%s%s\n%s",
		// TODO: migrate [s (style)] parameter to [lg.Color] once dark/light mode
		// is being integrated
		s.Bold(true).Render(header),
		s.Faint(true).Render("="),
		pipe,
	)
}

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
		cmds = append(cmds, validateScriptCmd(m.targetName, m.cnf))

	case InvalidScriptMsg:
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
		fmt.Fprintf(&b, "%s locating %s in ~/.mxs/scripts...\n",
			m.spinner.View(),
			ui.ScriptName.Render(m.targetName))

	case stateRunning:
		var logs string
		if len(m.output) == 0 {
			logs = ui.Subtle.Render("listening for output...")
		} else {
			logs = strings.Join(m.output, "\n")
		}

		fmt.Fprintf(&b, "%s%s\n%s",
			infoHeader("desc", ui.ScriptDescAndOutroMsgStyle, true),	
      ui.ScriptDescAndOutroMsgStyle.Render(m.scriptDesc),
      infoHeader("logs", ui.ScriptName, false),
		)
		b.WriteString(ui.RunningBoxStyle.Faint(true).Foreground(lg.Color(ui.Pink)).Render(logs))

		// Live 80ms stopwatch readout
		elapsed := time.Since(m.startTime).Round(10 * time.Millisecond)
		fmt.Fprintf(&b, "\n%sexecuting %s for %s",
			m.spinner.View(),
			ui.ScriptName.Render(m.targetName),
			ui.ScriptDescAndOutroMsgStyle.Render(elapsed.String()),
		)

	case stateFinished:
		if m.err != nil {
			if len(m.output) > 0 {
     	  fmt.Fprintf(&b, "%s%s\n%s",
  				infoHeader("desc", ui.ScriptDescAndOutroMsgStyle, true),	
          ui.ScriptDescAndOutroMsgStyle.Render(m.scriptDesc),
          infoHeader("logs", ui.Error, false),
   			)
				b.WriteString(ui.ErrBoxStyle.Foreground(lg.Color(ui.Red)).Render(strings.Join(m.output, "\n")))
				b.WriteString("\n")
			}

			b.WriteString(ui.Error.Render("✗"))
			b.WriteString(" executing ")
			b.WriteString(ui.Error.Render(m.targetName))
			b.WriteString(" failed, ")
			fmt.Fprintf(&b, "%s\n", m.err.Error())

		} else {
			if len(m.output) > 0 {
			  fmt.Fprintf(&b, "%s%s\n%s",
  				infoHeader("desc", ui.ScriptDescAndOutroMsgStyle, true),	
          ui.ScriptDescAndOutroMsgStyle.Render(m.scriptDesc),
          infoHeader("logs", ui.Success, false),
					)
				b.WriteString(ui.SuccessBoxStyle.Render(strings.Join(m.output, "\n")))
			}

			b.WriteString("\n")
			msg, err := ui.CreateSuccessMessage(fmt.Sprintf("executed %s for ", ui.Success.Render(m.targetName)))
			if err != nil {
				return ""
			}
			b.WriteString(msg)

			b.WriteString(
				ui.ScriptDescAndOutroMsgStyle.Render(m.duration.Round(time.Millisecond).String(), "\n"),
			)
		}
	}

	return b.String()
}
