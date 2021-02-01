package main

import (
	"fmt"
	"github.com/c-toast/directory-walker/filewalker"
	"github.com/mholt/archiver"
	"go.uber.org/zap"
	"os"
	"strings"
)

type unzipHandler struct{}

func (*unzipHandler) CanHandle(info *filewalker.FileInfoWrapper) bool {
	ext := info.Ext
	if ext == ".rar" || ext == ".zip" {
		return true
	}
	return false
}

func (*unzipHandler) Handle(info *filewalker.FileInfoWrapper) error {
	path := info.FullPath
	ext := info.Ext

	var unarchiver archiver.Unarchiver
	switch ext {
	case ".rar":
		unarchiver = archiver.NewRar()
	case ".zip":
		unarchiver = archiver.NewZip()
	default:
		return fmt.Errorf("the format of %s is not supported now", path)
	}

	//make a new dir to store unarchived file
	targetDir := strings.TrimSuffix(path, ext)
	err := os.Mkdir(targetDir, os.ModeDir)
	if err != nil {
		return err
	}

	logger.Info("unzipping file", zap.String("from", fmt.Sprintf("%s", path)),
		zap.String("to", fmt.Sprintf("%s", targetDir)))
	err = unarchiver.Unarchive(path, targetDir)

	return err
}
