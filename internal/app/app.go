/*
Package app wires together the application's
configuration and execution dependencies.
*/
package app

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/cfg"
	"github.com/yuriongit/mxs/internal/executor"
	"github.com/yuriongit/mxs/internal/scr"
)

// App holds MXS' core dependencies.
type App struct {
	HomePath  *string
	Cfg       *cfg.Cfg
	Scr       *scr.Scr
	Executor  *executor.Executor
	Validator *validator.Validate
}

// NewApp instantiates a new [App] struct.
//
// Returns the pointer to the created [App] struct.
func NewApp() (*App, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg := cfg.NewCfg(homePath)
	scr := scr.New()
	executor := executor.NewExecutor()

	validator := &validator.Validate{}

	return &App{
		&homePath,
		cfg,
		scr,
		executor,
		validator,
	}, nil
}
