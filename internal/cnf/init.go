/*
Package cnf provides all the functionality
for Xs's config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

import (
	"fmt"
	"os"
)

func (c *Cnf) InitCnfDir() error {
	directories := []struct {
		path string
		name string
	}{
		{path: c.FullCnfPath, name: fmt.Sprintf("~/%s", c.Dir)},
		{path: c.FullScriptsPath, name: fmt.Sprintf("~/%s", c.ScriptsPath)},
	}

	for _, directory := range directories {
		info, err := os.Stat(directory.path)

		switch {
		case err == nil:
			if !info.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", directory.name)
			}

			fmt.Printf("✓ Found existing %s directory\n", directory.name)

		case os.IsNotExist(err):
			if err := os.MkdirAll(directory.path, 0700); err != nil {
				return fmt.Errorf(
					"failed to create %s directory: %w",
					directory.name,
					err,
				)
			}

			fmt.Printf("✓ Created %s directory\n", directory.name)

		default:
			return fmt.Errorf("failed to check %s directory: %w", directory.name, err)
		}
	}

	return nil
}
