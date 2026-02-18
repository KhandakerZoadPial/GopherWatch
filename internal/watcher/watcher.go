package watcher

import (
	"io"
	"os"
)

type Watcher interface {
	Watch() (chan string, error)
}

type FileWatcher struct {
	FileName string
}

func (wh *FileWatcher) Watch() (chan string, error) {
	var file, err = os.Open(wh.FileName)

	if err != nil {
		return nil, err
	}

	var endLine = file.Seek(0, io.SeekEnd)

}
