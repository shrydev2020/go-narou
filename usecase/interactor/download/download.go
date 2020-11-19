package download

import (
	"context"
	"net/url"
	"regexp"
	"strconv"
	"time"

	metadataModel "narou/domain/metadata"
	"narou/domain/novel"
	query2 "narou/domain/query"
	"narou/infrastructure/log"
	"narou/interface/gateway/crawl"
	"narou/usecase/port"
	"narou/usecase/port/boudary/download"
)

type Interactor interface {
	Execute(uri metadataModel.URI) ([]string, port.ApplicationError)
}

type interactor struct {
	ctx           context.Context
	logger        log.Logger
	novelMetaRepo metadataModel.IRepository
	novelRp       novel.IRepository
	crawl         crawl.Crawler
	queryService  func(string) (query2.IQuery, error)
	outPutPort    download.OutputPorter
}

func NewDownloadInteractor(
	ctx context.Context,
	lg log.Logger,
	metaDataRepo metadataModel.IRepository,
	novelRepo novel.IRepository,
	outputPort download.OutputPorter,
	crawl crawl.Crawler,
	queryService func(string) (query2.IQuery, error)) Interactor {
	return &interactor{
		ctx:           ctx,
		logger:        lg,
		novelMetaRepo: metaDataRepo,
		novelRp:       novelRepo,
		crawl:         crawl,
		queryService:  queryService,
		outPutPort:    outputPort,
	}
}

func (uc *interactor) Execute(uri metadataModel.URI) ([]string, port.ApplicationError) {
	index, err := uc.crawl.Start(uri)
	if err != nil {
		return nil, port.NewPortError(err, port.CrawlerError)
	}

	topPage, err := uc.queryService(index)
	if err != nil {
		return nil, port.NewPortError(err, port.CrawlerError)
	}

	title := topPage.FindTitle()
	story := topPage.FindOverView()
	author := topPage.FindAuthor()
	subTitles := topPage.FindNumberOfEpisodes()

	meta, err := uc.novelMetaRepo.FindByTopURI(uri)
	if err != nil {
		return nil, port.NewPortError(err, port.RepositoryError)
	}
	if meta == nil {
		meta = metadataModel.NewMetaNovel(author, title, story, uri, subTitles)
	}

	_, err = uc.novelMetaRepo.Store(meta)
	if err != nil {
		return nil, port.NewPortError(err, port.RepositoryError)
	}
	err2 := uc.novelRp.Store(meta.SiteName,
		meta.Title,
		strconv.Itoa(0)+" "+"000　index.html",
		index)
	if err2 != nil {
		return nil, port.NewPortError(err, port.RepositoryError)
	}
	uc.logger.Info("start download")
	baseURI, _ := url.Parse(string(uri))
	subURLs := topPage.FindSubURIs()
	if meta.Length > 1 {
		_, err = uc.novelMetaRepo.Store(meta)
		if err != nil {
			return nil, port.NewPortError(err, port.RepositoryError)
		}
		// uc.outPutPort.OutPUtBoundary(novelOutputData)
		return uc.downloadSubs(baseURI, subURLs, meta)
	}
	topPage.FindBody()
	return nil, nil
}

func (uc *interactor) downloadSubs(baseURI *url.URL, subURLs []metadataModel.URI, meta *metadataModel.Novel) ([]string, port.ApplicationError) {
	downloadAt := time.Now()
	re := regexp.MustCompile("/")
	for i, sub := range subURLs {
		absURI, e := toAbsURL(baseURI, sub)
		if e != nil {
			uc.logger.Error("err occurred", "msg", e.Error())
			return nil, port.NewPortError(e, port.InvalidParam)
		}

		pageText, er := uc.crawl.Start(metadataModel.URI(absURI))
		if er != nil {
			return nil, port.NewPortError(er, port.CrawlerError)
		}

		query, err := uc.queryService(pageText)
		if err != nil {
			return nil, port.NewPortError(err, port.CrawlerError)
		}

		err2 := uc.novelRp.Store(meta.SiteName,
			meta.Title,
			strconv.Itoa(i+1)+" "+re.ReplaceAllString(query.FindEpisodeTitle(), ":")+".html",
			pageText)
		if err2 != nil {
			uc.logger.Error("error occurred when store index", "err occurred when store novel", err2.Error())
			return nil, port.NewPortError(er, port.CrawlerError)
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
