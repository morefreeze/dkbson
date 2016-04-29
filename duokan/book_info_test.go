package duokan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/juju/errors"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error [%s]", errors.ErrorStack(err))
	}
}

func TestFetchURL(t *testing.T) {
	p := NewBsonProxy(nil)
	urlList := []string{
		"http://www.duokan.com/reader/book_info/4479703547c34aba930ef5e754c69381/medium",
		"http://www.duokan.com/reader/page/5e5c9d902fbd49c8ae1860a161a83242/iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo",
	}
	for i := 0; i < len(urlList); i++ {
		expectData, err := ioutil.ReadFile(fmt.Sprintf("../results/test_url%d", i))
		checkErr(t, err)
		fName, err := p.FetchURL(urlList[i])
		checkErr(t, err)
		f, err := os.Open(fName)
		checkErr(t, err)
		data, err := ioutil.ReadAll(f)
		checkErr(t, err)
		if bytes.Compare(data, expectData) != 0 {
			t.Errorf("Unexpected data url[%s] expected[%q] got[%q]", urlList[i], expectData, data)
		}

	}
}

func TestGetBookInfo(t *testing.T) {
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	l := NewLibrarian(newLocalProxy())
	bInfo, err := l.GetBookInfo(bid)
	checkErr(t, err)
	expectBInfo := BookInfo{
		Title:    "跟着美剧《老友记》学英语",
		Pages:    make([]PageInfo, 60),
		Revision: "20151110.1",
	}
	if !compareBook(expectBInfo, bInfo) {
		t.Fatalf("Unexpected book info expected[%v] got[%v]", expectBInfo, bInfo)
	}
}

func TestIss(t *testing.T) {
	l := NewLibrarian(newLocalProxy())
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	iss := "iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo"
	jsAddr, err := l.iss2Js(bid, iss)
	checkErr(t, err)
	expectJs := "http://pages.read.duokan.com/mfsv2/download/s010/p01HBRCyEWdO/Bizwt2J9FkIiVJh.js"
	if expectJs != jsAddr {
		t.Fatalf("Unexpected js address\nexpected[%s]\ngot[%s]", expectJs, jsAddr)
	}
	// TODO: add status!=ok test
}

func TestGetPageContent(t *testing.T) {
	l := NewLibrarian(newLocalProxy())
	js := "http://pages.read.duokan.com/mfsv2/download/s010/p01WgwI7imAG/eZryvzKeQBfcSbb.js"
	page, err := l.getPageContent(js)
	checkErr(t, err)
	if 555 != len(page.Items) {
		t.Fatalf("Page item expected[%d] got[%d]", 555, len(page.Items))
	}
	content, err := page.GenerateContent()
	checkErr(t, err)
	expectContent, _ := ioutil.ReadFile("../results/test_js")
	if string(expectContent) != content {
		t.Fatalf("Unexpected page content\nexp[%s]\ngot[%s]", expectContent, content)
	}
}

func TestDecodeURL(t *testing.T) {
	l := NewLibrarian(newLocalProxy())
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	url := fmt.Sprintf("http://www.duokan.com/reader/book_info/%s/medium", bid)
	data, err := l.DecodeURL(url)
	checkErr(t, err)
	f, err := os.Open("../results/test_book_info")
	checkErr(t, err)
	expectData, err := ioutil.ReadAll(f)
	checkErr(t, err)
	checkErr(t, f.Close())
	if bytes.Compare(expectData, data) != 0 {
		t.Fatalf("Unexpected decode data expected[%q] got[%q]", expectData, data)
	}
}

func compareBook(lhs, rhs BookInfo) bool {
	lMap := bookInfoToMap(lhs)
	rMap := bookInfoToMap(rhs)
	return reflect.DeepEqual(lMap, rMap)
}

// bookInfoToMap dumps BookInfo to map, but only record length of array/slice.
func bookInfoToMap(in BookInfo) map[string]string {
	ret := make(map[string]string)
	v := reflect.ValueOf(in)
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		var v string
		switch f.Type().Kind() {
		case reflect.Slice, reflect.Array:
			v = fmt.Sprintf("%d", f.Len())
		default:
			v = fmt.Sprintf("%s", f.String())
		}
		ret[typ.Field(i).Name] = v
	}
	return ret
}

type localProxy struct {
}

func newLocalProxy() *localProxy {
	return &localProxy{}
}

func (p *localProxy) FetchURL(url string) (string, error) {
	var fileName string
	switch url {
	case "http://www.duokan.com/reader/book_info/5e5c9d902fbd49c8ae1860a161a83242/medium":
		fileName = "../tests/test_book_info"
	case "http://www.duokan.com/reader/page/5e5c9d902fbd49c8ae1860a161a83242/iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo":
		fileName = "../tests/test_iss"
	case "http://pages.read.duokan.com/mfsv2/download/s010/p01WgwI7imAG/eZryvzKeQBfcSbb.js":
		fileName = "../tests/test_js"
	case "http://pages.read.duokan.com/mfsv2/download/s010/p01dmt6Irale/RQvGquaD1rmKtdZ.js":
		fileName = "../tests/test_pic_js"
	default:
		return "", errors.Errorf("no such url[%s]", url)
	}
	f, _ := os.Open(fileName)
	data, _ := ioutil.ReadAll(f)
	f.Close()
	outF, _ := ioutil.TempFile(os.TempDir(), "duokan_test")
	outF.Write(data)
	outF.Close()
	return outF.Name(), nil
}
