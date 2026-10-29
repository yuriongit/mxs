// Package scr manages script data, creation, validation, and related operations.
package scr

import (
	"fmt"

	"github.com/yuriongit/mxs/internal/ui"
)

func renderListAllScriptsView(scriptsDir string) {
	fmt.Printf("%s%s\n",
		ui.CreateInfoHeader("scripts", ui.Primary, false),
		scriptsDir,
	)
}
