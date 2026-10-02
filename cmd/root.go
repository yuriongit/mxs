package cmd

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
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// ============================================================================
// Colors
// ============================================================================

var (
	grey   = "#626262"
	green  = "#9AFC8B"
	orange = "#FC8B9A"
	violet = "#8B9AFC"
)

// ============================================================================
// Colors
// ============================================================================

// ============================================================================
// Lip Gloss Styles
// ============================================================================

var (
	subtleStyle = lg.NewStyle().Faint(true)

	scriptNameStyle = lg.NewStyle().
			Foreground(lg.Color(violet)).
			Bold(true)

	outputBoxStyle = lg.NewStyle().
			Border(lg.RoundedBorder()).
			BorderForeground(lg.Color(violet)).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	successStyle = lg.NewStyle().
			Foreground(lg.Color(green)).
			Bold(true)

	errorStyle = lg.NewStyle().
			Foreground(lg.Color(orange)).
			Bold(true)

	spinnerStyle = lg.NewStyle().Foreground(lg.Color(violet))
)

// ============================================================================
// Application State Enum
// ============================================================================

type state int

const (
	stateSearching state = iota // 0: Finding script in ~/dev/scripts
	stateRunning                // 1: Streaming script output
	stateFinished               // 2: Completed execution or failed
)

// ============================================================================
// Bubble Tea Message Definitions
// ============================================================================

// Sent when script path is verified on disk
type scriptFoundMsg struct{ path string }

// Sent if script is missing or directory unreadable
type scriptNotFoundMsg struct{ err error }

// Sent each time stdout/stderr streams a new line
type outputLineMsg string

// Sent when the underlying OS process exits
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

func initialModel(targetName string, scriptArgs []string) model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = spinnerStyle

	return model{
		targetName: targetName,
		scriptArgs: scriptArgs,
		state:      stateSearching,
		spinner:    s,
		output:     make([]string, 0),
		sub:        make(chan tea.Msg),
	}
}

// Init triggers initial concurrent tasks when the program starts
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,              // Start animating the spinner
		findScriptCmd(m.targetName), // Start searching for the script file
	)
}

// ============================================================================
// Async Commands (tea.Cmd)
// ============================================================================

// Async command: Searches ~/dev/scripts for the specified target script
func findScriptCmd(name string) tea.Cmd {
	return func() tea.Msg {
		// Brief delay for smooth visual transition
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

// Async command: Spawns bash process and streams stdout/stderr via Go channel
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

		// Stream stdout & stderr lines into sub channel
		streamReader := func(r io.Reader) {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				sub <- outputLineMsg(scanner.Text())
			}
		}

		go streamReader(stdout)
		go streamReader(stderr)

		// Block until process exits
		err = cmd.Wait()
		return executionFinishedMsg{err: err}
	}
}

// Helper listener command that awaits the next message from the channel
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

	// Handle keyboard shortcuts
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	// Script located -> start execution state
	case scriptFoundMsg:
		m.scriptPath = msg.path
		m.state = stateRunning
		m.startTime = time.Now()
		cmds = append(cmds, runScriptCmd(m.scriptPath, m.scriptArgs, m.sub), waitForActivity(m.sub))

	// Script missing -> terminate
	case scriptNotFoundMsg:
		m.state = stateFinished
		m.err = msg.err
		return m, tea.Quit

	// Streamed line received from stdout/stderr
	case outputLineMsg:
		m.output = append(m.output, string(msg))
		if len(m.output) > 12 {
			m.output = m.output[len(m.output)-12:] // Maintain last 12 lines
		}
		cmds = append(cmds, waitForActivity(m.sub))

	// Process exited
	case executionFinishedMsg:
		m.state = stateFinished
		m.duration = time.Since(m.startTime)
		m.err = msg.err
		return m, tea.Quit

	// Advance spinner frame
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
			scriptNameStyle.Render(m.targetName)))

	case stateRunning:
		argInfo := ""
		if len(m.scriptArgs) > 0 {
			argInfo = fmt.Sprintf(" (args: %s)", strings.Join(m.scriptArgs, " "))
		}
		b.WriteString(fmt.Sprintf("%s Running %s%s...\n",
			m.spinner.View(),
			scriptNameStyle.Render(filepath.Base(m.scriptPath)),
			subtleStyle.Render(argInfo)))

		var logs string
		if len(m.output) == 0 {
			logs = subtleStyle.Render("Listening for output...")
		} else {
			logs = strings.Join(m.output, "\n")
		}
		b.WriteString(outputBoxStyle.Render(logs+"\n"))

	case stateFinished:
		if m.err != nil {
			b.WriteString(errorStyle.Render("✗ Execution failed\n"))
			if len(m.output) > 0 {
				b.WriteString(outputBoxStyle.Render(strings.Join(m.output, "\n")))
			}
			b.WriteString(subtleStyle.Render("\n"+m.err.Error()+"\n"))
		} else {
			b.WriteString(successStyle.Render("✓ Execution completed\n"))
			if len(m.output) > 0 {
				b.WriteString(outputBoxStyle.Render(strings.Join(m.output, "\n")))
			}
			b.WriteString(subtleStyle.Render(fmt.Sprintf("\nFinished in %s", m.duration.Round(time.Millisecond))+ "\n"))
		}
	}

	return b.String()
}

// ============================================================================
// Cobra Command Configuration
// ============================================================================

var rootCmd = &cobra.Command{
	Use:   "xs <script-name> [args...]",
	Short: "A CLI tool with AI capabilities for managing and executing scripts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		scriptArgs := args[1:]

		p := tea.NewProgram(initialModel(scriptName, scriptArgs))
		if _, err := p.Run(); err != nil {
			log.Fatal("Fatal error running xs", "err", err)
		}
	},
}

func Execute() {
	log.SetPrefix("xs")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
