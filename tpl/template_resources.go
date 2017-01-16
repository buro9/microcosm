// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	remoteURLLock = &remoteLock{m: make(map[string]*sync.Mutex)}
	resSleep      = time.Second * 2 // if JSON decoding failed sleep for n seconds before retrying
	resRetries    = 1               // number of retries to load the JSON from URL or local file system
)

type remoteLock struct {
	sync.RWMutex
	m map[string]*sync.Mutex
}

// URLLock locks an URL during download
func (l *remoteLock) URLLock(url string) {
	l.Lock()
	if _, ok := l.m[url]; !ok {
		l.m[url] = &sync.Mutex{}
	}
	l.Unlock() // call this Unlock before the next lock will be called. NFI why but defer doesn't work.
	l.m[url].Lock()
}

// URLUnlock unlocks an URL when the download has been finished. Use only in defer calls.
func (l *remoteLock) URLUnlock(url string) {
	l.RLock()
	defer l.RUnlock()
	if um, ok := l.m[url]; ok {
		um.Unlock()
	}
}

// getCacheFileID returns the cache ID for a string
func getCacheFileID(id string) string {
	return viper.GetString("cacheDir") + url.QueryEscape(id)
}

// resGetCache returns the content for an ID from the file cache or an error
// if the file is not found returns nil,nil
func resGetCache(id string, fs afero.Fs, ignoreCache bool) ([]byte, error) {
	if ignoreCache {
		return nil, nil
	}
	fID := getCacheFileID(id)
	isExists, err := helpers.Exists(fID, fs)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, nil
	}

	return afero.ReadFile(fs, fID)

}

// resWriteCache writes bytes to an ID into the file cache
func resWriteCache(id string, c []byte, fs afero.Fs, ignoreCache bool) error {
	if ignoreCache {
		return nil
	}
	fID := getCacheFileID(id)
	f, err := fs.Create(fID)
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to create file: " + fID)
	}
	defer f.Close()
	n, err := f.Write(c)
	if n == 0 {
		return errors.New("No bytes written to file: " + fID)
	}
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to write to file: " + fID)
	}
	return nil
}

func resDeleteCache(id string, fs afero.Fs) error {
	return fs.Remove(getCacheFileID(id))
}

// resGetRemote loads the content of a remote file. This method is thread safe.
func resGetRemote(url string, fs afero.Fs, hc *http.Client) ([]byte, error) {
	c, err := resGetCache(url, fs, viper.GetBool("ignoreCache"))
	if c != nil && err == nil {
		return c, nil
	}
	if err != nil {
		return nil, err
	}

	// avoid race condition with locks, block other goroutines if the current url is processing
	remoteURLLock.URLLock(url)
	defer func() { remoteURLLock.URLUnlock(url) }()

	// avoid multiple locks due to calling resGetCache twice
	c, err = resGetCache(url, fs, viper.GetBool("ignoreCache"))
	if c != nil && err == nil {
		return c, nil
	}
	if err != nil {
		return nil, err
	}

	jww.INFO.Printf("Downloading: %s ...", url)
	res, err := hc.Get(url)
	if err != nil {
		return nil, err
	}
	c, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = resWriteCache(url, c, fs, viper.GetBool("ignoreCache"))
	if err != nil {
		return nil, err
	}
	jww.INFO.Printf("... and cached to: %s", getCacheFileID(url))
	return c, nil
}

// resGetLocal loads the content of a local file
func resGetLocal(url string, fs afero.Fs) ([]byte, error) {
	filename := filepath.Join(viper.GetString("workingDir"), url)
	if e, err := helpers.Exists(filename, fs); !e {
		return nil, err
	}

	return afero.ReadFile(fs, filename)

}

// resGetResource loads the content of a local or remote file
func resGetResource(url string) ([]byte, error) {
	if url == "" {
		return nil, nil
	}
	if strings.Contains(url, "://") {
		return resGetRemote(url, hugofs.Source(), http.DefaultClient)
	}
	return resGetLocal(url, hugofs.Source())
}

// getJSON expects one or n-parts of a URL to a resource which can either be a local or a remote one.
// If you provide multiple parts they will be joined together to the final URL.
// GetJSON returns nil or parsed JSON to use in a short code.
func getJSON(urlParts ...string) interface{} {
	var v interface{}
	url := strings.Join(urlParts, "")

	for i := 0; i <= resRetries; i++ {
		c, err := resGetResource(url)
		if err != nil {
			jww.ERROR.Printf("Failed to get json resource %s with error message %s", url, err)
			return nil
		}

		err = json.Unmarshal(c, &v)
		if err != nil {
			jww.ERROR.Printf("Cannot read json from resource %s with error message %s", url, err)
			jww.ERROR.Printf("Retry #%d for %s and sleeping for %s", i, url, resSleep)
			time.Sleep(resSleep)
			resDeleteCache(url, hugofs.Source())
			continue
		}
		break
	}
	return v
}

// parseCSV parses bytes of CSV data into a slice slice string or an error
func parseCSV(c []byte, sep string) ([][]string, error) {
	if len(sep) != 1 {
		return nil, errors.New("Incorrect length of csv separator: " + sep)
	}
	b := bytes.NewReader(c)
	r := csv.NewReader(b)
	rSep := []rune(sep)
	r.Comma = rSep[0]
	r.FieldsPerRecord = 0
	return r.ReadAll()
}

// getCSV expects a data separator and one or n-parts of a URL to a resource which
// can either be a local or a remote one.
// The data separator can be a comma, semi-colon, pipe, etc, but only one character.
// If you provide multiple parts for the URL they will be joined together to the final URL.
// GetCSV returns nil or a slice slice to use in a short code.
func getCSV(sep string, urlParts ...string) [][]string {
	var d [][]string
	url := strings.Join(urlParts, "")

	var clearCacheSleep = func(i int, u string) {
		jww.ERROR.Printf("Retry #%d for %s and sleeping for %s", i, url, resSleep)
		time.Sleep(resSleep)
		resDeleteCache(url, hugofs.Source())
	}

	for i := 0; i <= resRetries; i++ {
		c, err := resGetResource(url)

		if err == nil && !bytes.Contains(c, []byte(sep)) {
			err = errors.New("Cannot find separator " + sep + " in CSV.")
		}

		if err != nil {
			jww.ERROR.Printf("Failed to read csv resource %s with error message %s", url, err)
			clearCacheSleep(i, url)
			continue
		}

		if d, err = parseCSV(c, sep); err != nil {
			jww.ERROR.Printf("Failed to parse csv file %s with error message %s", url, err)
			clearCacheSleep(i, url)
			continue
		}
		break
	}
	return d
}
