package internal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const fileUrl = "https://github.com/cjtoolkit/ignition/archive/master.tar.gz"

func getFileFromGitHub() *tar.Reader {
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

	gz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	return tar.NewReader(gz)
}
