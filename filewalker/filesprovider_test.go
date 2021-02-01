package filewalker_test

import (
	"fmt"
	"github.com/c-toast/directory-walker/filewalker"
	"github.com/c-toast/directory-walker/filewalker/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type testFileInfo struct {
	name  string
	isDir bool
}

func (i testFileInfo) Name() string       { return i.name }
func (i testFileInfo) IsDir() bool        { return i.isDir }
func (i testFileInfo) Size() int64        { return 0 }
func (i testFileInfo) Mode() os.FileMode  { return 0 }
func (i testFileInfo) ModTime() time.Time { return time.Time{} }
func (i testFileInfo) Sys() interface{}   { return nil }

func TestNext(t *testing.T) {
	r := &mocks.TestDirReader{}
	info1 := testFileInfo{"1", false}
	info2 := testFileInfo{"path2", true}
	info3 := testFileInfo{"3", false}
	info4 := testFileInfo{"error_path4", true}
	info5 := testFileInfo{"error_file5", false}
	info6 := testFileInfo{"6", false}
	//directory structure:
	//path1
	//	1
	//	path2
	//		3
	//	error_path4
	//		error_file5
	//	6
	r.On("ReadDir", "path1").Return([]os.FileInfo{info1, info2, info4, info6}, nil)
	r.On("ReadDir", "path1/path2").Return([]os.FileInfo{info3}, nil)
	r.On("ReadDir", "path1/error_path4").Return([]os.FileInfo{info5}, fmt.Errorf("path error"))

	p := filewalker.DefaultFilesProvider{}
	err := p.Init("path1", r)

	f, err := p.Next()
	assert.Equal(t, "1", f.FileName)
	assert.NoError(t, err)

	f, err = p.Next()
	assert.Equal(t, "3", f.FileName)
	assert.NoError(t, err)

	f, err = p.Next()
	assert.Equal(t, "path2", f.FileName)
	assert.NoError(t, err)

	f, err = p.Next()
	assert.Equal(t, "error_path4", f.FileName)
	assert.Error(t, err)

	f, err = p.Next()
	assert.Equal(t, "6", f.FileName)
	assert.NoError(t, err)

	f, err = p.Next()
	assert.Nil(t, nil, f)
	assert.NoError(t, err)
}
