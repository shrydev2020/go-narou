package download

import (
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"

	metadataModel "narou/domain/metadata"
	"narou/domain/novel"
	"narou/domain/text_query"
	"narou/interface/gateway/crawl"
)

type UseCase interface {
	Execute(uri metadataModel.URI) ([]string, error)
}

type usecase struct {
	logger        *slog.Logger
	novelMetaRepo metadataModel.IRepository
	novelRp       novel.IRepository
	crawl         crawl.Crawler
	queryService  func(string) (text_query.IQuery, error)
}

func NewDownloadUseCase(
	lg *slog.Logger,
	metaDataRepo metadataModel.IRepository,
	novelRepo novel.IRepository,
	crawl crawl.Crawler,
	queryService func(string) (text_query.IQuery, error)) UseCase {
	return &usecase{
		logger:        lg,
		novelMetaRepo: metaDataRepo,
		novelRp:       novelRepo,
		crawl:         crawl,
		queryService:  queryService,
	}
}

func (uc *usecase) Execute(uri metadataModel.URI) ([]string, error) {
	index, err := uc.crawl.Start(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start crawler")
	}

	topPage, err := uc.queryService(index)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query top page")
	}

	title := topPage.FindTitle()
	story := topPage.FindOverView()
	author := topPage.FindAuthor()
	subTitles := topPage.FindNumberOfEpisodes()

	meta, err := uc.novelMetaRepo.FindByTopURI(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find metadata")
	}
	if meta == nil {
		meta = metadataModel.NewMetaNovel(author, title, story, uri, subTitles)
	}

	_, err = uc.novelMetaRepo.Store(meta)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to store metadata")
	}
	err2 := uc.novelRp.Store(meta.SiteName,
		meta.Title,
		strconv.Itoa(0)+" "+"000 index.html",
		index)
	if err2 != nil {
		return nil, errors.Wrapf(err, "failed to store metadata")
	}
	uc.logger.Info("start download")
	baseURI, _ := url.Parse(string(uri))
	subURLs := topPage.FindSubURIs()
	if meta.Length > 1 {
		_, err = uc.novelMetaRepo.Store(meta)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to store metadata")
		}
		return uc.downloadSubs(baseURI, subURLs, meta)
	}
	topPage.FindBody()
	return nil, nil
}

func (uc *usecase) downloadSubs(baseURI *url.URL, subURLs []metadataModel.URI, meta *metadataModel.Novel) ([]string, error) {
	downloadAt := time.Now()
	re := regexp.MustCompile("/")
	for i, sub := range subURLs {
		absURI, e := toAbsURL(baseURI, sub)
		if e != nil {
			uc.logger.Error("err occurred", "msg", e.Error())
			return nil, errors.Wrapf(e, "failed to parse sub url")
		}

		pageText, er := uc.crawl.Start(metadataModel.URI(absURI))
		if er != nil {
			return nil, errors.Wrapf(er, "failed to start crawler")
		}

		query, err := uc.queryService(pageText)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start crawler")
		}

		err2 := uc.novelRp.Store(meta.SiteName,
			meta.Title,
			strconv.Itoa(i+1)+" "+re.ReplaceAllString(query.FindEpisodeTitle(), ":")+".html",
			pageText)
		if err2 != nil {
			uc.logger.Error("error occurred when store index", "err occurred when store novel", err2.Error())
			return nil, errors.Wrapf(err2, "failed to start crawler")
		}
		if _, err3 := uc.novelMetaRepo.StoreSub(&metadataModel.Sub{
			NovelID:      meta.ID,
			IndexID:      i,
			Href:         sub,
			Chapter:      query.FindChapterTitle(),
			Subtitle:     query.FindEpisodeTitle(),
			SubDate:      downloadAt,
			SubUpDatedAt: downloadAt,
			DownloadAt:   downloadAt}); err3 != nil {
			uc.logger.Error("err occurred when save sub data", "err", err3.Error())
		}
		uc.logger.Info(query.FindEpisodeTitle(), "total", len(subURLs), "current", i+1)
		// todo from config
		time.Sleep(getSec())
	}
	return nil, nil
}

const oneSec = 2

func getSec() time.Duration {
	return time.Second * oneSec
}

func toAbsURL(baseURL *url.URL, uri metadataModel.URI) (string, error) {
	parsedURL, err := url.Parse(string(uri))
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(parsedURL).String(), nil
}
