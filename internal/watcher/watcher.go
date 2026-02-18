package watcher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type Watcher interface {
	Watch() (<-chan string, error)
}

type FileWatcher struct {
	FileName string
}

func (fw *FileWatcher) run(file *os.File, ch chan string) {
	defer file.Close()
	defer close(ch)

	var reader = bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			fmt.Printf("Watcher Error On: %s:%v\n", fw.FileName, err)
			return
		}
		ch <- line
	}
}

func (fw *FileWatcher) Watch() (chan string, error) {
	var file, err = os.Open(fw.FileName)

	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		file.Close()
		return nil, err
	}

	var channel = make(chan string)
	go fw.run(file, channel)

	return channel, nil
}
