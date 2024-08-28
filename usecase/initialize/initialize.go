package initialize

import (
	"narou/domain/metadata"
)

type UseCase interface {
	Execute() error
}

type usecase struct {
	novelRepo metadata.IRepository
}

func NewInitUseCase(repository metadata.IRepository) UseCase {
	return &usecase{
		novelRepo: repository,
	}
}

// Execute execute db initialization.
func (uc *usecase) Execute() error {
	return uc.novelRepo.Initialize()
}
