package duokan

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/juju/errors"
)

// FileCookie read cookie from file.
type FileCookie struct {
	mu      sync.Mutex
	cookies map[string][]*http.Cookie
}

// NewFileCookie parse fileName to cookie like cookie.txt be exported.
func NewFileCookie(fileName string) (*FileCookie, error) {
	ret := &FileCookie{}
	ret.cookies = make(map[string][]*http.Cookie)
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		words := strings.Split(line, "\t")
		// cookie.txt has 7 fields.
		if len(words) < 7 {
			continue
		}
		u, err := url.Parse("http://" + words[0])
		if err != nil {
			return nil, errors.Trace(err)
		}
		ts, err := strconv.ParseFloat(words[4], 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		httpOnly, err := strconv.ParseBool(words[1])
		if err != nil {
			return nil, errors.Trace(err)
		}
		secure, err := strconv.ParseBool(words[3])
		if err != nil {
			return nil, errors.Trace(err)
		}

		cookie := &http.Cookie{
			Name:     words[5],
			Value:    words[6],
			Path:     words[2],
			Domain:   words[0],
			Expires:  time.Unix(int64(ts), int64(ts-float64(int(ts)))*int64(time.Second/time.Nanosecond)),
			HttpOnly: httpOnly,
			Secure:   secure,
		}
		ret.SetCookies(u, []*http.Cookie{cookie})
	}
	return ret, nil
}

// SetCookies set cookies. If u exists then add this cookie.
func (c *FileCookie) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.cookies[u.Host]; !ok {
		c.cookies[u.Host] = make([]*http.Cookie, 0)
	}
	c.cookies[u.Host] = append(c.cookies[u.Host], cookies...)
}

// Cookies return cookies.
func (c *FileCookie) Cookies(u *url.URL) []*http.Cookie {
	var ret []*http.Cookie

	for domain := u.Host; len(domain) > 0; {
		if cookies, ok := c.cookies[domain]; ok {
			ret = append(ret, cookies...)
		}
		domain = strings.TrimLeft(domain, ".")
		if pos := strings.IndexRune(domain, '.'); pos != -1 {
			domain = domain[pos:]
		} else {
			domain = ""
		}
	}
	return ret
}
