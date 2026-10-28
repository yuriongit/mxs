// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"fmt"

	"github.com/yuriongit/mxs/internal/ui"
)

func renderListAllScriptsView(scriptsDirName string, scriptsDir string) {
	fmt.Printf("%s\n%s\n",
		ui.ScriptsDirStyle.Render("~/"+scriptsDirName+"/"),
		scriptsDir,
	)
}
