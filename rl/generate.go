package main

import (
	"encoding/hex"
	"fmt"
	"encoding/json"
	"flag"
	"github.com/xfk/colly"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path"
	"strconv"
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
	}
	for _, ext := range interesting_files {
		if strings.Contains(uri, ext) {
			return true
		}
	}

	return false
}

func TempFileName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

type Path struct {
	XPaths [][]map[string]string
	Pages []url.URL
	Target string
}

func writeRes(out string, start string, xpaths [][]map[string]string, pages []url.URL, target string) {
	startUrl, _ := url.Parse(start)
	/*
	for i := 0; i < len(pages); i++ {
		//fmt.Printf("%+v\n", pages[i])
		//fmt.Printf("%+v\n", xpaths[i])
	}
	*/

	pages = append([]url.URL{*startUrl}, pages...)
	path := &Path{
		XPaths: xpaths,
		Pages: pages,
		Target: target,
	}
	
	file, err := json.MarshalIndent(path, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(out, file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	name := flag.String("name", "pointstone", "Spider Name")
	start_url := flag.String("url", "https://www.pointstone.com/download/", "Start URL")
	depth := flag.Int("depth", 10, "Max depth")
	credentials := flag.String("credentials", "", "Login credentials admin:admin")
	login := flag.String("login", "", "Login Page")
	flag.Parse()

	domain, _ := url.Parse(*start_url)

	dir := path.Join("../data/" + *name, strconv.Itoa(*depth))
	os.MkdirAll(dir, 0755)

	c := colly.NewCollector(
		colly.AllowedDomains(domain.Host),
		colly.MaxDepth(*depth),
		colly.UserAgent("2019RLCrawlAThon"),
	)

	if *credentials != "" {
		if *login == "" {
			*login = *start_url
		}

		cred := strings.Split(*credentials, ":")
		err := c.Post(
			*login,
			map[string]string{"username": cred[0], "password": cred[1]},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		xpath := getElementPath(e)

		//fmt.Printf("%#v\n", e.Request.Prev)
		//fmt.Printf("%#v\n", xpaths)

		if isInteresting(e.Request.AbsoluteURL(link)) {
			file := path.Join(dir, TempFileName())
			target_url := e.Request.AbsoluteURL(link)
			pages := e.Request.Prev
			xpaths := append(e.Request.XPath, xpath)

			writeRes(file, *start_url, xpaths, pages, target_url)
		} else {
			e.Request.Visit(e.Request.AbsoluteURL(link), xpath)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(*start_url)
}
