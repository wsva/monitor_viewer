package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	mlib "github.com/wsva/monitor_lib_go"
)

const (
	TTOC = `<li><a href="#sec-__ReplaceWithSection__">__ReplaceWithName__</a></li>
`

	TSectionH1 = `<h1 id="sec-__ReplaceWithSection__">__ReplaceWithName__</h1>
__ReplaceWithContent__
`

	TContentSingleLine = `<p>__ReplaceWithName__&emsp;
__ReplaceWithAddress__&emsp;
__ReplaceWithContent__&emsp;
<timestamp>__ReplaceWithTimestamp__</timestamp></p>
`

	TContentMultiLine = `<p>__ReplaceWithName__&emsp;
__ReplaceWithAddress__&emsp;
<timestamp>__ReplaceWithTimestamp__</timestamp></p>
<pre>__ReplaceWithContent__</pre>
`

	TContentSimplest = `<p>__ReplaceWithName__</p>
<pre>__ReplaceWithContent__</pre>
`
)

func getFinalHTML(htmlType, viewType string, mrList []mlib.MR) (string, error) {
	var templateFilename string
	switch htmlType {
	case HTMLTypeCurrent:
		templateFilename = "template/html/current.html"
	case HTMLTypeHistory:
		templateFilename = "template/html/history.html"
	}
	contentBytes, err := os.ReadFile(templateFilename)
	if err != nil {
		return "", err
	}
	regTOC := regexp.MustCompile(`__ReplaceWithTOC__`)
	regSection := regexp.MustCompile(`__ReplaceWithSection__`)
	tocAll, sectionAll := getMRListHTML(viewType, mrList)
	htmlContent := regTOC.ReplaceAllString(string(contentBytes), tocAll)
	htmlContent = regSection.ReplaceAllString(htmlContent, sectionAll)
	return htmlContent, nil
}

func getMRListHTML(viewType string, mrList []mlib.MR) (string, string) {
	mrMap := sortMRMap(getMRMapByViewType(viewType, mrList))
	var tocAll, sectionAll string
	regName := regexp.MustCompile(`__ReplaceWithName__`)
	regAddress := regexp.MustCompile(`__ReplaceWithAddress__`)
	regTimestamp := regexp.MustCompile(`__ReplaceWithTimestamp__`)
	regContent := regexp.MustCompile(`__ReplaceWithContent__`)
	regSection := regexp.MustCompile(`__ReplaceWithSection__`)
	for _, mt := range mtypeConfig.MTypeList {
		if _, exist := mrMap[mt.ID]; !exist {
			continue
		}
		toc := regSection.ReplaceAllString(TTOC, mt.ID)
		toc = regName.ReplaceAllString(toc, mt.Name)
		tocAll += toc
		sectionH1 := regSection.ReplaceAllString(TSectionH1, mt.ID)
		sectionH1 = regName.ReplaceAllString(sectionH1, mt.Name)
		var sectionContentAll string
		for _, v := range mrMap[mt.ID] {
			warning := v.GetWarning()
			var sectionContent string
			switch viewType {
			case ViewTypeNormal, ViewTypeWarning, ViewTypeSimple:
				if warning == "ok" {
					sectionContent = TContentSingleLine
					sectionContent = regContent.ReplaceAllString(sectionContent,
						warning)
				} else {
					sectionContent = TContentMultiLine
					sectionContent = regContent.ReplaceAllString(sectionContent,
						deleteNewLineAtEnd(warning))
				}
			case ViewTypeSimplest:
				sectionContent = TContentSimplest
				sectionContent = regContent.ReplaceAllString(sectionContent,
					deleteNewLineAtEnd(warning))
			case ViewTypeDetail:
				sectionContent = TContentMultiLine
				detail := ""
				if v.ErrorString != "" {
					detail = v.ErrorString
				} else {
					md, err := mlib.GetMD(v.MonitorType, v.DetailJSON)
					if err != nil {
						fmt.Println(v)
						fmt.Println(err)
						detail = fmt.Sprint(v, err)
					} else {
						detail = md.DetailString()
					}
				}
				sectionContent = regContent.ReplaceAllString(sectionContent, detail)
			}
			sectionContent = regName.ReplaceAllString(sectionContent, v.Name)
			sectionContent = regAddress.ReplaceAllString(sectionContent, v.Address)
			sectionContent = regTimestamp.ReplaceAllString(sectionContent, v.TimeString)
			sectionContent = regName.ReplaceAllString(sectionContent, v.Name)
			sectionContentAll += sectionContent
		}
		sectionAll += regContent.ReplaceAllString(sectionH1, sectionContentAll)
	}
	return tocAll, sectionAll
}

// get MR Map from mrList By ViewType
func getMRMapByViewType(viewType string, mrList []mlib.MR) map[string][]mlib.MR {
	mrMap := make(map[string][]mlib.MR)
	switch viewType {
	case ViewTypeNormal, ViewTypeDetail:
		for _, v := range mrList {
			mrMap[v.MonitorType] = append(mrMap[v.MonitorType], v)
		}
	case ViewTypeWarning:
		for _, v := range mrList {
			if v.GetWarning() != "ok" {
				mrMap[v.MonitorType] = append(mrMap[v.MonitorType], v)
			}
		}
	case ViewTypeSimple, ViewTypeSimplest:
		for _, v := range mrList {
			if v.GetWarning() != "ok" && !strings.Contains(v.Name, "预生产") {
				mrMap[v.MonitorType] = append(mrMap[v.MonitorType], v)
			}
		}
	}
	return mrMap
}

func sortMRMap(mrMap map[string][]mlib.MR) map[string][]mlib.MR {
	for k := range mrMap {
		sort.Slice(mrMap[k], func(i int, j int) bool {
			list := mrMap[k]
			if list[i].Name == list[j].Name {
				return list[i].Address < list[j].Address
			}
			return list[i].Name < list[j].Name
		})
	}
	return mrMap
}

func deleteNewLineAtEnd(source string) string {
	regNewLine := regexp.MustCompile(`\r|\n`)
	regBlankAtEnd := regexp.MustCompile(`\s+$`)
	lineList := regNewLine.Split(source, -1)
	result := ""
	for _, line := range lineList {
		line = regBlankAtEnd.ReplaceAllString(line, "")
		if len(line) > 0 {
			if result != "" {
				result += "\n"
			}
			result += line
		}
	}
	return result
}
