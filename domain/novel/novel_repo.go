package novel

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Store(novelType, title, chapterTitle, body string) error
	FindByNobelSiteAndTitle(novelSite, title string) []string
}
