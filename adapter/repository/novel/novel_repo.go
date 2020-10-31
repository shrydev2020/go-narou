package novel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"narou/infrastructure/storage"

	"narou/domain/novel"
	"narou/sdk/slice"
)

type repo struct {
	dist, subDist string
}

func NewRepository(manager storage.Manager) novel.IRepository {
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

	err = ioutil.WriteFile(path+"/"+chapter, []byte(body), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByNobelSiteAndTitle(novelType, title string) []string {
	files := dirWalk(r.dist + "/" + novelType + "/" + title + "/" + r.subDist)
	ret := make([]string, 0, len(files))

	for _, v := range files {
		b, err := ioutil.ReadFile(v)
		if err != nil {
			panic(err)
		}

		ret = append(ret, string(b))
	}

	return ret
}

func dirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string

	for _, file := range files {
		if file.IsDir() {
			if strings.Contains(file.Name(), ".html") {
				paths = append(paths, dirWalk(filepath.Join(dir, file.Name()))...)
			}

			continue
		}

		if strings.Contains(file.Name(), ".html") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return slice.SortStrings(paths)
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
