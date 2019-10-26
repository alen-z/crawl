package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"path"
)

type Path struct {
	XPaths [][]map[string]string
	Pages []url.URL
	Target string
}

type MPath struct {
	Path string
	Pos int
}

type Minimum struct {
	XPaths []map[string]string
	AllowedPaths []MPath
	RejectedPaths []MPath
}

func equalElement(a map[string]string, b map[string]string) bool {
	return a["tag"] == b["tag"] && a["id"] == b["id"] && a["class"] == b["class"]
}

func main() {
	name := flag.String("name", "pointstone", "Spider Name")
	depth := flag.Int("depth", 10, "Max depth")
	flag.Parse()

	dir := path.Join("../data/" + *name, strconv.Itoa(*depth))

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	
	/*
	atMin := &Minimum{
		XPaths: xpaths,
		AllowedPaths: allowed,
		RejectedPaths: rejected,
	}
	*/
	atMin := &Minimum{
		XPaths: make([]map[string]string, 0),
		AllowedPaths: make([]MPath, 0),
		RejectedPaths: make([]MPath, 0),
	}

	for _, file := range files {
		inFile, _ := ioutil.ReadFile(path.Join(dir, file.Name()))
		path := &Path{}
		json.Unmarshal([]byte(inFile), path)
		fmt.Println("file je: " + file.Name())
		for _, xpaths := range path.XPaths {
			for index, elements := range xpaths {
				if index > 3  || len(xpaths) < 3 {
					atMin.XPaths = append(atMin.XPaths, elements)
				}
			}
		}
		for i, page := range path.Pages {
			mpage := MPath{
				Path: page.String(),
				Pos: i,
			}
			//fmt.Println("page:")
			//fmt.Println(mpage)
			atMin.AllowedPaths = append(atMin.AllowedPaths, mpage)
		}
	}

	
	file, err := json.MarshalIndent(atMin, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	outFile := path.Join("../data/" + *name, *name + strconv.Itoa(*depth))
	//fmt.Println(outFile, file)
	err = ioutil.WriteFile(outFile, file, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
