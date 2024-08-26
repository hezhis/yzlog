package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var OpenNewFileByByDay CheckTimeToOpenNewFileFunc = func(lastOpenFileTime *time.Time) (string, bool) {
	if lastOpenFileTime == nil {
		return time.Now().Format(".01-02.log"), true
	}

	lastOpenYear, lastOpenMonth, lastOpenDay := lastOpenFileTime.Date()

	now := time.Now()
	nowYear, nowMonth, nowDay := now.Date()

	if lastOpenDay != nowDay || lastOpenMonth != nowMonth || lastOpenYear != nowYear {
		return time.Now().Format(".01-02.log"), true
	}

	return "", false
}

type FileWriter struct {
	*WriterConfig
	fp *os.File

	lastCheckIsFullAt int64
	isFileFull        bool

	openCurrentFileTime *time.Time

	bufCh           chan []byte
	flushSignCh     chan struct{}
	flushDoneSignCh chan error
}

func NewFileWriter(config WriterConfig) *FileWriter {
	config.BasePath = strings.TrimRight(config.BasePath, "/")
	if config.CheckFileFullIntervalSecs <= 0 {
		config.CheckFileFullIntervalSecs = 1
	}
	writer := &FileWriter{
		WriterConfig:    &config,
		bufCh:           make(chan []byte, config.ChanCapacity),
		flushSignCh:     make(chan struct{}),
		flushDoneSignCh: make(chan error),
	}

	go func() {
		err := writer.Loop()
		if err != nil {
			panic(err)
		}
	}()
	return writer
}

func (w *FileWriter) Write(content string) {
	select {
	case w.bufCh <- []byte(content):
	default:
		fmt.Println("log content cached buf full, lost:" + content)
	}
}

func (w *FileWriter) Sync() error {
	w.flushSignCh <- struct{}{}
	return <-w.flushDoneSignCh
}

func (w *FileWriter) Loop() error {
	doWriteMoreAsPossible := func(buf []byte) error {
		for {
			var moreBuf []byte
			select {
			case moreBuf = <-w.bufCh:
				buf = append(buf, moreBuf...)
			default:
			}

			if moreBuf == nil {
				break
			}
		}

		if len(buf) == 0 {
			return nil
		}

		if err := w.tryOpenNewFile(); err != nil {
			return err
		}

		if isFull, err := w.checkFileIsFull(); err != nil {
			return err
		} else if isFull {
			if err := w.rotate(); err != nil {
				return err
			}
		}

		bufLen := len(buf)
		var totalWrittenBytes int
		for {
			n, err := w.fp.Write(buf[totalWrittenBytes:])
			if err != nil {
				return err
			}
			totalWrittenBytes += n
			if totalWrittenBytes >= bufLen {
				break
			}
		}

		return nil
	}

	for {
		select {
		case buf := <-w.bufCh:
			if err := doWriteMoreAsPossible(buf); err != nil {
				return err
			}
		case _ = <-w.flushSignCh:
			if err := doWriteMoreAsPossible([]byte{}); err != nil {
				w.finishFlush(err)
				break
			}
			if err := w.fp.Sync(); err != nil {
				w.finishFlush(err)
				break
			}
			w.finishFlush(nil)
		}
	}
}

func (w *FileWriter) tryOpenNewFile() error {
	var err error
	fileName, ok := w.CheckTimeToOpenNewFile(w.openCurrentFileTime)
	if !ok {
		if w.fp == nil {
			return errors.New("get first file name failed")
		}

		return nil
	}

	if w.fp == nil {
		if _, err = os.Stat(w.BasePath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err = os.MkdirAll(w.BasePath, w.Perm); err != nil {
				return err
			}
		}
	}

	if w.fp, err = os.OpenFile(fmt.Sprintf("%s/%s%s", w.BasePath, w.BaseFileName, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, w.Perm); err != nil {
		return err
	}

	openFileTime := time.Now()
	w.openCurrentFileTime = &openFileTime
	w.isFileFull = false
	w.lastCheckIsFullAt = 0

	return nil
}

func (w *FileWriter) finishFlush(err error) {
	w.flushDoneSignCh <- err
}

func (w *FileWriter) checkFileIsFull() (bool, error) {
	if w.lastCheckIsFullAt != 0 && w.lastCheckIsFullAt+w.CheckFileFullIntervalSecs > time.Now().Unix() {
		return w.isFileFull, nil
	}

	fileInfo, err := w.fp.Stat()
	if err != nil {
		return false, err
	}

	w.isFileFull = fileInfo.Size() >= w.MaxFileSize
	w.lastCheckIsFullAt = time.Now().Unix()

	return w.isFileFull, nil
}

func (w *FileWriter) rotate() error {
	var fileName string
	if w.fp != nil {
		fileName = w.fp.Name()
	}
	if err := w.close(); err != nil {
		return err
	}
	if len(fileName) > 0 {
		if err := w.backup(fileName); err != nil {
			return err
		}
	}

	if err := w.tryOpenNewFile(); err != nil {
		return err
	}

	return nil
}

func (w *FileWriter) close() error {
	if w.fp == nil {
		return nil
	}

	err := w.fp.Close()
	w.fp = nil

	w.openCurrentFileTime = nil

	return err
}

func (w *FileWriter) backup(fileName string) error {
	return os.Rename(fileName, backUpName(fileName))
}

func backUpName(name string) string {
	const backupTimeFormat = "2006-01-02T15-04-05.000"

	dir := filepath.Dir(name)
	fileName := filepath.Base(name)
	ext := filepath.Ext(fileName)
	prefix := fileName[:len(fileName)-len(ext)]
	timestamp := time.Now().Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", prefix, timestamp, ext))
}
