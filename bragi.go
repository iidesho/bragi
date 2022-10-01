package bragi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	human  = log.New(os.Stdout, "", 0)
	humanf *os.File
	json   = log.New(os.Stdout, "", 0)
	jsonf  *os.File
	folder string
	prefix = "Default"
	level  = DEBUG
	ctx    context.Context
	cancel func()
)

type Level int

const (
	DEBUG Level = iota
	INFO
	NOTICE
	WARNING
	ERROR
	CRIT
)

func (l Level) String() string {
	return []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRIT"}[l]
}

func SetPrefix(p string) {
	prefix = p
}

func Closer() {
	humanf.Close()
	jsonf.Close()
	cancel()
}

func SetOutputFolder(path string) func() {
	ctx, cancel = context.WithCancel(context.Background())
	folder = path
	if !fileExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil
		}
	}
	jsonPath := path + "/json"
	if !fileExists(jsonPath) {
		err := os.MkdirAll(jsonPath, 0755)
		if err != nil {
			return nil
		}
	}
	var err error
	humanf, jsonf, err = newLogFiles(path, jsonPath)
	if err != nil {
		AddError(err).Error("unable to create new logfiles")
		return nil
	}
	human.SetOutput(humanf)
	json.SetOutput(jsonf)
	go func() {
		nextDay := time.Now().UTC().AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 1, time.UTC)
		rotateTicker := time.Tick(time.Second)
		rotateDayTicker := time.NewTicker(nextDay.Sub(nextDay))
		truncateTaleTicker := time.Tick(time.Second * 5)
		firstDay := true
		for {
			select {
			case <-ctx.Done():
				return
			case <-rotateTicker:
				rotate(path, jsonPath)
			case <-rotateDayTicker.C:
				if firstDay {
					firstDay = false
					rotateDayTicker.Reset(24 * time.Hour)
				}
				rotate(path, jsonPath)
			case <-truncateTaleTicker:
				truncateTale(path)
				truncateTale(jsonPath)
			}
		}
	}()
	return Closer
}

func SetLevel(l Level) {
	level = l
}

type logData struct {
	err   error
	level Level
}

func AddError(e error) logData {
	return logData{err: e, level: INFO}
}

type Stringer interface {
	String() string
}

var jsonEscaper = strings.NewReplacer(
	`"`, `\"`,
)

func (ld logData) format(s string) (human, json string) {
	human = time.Now().Format("15:04:05 MST")
	json = fmt.Sprintf(`{"@timestamp":"%s"`, time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
	var (
		function uintptr
		file     string
		line     int
	)
	path := strings.Split(file, "/")
	for i := 1; len(path) == 1 || path[len(path)-1] == "bragi.go"; i++ {
		function, file, line, _ = runtime.Caller(i)
		path = strings.Split(file, "/")
	}
	if ld.level == DEBUG || ld.level == CRIT {
		human = fmt.Sprintf("%s %s:%d/%s", human, path[len(path)-1], line, runtime.FuncForPC(function).Name())
	}
	json = fmt.Sprintf(`%s,"data":{"file":"%s","line":%d,"function":"%s"`, json, path[len(path)-1], line, runtime.FuncForPC(function).Name())
	if ld.err != nil {
		json = fmt.Sprintf(`%s,"error":"%s"`, json, jsonEscaper.Replace(ld.err.Error()))
	}
	json = fmt.Sprintf(`%s}`, json)
	human = fmt.Sprintf("%s [%s]%s", human, ld.level, s)
	json = fmt.Sprintf(`%s,"level":"%s","message":"%s"`, json, ld.level, jsonEscaper.Replace(s))
	if ld.err != nil {
		human = fmt.Sprintf("%s. Err: %v", human, ld.err)
	}
	json = fmt.Sprintf("%s}", json)
	return
}

func (ld logData) Print(a ...interface{}) {
	if level > ld.level {
		return
	}
	humanString, jsonString := ld.format(fmt.Sprint(a...))
	human.Print(humanString)
	if folder == "" {
		return
	}
	json.Print(jsonString)
}

func (ld logData) Printf(format string, a ...interface{}) {
	ld.Print(fmt.Sprintf(format, a...))
}

func Printf(format string, a ...interface{}) {
	AddError(nil).Printf(format, a...)
}

func (ld logData) Println(a ...interface{}) {
	if level > ld.level {
		return
	}
	humanString, jsonString := ld.format(fmt.Sprint(a...))
	human.Println(humanString)
	if folder == "" {
		return
	}
	json.Println(jsonString)
}

func Println(a ...interface{}) {
	AddError(nil).Println(a...)
}

func (ld logData) Debug(a ...interface{}) {
	ld.level = DEBUG
	ld.Print(a...)
}

func Debug(a ...interface{}) {
	AddError(nil).Debug(a...)
}

func (ld logData) Info(a ...interface{}) {
	ld.level = INFO
	ld.Print(a...)
}

func Info(a ...interface{}) {
	AddError(nil).Info(a...)
}

func (ld logData) Notice(a ...interface{}) {
	ld.level = NOTICE
	ld.Print(a...)
}

func Notice(a ...interface{}) {
	AddError(nil).Notice(a...)
}

func (ld logData) Warning(a ...interface{}) {
	ld.level = WARNING
	ld.Print(a...)
}

func Warning(a ...interface{}) {
	AddError(nil).Warning(a...)
}

func (ld logData) Error(a ...interface{}) {
	ld.level = ERROR
	ld.Print(a...)
}

func Error(a ...interface{}) {
	AddError(nil).Error(a...)
}

func (ld logData) Crit(a ...interface{}) {
	ld.level = CRIT
	ld.Print(a...)
}

func Crit(a ...interface{}) {
	AddError(nil).Crit(a...)
}

func (ld logData) Fatal(a ...interface{}) {
	ld.Crit(a...)
	panic("Exiting from call to fatal")
}

func Fatal(a ...interface{}) {
	AddError(nil).Fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a...)) // Tmp until i need more
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func newLogFiles(path, jsonPath string) (hf *os.File, jf *os.File, err error) {
	hf, err = os.OpenFile(fmt.Sprintf("%s/%s.log", path, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	jf, err = os.OpenFile(fmt.Sprintf("%s/%s.log", jsonPath, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		hf.Close()
		return
	}
	return
}

func rotate(path, jsonPath string) {
	jsonStat, err := jsonf.Stat()
	if err != nil {
		AddError(err).Error("unable to get json log file stats for rotation")
	}
	if jsonStat.Size() < 24*MB {
		return
	}
	tf := time.Now().Format("2006-01-02T15:04:05")
	oldName := humanf.Name()
	err = os.Rename(oldName, strings.Replace(oldName, ".log", fmt.Sprintf("-%s.log", tf), 1))
	if err != nil {
		AddError(err).Error("unable to move old human log file")
		return
	}
	oldName = jsonf.Name()
	err = os.Rename(oldName, strings.Replace(oldName, ".log", fmt.Sprintf("-%s.log", tf), 1))
	if err != nil {
		AddError(err).Error("unable to move old json log file")
		return
	}
	newHumanf, newJsonf, err := newLogFiles(path, jsonPath)
	if err != nil {
		AddError(err).Error("unable to create new logfiles")
		return
	}
	oldHumanf := humanf
	oldJsonf := jsonf
	humanf = newHumanf
	jsonf = newJsonf
	human.SetOutput(humanf)
	json.SetOutput(jsonf)
	oldHumanf.Close()
	oldJsonf.Close()
}

func truncateTale(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		AddError(err).Error("could not read dir for logs")
		return
	}
	numFiles := 0
	var oldestFile os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		numFiles++
		fi, err := file.Info()
		if err != nil {
			AddError(err).Error("cannot convert direntry to fileinfo in ls of log dir")
			continue
		}
		if oldestFile == nil || fi.ModTime().Before(oldestFile.ModTime()) {
			oldestFile = fi
		}
	}
	if numFiles >= 10 {
		err = os.Remove(fmt.Sprintf("%s/%s", path, oldestFile.Name()))
		if err != nil {
			AddError(err).Error("unable to remove old log file")
			return
		}
	}
}

const (
	B int64 = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)
