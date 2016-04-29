package duokan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

const (
	maxRetry      = 5
	retryInterval = 200 * time.Millisecond
	pageInterval  = 300 * time.Millisecond
)

// PageInfo duokan page info
type PageInfo struct {
	// Pid is iss format of bid.
	Pid string `json:"page_id"`
	Num int    `json:"page_number"`
}

// BookInfo duokan book info, including name, iss and so on.
type BookInfo struct {
	Title    string     `json:"title"`
	Pages    []PageInfo `json:"pages"` // Pages[i].Num == i+1
	Revision string     `json:"revision"`
	// chapters
}

// JsResp received response which contains page js address.
type JsResp struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// Proxy fetches url page interface.
type Proxy interface {
	// FetchURL fetches url content, store content in a file and return file name.
	FetchURL(string) (string, error)
}

// BsonProxy fetches url with cookie and decode it.
type BsonProxy struct {
	c *http.Client
}

// NewBsonProxy set jar as cookie.
func NewBsonProxy(jar http.CookieJar) *BsonProxy {
	return &BsonProxy{
		c: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				Dial: func(n, addr string) (net.Conn, error) {
					conn, err := net.Dial(n, addr)
					if err == nil {
						conn.SetDeadline(time.Now().Add(10 * time.Second))
					}
					return conn, err
				},
				DisableKeepAlives: true,
			},
			Timeout: 20 * time.Second,
		},
	}
}

// FetchURL fetches url and saves as a file, it includes retry and decoding bson.
func (p *BsonProxy) FetchURL(ref string) (string, error) {
	for tryCount := 0; tryCount < maxRetry; tryCount++ {
		if tryCount > 0 {
			log.Warnf("Retry after %d times", tryCount)
			time.Sleep(retryInterval)
		}
		resp, err := p.c.Get(ref)
		if err != nil {
			log.Error(err)
			continue
		}
		f, err := ioutil.TempFile(os.TempDir(), "duokan")
		if err != nil {
			log.Error(err)
			continue
		}
		// Why use ioutil.ReadAll doesn't work? It seems resp.Body can't Close,
		// so reader is blocking.
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Errorf("Get error url[%s], err[%s]", ref, err)
			continue
		}
		err = f.Close()
		if err != nil {
			log.Error(err)
			continue
		}
		return f.Name(), nil
	}
	return "", errors.Errorf("reach max retry time [%d]", maxRetry)
}

// Librarian This guy manager all books.
type Librarian struct {
	proxy Proxy
	books map[string]BookInfo // Now books is read only.
}

// GetBookInfo get book info by bid.
func (l *Librarian) GetBookInfo(bid string) (BookInfo, error) {
	if res, ok := l.books[bid]; ok {
		return res, nil
	}
	var res BookInfo
	url := fmt.Sprintf("http://www.duokan.com/reader/book_info/%s/medium", bid)
	jsonData, err := l.DecodeURL(url)
	if err != nil {
		return res, errors.Trace(err)
	}
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	err = dec.Decode(&res)
	if err != nil {
		return res, errors.Trace(err)
	}
	// TODO: make pages sorted and check some pages is missing.
	l.books[bid] = res
	return res, nil
}

// DecodeURL fetches url and decodes it through calling node dkbson.
func (l *Librarian) DecodeURL(ref string) ([]byte, error) {
	fileName, err := l.proxy.FetchURL(ref)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return runDecode(fileName)
}

// SaveBook saves book content as file by bid.
func (l *Librarian) SaveBook(bid, outFile string) error {
	bInfo, err := l.GetBookInfo(bid)
	if err != nil {
		return errors.Trace(err)
	}
	outF, err := os.Create(outFile)
	if err != nil {
		return errors.Trace(err)
	}
	for idx, page := range bInfo.Pages {
		log.Debugf("Fetching %d/%d page...", idx+1, len(bInfo.Pages))
		content, err := l.fetchPageContent(bid, page.Pid)
		_, err = outF.WriteString(content)
		if err != nil {
			return errors.Trace(err)
		}
		time.Sleep(pageInterval)
	}
	return nil
}

// fetchPageContent gets iss of bid and convert it to content.
func (l *Librarian) fetchPageContent(bid, iss string) (string, error) {
	js, err := l.iss2Js(bid, iss)
	if err != nil {
		return "", errors.Trace(err)
	}
	page, err := l.getPageContent(js)
	if err != nil {
		return "", errors.Trace(err)
	}
	content, err := page.GenerateContent()
	if err != nil {
		return "", errors.Trace(err)
	}
	return content, nil
}

// iss2Js convert iss to js address.
func (l *Librarian) iss2Js(bid, iss string) (string, error) {
	url := fmt.Sprintf("http://www.duokan.com/reader/page/%s/%s", bid, iss)
	jsonData, err := l.DecodeURL(url)
	if err != nil {
		return "", errors.Trace(err)
	}
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	var jsResp JsResp
	err = dec.Decode(&jsResp)
	if err != nil {
		return "", errors.Trace(err)
	}
	if jsResp.Status != "ok" {
		return "", errors.Errorf("get js response error[%s]", jsResp.Status)
	}
	return jsResp.URL, nil
}

// getPageContent get js address and extract content to string.
func (l *Librarian) getPageContent(js string) (*PageContent, error) {
	fileName, err := l.proxy.FetchURL(js)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// Truncate file content to only base64 code.
	modifyFile(fileName, func(data []byte) []byte {
		out := bytes.Split(data, []byte("'"))
		if len(out) < 2 {
			return nil
		}
		return out[1]
	})
	jsonData, err := runDecode(fileName)
	var page PageContent
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	err = dec.Decode(&page)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &page, nil
}

// NewLibrarian use proxy as getting url proxy.
func NewLibrarian(proxy Proxy) *Librarian {
	if proxy == nil {
		proxy = NewBsonProxy(nil)
	}
	return &Librarian{
		proxy: proxy,
		books: make(map[string]BookInfo),
	}
}

func runDecode(fileName string) ([]byte, error) {
	cmd := exec.Command("node", "./decode.js", fileName)
	// Get package main path and set it as work directory.
	if _, currentFile, _, ok := runtime.Caller(1); ok {
		cmd.Dir = path.Join(path.Dir(currentFile), "..")
	}
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Maskf(err, "stderr: %s", data)
	}
	err = os.Remove(fileName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

// modifyFile modifies file with fn and save back to fileName.
func modifyFile(fileName string, fn func([]byte) []byte) error {
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return errors.Trace(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Trace(err)
	}
	data = fn(data)
	if _, err = f.Seek(0, 0); err != nil {
		return errors.Trace(err)
	}
	n, err := f.Write(data)
	if err != nil {
		return errors.Trace(err)
	}
	if err = f.Truncate(int64(n)); err != nil {
		return errors.Trace(err)
	}
	if err = f.Close(); err != nil {
		return errors.Trace(err)
	}
	return nil
}
