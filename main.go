package main

import (
	"fmt"

	"github.com/morefreeze/dkbson/duokan"
)

func main() {
	jar, _ := duokan.NewFileCookie("./duokan/cookie.txt")
	proxy := duokan.NewDefaultProxy(jar)
	l := duokan.NewLibrarian(proxy)
	b, _ := l.GetBookInfo("4479703547c34aba930ef5e754c69381")
	fmt.Printf("%v", b)
}
