/*
Package app wires together the application's
configuration and execution dependencies.
*/
package app

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/yuriongit/mxs/internal/cnf"
	"github.com/yuriongit/mxs/internal/executor"
	"github.com/yuriongit/mxs/internal/scr"
)

type App struct {
	HomePath  *string
	Cnf       *cnf.Cnf
	Scr       *scr.Scr
	Executor  *executor.Executor
	Validator *validator.Validate
}

func NewApp() (*App, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cnf := cnf.NewCnf(homePath)
	scr := scr.New()
	executor := executor.NewExecutor()

	validator := &validator.Validate{}

	return &App{
		&homePath,
		cnf,
		scr,
		executor,
		validator,
	}, nil
}
