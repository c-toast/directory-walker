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
	Info os.FileInfo

	FileName string

	IsDir bool
	//only when the file is not dir, BasePath is meaningful
	BasePath string

	FullPath string

	Ext string
}

type directoryWalker struct {
	directoryPath string

	provider filesProvider

	handlersList []fileHandler

	//simple implementation of error record, should be replaced by writer
	errList []error
}

func (w *FileInfoWrapper) init(basePath string, info os.FileInfo) {
	w.Info = info
	w.IsDir = info.IsDir()
	if !w.IsDir {
		w.BasePath = basePath
		w.FullPath = filepath.Join(basePath, info.Name())
		w.Ext = filepath.Ext(info.Name())
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
