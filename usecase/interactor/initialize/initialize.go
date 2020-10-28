package initialize

import (
	"narou/adapter/repository/metadata"
	metadata2 "narou/domain/metadata"
	"narou/infrastructure/database"
)

type Interactor interface {
	Execute() error
}

type interactor struct {
	novelRepo metadata2.IRepository
}

func NewInitializeInteractor(db database.DBM) Interactor {
	return &interactor{
		novelRepo: metadata.NewRepository(db),
	}
}

// Execute execute db initialization.
func (uc *interactor) Execute() error {
	return uc.novelRepo.Initialize()
}
