/*
Package execui provides the execution screen for locating scripts,
streaming their output, and displaying their completion status.

The screens directory contains terminal user interfaces used by the application.
*/
package execui

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

var (
	blue             = "#6eb6ff"
	pipe             = "│"
	RegularTextColor = "#f3f3f3"

	SubMsgStyle                = lg.NewStyle().Faint(true).PaddingLeft(2)
	ScriptDescAndOutroMsgStyle = lg.NewStyle().Foreground(lg.Color(blue))
	RegularTextStyle           = lg.NewStyle().Foreground(lg.Color(RegularTextColor))
)

// addFormattedScriptMsg formats each incoming "message" according to it's order.
//
// The script description and closing message (the first and last message) are
// formatted differently from the remainder of messages (logs). The script's
// description and closing message and colored in blue.
//
// Additionally, if the current message is not the description or the closing
// message, that message will be structured and styled as a "log" before being
// appended. Each log is fainted as well and includes a pipe wall to it's left.
func addFormattedScriptMsg(m *model, msg *outputLineMsg) {
	isOutroMsg := false
	currMsgAmt := len(m.output)
	// Check for the closing message
	if currMsgAmt > 0 {
		isOutroMsg = strings.Contains(string(*msg), "✓ ")
	}
	// If the current message is the first message, the Logs= header will
	// be formatted and appended first before any logs appear.
	//
	// Additionally, if the current message is not the first, nor the last, that
	// message will be formatted as a "log" before being appended.
	//
	// Currently, a description message and closing message is needed for this
	// feature to operate as intended.
	//
	// ATTENTION: This implementation does not support closing messages within an
	// error state. Furthermore, the detection of a closing message relies on a
	// checkmark being included.
	currLogAmt := len(m.output)
	if currLogAmt != 0 && !isOutroMsg {
		// If it's the first message, append the "Msg" header
		if currLogAmt == 1 {
			m.output = append(m.output, RegularTextStyle.Faint(true).Render("  Logs="))
		}
		// If not, it's a regular message. Append the pipe and message
		m.output = append(m.output, SubMsgStyle.Render(fmt.Sprint(pipe, " ", string(*msg))))
	} else {
		// Format the script description and last message
		m.output = append(m.output, ScriptDescAndOutroMsgStyle.Render(string(*msg)))
	}
}

func InitialModel(targetName string, scriptArgs []string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.Spinner

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
		addFormattedScriptMsg(&m, &msg)
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
		fmt.Fprintf(&b, "%s Locating %s in ~/.mxs/scripts...\n",
			m.spinner.View(),
			ui.ScriptName.Render(m.targetName))

	case stateRunning:
		var logs string
		if len(m.output) == 0 {
			logs = ui.Subtle.Render("Listening for output...")
		} else {
			logs = strings.Join(m.output, "\n")
		}

		runningBoxStyle := ui.BaseOutputBox.BorderForeground(lg.Color(ui.Pink))
		b.WriteString(runningBoxStyle.Faint(true).Foreground(lg.Color(ui.Pink)).Render(logs))

		// Live 80ms stopwatch readout
		elapsed := time.Since(m.startTime).Round(10 * time.Millisecond)
		fmt.Fprintf(&b, "\n%sExecuting %s for %s",
			m.spinner.View(),
			ui.ScriptName.Render(filepath.Base(m.scriptPath)),
			ScriptDescAndOutroMsgStyle.Render(elapsed.String()),
		)

	case stateFinished:
		if m.err != nil {
			if len(m.output) > 0 {
				errBoxStyle := ui.BaseOutputBox.BorderForeground(lg.Color(ui.Red))
				b.WriteString(errBoxStyle.Foreground(lg.Color(ui.Red)).Render(strings.Join(m.output, "\n")))
				b.WriteString("\n")
			}

			b.WriteString(
				ui.Error.Render(fmt.Sprintf("✗ Executing %s ", ui.ScriptName.Render(filepath.Base(m.scriptPath)))),
			)
			b.WriteString(ui.Error.Render(fmt.Sprintf("failed, %s\n", m.err.Error())))

		} else {
			if len(m.output) > 0 {
				successBoxStyle := ui.BaseOutputBox.BorderForeground(lg.Color(ui.Lime)).Foreground(lg.Color(ui.Lime))
				b.WriteString(successBoxStyle.Render(strings.Join(m.output, "\n")))
			}

			b.WriteString(ui.Success.Render("\n✓ Executed "))
			b.WriteString(ui.ScriptName.Render(filepath.Base(m.scriptPath)))
			b.WriteString(ui.Success.Render(" for "))
			b.WriteString(
				ScriptDescAndOutroMsgStyle.Bold(true).Render(m.duration.Round(time.Millisecond).String(), "\n"),
			)
		}
	}

	return b.String()
}
