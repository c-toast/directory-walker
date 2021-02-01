package main

import (
	"fmt"
	"github.com/c-toast/directory-walker/filewalker"
	"go.uber.org/zap"
	"os"
)

type removeHandler struct{}

func (*removeHandler) CanHandle(info *filewalker.FileInfoWrapper) bool {
	ext := info.Ext
	if info.IsDir || ext==".rar" ||ext==".zip"{
		return true
	}
	return false
}

func (*removeHandler) Handle(info *filewalker.FileInfoWrapper) error {
	path:=info.FullPath
	logger.Info("removing file", zap.String("filename", fmt.Sprintf("%s", path)), )
	err:=os.Remove(path)
	return err
}
