package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DSiDLData struct {
	Dsidl  int               `json:"dsidl"`
	Script map[string]string `json:"script"`
}

func main() {
	// create url flag
	url := flag.String("url", "", "url to root of webserver where files are hosted")
	dir := flag.String("dir", "", "directory of all files")
	flag.Parse()

	// if dir is empty, use current dir
	if *dir == "" {
		// if arg 1 exists, use it as dir
		if len(os.Args) > 1 {
			*dir = os.Args[1]
		} else {
			*dir = "./"
		}
	}

	// if url is empty, then fill it with default example.com
	if *url == "" {
		if _, err := os.Stat(*dir + ".url"); !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(*dir + ".url")
			if err != nil {
				fmt.Println("File reading error", err)
				return
			}
			dataStr := string(data)
			if !strings.HasSuffix(dataStr, "/") {
				dataStr = dataStr + "/"
			}
			*url = dataStr
		} else {
			*url = "http://example.com/"
		}

	}

	// create a new dsidl struct
	dsidl := DSiDLData{}
	dsidl.Dsidl = 1
	dsidl.Script = make(map[string]string)

	// walk the directory and add all files to the dsidl.Script map
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, *dir)

		if path == "" || strings.Contains(path, "dsidl.json") || info.Name() == ".url" {
			return nil
		}
		// ignore .git and dsidl.json
		if info.IsDir() && info.Name() == "dsidl.json" || info.Name() == ".git" {
			return filepath.SkipDir
		}

		// add file to dsidl.Script map
		dsidl.Script[path] = *url + path
		return nil
	})

	// if dsidl.json exists, delete it
	if _, err := os.Stat(*dir + "/dsidl.json"); err == nil {
		os.Remove(*dir + "/dsidl.json")
	}
	// write the dsidl struct to a json file in dir
	jsonFile, err := os.Create(*dir + "/dsidl.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	//convert the dsidl struct to json
	jsonData, err := json.Marshal(dsidl)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile(*dir+"/dsidl.json", []byte(jsonData), 0644)
	fmt.Println(string(jsonData))
}
