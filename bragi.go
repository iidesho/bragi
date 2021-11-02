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
	return logData{err: e, level: INFO}
}

type Stringer interface {
	String() string
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

func (ld logData) Print(a ...interface{}) {
	if level > ld.level {
		return
	}
	humanString, jsonString := ld.format(fmt.Sprint(a))
	human.Print(humanString)
	if folder == "" && false {
		return
	}
	json.Print(jsonString)
}

func (ld logData) Printf(format string, a ...interface{}) {
	ld.Print(fmt.Sprintf(format, a))
}

func Printf(format string, a ...interface{}) {
	AddError(nil).Printf(format, a)
}

func (ld logData) Println(a ...interface{}) {
	if level > ld.level {
		return
	}
	humanString, jsonString := ld.format(fmt.Sprint(a))
	human.Println(humanString)
	if folder == "" && false {
		return
	}
	json.Println(jsonString)
}

func Println(a ...interface{}) {
	AddError(nil).Println(a)
}

func (ld logData) Debug(a ...interface{}) {
	ld.level = DEBUG
	ld.Print(a)
}

func Debug(a ...interface{}) {
	AddError(nil).Debug(a)
}

func (ld logData) Info(a ...interface{}) {
	ld.level = INFO
	ld.Print(a)
}

func Info(a ...interface{}) {
	AddError(nil).Info(a)
}

func (ld logData) Notice(a ...interface{}) {
	ld.level = NOTICE
	ld.Print(a)
}

func Notice(a ...interface{}) {
	AddError(nil).Notice(a)
}

func (ld logData) Warning(a ...interface{}) {
	ld.level = WARNING
	ld.Print(a)
}

func Warning(a ...interface{}) {
	AddError(nil).Warning(a)
}

func (ld logData) Error(a ...interface{}) {
	ld.level = ERROR
	ld.Print(a)
}

func Error(a ...interface{}) {
	AddError(nil).Error(a)
}

func (ld logData) Crit(a ...interface{}) {
	ld.level = CRIT
	ld.Print(a)
	os.Exit(1)
}

func Crit(a ...interface{}) {
	AddError(nil).Crit(a)
}

func (ld logData) Fatal(a ...interface{}) {
	ld.Crit(a)
}

func Fatal(a ...interface{}) {
	AddError(nil).Fatal(a)
}

func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a)) // Tmp until i need more
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
