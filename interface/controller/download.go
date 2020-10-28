package controller

import (
	"narou/domain/metadata"
	"narou/infrastructure/log"
	downloadIc "narou/usecase/interactor/download"
)

type Downloader interface {
	Execute(args []string) ([]string, error)
}

type download struct {
	it     downloadIc.Interactor
	logger log.Logger
}

func NewDownloadController(it downloadIc.Interactor, logger log.Logger) Downloader {
	return &download{
		it:     it,
		logger: logger,
	}
}

func (d download) Execute(args []string) ([]string, error) {
	d.logger.Info("Execute -- start")
	defer d.logger.Info("Execute -- end")

	uri := metadata.URI(args[0])
	ret, err := d.it.Execute(uri)

	if err != nil {
		return nil, err
	}

	return ret, nil
}
