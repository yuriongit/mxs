/*
Package execUI provides the execution screen for locating scripts,
streaming their output, and displaying their completion status.

The screens directory contains terminal user interfaces used by the application.
*/
package execUI

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
	"github.com/yuriongit/xs/internal/ui/styles"
)

// ============================================================================
// Application State Enum
// ============================================================================

type state int

const (
	stateSearching state = iota // 0: Finding script in ~/.xs/scripts
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
	targetName string        // Name provided by user (e.g. "gitupd")
	scriptArgs []string      // Arguments forwarded to script
	scriptPath string        // Resolved path on disk
	state      state         // Current lifecycle state
	spinner    spinner.Model // Bubbles spinner component
	output     []string      // Streamed stdout/stderr log buffer
	err        error         // Execution error, if any
	startTime  time.Time     // Timer start
	duration   time.Duration // Total elapsed run time
	sub        chan tea.Msg  // Async channel for streaming output
}

func InitialModel(targetName string, scriptArgs []string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Spinner

	return model{
		targetName: targetName,
		scriptArgs: scriptArgs,
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
func Run(targetName string, scriptArgs []string) error {
	// Pre-scroll 20 lines so the terminal never scrolls up during rendering
	fmt.Print(strings.Repeat("\n", 20) + "\033[20A")

	p := tea.NewProgram(InitialModel(targetName, scriptArgs))
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

		scriptsDir := filepath.Join(home, ".xs", "scripts")
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

func runScriptCmd(scriptPath string, args []string, sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		execArgs := append([]string{scriptPath}, args...)
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
		m.startTime = time.Now()
		cmds = append(cmds,
			runScriptCmd(m.scriptPath, m.scriptArgs, m.sub),
			waitForActivity(m.sub),
			tickCmd(), // Start the 80ms stopwatch ticker
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
		m.output = append(m.output, string(msg))
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
		fmt.Fprintf(&b, "%s Locating %s in ~/.xs/scripts...\n",
			m.spinner.View(),
			styles.ScriptName.Render(m.targetName))

	case stateRunning:
		argInfo := ""
		if len(m.scriptArgs) > 0 {
			argInfo = fmt.Sprintf(" (args: %s)", strings.Join(m.scriptArgs, " "))
		}

		var logs string
		if len(m.output) == 0 {
			logs = styles.Subtle.Render("Listening for output...")
		} else {
			logs = strings.Join(m.output, "\n")
		}

		runningBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.Purple))
		b.WriteString(runningBoxStyle.Faint(true).Foreground(lg.Color(styles.Purple)).Render(logs + "\n..."))

		// Live 80ms stopwatch readout
		elapsed := time.Since(m.startTime).Round(10 * time.Millisecond)
		fmt.Fprintf(&b, "\n%s Running %s%s, time: %s",
			m.spinner.View(),
			styles.ScriptName.Render(filepath.Base(m.scriptPath)),
			styles.Subtle.Render(argInfo),
			elapsed,
		)

	case stateFinished:
		if m.err != nil {
			if len(m.output) > 0 {
				errBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.Red))
				b.WriteString(errBoxStyle.Foreground(lg.Color(styles.Red)).Render(strings.Join(m.output, "\n")))
			}
			b.WriteString(styles.Subtle.Render(m.err.Error()))
			b.WriteString(styles.Error.Render("\n✗ Execution failed\n"))
		} else {
			if len(m.output) > 0 {
				successBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.Green)).Foreground(lg.Color(styles.Green))
				b.WriteString(successBoxStyle.Render(strings.Join(m.output, "\n")))
			}

			b.WriteString(styles.Success.Render(fmt.Sprintf("\n✓ Execution completed, finished in %s\n", m.duration.Round(time.Millisecond))))
		}
	}

	return b.String()
}
