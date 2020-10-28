package metadata

import (
	"gorm.io/gorm"

	"narou/domain/metadata"
	"narou/infrastructure/database"
)

type repo struct {
	db database.DBM
}

func NewRepository(db database.DBM) metadata.IRepository {
	return &repo{
		db: db,
	}
}

// FindByURI .
func (r *repo) FindByTopURI(uri metadata.URI) (*metadata.Novel, error) {
	var ret metadata.Novel
	err := r.db.Model(&metadata.Novel{}).
		Where("top_url = ?", uri).
		First(&ret).Error()

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}

		return nil, err
	}

	return &ret, nil
}

// Initialize init metadata db.
func (r *repo) Initialize() error {
	if err := r.db.AutoMigrate(&metadata.Novel{}); err != nil {
		return err
	}

	if err := r.db.AutoMigrate(&metadata.Sub{}); err != nil {
		return err
	}

	return nil
}

// Store store metadata.
func (r *repo) Store(n *metadata.Novel) error {
	return r.db.Save(n).Error()
}

// StoreSubs store sub meta data.
func (r *repo) StoreSubs(subs []metadata.Sub) ([]metadata.Sub, error) {
	r.db.Create(&subs)
	return subs, r.db.Create(&subs).Error()
}
