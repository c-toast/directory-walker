package filewalker_test

import (
	"github.com/c-toast/directory-walker/filewalker"
	"github.com/c-toast/directory-walker/filewalker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestWalk(t *testing.T) {
	p := &mocks.TestFilesProvider{}
	p.On("Init", mock.Anything).Return(nil)

	p.On("Next", mock.Anything).Return(&filewalker.FileInfoWrapper{FileName: "1"}, nil).Times(1)
	p.On("Next", mock.Anything).Return(&filewalker.FileInfoWrapper{FileName: "2"}, nil).Times(1)
	p.On("Next", mock.Anything).Return(nil, nil).Times(1)

	walker := filewalker.New()
	err := walker.Init("path", p)
	assert.NoError(t, err)

	h1 := &mocks.TestFileHandler{}
	walker.RegisterHandler(h1)

	h1.On("CanHandle", mock.Anything).Return(true)
	h1.On("Handle", mock.Anything).Return(nil)

	walker.Walk()
	h1.AssertNumberOfCalls(t, "CanHandle", 2)
	h1.AssertNumberOfCalls(t, "Handle", 2)
}
