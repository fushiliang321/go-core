package logger

import (
	"fmt"
	loggerConfig "github.com/fushiliang321/go-core/config/logger"
	"github.com/fushiliang321/go-core/helper/file"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

type (
	handlerWriter struct {
		file   io.Writer
		stdout *os.File
		sync.Mutex
	}

	logFile struct {
		dirPath               string
		fileName              string
		file                  *os.File
		writeInterval         time.Duration //日志文件写入间隔
		lumberjack            *loggerConfig.Lumberjack
		lumberjackMaxSizeByte int64
		writeChan             chan []byte
		buffer                []byte
		sync.RWMutex
	}
)

const LogFileSuffix = ".log"

func (w *handlerWriter) setLogFile(file io.Writer) {
	w.file = file
}

func (w *handlerWriter) setStdout(stdout *os.File) {
	w.stdout = stdout
}

func (w *handlerWriter) Write(b []byte) (n int, err error) {
	w.Lock()
	defer func() {
		w.Unlock()
		if _err := recover(); _err != nil {
			slog.Error(fmt.Sprint(_err))
		}
	}()

	if w.stdout != nil {
		n, err = w.stdout.Write(b)
		if err != nil {
			return 0, err
		}
	}

	if w.file != nil {
		n, err = w.file.Write(b)
		if err != nil {
			return 0, err
		}
	}
	return
}

func (f *logFile) open() error {
	f.writeChan = make(chan []byte, 100)
	if f.lumberjack.MaxSize > 0 {
		f.lumberjackMaxSizeByte = int64(f.lumberjack.MaxSize) * 1024 * 1024
	}
	err := os.MkdirAll(f.dirPath, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	_file, err := os.OpenFile(f.dirPath+"/"+f.fileName+LogFileSuffix, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	f.file = _file
	go f.writeJoint()
	return nil
}

func (f *logFile) Write(b []byte) (n int, err error) {
	f.writeChan <- b
	return len(b), err
}

// 拼接写入的数据
func (f *logFile) writeJoint() {
	if f.writeInterval == 0 {
		//实时写入
		for b := range f.writeChan {
			f.writeHandle(b)
		}
		return
	}
	f.buffer = []byte{}
	time.AfterFunc(f.writeInterval, func() {
		//间隔执行写入
		f.RLock()
		if len(f.buffer) == 0 {
			f.RUnlock()
			return
		}
		f.RUnlock()
		f.Lock()
		buffer := f.buffer
		f.buffer = []byte{}
		f.Unlock()
		f.writeHandle(buffer)
	})

	for b := range f.writeChan {
		f.Lock()
		f.buffer = append(f.buffer, b...)
		f.Unlock()
	}
}

func (f *logFile) writeHandle(b []byte) {
	if f.lumberjackMaxSizeByte <= 0 {
		f.file.Write(b)
		return
	}

	stat, err := f.file.Stat()
	if err != nil {
		return
	}
	fileSize := stat.Size() + int64(len(b))
	if fileSize > f.lumberjackMaxSizeByte {
		f.lumberjackHandle()
	}
	f.file.Write(b)
}

// 切割文件
func (f *logFile) lumberjackHandle() {
	var oldB []byte
	n, err := f.file.Read(oldB)
	if err != nil {
		return
	}
	if err = f.file.Truncate(int64(n)); err != nil {
		return
	}

	go func() {
		defer func() {
			if _err := recover(); _err != nil {
				slog.Error(fmt.Sprint(_err))
			}
		}()
		var (
			fileName     = f.lumberjack.FileNameFormat()
			fileFullPath string
			suffix       int
		)
		for {
			if suffix > 0 {
				//需要添加后缀
				fileName = fileName + "(" + strconv.Itoa(suffix) + ")"
			}
			fileFullPath = f.dirPath + "/" + fileName + LogFileSuffix
			_, err = os.Stat(fileFullPath)
			if os.IsExist(err) {
				//文件已存在
				suffix++
			} else {
				//文件不存在
				break
			}
		}
		_file, err := os.Create(fileFullPath)
		if err != nil {
			return
		}
		_file.Write(oldB)
		_file.Close()

		if f.lumberjack.MaxBackups > 0 || f.lumberjack.MaxAge > 0 {
			sortFile, err := file.SortFile(f.dirPath)
			if err != nil {
				return
			}
			var files []*os.FileInfo
			for _, info := range sortFile {
				if path.Ext((*info).Name()) == LogFileSuffix {
					files = append(files, info)
				}
			}
			if f.lumberjack.MaxBackups > 0 {
				//限制保留旧文件的最大个数
				surplus := len(files) - f.lumberjack.MaxBackups
				if surplus > 0 {
					for i := 0; i < surplus; i++ {
						os.Remove(f.dirPath + "/" + (*files[i]).Name())
					}
					files = files[surplus:]
				}
			}
			if f.lumberjack.MaxAge > 0 {
				//限制保留旧文件的最大天数
				now := time.Now()
				for _, fileItem := range files {
					if int(now.Sub((*fileItem).ModTime()).Hours()/24) > f.lumberjack.MaxAge {
						os.Remove(f.dirPath + "/" + (*fileItem).Name())
					}
				}
			}
		}
	}()
}
