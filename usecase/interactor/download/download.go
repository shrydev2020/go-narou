package download

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"narou/adapter/query"
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
	logger        log.Logger
	novelMetaRepo metadataModel.IRepository
	novelRp       novel.IRepository
	crawl         crawl.Crawler
	queryService  func(string) (query2.IQuery, error)
	outPutPort    download.OutputPorter
}

func NewDownloadInteractor(
	lg log.Logger,
	metaDataRepo metadataModel.IRepository,
	novelRepo novel.IRepository,
	outputPort download.OutputPorter,
	crawl crawl.Crawler) Interactor {
	return &interactor{
		logger:        lg,
		novelMetaRepo: metaDataRepo,
		novelRp:       novelRepo,
		crawl:         crawl,
		queryService:  query.New,
		outPutPort:    outputPort,
	}
}

func (uc *interactor) Execute(uri metadataModel.URI) ([]string, port.ApplicationError) {
	meta, err := uc.novelMetaRepo.FindByTopURI(uri)
	if err != nil {
		return nil, port.NewPortError(err, port.RepositoryError)
	}

	index, err := uc.crawl.Start(uri)
	if err != nil {
		return nil, port.NewPortError(err, port.CrawlerError)
	}

	topPage, err := query.New(index)
	if err != nil {
		return nil, port.NewPortError(err, port.CrawlerError)
	}

	title := topPage.FindTitle()
	story := topPage.FindOverView()
	author := topPage.FindAuthor()
	subTitles := topPage.FindNumberOfEpisodes()
	subURLs := topPage.FindSubURIs()

	if meta == nil {
		meta = metadataModel.NewMetaNovel(author, title, story, uri, subTitles)
	}

	err = uc.novelMetaRepo.Store(meta)
	// err = uc.novelMetaRepo.StoreSubs(subs)
	if err != nil {
		return nil, port.NewPortError(err, port.RepositoryError)
	}

	baseURI, _ := url.Parse(string(uri))

	fmt.Print("start download\n")
	for i, sub := range subURLs {
		u, e := toAbsURL(baseURI, sub)
		if e != nil {
			uc.logger.Error(e.Error())
			return nil, port.NewPortError(e, port.InvalidParam)
		}

		pageText, er := uc.crawl.Start(metadataModel.URI(u))
		if er != nil {
			return nil, port.NewPortError(er, port.CrawlerError)
		}

		d, err := query.New(pageText)
		if err != nil {
			return nil, port.NewPortError(err, port.CrawlerError)
		}

		err2 := uc.novelRp.Store(meta.SiteName,
			title,
			strconv.Itoa(i+1)+" "+d.FindEpisodeTitle()+".html",
			pageText)
		if err2 != nil {
			uc.logger.Error("error occurred when store index", "err occurred when store novel", err2.Error())
			return nil, port.NewPortError(er, port.CrawlerError)
		}

		fmt.Printf("%s	(%d/%d) done\n", d.FindEpisodeTitle(), i+1, len(subURLs))
		// todo from config
		time.Sleep(getSec())
	}
	// uc.outPutPort.OutPUtBoundary(novelOutputData)
	return nil, nil
}

const oneSec = 1

func getSec() time.Duration {
	return time.Second * oneSec
}

func toAbsURL(baseURL *url.URL, weburl metadataModel.URI) (string, error) {
	parsedURL, err := url.Parse(string(weburl))

	if err != nil {
		return "", err
	}

	absURL := baseURL.ResolveReference(parsedURL)

	return absURL.String(), nil
}
