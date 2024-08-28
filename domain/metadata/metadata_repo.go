package metadata

import (
	"gorm.io/gorm"

	"narou/infrastructure/database"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Initialize() error
	Store(*Novel) (*Novel, error)
	StoreSub(*Sub) (*Sub, error)
	FindByTopURI(uri URI) (*Novel, error)
	FindALL() ([]Novel, error)
}
type repo struct {
	db database.DBM
}

func NewRepository(db database.DBM) IRepository {
	return &repo{
		db: db,
	}
}

// FindByURI .
func (r *repo) FindByTopURI(uri URI) (*Novel, error) {
	var ret Novel
	err := r.db.Model(&Novel{}).
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
	if err := r.db.AutoMigrate(&Novel{}); err != nil {
		return err
	}

	if err := r.db.AutoMigrate(&Sub{}); err != nil {
		return err
	}

	return nil
}

// Store store metadata.
func (r *repo) Store(n *Novel) (*Novel, error) {
	d := r.db.Save(n)
	return n, d.Error()
}

// StoreSubs store sub meta data.
func (r *repo) StoreSub(sub *Sub) (*Sub, error) {
	var tmp []Sub
	var err error
	r.db.Where("novel_id = ?", sub.NovelID).
		Where("index_id = ?", sub.IndexID).Find(&tmp)

	if len(tmp) == 0 {
		err = r.db.Create(sub).Error()
		return sub, err
	}

	err = r.db.
		Where("novel_id = ?", sub.NovelID).
		Where("index_id = ?", sub.IndexID).Save(sub).Error()

	return sub, err
}

func (r *repo) FindALL() ([]Novel, error) {
	var tmp []Novel
	if err := r.db.Preload("Subs").Find(&tmp).Error(); err != nil {
		panic(err)
	}
	return tmp, nil
}
