/*
Package app wires together the application's
configuration and execution dependencies.
*/
package app

import (
	"os"

	"github.com/yuriongit/xs/internal/cnf"
	"github.com/yuriongit/xs/internal/executor"
)

type App struct {
	HomePath *string
	Cnf      *cnf.Cnf
	Executor *executor.Executor
}

func NewApp() (*App, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cnf := cnf.NewCnf(homePath)
	executor := executor.NewExecutor()

	return &App{
		&homePath,
		cnf,
		executor,
	}, nil
}
