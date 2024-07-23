package content

import (
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
)

func NewFilepath(getPath func() string) Handler {
	return newFS(&filepathFS{
		getPath: getPath,
	})
}

var _ fsContent = &filepathFS{}

type filepathFS struct {
	getPath func() string
}

func (f *filepathFS) ToFileServer(basePaths ...string) http.Handler {
	root := f.getPath()
	if root == "" {
		return http.NotFoundHandler()
	}
	path := filepath.Join(append([]string{string(root)}, basePaths...)...)
	return http.FileServer(http.Dir(path))
}

func (f *filepathFS) Open(name string) (fs.File, error) {
	root := f.getPath()
	if root == "" {
		return nil, errors.New("filepath fs is not ready")
	}
	return http.Dir(root).Open(name)
}

func (f *filepathFS) Refresh() error {
	return nil
}
