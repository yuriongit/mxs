package execScreen

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
	"github.com/yuriongit/xs/styles"
)

// ============================================================================
// Application State Enum
// ============================================================================

type state int

const (
	stateSearching state = iota // 0: Finding script in ~/.xs/scripts
	stateRunning                // 1: Streaming script output (Lime box)
	stateFinished               // 2: Completed execution or failed (Lime / Red box)
)

// ============================================================================
// Bubble Tea Message Definitions
// ============================================================================

type scriptFoundMsg struct{ path string }
type scriptNotFoundMsg struct{ err error }
type outputLineMsg string
type executionFinishedMsg struct{ err error }

// ============================================================================
// Bubble Tea Model
// ============================================================================

type model struct {
	targetName string        // Name provided by user (e.g. "gitupd")
	scriptArgs []string      // Arguments forwarded to script
	scriptPath string        // Resolved path on disk
	state      state         // Current lifecycle state
	spinner    spinner.Model // Bubbles spinner component
	output     []string      // Ring buffer holding streamed logs
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
		cmds = append(cmds, runScriptCmd(m.scriptPath, m.scriptArgs, m.sub), waitForActivity(m.sub))

	case scriptNotFoundMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit

	case outputLineMsg:
		m.output = append(m.output, string(msg))
		if len(m.output) > 12 {
			m.output = m.output[len(m.output)-12:]
		}
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
		b.WriteString(fmt.Sprintf("%s Locating %s in ~/.xs/scripts...\n",
			m.spinner.View(),
			styles.ScriptName.Render(m.targetName)))

	case stateRunning:
		argInfo := ""
		if len(m.scriptArgs) > 0 {
			argInfo = fmt.Sprintf(" (args: %s)", strings.Join(m.scriptArgs, " "))
		}
		b.WriteString(fmt.Sprintf("%s Running %s%s...\n",
			m.spinner.View(),
			styles.ScriptName.Render(filepath.Base(m.scriptPath)),
			styles.Subtle.Render(argInfo)))

		var logs string
		if len(m.output) == 0 {
			logs = styles.Subtle.Render("Listening for output...")
		} else {
			logs = strings.Join(m.output, "\n")
		}

		// Live box border is Lime while script is in progress
		runningBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.LightPink))
		b.WriteString(runningBoxStyle.Render(logs + "\n"))

	case stateFinished:
		if m.err != nil {
			b.WriteString(styles.Error.Render("✗ Execution failed\n"))
			if len(m.output) > 0 {
				// Box border turns Red on failure
				errBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.Red))
				b.WriteString(errBoxStyle.Render(strings.Join(m.output, "\n")))
			}
			b.WriteString(styles.Subtle.Render("\n" + m.err.Error() + "\n"))
		} else {
			b.WriteString(styles.Success.Render("✓ Execution completed\n"))
			if len(m.output) > 0 {
				// Box border turns Lime on success
				successBoxStyle := styles.BaseOutputBox.BorderForeground(lg.Color(styles.Green))
				b.WriteString(successBoxStyle.Render(strings.Join(m.output, "\n")))
			}
			b.WriteString(
				styles.Subtle.Render(fmt.Sprintf("\nFinished in %s", m.duration.Round(time.Millisecond)) + "\n"),
			)
		}
	}

	return b.String()
}
