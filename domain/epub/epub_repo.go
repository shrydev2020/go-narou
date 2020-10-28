package epub

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Store(model IEpub, siteName string) error
}
