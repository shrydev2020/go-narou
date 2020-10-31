package epub

import (
	"narou/adapter/repository"
	"narou/domain/epub"
	"narou/infrastructure/storage"
)

type repo struct {
	dist string
}

func NewRepository(manager storage.Manager) epub.IRepository {
	return &repo{
		dist: manager.GetDist(),
	}
}

func (r *repo) Store(model epub.IEpub, siteName string) error {
	if err := model.Write(r.dist + "/" + siteName + "/" + model.GetTitle() + "/" + model.GetTitle() + ".epub"); err != nil {
		return repository.ErrorStorage
	}
	return nil
}
