package bragi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func StartRotate(done <-chan func()) {
	ticker := time.NewTicker(getNextTick())
	ticker2 := time.NewTicker(4 * time.Second) // time.Minute * 5)
	go func() {                                //Handle panicks
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				rotateLog()
				ticker.Reset(getNextTick())
			case <-ticker2.C:
				hstat, herr := humanf.Stat()
				jstat, jerr := jsonf.Stat()
				if herr != nil {
					AddError(herr).Warning("Could not get human file stats while checking if it should be rotated")
				}
				if jerr != nil {
					AddError(jerr).Warning("Could not get json file stats while checking if it should be rotated")
				}
				if !(herr == nil && hstat.Size() > 11<<20 || jerr == nil && jstat.Size() > 11<<20) { // Bitshifting to megabyte so i dont have to write the whole number
					continue // Continuing if bouth files are smaller than 11MB
				}
				rotateLog()
			}
		}
	}()
}

func getNextTick() time.Duration {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(24*time.Hour - 5*time.Microsecond).Sub(now)
}

func rotateLog() {
	newFilePrefix := fmt.Sprintf("%s-%s", prefix, time.Now().Format("2006.01.02"))
	stat, err := humanf.Stat()
	if err != nil {
		AddError(err).Warning("Could not get human file stats while rotating logs")
	} else if stat.Size() == 0 {
		Info("Logs did not rotate because human file size was zero 0")
		return
	}
	err = os.Rename(fmt.Sprintf("%s/%s.log", folder, prefix), fmt.Sprintf("%s/%s.%d.log", folder, newFilePrefix, numLogfiles(folder, newFilePrefix)))
	if err != nil {
		AddError(err).Error("Moving human readable log failed while rotating logs")
		return
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.log", folder, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	human.SetOutput(f)
	humanf.Close()
	humanf = f
	stat, err = jsonf.Stat()
	if err != nil {
		AddError(err).Warning("Could not get json file stats while rotating logs")
	} else if stat.Size() == 0 {
		Info("Json logs did not rotate because json file size was zero 0")
		return
	}
	jsonFolder := folder + "/json"
	err = os.Rename(fmt.Sprintf("%s/%s.log", jsonFolder, prefix), fmt.Sprintf("%s/%s.%d.log", jsonFolder, newFilePrefix, numLogfiles(jsonFolder, newFilePrefix)))
	if err != nil {
		AddError(err).Error("Moving json log failed while rotating logs")
		return
	}
	jf, err := os.OpenFile(fmt.Sprintf("%s/%s.log", jsonFolder, prefix), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		f.Close()
		return
	}
	json.SetOutput(jf)
	jsonf.Close()
	jsonf = jf
}

func numLogfiles(dir, prefix string) (num int) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".log") {
			continue
		}
		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}
		num++
	}
	return
}
