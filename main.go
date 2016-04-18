package main

import (
	"fmt"

	"github.com/morefreeze/dkbson/duokan"
)

func main() {
	l := duokan.NewLibrarian(nil)
	b, _ := l.GetBookInfo("0c9919558af14bcb9c38ec22f3885b78")
	fmt.Printf("%v", b)
}
