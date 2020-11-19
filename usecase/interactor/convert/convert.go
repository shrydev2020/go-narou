package convert

import (
	"context"
	"fmt"
	"regexp"

	"narou/adapter/query/narou"

	"github.com/logrusorgru/aurora"

	"narou/config"
	"narou/domain/epub"
	metadataModel "narou/domain/metadata"
	"narou/domain/novel"
	narouQuery "narou/domain/query"
	"narou/infrastructure/log"
	"narou/usecase/port"
)

type Interactor interface {
	Execute(uri metadataModel.URI) port.ApplicationError
}

type interactor struct {
	ctx          context.Context
	logger       log.Logger
	novelMetaRp  metadataModel.IRepository
	novelRp      novel.IRepository
	epubRp       epub.IRepository
	queryService func(string) (narouQuery.IQuery, error)
	cfg          config.IConfigure
}

func NewConvertInteractor(
	ctx context.Context,
	lg log.Logger,
	metaDataRepo metadataModel.IRepository,
	novelRepo novel.IRepository,
	epubRp epub.IRepository,
	cfg config.IConfigure) Interactor {

	return &interactor{
		ctx:          ctx,
		logger:       lg,
		novelMetaRp:  metaDataRepo,
		novelRp:      novelRepo,
		queryService: narou.New,
		epubRp:       epubRp,
		cfg:          cfg,
	}
}

func (uc *interactor) Execute(uri metadataModel.URI) port.ApplicationError {
	defer uc.logger.Info("convert done", "", aurora.Cyan("Convert Done"))

	meta, err := uc.novelMetaRp.FindByTopURI(uri)
	if err != nil {
		uc.logger.Error(err.Error())
		return port.NewPortError(err, port.RepositoryError)
	}

	data := convertText2epubData(uc.novelRp.FindByNobelSiteAndTitle(meta.SiteName, meta.Title), uri)

	epubModel, err := uc.newEpubData(epub.NewEpub(uc.ctx, uc.cfg).New(meta.Author, meta.Title), data)

	if err != nil {
		return port.NewPortError(err, port.EpubError)
	}

	if err := uc.epubRp.Store(epubModel, meta.SiteName); err != nil {
		return port.NewPortError(err, port.RepositoryError)
	}

	return nil
}

type epubSection struct {
	chapterTitle string
	subTitle     string
	body         string
}

func (uc *interactor) newEpubData(model epub.IEpub, data []epubSection) (epub.IEpub, error) {
	cssPath, err := model.AddCSS("preset/vertical.css", "")
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, port.NewPortError(err, port.EpubError)
	}

	_, err2 := model.AddFont("preset/DMincho.ttf", "")
	if err2 != nil {
		uc.logger.Error(err2.Error())
		return nil, port.NewPortError(err, port.EpubError)
	}

	for _, d := range data {
		_, err3 := model.AddSection(d.body, d.subTitle, "", cssPath)
		if err3 != nil {
			uc.logger.Error(err3.Error())
			return nil, port.NewPortError(err, port.EpubError)
		}
	}

	return model, nil
}

func convertText2epubData(texts []string, uri metadataModel.URI) []epubSection {
	var data []epubSection
	for _, txt := range texts {
		d, _ := narou.New(txt)
		overView := d.FindOverView()
		if len(overView) > 0 {
			data = append(data, epubSection{
				chapterTitle: "あらすじ",
				subTitle:     "あらすじ",
				body:         "<hr/>" + overView + `<br/>Refer: <a href="` + string(uri) + `">` + string(uri) + `</a><hr/>`,
			})
			continue
		}
		data = append(data, epubSection{
			chapterTitle: d.FindChapterTitle(),
			subTitle:     d.FindEpisodeTitle(),
			body: addCSSClass(getEpisodeTitle(d) +
				getPrefaceCSS(d) +
				d.FindBody() +
				getAfterWordCSS(d)),
		})
	}

	ret := append([]epubSection{newCoverPage(texts[0])}, data...)
	return ret
}

func newCoverPage(txt string) epubSection {
	d, _ := narou.New(txt)
	return epubSection{
		chapterTitle: "",
		subTitle:     d.FindTitle(),
		body:         addCoverCSS(d.FindTitle()) + addAuthorCSS("author "+d.FindAuthor()),
	}
}

func addAuthorCSS(author string) string {
	return `<div id="author" class="author">` + author + `</div>`
}

func getPrefaceCSS(d narouQuery.IQuery) string {
	p := d.FindPreface()
	if len(p) == 0 {
		return ""
	}
	return `<div id="preface=%s" class="episode-preface">` + p + `</div>`
}

func getAfterWordCSS(d narouQuery.IQuery) string {
	aw := d.FindAfterword()
	if len(aw) == 0 {
		return ""
	}
	return `<div id="afterword-%s" class="episode-afterword">` + aw + `</div>`
}

func getEpisodeTitle(d narouQuery.IQuery) string {
	return `<div id="episode-title" class="episode-title">` +
		digitPatternTil3.ReplaceAllString(d.FindEpisodeTitle(),
			`<span class="text-combine">${key}</span>`) +
		`</div>`
}

var digitPatternTil2 = regexp.MustCompile(`(?P<key>\b(\d{1,2})\b)`)
var digitPatternTil3 = regexp.MustCompile(`(?P<key>\b(\d{1,3})\b)`)
var exclamationAndQuestionPattern = regexp.MustCompile(`(?P<key>(!!|!\?|\?|!|！！|！？|？|！))`)

func addCoverCSS(cover string) string {
	return `<div class="cover">` + cover + `</div>`
}

func addCSSClass(body string) string {
	tmp := digitPatternTil2.ReplaceAllString(
		body, fmt.Sprintf(`<span class="text-combine">${key}</span>`))
	return exclamationAndQuestionPattern.ReplaceAllString(
		tmp, `<span class="text-combine">${key}</span>`)
}
