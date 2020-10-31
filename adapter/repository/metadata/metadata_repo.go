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
		Where("top_url like ?", "%"+uri+"%").
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
func (r *repo) Store(n *metadata.Novel) (*metadata.Novel, error) {
	d := r.db.Save(n)
	return n, d.Error()
}

// StoreSubs store sub meta data.
func (r *repo) StoreSub(sub *metadata.Sub) (*metadata.Sub, error) {
	var tmp []metadata.Sub
	d := r.db.
		Where("novel_id = ?", sub.NovelID).
		Where("index_id = ?", sub.IndexID).Find(&tmp)
	if len(tmp) == 0 {
		d = r.db.Create(sub)
	} else {
		d = r.db.
			Where("novel_id = ?", sub.NovelID).
			Where("index_id = ?", sub.IndexID).Save(sub)
	}
	return sub, d.Error()
}
