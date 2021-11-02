package bragi

import (
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
	json   = log.New(os.Stdout, "", 0)
	folder string
	prefix = "Default"
	level  = Level(0)
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
	return []string{"Debug", "Info", "Notice", "Warning", "Error", "Crit"}[l]
}

func SetPrefix(p string) {
	prefix = p
}

func SetOutputFolder(path string) func() {
	folder = path
	if !fileExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil
		}
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.log", path, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	human = log.New(f, "", 0)
	jsonPath := path + "/json"
	if !fileExists(jsonPath) {
		err := os.MkdirAll(jsonPath, 0755)
		if err != nil {
			f.Close()
			return nil
		}
	}
	jsonf, err := os.OpenFile(fmt.Sprintf("%s/%s.log", jsonPath, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		f.Close()
		return nil
	}
	json.SetOutput(jsonf)
	return func() {
		f.Close()
		jsonf.Close()
	}
}

type logData struct {
	err   error
	level Level
}

func AddError(e error) logData {
	return logData{err: e}
}

func (ld logData) format(s string) (human, json string) {
	human = time.Now().Format("15:04:05 MST")
	json = fmt.Sprintf(`{"time":"%s"`, time.Now().Format("15:04:05 MST"))
	if ld.level == DEBUG {
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
		human = fmt.Sprintf("%s %s:%d/%s", human, path[len(path)-1], line, runtime.FuncForPC(function).Name())
		json = fmt.Sprintf(`%s,"file":"%s","line":%d,"function":"%s"`, json, path[len(path)-1], line, runtime.FuncForPC(function).Name())
	}
	human = fmt.Sprintf("%s [%s]%s", human, ld.level, s)
	json = fmt.Sprintf(`%s,"level":"%s","text":"%s"`, json, ld.level, s)
	if ld.err != nil {
		human = fmt.Sprintf("%s. Err: %v", human, ld.err)
		json = fmt.Sprintf(`%s,"error":"%v"`, json, ld.err)
	}
	json = fmt.Sprintf("%s}", json)
	return
}

func (ld logData) Print(s string) {
	if level > ld.level {
		return
	}
	humanString, jsonString := ld.format(s)
	human.Printf("%s\n", humanString)
	if folder == "" && false {
		return
	}
	json.Printf("%s\n", jsonString)
}

func (ld logData) Debug(s string) {
	ld.level = DEBUG
	ld.Print(s)
}

func Debug(s string) {
	AddError(nil).Debug(s)
}

func (ld logData) Info(s string) {
	ld.level = INFO
	ld.Print(s)
}

func Info(s string) {
	AddError(nil).Info(s)
}

func (ld logData) Notice(s string) {
	ld.level = NOTICE
	ld.Print(s)
}

func Notice(s string) {
	AddError(nil).Notice(s)
}

func (ld logData) Warning(s string) {
	ld.level = WARNING
	ld.Print(s)
}

func Warning(s string) {
	AddError(nil).Warning(s)
}

func (ld logData) Error(s string) {
	ld.level = ERROR
	ld.Print(s)
}

func Error(s string) {
	AddError(nil).Error(s)
}

func (ld logData) Crit(s string) {
	ld.level = CRIT
	ld.Print(s)
	os.Exit(1)
}

func Crit(s string) {
	AddError(nil).Crit(s)
}

func (ld logData) Fatal(s string) {
	ld.Crit(s)
}

func Fatal(s string) {
	AddError(nil).Fatal(s)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
