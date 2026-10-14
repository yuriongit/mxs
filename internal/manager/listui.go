// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/yuriongit/mxs/internal/ui"
)

var (
	scriptsDirStyle = lg.NewStyle().
		Foreground(lg.Color(ui.Pink)).
		PaddingLeft(0).
		Bold(true)
)

func renderListAllScriptsView(scriptsDirName string, scriptsDir string) {
	fmt.Printf("%s\n%s\n",
		scriptsDirStyle.Render("~/"+scriptsDirName+"/"),
		scriptsDir,
	)
}
