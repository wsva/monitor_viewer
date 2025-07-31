package main

import (
	"encoding/json"
	"os"
	"path"

	wl_fs "github.com/wsva/lib_go/fs"
	mlib "github.com/wsva/monitor_lib_go"
)

const (
	HTMLTypeCurrent = "current"
	HTMLTypeHistory = "history"

	ViewTypeNormal   = "normal"
	ViewTypeWarning  = "warning"
	ViewTypeSimple   = "simple"
	ViewTypeSimplest = "simplest"
	ViewTypeDetail   = "detail"
)

// MainConfig comment
type MainConfig struct {
	Port         string              `json:"Port"`
	Language     string              `json:"Language"`
	DirectoryLog string              `json:"DirectoryLog"`
	FilterList   []mlib.FilterRegexp `json:"FilterList"`
}

var (
	MainConfigFile = "viewer_config.json"
)

var mainConfig MainConfig

var latestMessage []mlib.MR

func initGlobals() error {
	basepath, err := wl_fs.GetExecutableFullpath()
	if err != nil {
		return err
	}

	MainConfigFile = path.Join(basepath, MainConfigFile)

	contentBytes, err := os.ReadFile(MainConfigFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contentBytes, &mainConfig)
	if err != nil {
		return err
	}

	if !path.IsAbs(mainConfig.DirectoryLog) {
		mainConfig.DirectoryLog = path.Join(basepath, mainConfig.DirectoryLog)
	}
	err = wl_fs.CheckDirectoryExistAndCreateIfNot(mainConfig.DirectoryLog)
	if err != nil {
		return err
	}

	return nil
}
