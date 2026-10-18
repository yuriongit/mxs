// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

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
