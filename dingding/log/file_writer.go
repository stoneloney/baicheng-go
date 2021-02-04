package log

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

//const logFileNameFormat = "%s/%s.%4d-%02d-%02d.log"
const logFileNameFormat = "%s/%s.log"

// FileWriter 日志实现Writer
type FileWriter struct {
	maxSize  int64
	maxNum   int
	fileName string
	filePath string
	file     *os.File
	writer   io.Writer
	mu       sync.Mutex
	ch       chan []byte
}

// NewFileWriter 新建一个日志writer，并启动三个goroutine来 rotate, check, flush
func NewFileWriter(filePath string, fileName string, maxSize int64, maxNum int) (io.Writer, error) {
	//y, m, d := time.Now().Date()
	_, e := os.Stat(filePath)
	if e != nil && os.IsNotExist(e) {
		e = os.Mkdir(filePath, os.ModePerm)
		if e != nil {
			panic(e)
		}
	}
	//path := fmt.Sprintf(logFileNameFormat, filePath, fileName, y, m, d)
	path := fmt.Sprintf(logFileNameFormat, filePath, fileName)
	file, e := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if e != nil {
		return nil, e
	}
	writer := &FileWriter{fileName: fileName, filePath: path, file: file, writer: file, ch: make(chan []byte, 256), maxSize: maxSize, maxNum: maxNum}
	// set log output
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(file)

	//go writer.rotate()
	//go writer.flush()
	go writer.check()

	return writer, nil
}

// Write 异步channel写日志
func (w *FileWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	select {
	case w.ch <- buf:
		return len(buf), nil
	default:
		return 0, errors.New("chan full, drop")
	}
}

// check 每分钟检查一下日志文件是否存在，运维误删log文件但是进程一直在打日志，fd会一直存在，需要关闭。超过maxSize自动rotate
func (w *FileWriter) check() {
	for {
		time.Sleep(time.Minute)

		w.mu.Lock()
		fileInfo, err := os.Stat(w.filePath)
		if os.IsNotExist(err) {
			file, e := os.OpenFile(w.filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if e == nil {
				w.file.Close()
				w.file = file
				w.writer = file
			}
			w.mu.Unlock()
			continue
		}
		if w.maxSize > 0 && fileInfo.Size() > w.maxSize {
			name := w.filePath + ".full."
			files, _ := ioutil.ReadDir(path.Dir(w.filePath))
			var minNum = 1000000
			var maxNum = 0
			var totalNum = 0
			for _, f := range files {
				if strings.Contains(f.Name(), name) {
					totalNum++
					s := strings.Split(f.Name(), ".")
					if len(s) > 4 {
						n, _ := strconv.Atoi(s[3])
						if n > maxNum {
							maxNum = n
						}
						if n < minNum {
							minNum = n
						}
					}
				}
			}
			w.file.Close()
			//rename log file
			name = fmt.Sprintf("%s.full.%d.log", w.filePath, maxNum+1)
			err := os.Rename(w.filePath, name)
			if err != nil {
				fmt.Printf("rename file path:%s fail:%s\n", w.filePath, err)
			}
			file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("open file path:%s fail:%s\n", w.filePath, err)
			}
			if err == nil {
				w.file = file
				w.writer = file
			}
			// set log output
			log.SetOutput(file)

			if totalNum >= w.maxNum {
				//remove oldest log file
				name = fmt.Sprintf("%s.full.%d.log", w.filePath, minNum)
				err := os.Remove(name)
				if err != nil {
					fmt.Printf("remove file path:%s fail:%s\n", name, err)
				}
			}
		}
		w.mu.Unlock()
	}
}

// rotate 按天更新日志文件名
/*
func (w *FileWriter) rotate() {
	for {
		now := time.Now()
		y, m, d := now.Add(24 * time.Hour).Date()
		nextDay := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		tm := time.NewTimer(time.Duration(nextDay.UnixNano() - now.UnixNano() - 100))
		<-tm.C
		w.mu.Lock()
		path := fmt.Sprintf(logFileNameFormat, w.fileName, y, m, d)
		file, e := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if e == nil {
			w.file.Close()
			w.file = file
			w.writer = file
			w.filePath = path
		}
		w.mu.Unlock()
	}
}
*/

// flush 刷新日志到磁盘中
func (w *FileWriter) flush() {
	for {
		log := <-w.ch
		w.mu.Lock()
		w.writer.Write(log)
		w.mu.Unlock()
	}
}
