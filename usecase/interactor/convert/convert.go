package convert

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"

	"narou/config"
	"narou/domain/epub"
	metadataModel "narou/domain/metadata"
	"narou/domain/novel"
	"narou/domain/query"
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
	queryService func(string) (query.IQuery, error)
	cfg          config.IConfigure
}

func NewConvertInteractor(
	ctx context.Context,
	lg log.Logger,
	metaDataRepo metadataModel.IRepository,
	novelRepo novel.IRepository,
	epubRp epub.IRepository,
	cfg config.IConfigure,
	queryService func(string) (query.IQuery, error)) Interactor {

	return &interactor{
		ctx:          ctx,
		logger:       lg,
		novelMetaRp:  metaDataRepo,
		novelRp:      novelRepo,
		queryService: queryService,
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

	data := uc.convertText2epubData(uc.novelRp.FindByNobelSiteAndTitle(meta.SiteName, meta.Title), uri)

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

func (uc *interactor) convertText2epubData(texts []string, uri metadataModel.URI) []epubSection {
	var data []epubSection
	for i, txt := range texts {
		d, _ := uc.queryService(txt)
		overView := d.FindOverView()
		if i == 0 && len(overView) > 0 {
			data = append(data, epubSection{
				chapterTitle: "あらすじ",
				subTitle:     "あらすじ",
				body:         "<hr />" + overView + `<br /> Refer: <a href="` + string(uri) + `">` + string(uri) + `</a>`,
			})
			continue
		}
		data = append(data, epubSection{
			chapterTitle: d.FindChapterTitle(),
			subTitle:     d.FindEpisodeTitle(),
			body: addEpisodeTitleCSS(d.FindEpisodeTitle()) +
				newBodyWithCSS(
					newPrefaceWithCSS(d.FindPreface())+
						replaceBodyHTML2Text(d.FindBody())+
						newAfterWordWithCSS(d.FindAfterword())),
		})
	}
	tmp := uc.newIndexPage(texts)
	ret := append([]epubSection{uc.newCoverPage(texts[0]), tmp}, data...)
	return ret
}

func (uc *interactor) newIndexPage(texts []string) epubSection {
	body := `<div id="index-page"><ol>`
	for i, text := range texts {
		d, _ := uc.queryService(text)
		if i == 0 {
			// i == 0 cover, 1==index, 2 == あらすじ
			body += fmt.Sprintf(`<li><a href="section%04s.xhtml">`, strconv.Itoa(i+2)) + "INDEX" + `</a></li>`
			body += fmt.Sprintf(`<li><a href="section%04s.xhtml">`, strconv.Itoa(i+3)) + "あらすじ" + `</a></li>`
			continue
		}

		body += fmt.Sprintf(`<li><a href="section%04s.xhtml">`, strconv.Itoa(i+3)) + d.FindEpisodeTitle() + `</a></li>`
	}
	return epubSection{
		chapterTitle: "",
		subTitle:     "index",
		body:         body + "</ol></div>",
	}
}

func (uc *interactor) newCoverPage(txt string) epubSection {
	d, _ := uc.queryService(txt)
	return epubSection{
		chapterTitle: "",
		subTitle:     d.FindTitle(),
		body:         addCoverCSS(d.FindTitle()) + newAuthorWithCSS("author "+d.FindAuthor()),
	}
}

var digitPatternTil2 = regexp.MustCompile(`(?P<key>(\b\d{1,2}\b))`) // `Word boundary {1 or 2 digit} Word boundary` will hit
var digitPatternTil3 = regexp.MustCompile(`(?P<key>\b(\d{1,3})\b)`)
var singleQuotePattern = regexp.MustCompile(`"(.+?)"`)
var exclamationAndQuestionPattern = regexp.MustCompile(`(?P<key>(!!|!\?|\?|!|！！|！？|？|！))`)

func addCoverCSS(cover string) string {
	return `<div class="cover">` + cover + `</div>`
}

func newBodyWithCSS(body string) string {
	tmp := digitPatternTil2.ReplaceAllString(
		body, `<span class="text-combine">${key}</span>`)
	return exclamationAndQuestionPattern.ReplaceAllString(
		tmp, `<span class="text-combine">${key}</span>`)
}

func newAuthorWithCSS(author string) string {
	return `<div id="author" class="author">` + author + `</div>`
}

func newPrefaceWithCSS(preface string) string {
	if len(preface) == 0 {
		return ""
	}
	return `<div id="preface=%s" class="episode-preface">` + preface + `</div>`
}

func newAfterWordWithCSS(aw string) string {
	if len(aw) == 0 {
		return ""
	}
	return `<div id="afterword-%s" class="episode-afterword">` + aw + `</div>`
}

func addEpisodeTitleCSS(title string) string {
	return `<div id="episode-title" class="episode-title">` +
		digitPatternTil3.ReplaceAllString(title, `<span class="text-combine">${key}</span>`) +
		`</div>`
}

func replaceBodyHTML2Text(s string) string {
	s = strings.ReplaceAll(s, "&#34;", `"`) // &#34; -> "
	// FIXME &#XX;
	s = convertDoubleQuote2DoubleInDesignDoubleQuote(s) // "aaa" -> 『aaa』
	return s
}

func convertDoubleQuote2DoubleInDesignDoubleQuote(s string) string {
	return singleQuotePattern.ReplaceAllString(s, `『$1』`)
}
