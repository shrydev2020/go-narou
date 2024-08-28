package convert

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/logrusorgru/aurora"

	"narou/config"
	"narou/domain/epub"
	"narou/domain/metadata"
	"narou/domain/novel"
	"narou/domain/text_query"
)

type UseCase interface {
	Execute(ctx context.Context, uri metadata.URI) error
}

type convertUseCase struct {
	logger       *slog.Logger
	novelMetaRp  metadata.IRepository
	novelRp      novel.IRepository
	epubRp       epub.IRepository
	queryService func(string) (text_query.IQuery, error)
	cfg          config.IConfigure
}

func NewConvertUseCase(
	lg *slog.Logger,
	metaDataRepo metadata.IRepository,
	novelRepo novel.IRepository,
	epubRp epub.IRepository,
	cfg config.IConfigure,
	queryService func(string) (text_query.IQuery, error)) UseCase {

	return &convertUseCase{
		logger:       lg,
		novelMetaRp:  metaDataRepo,
		novelRp:      novelRepo,
		queryService: queryService,
		epubRp:       epubRp,
		cfg:          cfg,
	}
}

func (uc *convertUseCase) Execute(ctx context.Context, uri metadata.URI) error {
	defer uc.logger.Info("convert done", "", aurora.Cyan("Convert Done"))

	meta, err := uc.novelMetaRp.FindByTopURI(uri)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return errors.Wrapf(err, "find by top uri %s", uri)
	}
	s, err := uc.novelRp.FindByNobelSiteAndTitle(meta.SiteName, meta.Title)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return errors.Wrapf(err, "find by site %s and title %s", meta.SiteName, meta.Title)
	}
	data := uc.convertText2epubData(s, uri)

	epubModel, err := uc.newEpubData(ctx, epub.NewEpub(uc.cfg).New(ctx, meta.Author, meta.Title), data)
	if err != nil {
		return errors.Wrap(err, "create epub")
	}

	if err := uc.epubRp.Store(ctx, epubModel, meta.SiteName); err != nil {
		return errors.Wrap(err, "store epub")
	}
	return nil
}

type epubSection struct {
	chapterTitle string
	subTitle     string
	body         string
}

func (uc *convertUseCase) newEpubData(ctx context.Context, model epub.IEpub, data []epubSection) (epub.IEpub, error) {
	cssPath, err := model.AddCSS(ctx, "preset/vertical.css", "")
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, errors.Newf("AddCSS error:%v", err)
	}

	_, err2 := model.AddFont(ctx, "preset/DMincho.ttf", "")
	if err2 != nil {
		uc.logger.Error(err2.Error())
		return nil, errors.Newf("AddFont error:%v", err)
	}

	for _, d := range data {
		_, err3 := model.AddSection(ctx, d.body, d.subTitle, "", cssPath)
		if err3 != nil {
			uc.logger.Error(err3.Error())
			return nil, errors.Newf("AddSection error:%v", err)
		}
	}
	return model, nil
}

func (uc *convertUseCase) convertText2epubData(texts []string, uri metadata.URI) []epubSection {
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

func (uc *convertUseCase) newIndexPage(texts []string) epubSection {
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

func (uc *convertUseCase) newCoverPage(txt string) epubSection {
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
