package novel

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"

	"narou/infrastructure/storage"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
type IRepository interface {
	Store(novelType, title, chapterTitle, body string) error
	FindByNobelSiteAndTitle(novelSite, title string) ([]string, error)
}
type repo struct {
	dist, subDist string
}

func NewRepository(manager storage.Manager) IRepository {
	return &repo{
		dist:    manager.GetDist(),
		subDist: manager.GetSubDist(),
	}
}

func (r *repo) Store(nobelType, title, chapter, body string) error {
	path, err := r.makeDirIfNoExist(nobelType, title)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+chapter, []byte(body), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByNobelSiteAndTitle(novelType, title string) ([]string, error) {
	files, err := dirWalk(r.dist + "/" + novelType + "/" + title + "/" + r.subDist)
	if err != nil {
		return nil, errors.Wrapf(err, "dirWailErr:%s", files)
	}
	ret := make([]string, 0, len(files))

	for _, v := range files {
		b, err := os.ReadFile(v)
		if err != nil {
			return nil, errors.Wrapf(err, "dirWailErr:%s", files)
		}

		ret = append(ret, string(b))
	}

	return ret, nil
}

func dirWalk(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			subDirPaths, err := dirWalk(fullPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subDirPaths...)
		} else if strings.HasSuffix(file.Name(), ".html") {
			paths = append(paths, fullPath)
		}
	}
	return paths, nil
}

func (r *repo) makeDirIfNoExist(novelType, title string) (string, error) {
	path := r.dist + "/" + novelType + "/" + title + "/" + r.subDist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}
