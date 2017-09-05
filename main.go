package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/morefreeze/dkbson/duokan"
	"github.com/ngaut/log"
)

func main() {
	// If cookie.txt missing, it will get book without login info.
	jar, _ := duokan.NewFileCookie("./duokan/cookie.txt")
	proxy := duokan.NewBsonProxy(jar)
	l := duokan.NewLibrarian(proxy)
	bid := "837222cb5b3f48428b57a29869d7bd30"
	b, err := l.GetBookInfo(bid)
	if err != nil {
		log.Errorf("%s", errors.ErrorStack(err))
		return
	}
	log.Debugf("%d", len(b.Pages))
	if err := l.SaveBook(bid, fmt.Sprintf("%s.txt", bid)); err != nil {
		log.Errorf("%s", errors.ErrorStack(err))
		return
	}
}
