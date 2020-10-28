package epub

import (
	"narou/adapter/repository"
	"narou/domain/epub"
)

type repo struct {
	dist string
}

func NewRepository(dist string) epub.IRepository {
	return &repo{
		dist: dist,
	}
}

func (r *repo) Store(model epub.IEpub, siteName string) error {
	if err := model.Write(r.dist + "/" + siteName + "/" + model.GetTitle() + "/" + model.GetTitle() + ".epub"); err != nil {
		return repository.ErrorStorage
	}
	return nil
}
