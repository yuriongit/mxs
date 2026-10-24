// Package manager provides managerial (CRUD) operations
// for scripts.
package manager

import (
	"fmt"

	"github.com/yuriongit/mxs/internal/ui"
)

func renderListAllScriptsView(scriptsDir string) {
	fmt.Printf("%s%s\n",
		ui.CreateInfoHeader("scripts", ui.ScriptsDirStyle, false),
		scriptsDir,
	)
}
