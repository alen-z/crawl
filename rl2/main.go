package main

import (
	//"fmt"
	"github.com/xfk/colly"
	"golang.org/x/net/html"
	//"net/url"
	//"strings"
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

func main() {
	xpaths := make(map[string][]map[string]string)

	c := colly.NewCollector(
		colly.AllowedDomains("pointstone.com", "www.pointstone.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		xpath := getElementPath(e)
		xpaths[e.Request.URL.String() + e.Request.AbsoluteURL(link)] = xpath

		e.Request.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.pointstone.com/download/")
}
