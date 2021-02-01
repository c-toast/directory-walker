package filewalker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//dirReader is designed to be an interface for unit test convenience
type dirReader interface {
	ReadDir(path string) ([]os.FileInfo, error)
}

type DefaultdirReader struct{}

func (DefaultdirReader) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

//filesProvider determines the order of files to return.
//default implementation return inorder sequence of files
type filesProvider interface {
	//Next will provide file only after the Init have been called
	Init(dirPath string, dirReaderImpl ...interface{}) error

	Next() (*FileInfoWrapper, error)
}

type DefaultFilesProvider struct {
	dirPath string

	subProvider *DefaultFilesProvider

	reader dirReader

	fileInfoLists []os.FileInfo
}

func (p *DefaultFilesProvider) Init(dirPath string, dirReaderImpl ...interface{}) error {
	if len(dirReaderImpl) > 1 {
		return fmt.Errorf("provide more than one extra args to Init")
	}
	if len(dirReaderImpl) == 1 {
		ok := false
		if p.reader, ok = dirReaderImpl[0].(dirReader); !ok {
			return fmt.Errorf("extra arg is not a implementation of dirReader")
		}
	} else {
		p.reader = DefaultdirReader{}
	}

	p.dirPath = dirPath
	p.subProvider = nil
	var err error = nil
	p.fileInfoLists, err = p.reader.ReadDir(dirPath)
	return err
}

func (p *DefaultFilesProvider) wrapFirstFileInfo() *FileInfoWrapper {
	w := &FileInfoWrapper{}
	w.init(p.dirPath, p.fileInfoLists[0])
	p.fileInfoLists = p.fileInfoLists[1:]
	return w
}

//subDirNext return the dir if all the file in the dir have been visited
//return error if fail to open directory
func (p *DefaultFilesProvider) subDirNext() (*FileInfoWrapper, error) {
	if p.subProvider == nil {
		p.subProvider = &DefaultFilesProvider{}
		err := p.subProvider.Init(path.Join(p.dirPath, p.fileInfoLists[0].Name()), p.reader)
		//if err happened, return the dir and the error
		if err != nil {
			return p.wrapFirstFileInfo(), err
		}
	}

	w, _ := p.subProvider.Next()
	if w == nil {
		//if all the file in dir have been walked, reset the subProvider and return the dir
		p.subProvider=nil
		w = p.wrapFirstFileInfo()
	}
	return w, nil
}

//Next return nil FileInfoWrapper if all the file have been visited.
func (p *DefaultFilesProvider) Next() (*FileInfoWrapper, error) {
	if len(p.fileInfoLists) == 0 {
		return nil, nil
	}
	if p.fileInfoLists[0].IsDir() {
		return p.subDirNext()
	}
	return p.wrapFirstFileInfo(), nil
}
