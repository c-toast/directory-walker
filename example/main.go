package main

import (
	"fmt"
	"github.com/c-toast/directory-walker/filewalker"
)

type printHandler struct{}

func (*printHandler) CanHandle(info *filewalker.FileInfoWrapper) bool {
	return true
}

func (*printHandler) Handle(info *filewalker.FileInfoWrapper) error {
	fmt.Printf("%s\n", info.FileName)
	return nil
}

func main() {
	walker := filewalker.New()
	err := walker.Init("./")
	if err != nil {
		fmt.Printf("%v", err)
	}
	walker.RegisterHandler(&printHandler{})
	walker.Walk()
}
