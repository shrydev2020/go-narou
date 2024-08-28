package epub

import (
	"context"

	"github.com/cockroachdb/errors"

	"narou/infrastructure/storage"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Store(ctx context.Context, model IEpub, siteName string) error
}

type repo struct {
	dist string
}

func NewRepository(manager storage.Manager) IRepository {
	return &repo{
		dist: manager.GetDist(),
	}
}

func (r *repo) Store(ctx context.Context, model IEpub, siteName string) error {
	if err := model.Write(ctx, r.dist+"/"+siteName+"/"+model.GetTitle(ctx)+"/"+model.GetTitle(ctx)+".epub"); err != nil {
		return errors.Wrapf(err, `epub store fail:%s`, model.GetTitle(ctx))
	}
	return nil
}
