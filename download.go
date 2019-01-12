package main

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// DownloadFile makes a GET request to download a photo/video etc. to a file
func DownloadFile(url, dir, name string) error {
	logrus.Debugf("downloading %s to %s/%s", url, dir, name)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("could not download file")
	}
	if res.ContentLength == 0 {
		logrus.Debugf("tried to download %s, but got nothing", url)
		return nil
	}
	defer res.Body.Close()

	f, err := os.Create(dir + "/" + name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return nil
}
