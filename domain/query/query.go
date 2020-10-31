package query

import "narou/domain/metadata"

type IQuery interface {
	FindTitle() string
	FindAuthor() string
	FindOverView() string
	FindNumberOfEpisodes() int
	FindSubURIs() []metadata.URI
	FindChapterTitle() string
	FindEpisodeTitle() string
	FindPreface() string
	FindAfterword() string
	FindBody() string
}
