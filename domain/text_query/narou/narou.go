package narou

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"narou/domain/metadata"
	"narou/domain/text_query"
)

type narouQuery struct {
	d *goquery.Document
}

func NewNarouQuery(html string) (text_query.IQuery, error) {
	cls := io.NopCloser(strings.NewReader(html))
	d, err := goquery.NewDocumentFromReader(cls)

	if err != nil {
		return &narouQuery{}, err
	}

	return &narouQuery{
		d: d,
	}, nil
}

// FindTitle return Novel Tile.
func (n *narouQuery) FindTitle() string {
	return n.d.Find(".novel_title").Text()
}

// FindAuthor return author name.
func (n *narouQuery) FindAuthor() string {
	return n.d.Find(".novel_writername > a").Text()
}

// FindOverView return novel's over view.
func (n *narouQuery) FindOverView() string {
	b := ""

	n.d.Find("div#novel_ex").Each(func(i int, selection *goquery.Selection) {
		tmp, _ := selection.Html()
		if tmp == "<br/>" {
			return
		}
		b += "<p>" + tmp + "</p>"
	})
	return b
}

// FindNumberOfEpisodes return how many episodes are there.
func (n *narouQuery) FindNumberOfEpisodes() int {
	ret := 0

	n.d.Find(".novel_sublist2").
		Each(func(i int, selection *goquery.Selection) {
			ret++
		})

	return ret
}

// FindSubURIs return episode uri like "/n5378gc/"
func (n *narouQuery) FindSubURIs() []metadata.URI {
	var subs []metadata.URI

	n.d.Find(".novel_sublist2").
		Each(func(i int, selection *goquery.Selection) {
			uri, _ := selection.Find(".subtitle  a").Attr("href")
			subs = append(subs, metadata.URI(uri))
		})

	return subs
}

// FindChapterTitle return chapter tile.
func (n *narouQuery) FindChapterTitle() string {
	return n.d.Find(".chapter_title").Text()
}

// FindEpisodeTitle return each episode title.
func (n *narouQuery) FindEpisodeTitle() string {
	return n.d.Find(".novel_subtitle").Text()
}

// FindBody return novel body.
func (n *narouQuery) FindBody() string {
	b := ""

	n.d.Find("div#novel_honbun").
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
func (n *narouQuery) FindPreface() string {
	b := ""
	n.d.Find(".novel_view#novel_p").Each(func(_ int, selection *goquery.Selection) {
		b += selection.Text()
	})
	return b
}

// FindAfterword return episode afterword.
func (n *narouQuery) FindAfterword() string {
	b := ""
	n.d.Find(".novel_view#novel_a").Each(func(_ int, selection *goquery.Selection) {
		b += selection.Text()
	})
	return b
}
