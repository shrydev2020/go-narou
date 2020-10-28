package controller

import (
	"narou/infrastructure/database"
	"narou/infrastructure/log"
	"narou/usecase/interactor/initialize"
)

type Initializer interface {
	Execute() error
}

type initializer struct {
	it initialize.Interactor
	lg log.Logger
	db database.DBM
}

func NewInitializeController(it initialize.Interactor, logger log.Logger, dbm database.DBM) Initializer {
	return &initializer{
		it: it,
		lg: logger,
		db: dbm,
	}
}

// Execute handle init db.
func (i initializer) Execute() error {
	i.lg.Info("Execute start")
	defer i.lg.Debug("Execute end")

	return i.it.Execute()
}
