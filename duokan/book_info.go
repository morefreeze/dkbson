package duokan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

const (
	maxFetch int = 5
)

// PageInfo duokan page info
type PageInfo struct {
	Pid string `json:"page_id"`
	Num int    `json:"page_number"`
}

// BookInfo duokan book info, including name, iss and so on.
type BookInfo struct {
	Title    string      `json:"title"`
	Pages    []*PageInfo `json:"pages"`
	Revision string      `json:"revision"`
	// chapters
}

// Proxy get url page interface.
type Proxy interface {
	getURL(string) ([]byte, error)
}

type DefaultProxy struct {
	c *http.Client
}

// NewDefaultProxy set jar as cookie.
func NewDefaultProxy(jar http.CookieJar) *DefaultProxy {
	return &DefaultProxy{
		c: &http.Client{
			Jar: jar,
		},
	}
}

func (p *DefaultProxy) getURL(url string) ([]byte, error) {
	for tryCount := 0; tryCount < maxFetch; tryCount++ {
		resp, err := p.c.Get(url)
		if err != nil {
			log.Error(err)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			continue
		}
		return data, nil
	}
	return nil, errors.Errorf("reach max retry time [%d]", maxFetch)
}

// Librarian This guy manager all books.
type Librarian struct {
	proxy Proxy
}

// GetBookInfo get book info by bid
func (l *Librarian) GetBookInfo(bid string) (BookInfo, error) {
	var res BookInfo
	url := fmt.Sprintf("http://www.duokan.com/reader/book_info/%s/medium", bid)
	data, err := l.proxy.getURL(url)
	if err != nil {
		return res, errors.Trace(err)
	}
	cmd := exec.Command("node", "decode.js", "/dev/stdin")
	sin, err := cmd.StdinPipe()
	if err != nil {
		return res, errors.Trace(err)
	}
	sin.Write(data)
	sin.Close()
	jsonData, err := cmd.Output()
	if err != nil {
		return res, errors.Trace(err)
	}
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	err = dec.Decode(&res)
	if err != nil {
		return res, errors.Trace(err)
	}
	return res, nil
}

// NewLibrarian use proxy as getting url proxy.
func NewLibrarian(proxy Proxy) *Librarian {
	if proxy == nil {
		proxy = NewDefaultProxy(nil)
	}
	return &Librarian{
		proxy: proxy,
	}
}
