package crawl

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"narou/domain/metadata"
	"narou/infrastructure/log"
)

type Crawler interface {
	Start(uri metadata.URI) (string, error)
}

type crawl struct {
	lg log.Logger
}

func NewCrawler(lg log.Logger) Crawler {
	return &crawl{
		lg: lg,
	}
}

// Start crawl return html and err
func (c crawl) Start(uri metadata.URI) (string, error) {
	res, err := http.Get(string(uri))
	if err != nil {
		return "", err
	}
	defer handleDefer(c.lg, res.Body.Close)

	if res.StatusCode != 200 {
		c.lg.Error("status code error", "code", res.StatusCode, "status", res.Status)
	}
	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		c.lg.Error("err occurred", "err", err)
	}
	buf := bytes.NewBuffer(body)
	return buf.String(), nil
}

func handleDefer(lg log.Logger, c func() error) {
	if err := c(); err != nil {
		lg.Error(err.Error())
	}
}
