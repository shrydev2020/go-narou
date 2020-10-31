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

	data := convertText2epubData(uc.novelRp.FindByNobelSiteAndTitle(meta.SiteName, meta.Title))

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
	chapterMap := map[string]struct{}{}
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
		if _, found := chapterMap[d.chapterTitle]; !found {
			chapterMap[d.chapterTitle] = struct{}{}
			// add chapter page
			_, err2 := model.AddSection(d.chapterTitle, d.chapterTitle, "", cssPath)
			if err2 != nil {
				uc.logger.Error(err2.Error())
				return nil, port.NewPortError(err, port.EpubError)
			}
			continue
		}
		_, err := model.AddSection(d.body, d.subTitle, "", cssPath)
		if err != nil {
			uc.logger.Error(err.Error())
			return nil, port.NewPortError(err, port.EpubError)
		}
	}

	return model, nil
}

func convertText2epubData(texts []string) []epubSection {
	currentChTitle := ""
	var data []epubSection
	for _, txt := range texts {
		d, _ := query.New(txt)
		newChTitle := d.FindChapterTitle()
		// make ch title page
		if currentChTitle == "" ||
			currentChTitle != newChTitle {

			currentChTitle = newChTitle
			data = append(data, epubSection{
				chapterTitle: currentChTitle,
				subTitle:     "",
				body:         "",
			})
		}

		data = append(data, epubSection{
			chapterTitle: currentChTitle,
			subTitle:     d.FindEpisodeTitle(),
			body:         addCSSClass(d.FindBody()),
		})
	}
	return data
}

var digitPattern = regexp.MustCompile(`(?P<key>\b(\d{1,2})\b)`)
var exclamationAndQuestionPattern = regexp.MustCompile(`(?P<key>(!!|!\?|\?|!|！！|！？|？|！))`)
var threePoint = regexp.MustCompile(`(?P<key>(…))`)

func addCSSClass(body string) string {
	return threePoint.ReplaceAllString(
		exclamationAndQuestionPattern.ReplaceAllString(
			digitPattern.ReplaceAllString(
				body, `<span class="text-combine">${key}</span>`),
			`<span class="text-combine">${key}</span>`),
		`<span class="text-three-point-reader">${key}</span>`)
}
