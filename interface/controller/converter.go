package controller

import (
	"narou/domain/metadata"
	"narou/infrastructure/log"
	"narou/usecase/interactor/convert"
)

type Converter interface {
	Execute(args []string) error
}

type converter struct {
	it     convert.Interactor
	logger log.Logger
}

func NewConvertController(it convert.Interactor, logger log.Logger) Converter {
	return &converter{
		it:     it,
		logger: logger,
	}
}

func (c *converter) Execute(args []string) error {
	c.logger.Info("Execute -- start")
	defer c.logger.Info("Execute -- end")

	uri := metadata.URI(args[0])
	if err := c.it.Execute(uri); err != nil {
		c.logger.Error(err.Error())
	}

	return nil
}
