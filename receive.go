package main

import (
	"path"
	"time"

	mlib "github.com/wsva/monitor_lib_go"
)

func receive(body []byte) {
	latestMessageAll, err := mlib.GetMRListFromJSON(body)
	if err != nil {
		return
	}
	var latestMessageFiltered []mlib.MR
	for _, v := range latestMessageAll {
		if mlib.GetFilterResult(mainConfig.FilterList, v) {
			latestMessageFiltered = append(latestMessageFiltered, v)
		}
	}
	latestMessage = latestMessageFiltered
}

func writeToLog() {
	logPath := mainConfig.DirectoryLog
	err := mlib.WriteMRListToFile(latestMessage,
		path.Join(logPath, time.Now().Format("20060102_150405")+".json"))
	if err != nil {
		return
	}
}
