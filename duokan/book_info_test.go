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

func TestGetURL(t *testing.T) {
	p := NewDefaultProxy(nil)
	expectData := []byte("AiBzdGF0dXMAAm9rIHVybABQaHR0cDovL3BhZ2VzLnJlYWQuZHVva2FuLmNvbS9tZnN2Mi9kb3dubG9hZC9zMDEwL3AwMUhCUkN5RVdkTy9CaXp3dDJKOUZrSWlWSmguanM=")
	fName, err := p.GetURL("http://www.duokan.com/reader/page/5e5c9d902fbd49c8ae1860a161a83242/iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo")
	checkErr(t, err)
	f, err := os.Open(fName)
	checkErr(t, err)
	data, err := ioutil.ReadAll(f)
	checkErr(t, err)
	if bytes.Compare(data, expectData) != 0 {
		t.Errorf("Unexpected data expected[%q] got[%q]", expectData, data)
	}
}

func TestGetBookInfo(t *testing.T) {
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	l := &Librarian{
		proxy: newLocalProxy(),
	}
	bInfo, err := l.GetBookInfo(bid)
	checkErr(t, err)
	expectBInfo := BookInfo{
		Title:    "跟着美剧《老友记》学英语",
		Pages:    make([]*PageInfo, 60),
		Revision: "20151110.1",
	}
	if !compareBook(expectBInfo, bInfo) {
		t.Fatalf("Unexpected book info expected[%v] got[%v]", expectBInfo, bInfo)
	}
}

func TestIss(t *testing.T) {
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	iss := "iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo"
	l := &Librarian{
		proxy: newLocalProxy(),
	}
	jsAddr, err := l.iss2Js(bid, iss)
	checkErr(t, err)
	expectJs := "http://pages.read.duokan.com/mfsv2/download/s010/p01HBRCyEWdO/Bizwt2J9FkIiVJh.js"
	if expectJs != jsAddr {
		t.Fatalf("Unexpected js address expected[%s] got[%s]", expectJs, jsAddr)
	}
	// TODO: add status!=ok test
}

func TestGetPageContent(t *testing.T) {
	js := "http://pages.read.duokan.com/mfsv2/download/s010/p01HBRCyEWdO/Bizwt2J9FkIiVJh.js"
	l := &Librarian{
		proxy: newLocalProxy(),
	}
	content, err := l.getPageContent(js)
	checkErr(t, err)
	_ = content
}

func TestDecodeURL(t *testing.T) {
	bid := "5e5c9d902fbd49c8ae1860a161a83242"
	l := &Librarian{
		proxy: newLocalProxy(),
	}
	url := fmt.Sprintf("http://www.duokan.com/reader/book_info/%s/medium", bid)
	data, err := l.DecodeURL(url)
	checkErr(t, err)
	f, err := os.Open("../tests/test_book_info.ret")
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

func (p *localProxy) GetURL(url string) (string, error) {
	var fileName string
	switch url {
	case "http://www.duokan.com/reader/book_info/5e5c9d902fbd49c8ae1860a161a83242/medium":
		fileName = "../tests/test_book_info"
	case "http://www.duokan.com/reader/page/5e5c9d902fbd49c8ae1860a161a83242/iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo":
		fileName = "../tests/test_iss"
	case "http://pages.read.duokan.com/mfsv2/download/s010/p01HBRCyEWdO/Bizwt2J9FkIiVJh.js":
		fileName = "../tests/test_js"
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
