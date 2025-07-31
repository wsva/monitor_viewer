package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func getHistoryIndex(start, end, schemaandhost string) (string, error) {
	reg := regexp.MustCompile(`^\d{8}_\d{6}$`)
	if !reg.MatchString(start) || !reg.MatchString(end) {
		return "", errors.New("invalid start or end format, please use 20210101_000000")
	}

	entryList, err := os.ReadDir(mainConfig.DirectoryLog)
	if err != nil {
		return "", err
	}
	var logfileList []string
	for _, v := range entryList {
		if v.IsDir() {
			continue
		}
		if v.Name() >= start && v.Name() <= end {
			logfileList = append(logfileList, v.Name())
		}
	}
	if logfileList == nil {
		return "", errors.New("no history found between start and end")
	}
	sort.Slice(logfileList, func(i, j int) bool {
		return logfileList[i] < logfileList[j]
	})

	var result strings.Builder
	for _, v := range logfileList {
		result.WriteString(
			fmt.Sprintf("<a href=\"%s/history?logfile=%s&type=warning\">%s</a><br/>\n",
				schemaandhost, v, v))
	}
	contentBytes, err := ioutil.ReadFile("template/html/historyindex.html")
	if err != nil {
		return "", err
	}
	reg = regexp.MustCompile(`__ReplaceWithSection__`)
	return reg.ReplaceAllString(string(contentBytes), result.String()), nil
}

//get a logfile nearest in 10 minites
func getLogfileNear(logtime time.Time) (string, error) {
	logtimeList := []string{logtime.Format("20060102_1504")}
	for i := time.Duration(1); i < 11; i++ {
		t1 := logtime.Add(i * time.Minute)
		t2 := logtime.Add(-1 * i * time.Minute)
		logtimeList = append(logtimeList, t1.Format("20060102_1504"))
		logtimeList = append(logtimeList, t2.Format("20060102_1504"))
	}

	logtimeRegexp := strings.Join(logtimeList, "|")
	reg := regexp.MustCompile(logtimeRegexp)

	var logfilelist []string
	err := filepath.Walk(mainConfig.DirectoryLog, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() && reg.MatchString(path) {
			logfilelist = append(logfilelist, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	for _, v1 := range logtimeList {
		for _, v2 := range logfilelist {
			if strings.Contains(v2, v1) {
				return v2, nil
			}
		}
	}

	return "", errors.New("no logfile found")
}
