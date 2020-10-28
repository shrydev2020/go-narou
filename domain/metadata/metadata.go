package metadata

import (
	"strings"
	"time"
)

type Novel struct {
	ID              int `gorm:"primary_key,column:id"`
	Author          string
	Title           string
	FileTitle       string
	TopUrl          URI `gorm:"column:top_url"`
	SiteName        string
	Story           string
	NovelType       Kind
	End             bool
	LastUpdate      time.Time
	NewArrivalsDate time.Time
	UseSubdirectory bool
	GeneralFirstUp  time.Time
	NovelUpdatedAt  time.Time
	GeneralKastUp   time.Time
	Length          int
	Suspend         bool
	GeneralAllNo    int
	LastCheckAt     time.Time
	Sub             []Sub // TODO foreign key
}

func (n *Novel) TableName() string {
	return "novels"
}

type Sub struct {
	NovelID      int
	Index        int
	Href         URI
	Chapter      string
	Subtitle     string
	SubDate      time.Time
	SubUpDatedAt time.Time
	DownloadAt   time.Time
}

func (n *Sub) TableName() string {
	return "subs"
}

type URI string

type Kind int

const (
	Series Kind = iota + 1
	SS
)

func GetSiteName(uri URI) string {
	ur := string(uri)
	if strings.Contains(ur, "https://ncode.syosetu.com/") {
		return "小説家になろう"
	} else if strings.Contains(ur, "https://syosetu.org/novel/") {
		return "ハーメルン"
	}
	return "Nof Supported"
}

func GetID(uri URI) string {
	eplaced1 := strings.Replace(string(uri), "https://ncode.syosetu.com/", "", -1)
	eplaced1 = strings.Replace(string(uri), "/", "", -1)
	return eplaced1
}

func NewMetaNovel(author, title, outline string, uri URI, length int) *Novel {
	return &Novel{
		ID:        0,
		Author:    author,
		Title:     title,
		Story:     outline,
		FileTitle: GetID(uri) + "-" + title,
		TopUrl:    uri,
		SiteName:  GetSiteName(uri),
		// TODO fill
		NovelType:       0,
		End:             false,
		LastUpdate:      time.Time{},
		NewArrivalsDate: time.Time{},
		UseSubdirectory: false,
		GeneralFirstUp:  time.Time{},
		NovelUpdatedAt:  time.Time{},
		GeneralKastUp:   time.Time{},
		Length:          length,
		Suspend:         false,
		GeneralAllNo:    0,
		LastCheckAt:     time.Time{},
		Sub:             nil,
	}
}

func NewSubs(n *Novel, i int, subURI URI) Sub {
	return Sub{
		NovelID:      n.ID,
		Index:        i + 1,
		Href:         subURI,
		Chapter:      "",
		Subtitle:     "",
		SubDate:      time.Time{},
		SubUpDatedAt: time.Time{},
		DownloadAt:   time.Time{},
	}
}
