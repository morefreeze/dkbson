package duokan

import "testing"

func TestFileCookie(t *testing.T) {
	jar, err := NewFileCookie("./cookie.txt")
	checkErr(t, err)
	t.Logf("%v", jar.cookies)
}
