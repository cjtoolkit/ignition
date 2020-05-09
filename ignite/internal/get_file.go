package internal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	fileUrl = "https://github.com/cjtoolkit/ignition/archive/master.tar.gz"
	expires = 5 * time.Minute
)

func getFile() *tar.Reader {
	r, err := getFromCache()
	if err == nil {
		return r
	}

	return getFileFromGitHubByHttp()
}

func getFileFromGitHubByHttp() *tar.Reader {
	client := &http.Client{Timeout: 30 * time.Second}

	res, err := client.Get(fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal("Unable to download the file.")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	saveToCache(b)

	return createReader(b)
}

func createReader(b []byte) *tar.Reader {
	gz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	return tar.NewReader(gz)
}

func saveToCache(b []byte) {
	cacheDir, done := getCacheDir()
	if done {
		return
	}
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			log.Println(err)
			return
		}
	}
	file, err := os.Create(cacheDir + filepath.FromSlash("/cache.tar.gz"))
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		log.Println(err)
		return
	}
}

func getCacheDir() (string, bool) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return "", true
	}
	cacheDir += filepath.FromSlash("/ignite")
	return cacheDir, false
}

func getFromCache() (*tar.Reader, error) {
	cacheDir, done := getCacheDir()
	if done {
		return nil, fmt.Errorf("Could not find cache directory")
	}

	file, err := os.Open(cacheDir + filepath.FromSlash("/cache.tar.gz"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() > stat.ModTime().Add(expires).Unix() {
		return nil, fmt.Errorf("Expired")
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return createReader(b), nil
}
