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
)

type App struct {
	HomePath  *string
	Cnf       *cnf.Cnf
	Executor  *executor.Executor
	Validator *validator.Validate
}

func NewApp() (*App, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cnf := cnf.NewCnf(homePath)
	executor := executor.NewExecutor()

	validator := &validator.Validate{}

	return &App{
		&homePath,
		cnf,
		executor,
		validator,
	}, nil
}
