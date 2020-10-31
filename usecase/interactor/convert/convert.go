package convert

import (
	"context"
	"regexp"

	"github.com/logrusorgru/aurora"

	"narou/adapter/query"
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
		queryService: query.New,
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

	_, err = model.AddFont("preset/DMincho.ttf", "")
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, port.NewPortError(err, port.EpubError)
	}

	for _, d := range data {
		_, err := model.AddSection(d.body, d.subTitle, "", cssPath)
		if err != nil {
			uc.logger.Error(err.Error())
			return nil, port.NewPortError(err, port.EpubError)
		}
	}

	return model, nil
}

func convertText2epubData(texts []string, uri metadataModel.URI) []epubSection {
	currentChTitle := ""
	var data []epubSection
	for _, txt := range texts {
		d, _ := query.New(txt)
		overView := d.FindOverView()
		if len(overView) > 0 {
			data = append(data, epubSection{
				chapterTitle: "",
				subTitle:     d.FindTitle(),
				body:         "<hr>" + overView + "\nRefer:<br/> <a href=\"" + string(uri) + "\" >" + string(uri) + "</a>" + "<hr>",
			})
			continue
		}
		data = append(data, epubSection{
			chapterTitle: currentChTitle,
			subTitle:     d.FindEpisodeTitle(),
			body: addCSSClass(getEpisodeTitle(d) +
				getPrefaceCSS(d) +
				d.FindBody() +
				getAfterWordCSS(d)),
		})
	}
	return data
}

func getPrefaceCSS(d narouQuery.IQuery) string {
	if len(d.FindPreface()) == 0 {
		return ""
	}
	return `<div class="episode-preface ">` + d.FindPreface() + `</div>`
}

func getAfterWordCSS(d narouQuery.IQuery) string {
	if len(d.FindPreface()) == 0 {
		return ""
	}
	return `<div class="episode-afterword ">` + d.FindPreface() + `</div>`
}

func getEpisodeTitle(d narouQuery.IQuery) string {
	return `<div class="episode-title">` +
		digitPatternTil3.ReplaceAllString(d.FindEpisodeTitle(),
			`<span class="text-combine">${key}</span>`) +
		`</div>`
}

var digitPatternTil2 = regexp.MustCompile(`(?P<key>\b(\d{1,2})\b)`)
var digitPatternTil3 = regexp.MustCompile(`(?P<key>\b(\d{1,3})\b)`)
var exclamationAndQuestionPattern = regexp.MustCompile(`(?P<key>(!!|!\?|\?|!|！！|！？|？|！))`)

func addCSSClass(body string) string {
	return exclamationAndQuestionPattern.ReplaceAllString(
		digitPatternTil2.ReplaceAllString(
			body, `<span class="text-combine">${key}</span>`),
		`<span class="text-combine">${key}</span>`)
}
