package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ZipFile struct {
	Folders   []string
	Data      map[string]*Folder
	Timestamp time.Time
}

type Folder struct {
	Name  string
	Path  string
	Files []*File
}

type File struct {
	Name     string
	Basename string
	Path     string
	Content  string
}

func getData(input string) (*ZipFile, error) {
	var zipData map[string]string
	var err error

	if strings.HasPrefix(input, "http") {
		// read zip content from url
		zipData, err = readZipFromURL(input)
	} else {
		// read zip content from local file
		zipData, err = readZipFromFile(input)
	}
	if err != nil {
		return nil, err
	}

	zipFile := &ZipFile{
		Timestamp: time.Now(),
		Data:      make(map[string]*Folder),
	}

	// get a 'set' of all paths
	paths := make(map[string]interface{})
	for key, _ := range zipData {
		paths[path.Dir(key)] = nil
	}

	// collect and sort all paths
	var folders []string
	for folder, _ := range paths {
		folders = append(folders, folder)
	}
	sort.Sort(sort.StringSlice(folders))
	zipFile.Folders = folders

	// store root folder/path
	root := zipFile.Folders[0]

	// strip away root folder/path from all folders/paths
	for i := range zipFile.Folders {
		zipFile.Folders[i] = strings.TrimPrefix(zipFile.Folders[i], root)
	}

	// go through all paths, collect their files and create Folder/File structure
	for _, p := range zipFile.Folders {
		folder := &Folder{
			Name: path.Base(p),
			Path: p,
		}

		// collect all files for this path
		var files []*File
		for key, value := range zipData {
			if strings.TrimPrefix(path.Dir(key), root) == p {
				basename := filepath.Base(key)
				file := &File{
					Name:     path.Base(key),
					Basename: strings.TrimSuffix(basename, filepath.Ext(basename)),
					Path:     p,
					Content:  value,
				}
				files = append(files, file)
			}
		}
		folder.Files = files

		// store Folder
		zipFile.Data[p] = folder
	}

	return zipFile, nil
}

func readZipFromURL(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return readZip(bytes.NewReader(data), resp.ContentLength)
}

func readZipFromFile(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return readZip(bytes.NewReader(data), fi.Size())
}

func readZip(data *bytes.Reader, size int64) (map[string]string, error) {
	r, err := zip.NewReader(data, size)
	if err != nil {
		return nil, err
	}

	// store data in a map of {filename:content}
	contents := make(map[string]string, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name, ".md") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		contents[f.Name] = string(data)
	}
	return contents, nil
}
