package duokan

import "testing"

func TestGetPageContentWithPic(t *testing.T) {
	js := "http://pages.read.duokan.com/mfsv2/download/s010/p01dmt6Irale/RQvGquaD1rmKtdZ.js"
	l := &Librarian{
		proxy: newLocalProxy(),
	}
	page, err := l.getPageContent(js)
	checkErr(t, err)
	if 270 != len(page.Items) {
		t.Fatalf("Page item expected[%d] got[%d]", 270, len(page.Items))
	}
	// TODO: check picture.
}
