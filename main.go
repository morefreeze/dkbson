package main

import (
	"github.com/juju/errors"
	"github.com/morefreeze/dkbson/duokan"
	"github.com/ngaut/log"
)

func main() {
	// If cookie.txt missing, it will get book without login info.
	jar, _ := duokan.NewFileCookie("./duokan/cookie.txt")
	proxy := duokan.NewBsonProxy(jar)
	l := duokan.NewLibrarian(proxy)
	bid := "4479703547c34aba930ef5e754c69381"
	b, err := l.GetBookInfo(bid)
	if err != nil {
		log.Errorf("%s", errors.ErrorStack(err))
		return
	}
	log.Debugf("%d", len(b.Pages))
	//if err := l.SaveBook(bid, fmt.Sprintf("%s.txt", bid)); err != nil {
	//log.Errorf("%s", errors.ErrorStack(err))
	//return
	//}
}
