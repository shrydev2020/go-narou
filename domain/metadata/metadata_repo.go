package metadata

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Initialize() error
	Store(*Novel) error
	StoreSubs([]Sub) ([]Sub, error)
	FindByTopURI(uri URI) (*Novel, error)
}
