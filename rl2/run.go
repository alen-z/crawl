package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"github.com/xfk/colly"
	"golang.org/x/net/html"
	"io/ioutil"
	"path"
	"strings"
)

func getElementPath(e *colly.HTMLElement) []map[string]string {
	var p []map[string]string

	s := e.DOM
	for {
		tmp := make(map[string]string)
		if id, exists := s.Attr("id"); exists {
			tmp["id"] = id
		} else {
			tmp["id"] = ""
		}
		if class, exists := s.Attr("class"); exists {
			tmp["class"] = class
		} else {
			tmp["class"] = ""
		}

		if len(s.Nodes) > 0 && s.Nodes[0].Type == html.ElementNode {
			tmp["tag"] = s.Nodes[0].Data
		}

		p = append(p, tmp)

		if tmp["tag"] == "html" {
			break
		}

		s = s.Parent()
	}

	/*
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	*/

	return p
}

type Path struct {
	Path string
	Pos int
}

type Minimum struct {
	XPaths []map[string]string
	AllowedPaths []Path
	RejectedPaths []Path
}

func isInteresting(uri string) bool {
        interesting_files := []string{
                ".exe",
                ".dmg",
                ".tar.gz",
                ".zip",
                ".tgz",
                ".apk",
                ".bin",
                ".bat",
                ".msi",
        }
        for _, ext := range interesting_files {
                if strings.Contains(uri, ext) {
                        return true
                }
        }

        return false
}

func equalElement(a map[string]string, b map[string]string) bool {
	return a["tag"] == b["tag"] && a["id"] == b["id"] && a["class"] == b["class"]
}

func hasMinElement(atMin *Minimum, e *colly.HTMLElement) bool {
	xpath := getElementPath(e)

	for _, minElement := range atMin.XPaths {
		for _, element := range xpath {
			if equalElement(minElement, element) {
				return true
			}
		}
	}

	return false
}

func hasMinPath(atMin *Minimum, path string) bool {
	pathSegments := strings.Split(path, "/")
	for i, seg := range pathSegments {
		for _, minSeg := range atMin.RejectedPaths {
			if seg == minSeg.Path {
				return false
			}
		}
		for _, minSeg := range atMin.AllowedPaths {
			if seg == minSeg.Path && i == minSeg.Pos {
				return true
			}
		}
	}
	return false
}

func shouldFollow(atMin *Minimum, e *colly.HTMLElement) bool {
	return hasMinElement(atMin, e) && hasMinPath(atMin, e.Request.URL.Path)
}

func main() {
	name := flag.String("name", "pointstone", "Spider Name")
	start_url := flag.String("url", "https://www.pointstone.com/download/", "Start URL")
	flag.Parse()

	file, _ := ioutil.ReadFile(path.Join("../data/" + *name, *name))
	atMin := &Minimum{}
	json.Unmarshal([]byte(file), atMin)
	
	c := colly.NewCollector(
		colly.AllowedDomains("pointstone.com"),
	)
	
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")


		if isInteresting(e.Request.AbsoluteURL(link)) {
			fmt.Println(e.Request.AbsoluteURL(link))
		} else if shouldFollow(atMin, e) {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(*start_url)
}
