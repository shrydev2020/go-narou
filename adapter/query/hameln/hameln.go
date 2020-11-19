package hameln

import (
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"narou/domain/metadata"
	"narou/domain/query"
)

type hamelnQuery struct {
	d *goquery.Document
}

func New(html string) (query.IQuery, error) {
	cls := ioutil.NopCloser(strings.NewReader(html))
	d, err := goquery.NewDocumentFromReader(cls)

	if err != nil {
		return &hamelnQuery{}, err
	}

	return &hamelnQuery{
		d: d,
	}, nil
}

// FindTitle return Novel Tile.
func (n *hamelnQuery) FindTitle() string {
	return n.d.Find("span[itemprop=name]").Text()
}

// FindAuthor return author name.
func (n *hamelnQuery) FindAuthor() string {
	return n.d.Find("span[itemprop=author]").Text()
}

// FindOverView return novel's over view.
func (n *hamelnQuery) FindOverView() string {
	story := ""
	n.d.Find("div.ss").Each(
		func(i int, selection *goquery.Selection) {
			if i == 1 {
				story = selection.Text()
			}
		})
	return story
}

// FindNumberOfEpisodes return how many episodes are there.
func (n *hamelnQuery) FindNumberOfEpisodes() int {
	ret := 0

	n.d.Find("tbody > tr > td > a").
		Each(func(i int, selection *goquery.Selection) {
			ret++
		})

	return ret
}

// FindSubURIs return episode uri like "/n5378gc/"
func (n *hamelnQuery) FindSubURIs() []metadata.URI {
	var subs []metadata.URI

	n.d.Find("tbody > tr > td > a").
		Each(func(i int, selection *goquery.Selection) {
			uri, _ := selection.Attr("href")
			subs = append(subs, metadata.URI(uri))
		})

	return subs
}

// FindChapterTitle return chapter tile.
func (n *hamelnQuery) FindChapterTitle() string {
	title := ""
	n.d.Find("span[style=\"font-size:120%\"]").Each(
		func(i int, selection *goquery.Selection) {
			if i == 1 {
				title = selection.Text()
			}
		})
	return title
}

// FindEpisodeTitle return each episode title.
func (n *hamelnQuery) FindEpisodeTitle() string {
	title := ""
	n.d.Find("span[style=\"font-size:120%\"]").Each(
		func(i int, selection *goquery.Selection) {
			if i == 1 {
				title = selection.Text()
			}
		})
	return title
}

// FindBody return novel body.
func (n *hamelnQuery) FindBody() string {
	b := ""

	n.d.Find("div#honbun").
		Find("p").
		Each(func(i int, selection *goquery.Selection) {
			tmp, _ := selection.Html()
			if tmp == "<br/>" {
				return
			}
			b += "<p>" + tmp + "</p>"
		})

	return b
}

// FindPreface return episode preface.
func (n *hamelnQuery) FindPreface() string {
	b := ""
	b = n.d.Find("#maegaki").Text()
	return b
}

// FindAfterword return episode afterword.
func (n *hamelnQuery) FindAfterword() string {
	b := ""
	n.d.Find(".novel_view#novel_a").Each(func(_ int, selection *goquery.Selection) {
		b += selection.Text()
	})
	return b
}
