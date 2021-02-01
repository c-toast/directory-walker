package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unzip/filehandler"
	"unzip/filewalker"
)

var testPath="D:\\Workspace\\drawing\\proko解剖\\1、the basics（1—4）"
var targetDir ="D:\\Workspace\\drawing\\proko解剖\\1、the basics（1—4）"

type unzipOverride struct{
	filehandler.UnzipHandler
}

type moveOverride struct{
	filehandler.MoveHandler
}

type removeOverride struct{
	filehandler.RemoveHandler
}

func (unzipOverride)CanHandle(basePath string,info os.FileInfo,err error) bool{
	ext:=filepath.Ext(info.Name())
	if ext==".rar" ||ext==".zip"{
		return true
	}
	return false
}

func (removeOverride)CanHandle(basePath string,info os.FileInfo,err error) bool{
	ext:=filepath.Ext(info.Name())
	if info.IsDir() || ext==".rar" ||ext==".zip"{
		return true
	}
	return false
}

func (moveOverride)CanHandle(basePath string,info os.FileInfo,err error) bool{
	ext:=filepath.Ext(info.Name())
	if ext==".rar" ||ext==".zip"{
		return false
	}
	return true
}

func (moveOverride) GeneratePath(originalPath string)string{
	return generateNewFilePath(originalPath)
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
//add the identifier ahead of the original file basename and put the file into targetDir
func generateNewFilePath(originPath string)string{
	var buf bytes.Buffer
	originFileBase:=filepath.Base(originPath)
	identifier:= GetIdentifier(originFileBase)
	//if the basename already contains the identifier, just join it
	if identifier!=""{
		buf.WriteString(filepath.Join(targetDir,originFileBase))
	} else{
		identifier= GetIdentifier(originPath)
		//if the identifier does not exist, we have no way but just use the original file base name
		if identifier==""{
			buf.WriteString(filepath.Join(targetDir,originFileBase))
		}else{//else, we add the identifier ahead of the file base name
			buf.WriteString(filepath.Join(targetDir,identifier))
			buf.WriteString(" ")
			buf.WriteString(originFileBase)
		}
	}
	return buf.String()
}

func main(){
	walker:=filewalker.New(nil)
	unzip:=filehandler.HandlerImpl{Impl: unzipOverride{}}
	remove:=filehandler.HandlerImpl{Impl: removeOverride{}}
	move:=filehandler.HandlerImpl{Impl: moveOverride{}}

	walker.RegisterHandler(unzip)
	walker.RegisterHandler(remove)
	walker.RegisterHandler(move)

	err:=walker.Walk(testPath)
	fmt.Printf("error:%v",err)
}