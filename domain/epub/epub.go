package epub

import (
	"context"

	"narou/config"

	"github.com/bmaupin/go-epub"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IEpub interface {
	New(ctx context.Context, author, title string) IEpub
	AddCSS(ctx context.Context, source string, internalFilename string) (string, error)
	AddFont(ctx context.Context, source string, internalFilename string) (string, error)
	AddImage(ctx context.Context, source string, imageFilename string) (string, error)
	AddSection(ctx context.Context, body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error)
	Author(ctx context.Context) string
	Identifier(ctx context.Context) string
	Lang(ctx context.Context) string
	Description(ctx context.Context) string
	Ppd(ctx context.Context) string
	SetAuthor(ctx context.Context, author string)
	SetCover(ctx context.Context, internalImagePath string, internalCSSPath string)
	SetIdentifier(ctx context.Context, identifier string)
	SetLang(ctx context.Context, lang string)
	SetDescription(ctx context.Context, desc string)
	SetPpd(ctx context.Context, direction string)
	SetTitle(ctx context.Context, title string)
	GetTitle(ctx context.Context) string
	Write(ctx context.Context, path string) error
}

type innerEP struct {
	e   *epub.Epub
	cfg config.IConfigure
}

func NewEpub(cfg config.IConfigure) IEpub {
	return &innerEP{
		cfg: cfg,
	}
}
func (ep *innerEP) New(ctx context.Context, author, title string) IEpub {
	e := epub.NewEpub(title)
	e.SetAuthor(author)
	e.SetIdentifier(author + "-" + title)

	lang, ppd := ep.cfg.GetEpubSetting()
	e.SetLang(lang)
	e.SetPpd(ppd)

	return &innerEP{e: e}
}

func (ep *innerEP) AddCSS(_ context.Context, source string, internalFilename string) (string, error) {
	return ep.e.AddCSS(source, internalFilename)
}

func (ep *innerEP) AddFont(_ context.Context, source string, internalFilename string) (string, error) {
	return ep.e.AddFont(source, internalFilename)
}

func (ep *innerEP) AddImage(_ context.Context, source string, imageFilename string) (string, error) {
	return ep.e.AddImage(source, imageFilename)
}

func (ep *innerEP) AddSection(_ context.Context, body string, sectionTitle string, internalFilename string, internalCSSPath string) (
	string, error) {
	return ep.e.AddSection(body, sectionTitle, internalFilename, internalCSSPath)
}

// Author returns the author of the IEpub.
func (ep *innerEP) Author(_ context.Context) string {
	return ep.e.Author()
}

// Identifier returns the unique identifier of the IEpub.
func (ep *innerEP) Identifier(_ context.Context) string {
	return ep.e.Identifier()
}

// Lang returns the language of the IEpub.
func (ep *innerEP) Lang(_ context.Context) string {
	return ep.e.Lang()
}

// Description returns the description of the IEpub.
func (ep *innerEP) Description(_ context.Context) string {
	return ep.e.Description()
}

// Ppd returns the page progression direction of the IEpub.
func (ep *innerEP) Ppd(_ context.Context) string {
	return ep.e.Ppd()
}

// SetAuthor sets the author of the IEpub.
func (ep *innerEP) SetAuthor(_ context.Context, author string) {
	ep.e.SetAuthor(author)
}

func (ep *innerEP) SetCover(_ context.Context, internalImagePath string, internalCSSPath string) {
	ep.e.SetCover(internalImagePath, internalCSSPath)
}

func (ep *innerEP) SetIdentifier(_ context.Context, identifier string) {
	ep.e.SetIdentifier(identifier)
}

// SetLang sets the language of the IEpub.
func (ep *innerEP) SetLang(_ context.Context, lang string) {
	ep.e.SetLang(lang)
}

// SetDescription sets the description of the IEpub.
func (ep *innerEP) SetDescription(_ context.Context, desc string) {
	ep.e.SetDescription(desc)
}

// SetPpd sets the page progression direction of the IEpub.
func (ep *innerEP) SetPpd(_ context.Context, direction string) {
	ep.e.SetPpd(direction)
}

// SetTitle sets the title of the IEpub.
func (ep *innerEP) SetTitle(_ context.Context, title string) {
	ep.e.SetTitle(title)
}

// GetTitle returns the title of the IEpub.
func (ep *innerEP) GetTitle(_ context.Context) string {
	return ep.e.Title()
}

func (ep *innerEP) Write(_ context.Context, path string) error {
	return ep.e.Write(path)
}
