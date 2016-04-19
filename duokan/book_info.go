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
	"time"

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

// JsResp received response which contains page js address.
type JsResp struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// Item each item is a word or a picture in a page.
type Item struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Pos  []int  `json:"pos"`
	Type string `json:"type"`
	Char string `json:"char"`
}

// PageContent contains many items which represent element in page.
type PageContent struct {
	Items []*Item
}

// Proxy get url page interface.
type Proxy interface {
	// GetURL get url string return file name.
	GetURL(string) (string, error)
}

// DefaultProxy default proxy only keep cookie.
type DefaultProxy struct {
	c *http.Client
}

// NewDefaultProxy set jar as cookie.
func NewDefaultProxy(jar http.CookieJar) *DefaultProxy {
	return &DefaultProxy{
		c: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				Dial: func(n, addr string) (net.Conn, error) {
					conn, err := net.Dial(n, addr)
					if err == nil {
						conn.SetDeadline(time.Now().Add(1 * time.Second))
					}
					return conn, err
				},
				DisableKeepAlives: true,
			},
		},
	}
}

// GetURL get url and save as a file.
func (p *DefaultProxy) GetURL(ref string) (string, error) {
	for tryCount := 0; tryCount < maxFetch; tryCount++ {
		if tryCount > 0 {
			log.Warn(tryCount)
			time.Sleep(time.Second)
		}
		resp, err := p.c.Get(ref)
		if err != nil {
			log.Error(err)
			continue
		}
		defer resp.Body.Close()
		f, err := ioutil.TempFile(os.TempDir(), "duokan")
		if err != nil {
			log.Error(err)
			continue
		}
		// Why use ioutil.ReadAll doesn't work? It seems resp.Body can't Close,
		// so reader is blocking.
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Error(err)
			continue
		}
		err = f.Close()
		if err != nil {
			log.Error(err)
			continue
		}
		return f.Name(), nil
	}
	return "", errors.Errorf("reach max retry time [%d]", maxFetch)
}

// Librarian This guy manager all books.
type Librarian struct {
	proxy Proxy
}

// GetBookInfo get book info by bid.
func (l *Librarian) GetBookInfo(bid string) (BookInfo, error) {
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
	return res, nil
}

// DecodeURL fetch url and decode it with dkbson.
func (l *Librarian) DecodeURL(ref string) ([]byte, error) {
	fileName, err := l.proxy.GetURL(ref)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return runDecode(fileName)
}

// GetBookContent save book content as file by bid.
func (l *Librarian) GetBookContent(bid, outFile string) error {
	bInfo, err := l.GetBookInfo(bid)
	if err != nil {
		return errors.Trace(err)
	}
	outF, err := os.Create(outFile)
	if err != nil {
		return errors.Trace(err)
	}
	for _, page := range bInfo.Pages {
		content, err := l.FetchPageContent(bid, page.Pid)
		_, err = outF.WriteString(content)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// FetchPageContent get iss of bid and convert it to content.
func (l *Librarian) FetchPageContent(bid, iss string) (string, error) {
	js, err := l.iss2Js(bid, iss)
	if err != nil {
		return "", errors.Trace(err)
	}
	content, err := l.getPageContent(js)
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
func (l *Librarian) getPageContent(js string) (string, error) {
	fileName, err := l.proxy.GetURL(js)
	if err != nil {
		return "", errors.Trace(err)
	}
	// Truncate file content to only base64 code.
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return "", errors.Trace(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.Trace(err)
	}
	if _, err = f.Seek(0, 0); err != nil {
		return "", errors.Trace(err)
	}
	data = data[len("duokan_page('") : len(data)-4]
	if _, err = f.Write(data); err != nil {
		return "", errors.Trace(err)
	}
	if err = f.Truncate(int64(len(data))); err != nil {
		return "", errors.Trace(err)
	}
	if err = f.Close(); err != nil {
		return "", errors.Trace(err)
	}
	jsonData, err := runDecode(fileName)
	var page PageContent
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	err = dec.Decode(&page)
	if err != nil {
		return "", errors.Trace(err)
	}
	return constructPage(page)
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

func runDecode(fileName string) ([]byte, error) {
	cmd := exec.Command("node", "../decode.js", fileName)
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

func constructPage(page PageContent) (string, error) {
	return "", nil
}
