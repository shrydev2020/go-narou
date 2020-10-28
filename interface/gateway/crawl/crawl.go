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
	l log.Logger
}

func NewCrawler(l log.Logger) Crawler {
	return &crawl{
		l: l,
	}
}

// Start クロール.
// @return html and err
func (c crawl) Start(uri metadata.URI) (string, error) {
	res, err := http.Get(string(uri))
	if err != nil {
		return "", err
	}
	defer handleDefer(c.l, res.Body.Close)

	if res.StatusCode != 200 {
		c.l.Error("status code error", "code", res.StatusCode, "status", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	buf := bytes.NewBuffer(body)
	return buf.String(), nil
}

func handleDefer(l log.Logger, c func() error) {
	if err := c(); err != nil {
		l.Error(err.Error())
	}
}
