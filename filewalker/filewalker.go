package filewalker

import (
	"fmt"
	"os"
	"path/filepath"
)

//directoryWalker and fileHandler follow observer pattern. When directoryWalker Walk a
//file(or directory), it will notify all fileHandler that have registered to
//handle the file
type fileHandler interface {
	CanHandle(info *FileInfoWrapper) bool

	Handle(info *FileInfoWrapper) error
}

type FileInfoWrapper struct {
	info os.FileInfo

	FileName string

	isDir bool
	//only when the file is not dir, basePath is meaningful
	basePath string

	fullPath string

	ext string
}

type directoryWalker struct {
	directoryPath string

	provider filesProvider

	handlersList []fileHandler

	//simple implementation of error record, should be replaced by writer
	errList []error
}

func (w *FileInfoWrapper) init(basePath string, info os.FileInfo) {
	w.info = info
	w.isDir = info.IsDir()
	if !w.isDir {
		w.basePath = basePath
		w.fullPath = filepath.Join(basePath, info.Name())
		w.ext = filepath.Ext(info.Name())
	}
	w.FileName = info.Name()
}

func (walker *directoryWalker) RegisterHandler(h fileHandler) int {
	walker.handlersList = append(walker.handlersList, h)
	return len(walker.handlersList)
}

func (walker *directoryWalker) RegisterHandlers(handlersList []fileHandler) int {
	walker.handlersList = append(walker.handlersList, handlersList...)
	return len(walker.handlersList)
}

func (walker *directoryWalker) clearALLHandlers() {
	walker.handlersList = []fileHandler{}
}

func (walker *directoryWalker) handleFile(w *FileInfoWrapper) {
	for _, h := range walker.handlersList {
		if h.CanHandle(w) {
			err := h.Handle(w)
			if err != nil {
				walker.errList = append(walker.errList, err)
			}
		}
	}
}

func (walker *directoryWalker) Walk() {
	p := walker.provider
	for true {
		fileWrapper, err := p.Next()
		if err != nil {
			walker.errList = append(walker.errList, err)
		}
		if fileWrapper == nil {
			break
		}
		walker.handleFile(fileWrapper)
	}
}

func (walker *directoryWalker) Init(dirPath string, providerImpl ...interface{}) error {
	if len(providerImpl) == 1 {
		if _, ok := providerImpl[0].(filesProvider); !ok {
			return fmt.Errorf("extra arg is not a implementation of fileProvider")
		}
		walker.provider = providerImpl[0].(filesProvider)
	}
	err := walker.provider.Init(dirPath)
	return err
}

func New() *directoryWalker {
	w := directoryWalker{}
	w.provider = &DefaultFilesProvider{}
	return &w
}
