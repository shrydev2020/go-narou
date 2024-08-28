package metadata

import (
	"context"
	"log/slog"

	"github.com/cockroachdb/errors"

	"narou/domain/metadata"
	"narou/domain/text_query"
)

type UseCase interface {
	Execute(ctx context.Context) ([]metadata.Novel, error)
}

type usecase struct {
	novelMetaRepo metadata.IRepository
	queryService  func(string) (text_query.IQuery, error)
	logger        *slog.Logger
}

func NewMetaDataListUseCase(
	lg *slog.Logger,
	metaDataRepo metadata.IRepository,
	queryService func(string) (text_query.IQuery, error),
) UseCase {
	return &usecase{
		logger:        lg,
		novelMetaRepo: metaDataRepo,
		queryService:  queryService,
	}
}

func (uc *usecase) Execute(ctx context.Context) ([]metadata.Novel, error) {
	data, err := uc.novelMetaRepo.FindALL()
	if err != nil {
		return nil, errors.Wrapf(err, "novel metadata list failed")
	}
	return data, nil
}
