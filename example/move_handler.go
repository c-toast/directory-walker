package main

import (
	"bytes"
	"fmt"
	"github.com/c-toast/directory-walker/filewalker"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
)



type moveHandler struct{}

func (*moveHandler) CanHandle(info *filewalker.FileInfoWrapper) bool {
	ext := info.Ext
	if ext == ".rar" || ext == ".zip" {
		return false
	}
	return true
}

func (*moveHandler) Handle(info *filewalker.FileInfoWrapper) error {
	path:=info.FullPath

	targetName:=generateNewFilePath(path)
	logger.Info("moving file", zap.String("from", fmt.Sprintf("%s", path)),
		zap.String("to", fmt.Sprintf("%s", moveDir)))
	err:=os.Rename(path,targetName)
	return err
}

//return identifier such as 1a 1b 2a 2b
func GetIdentifier(path string)string{
	re:=regexp.MustCompile(`[1-9][a-z]`)
	res:=re.FindAllString(path,-1)
	if res!=nil{
		return res[len(res)-1]
	}
	return ""
}
//add the identifier ahead of the original file basename and put the file into moveDir
func generateNewFilePath(originPath string)string{
	var buf bytes.Buffer
	originFileBase:=filepath.Base(originPath)
	identifier:= GetIdentifier(originFileBase)
	//if the basename already contains the identifier, just join it
	if identifier!=""{
		buf.WriteString(filepath.Join(moveDir,originFileBase))
	} else{
		identifier= GetIdentifier(originPath)
		//if the identifier does not exist, we have no way but just use the original file base name
		if identifier==""{
			buf.WriteString(filepath.Join(moveDir,originFileBase))
		}else{//else, we add the identifier ahead of the file base name
			buf.WriteString(filepath.Join(moveDir,identifier))
			buf.WriteString(" ")
			buf.WriteString(originFileBase)
		}
	}
	return buf.String()
}
