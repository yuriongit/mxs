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
	"sync"
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
	cnf             *cnf.Cnf
	targetName      string   // Name provided by user (e.g. "gitupd")
	scriptArgs      []string // Arguments forwarded to script
	scriptPath      string   // Resolved path on disk
	scriptDesc      string
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
	m.output = append(m.output, ui.SubMsgStyle.Render(string(msg)))
}

var ScriptDescStyle = ui.ScriptDescAndOutroMsgStyle

func InitialModel(c *cnf.Cnf, targetName string, scriptArgs []string) model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = ui.Spinner

	ellipsisSpinner := spinner.New()
	ellipsisSpinner.Spinner = spinner.Ellipsis
	ellipsisSpinner.Style = ui.Subtle

	return model{
		cnf:             c,
		targetName:      targetName,
		scriptArgs:      scriptArgs,
		scriptPath:      "",
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
		findScriptCmd(m.cnf, m.targetName),
	)
}

// Run runs the script execution UI

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

func findScriptCmd(c *cnf.Cnf, name string) tea.Cmd {
	return func() tea.Msg {
		// time.Sleep(250 * time.Millisecond)
		path := filepath.Join(c.FullScriptsPath, name+".sh")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return scriptNotFoundMsg{
				err: fmt.Errorf("script %s not found", name),
			}
		}
		return scriptFoundMsg{path: path}
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
		if scriptPath == "" {
			return executionFinishedMsg{
				err: fmt.Errorf("cannot execute script: script path is empty"),
			}
		}

		var cleanArgs []string
		for _, arg := range args {
			if arg != "" {
				cleanArgs = append(cleanArgs, arg)
			}
		}

		execArgs := append([]string{scriptPath}, cleanArgs...)

		// #nosec G204 -- script name without the file extension and script arguments
		cmd := exec.Command(execArgs[0], execArgs[1:]...)
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

		var wg sync.WaitGroup
		wg.Add(2)

		scanStream := func(r io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				sub <- outputLineMsg(scanner.Text())
			}
		}

		go scanStream(stdout)
		go scanStream(stderr)

		// Wait for both output streams to flush completely
		wg.Wait()

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
		cmds = append(cmds, validateScriptCmd(m.targetName, m.cnf))
	case tickMsg:
		// Re-trigger tickCmd as long as execution is running
		if m.state == stateRunning {
			// if len(m.output) == 1 {
			// 	m.output = append(m.output, "logs")
			// b.WriteString("=)
			// }
			cmds = append(cmds, tickCmd())
		}
	case scriptNotFoundMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit
	case InvalidScriptMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit
	case validScriptMsg:
		m.scriptDesc = msg.scriptDesc
		m.state = stateRunning
		m.startTime = time.Now()
		cmds = append(cmds,
			runScriptCmd(m.scriptPath, m.scriptArgs, m.sub),
			waitForActivity(m.sub),
			tickCmd(), // Start the stopwatch
		)

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
			ui.ScriptName.Render(m.targetName),
			m.ellipsisSpinner.View(),
		)

	case stateRunning:
		var logs string
		if len(m.output) == 0 {
			logs = m.ellipsisSpinner.View()
		} else {
			logs = ui.ScriptName.Render(strings.Join(m.output, "\n"))
		}
		fmt.Fprintf(
			&b,
			"%s%s\n%s",
			ui.CreateInfoHeader("desc", ScriptDescStyle, true),
			ScriptDescStyle.Faint(true).Render(m.scriptDesc),
			ui.CreateInfoHeader("logs", ui.ScriptName, false),
		)
		b.WriteString(ui.RunningBoxStyle.Render(logs))

		// Live 80ms stopwatch readout
		elapsed := time.Since(m.startTime).Round(10 * time.Millisecond)
		fmt.Fprintf(&b, "\n%s executing %s for %s",
			m.spinner.View(),
			ui.ScriptName.Bold(true).Render(m.targetName),
			ui.ScriptDescAndOutroMsgStyle.Render(elapsed.String()),
		)

	case stateFinished:
		if m.err != nil {
			if len(m.output) >= 1 {
				ui.CreateInfoHeader("desc", ScriptDescStyle, true)
				fmt.Fprintf(
					&b,
					"%s%s\n%s",
					ui.CreateInfoHeader("desc", ScriptDescStyle, true),
					ScriptDescStyle.Faint(true).Render(m.scriptDesc),
					ui.CreateInfoHeader("logs", ui.Error, false),
				)
				b.WriteString(ui.ErrBoxStyle.
					Foreground(ui.Error.GetForeground()).
					Render(strings.Join(m.output, "\n")),
				)
				b.WriteString(ui.Error.Render("\n"))
			}
			b.WriteString(ui.Error.Render("✗"))
			b.WriteString(" executing ")
			b.WriteString(ui.Error.Render(m.targetName))
			b.WriteString(" failed, ")
			fmt.Fprintf(&b, "%s\n", m.err.Error())
		} else {
			if len(m.output) >= 1 {
				fmt.Fprintf(
					&b,
					"%s%s\n%s",
					ui.CreateInfoHeader("desc", ScriptDescStyle, true),
					ScriptDescStyle.Faint(true).Render(m.scriptDesc),
					ui.CreateInfoHeader("logs", ui.Success, false),
				)
				b.WriteString(ui.SuccessBoxStyle.
					Foreground(lg.Color(ui.Lime)).
					Render(strings.Join(m.output, "\n")),
				)
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
