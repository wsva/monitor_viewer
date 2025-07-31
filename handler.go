package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	wl_net "github.com/wsva/lib_go/net"
	mlib "github.com/wsva/monitor_lib_go"
)

func handleRoot(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/current", http.StatusFound)
	}
	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		go receive(body)
		go writeToLog()
	}
}

func handleNavigation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		contentBytes, err := os.ReadFile("template/html/navigation.html")
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			content := strings.ReplaceAll(string(contentBytes),
				"__ReplaceWithSchemaHost__", wl_net.GetSchemaAndHost(r))
			io.WriteString(w, content)
		}
	}
}

func handleCurrent(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		query := r.URL.Query()
		queryType := ""
		if len(query["type"]) > 0 {
			queryType = query["type"][0]
		}
		switch queryType {
		case ViewTypeWarning, ViewTypeSimple, ViewTypeSimplest, ViewTypeDetail:
			html, err := getFinalHTML(HTMLTypeCurrent, queryType, latestMessage)
			if err != nil {
				io.WriteString(w, err.Error())
			} else {
				io.WriteString(w, html)
			}
			return
		default:
			html, err := getFinalHTML(HTMLTypeCurrent, ViewTypeNormal, latestMessage)
			if err != nil {
				io.WriteString(w, err.Error())
			} else {
				io.WriteString(w, html)
			}
			return
		}
	}
}

func handleHistory(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		query := r.URL.Query()

		var queryStart, queryEnd, queryNear, queryLogfile, queryType string

		if len(query["near"]) > 0 {
			queryNear = query["near"][0]
		}
		if len(query["start"]) > 0 {
			queryStart = query["start"][0]
		}
		if len(query["end"]) > 0 {
			queryEnd = query["end"][0]
		}
		if len(query["logfile"]) > 0 {
			queryLogfile = query["logfile"][0]
		}
		if len(query["type"]) > 0 {
			queryType = query["type"][0]
		}

		//if logfile is specified, then show it
		if queryLogfile != "" {
			mrList, err := mlib.GetMRListFromFile(
				path.Join(mainConfig.DirectoryLog, queryLogfile))
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			switch queryType {
			case ViewTypeWarning, ViewTypeSimple, ViewTypeSimplest, ViewTypeDetail:
				html, err := getFinalHTML(HTMLTypeHistory, queryType, mrList)
				if err != nil {
					io.WriteString(w, err.Error())
				} else {
					io.WriteString(w, html)
				}
			default:
				html, err := getFinalHTML(HTMLTypeHistory, ViewTypeNormal, mrList)
				if err != nil {
					io.WriteString(w, err.Error())
				} else {
					io.WriteString(w, html)
				}

			}
			return
		}

		//if near specified, then use it
		if queryNear != "" {
			logtime, err := time.Parse("20060102_1504", queryNear)
			if err != nil {
				io.WriteString(w, "parse logtime error")
				return
			}
			logfile, err := getLogfileNear(logtime)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			logfile = path.Base(logfile)
			http.Redirect(w, r,
				wl_net.GetSchemaAndHost(r)+"/history?logfile="+logfile+"&type=warning", http.StatusFound)
			return
		}

		//if start and end are specified, then use it
		if queryStart != "" && queryEnd != "" {
			html, err := getHistoryIndex(queryStart, queryEnd, wl_net.GetSchemaAndHost(r))
			if err != nil {
				io.WriteString(w, err.Error())
			} else {
				io.WriteString(w, html)
			}
			return
		}

		io.WriteString(w, "wrong usage")
	}
}

func handleMRList(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "GET" {
		query := r.URL.Query()

		if len(query["latest"]) > 0 {
			jsonBytes, _ := json.Marshal(latestMessage)
			w.Write(jsonBytes)
			return
		}

		var queryNear string
		if len(query["near"]) > 0 {
			queryNear = query["near"][0]
		} else {
			io.WriteString(w, "near(20210801_1030) or latest is needed")
			return
		}

		logtime, err := time.Parse("20060102_1504", queryNear)
		if err != nil {
			io.WriteString(w, "parse logtime error")
			return
		}

		logfile, err := getLogfileNear(logtime)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		contentBytes, err := os.ReadFile(logfile)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		w.Write(contentBytes)
	}
}
