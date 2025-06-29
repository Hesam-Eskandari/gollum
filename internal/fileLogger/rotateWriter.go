package fileLogger

import (
	"os"
	"sync"
	"time"
)

type rotateWriter struct {
	lock     sync.Mutex
	filename string
	file     *os.File
	config   Config
	stop     chan struct{}
}

func newRotateWriter(filename string, config Config) (*rotateWriter, error) {
	config.CheckPeriod = max(min(config.CheckPeriod, maxCheckPeriod), minCheckPeriod)
	config.PurgePeriod = max(min(config.PurgePeriod, maxPurgePeriod), minPurgePeriod)
	config.MaxFileSizeByte = max(min(config.MaxFileSizeByte, maxFileSizeByte), minFileSizeByte)

	w := &rotateWriter{
		filename: filename,
		config:   config,
		stop:     make(chan struct{}, 1),
	}
	err := w.rotate()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.file.Write(output)
}

func (w *rotateWriter) Check() {
	ticker := time.NewTicker(w.config.CheckPeriod)
	defer ticker.Stop()
	lastCheck := time.Now()
loop:
	for {
		info, err := os.Stat(w.filename)
		if (err == nil && info.Size() >= w.config.MaxFileSizeByte) || time.Since(lastCheck) >= w.config.PurgePeriod {
			_ = w.rotate()
			lastCheck = time.Now()
		}
		select {
		case <-ticker.C:
		case <-w.stop:
			break loop
		}
	}
}

func (w *rotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.file != nil {
		err = w.file.Close()
		w.file = nil
		if err != nil {
			return
		}
	}
	if _, err = os.Stat(w.filename); err == nil {
		err = os.Rename(w.filename, getFilepath(w.filename))
		if err != nil {
			return
		}
	}
	w.file, err = os.Create(w.filename)
	return
}

func (w *rotateWriter) destroy() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	defer close(w.stop)
	w.stop <- struct{}{}
	if w.file != nil {
		err = w.file.Close()
		w.file = nil
	}
	return
}
