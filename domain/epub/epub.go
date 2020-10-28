package epub

import (
	"context"

	"narou/config"

	"github.com/bmaupin/go-epub"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IEpub interface {
	New(author, title string) IEpub
	AddCSS(source string, internalFilename string) (string, error)
	AddFont(source string, internalFilename string) (string, error)
	AddImage(source string, imageFilename string) (string, error)
	AddSection(body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error)
	Author() string
	Identifier() string
	Lang() string
	Description() string
	Ppd() string
	SetAuthor(author string)
	SetCover(internalImagePath string, internalCSSPath string)
	SetIdentifier(identifier string)
	SetLang(lang string)
	SetDescription(desc string)
	SetPpd(direction string)
	SetTitle(title string)
	GetTitle() string
	Write(path string) error
}

type innerEP struct {
	e   *epub.Epub
	ctx context.Context
	cfg config.IConfigure
}

func NewEpub(ctx context.Context, cfg config.IConfigure) IEpub {
	return &innerEP{
		ctx: ctx,
		cfg: cfg,
	}
}
func (ep *innerEP) New(author, title string) IEpub {
	e := epub.NewEpub(title)
	e.SetAuthor(author)
	e.SetIdentifier(author + "-" + title)

	lang, ppd := ep.cfg.GetEpubSetting()
	e.SetLang(lang)
	e.SetPpd(ppd)

	return &innerEP{e: e}
}

func (ep *innerEP) AddCSS(source string, internalFilename string) (string, error) {
	return ep.e.AddCSS(source, internalFilename)
}

func (ep *innerEP) AddFont(source string, internalFilename string) (string, error) {
	return ep.e.AddFont(source, internalFilename)
}

func (ep *innerEP) AddImage(source string, imageFilename string) (string, error) {
	return ep.e.AddImage(source, imageFilename)
}

func (ep *innerEP) AddSection(body string, sectionTitle string, internalFilename string, internalCSSPath string) (
	string, error) {
	return ep.e.AddSection(body, sectionTitle, internalFilename, internalCSSPath)
}

// Author returns the author of the IEpub.
func (ep *innerEP) Author() string {
	return ep.e.Author()
}

// Identifier returns the unique identifier of the IEpub.
func (ep *innerEP) Identifier() string {
	return ep.e.Identifier()
}

// Lang returns the language of the IEpub.
func (ep *innerEP) Lang() string {
	return ep.e.Lang()
}

// Description returns the description of the IEpub.
func (ep *innerEP) Description() string {
	return ep.e.Description()
}

// Ppd returns the page progression direction of the IEpub.
func (ep *innerEP) Ppd() string {
	return ep.e.Ppd()
}

// SetAuthor sets the author of the IEpub.
func (ep *innerEP) SetAuthor(author string) {
	ep.e.SetAuthor(author)
}

func (ep *innerEP) SetCover(internalImagePath string, internalCSSPath string) {
	ep.e.SetCover(internalImagePath, internalCSSPath)
}

func (ep *innerEP) SetIdentifier(identifier string) {
	ep.e.SetIdentifier(identifier)
}

// SetLang sets the language of the IEpub.
func (ep *innerEP) SetLang(lang string) {
	ep.e.SetLang(lang)
}

// SetDescription sets the description of the IEpub.
func (ep *innerEP) SetDescription(desc string) {
	ep.e.SetDescription(desc)
}

// SetPpd sets the page progression direction of the IEpub.
func (ep *innerEP) SetPpd(direction string) {
	ep.e.SetPpd(direction)
}

// SetTitle sets the title of the IEpub.
func (ep *innerEP) SetTitle(title string) {
	ep.e.SetTitle(title)
}

// GetTitle returns the title of the IEpub.
func (ep *innerEP) GetTitle() string {
	return ep.e.Title()
}

func (ep *innerEP) Write(path string) error {
	return ep.e.Write(path)
}
