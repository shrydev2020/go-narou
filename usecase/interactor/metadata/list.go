package metadata

import (
	"context"

	"narou/adapter/query"
	metadataModel "narou/domain/metadata"
	query2 "narou/domain/query"
	"narou/infrastructure/log"
	"narou/usecase/port"
	"narou/usecase/port/boudary/download"
)

type Interactor interface {
	Execute() ([]metadataModel.Novel, port.ApplicationError)
}

type interactor struct {
	ctx           context.Context
	logger        log.Logger
	novelMetaRepo metadataModel.IRepository
	queryService  func(string) (query2.IQuery, error)
	outPutPort    download.OutputPorter
}

func NewMetaDataListInteractor(
	ctx context.Context,
	lg log.Logger,
	metaDataRepo metadataModel.IRepository,
	outputPort download.OutputPorter) Interactor {
	return &interactor{
		ctx:           ctx,
		logger:        lg,
		novelMetaRepo: metaDataRepo,
		queryService:  query.New,
		outPutPort:    outputPort,
	}
}

func (uc *interactor) Execute() ([]metadataModel.Novel, port.ApplicationError) {
	data, err := uc.novelMetaRepo.FindALL()
	if err != nil {
		return nil, port.NewPortError(err, port.UnHandledError)
	}
	return data, nil
}
