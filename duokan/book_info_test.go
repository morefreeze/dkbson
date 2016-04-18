package duokan

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error [%s]", err)
	}
}

func TestGetURL(t *testing.T) {
	p := newDefaultProxy(nil)
	expectData := []byte("AiBzdGF0dXMAAm9rIHVybABQaHR0cDovL3BhZ2VzLnJlYWQuZHVva2FuLmNvbS9tZnN2Mi9kb3dubG9hZC9zMDEwL3AwMUhCUkN5RVdkTy9CaXp3dDJKOUZrSWlWSmguanM=")
	data, err := p.getURL("http://www.duokan.com/reader/page/5e5c9d902fbd49c8ae1860a161a83242/iss_0060008_-jUCuT3F0RNRz8f5rM8fuiHqApo")
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
		t.Errorf("Unexpected book info expected[%v] got[%v]", expectBInfo, bInfo)
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
	re *regexp.Regexp
}

func newLocalProxy() *localProxy {
	patt := "http://www.duokan.com/reader/book_info/([^/]+)/medium"
	return &localProxy{
		re: regexp.MustCompile(patt),
	}
}

func (p *localProxy) getURL(url string) ([]byte, error) {
	out := p.re.FindAllStringSubmatch(url, 1)
	if len(out) < 1 {
		return []byte(""), nil
	}
	switch out[0][1] {
	case "5e5c9d902fbd49c8ae1860a161a83242":
		return []byte("CyB0aXRsZQAk6Lef552A576O5Ymn44CK6ICB5Y+L6K6w44CL5a2m6Iux6K+tMHRyaWFsX21vZGUAASB0cmFpdAAGbWVkaXVtAm51bWJlcl9vZl9wYWdlcwC6AUFwYWdlX3NpemUAAgKbAgJ5AzBvd25fYm9vawAAQWNoYXB0ZXJzAC9AA0FwYWdlX3JhbmdlAAIBAQEBAW51bWJlcgABIGNpZAAkNWEwNGQxNGYtZjRjNC00NDljLTljOTQtZWE0NmRhYTg4N2I0QANBcGFnZV9yYW5nZQACAQIBAgFudW1iZXIAAiBjaWQAJGM2MGVlMGFiLTA0ZDYtNDJlMS05NzM4LTNjMmQ2YTAzNGViMkADQXBhZ2VfcmFuZ2UAAgEDAQUBbnVtYmVyAAMgY2lkACQ3OTk5ZGYxYi02MDI3LTQ3MmQtYmIyNi1jZTdlNzdkZTc1M2JAA0FwYWdlX3JhbmdlAAIBBgEGAW51bWJlcgAEIGNpZAAHMDAwMmI3NkADQXBhZ2VfcmFuZ2UAAgEHAQcBbnVtYmVyAAUgY2lkACQxNThjODc4NC1iYzIxLTQyNDMtODM1Mi1jZTk2OTVkNmMxNGVAA0FwYWdlX3JhbmdlAAIBCAEUAW51bWJlcgAGIGNpZAAkZGE2Nzg4ZjgtYzEwOC00NjhjLTlhYzMtMjlmMzQ1ZjMzYmZhQANBcGFnZV9yYW5nZQACARUBIQFudW1iZXIAByBjaWQAJDk2NWEzZjgyLWYwNTQtNDZhMS04MWVkLTdiNjE5YmY3MDZhMkADQXBhZ2VfcmFuZ2UAAgEiAS4BbnVtYmVyAAggY2lkACQ4ZjAzNGMxZi05MzU3LTRkZDItOWZhZi1lZTk4ZTAyODA5NzFAA0FwYWdlX3JhbmdlAAIBLwE7AW51bWJlcgAJIGNpZAAkZjRiNmZjMTQtNzBkYy00YjcxLWI1NDQtMjBhMzhmNzAzYjhiQANBcGFnZV9yYW5nZQACATwBSAFudW1iZXIACiBjaWQAJDQ4NmE0NzUzLTg4YmQtNDlhYi05YTViLWIzNzc3M2Q1ZGM1Y0ADQXBhZ2VfcmFuZ2UAAgFJAUkBbnVtYmVyAAsgY2lkAAc0ZThjYTc5QANBcGFnZV9yYW5nZQACAUoBSgFudW1iZXIADCBjaWQAJGIyNzM2NDk1LTJmZTYtNDI1MS1iZGEzLWU5MTY5ZTk4MWU4NEADQXBhZ2VfcmFuZ2UAAgFLAVgBbnVtYmVyAA0gY2lkACQxNDVmZGQ2MC1mYWE1LTQwNjQtOTlmMS00NTVkY2Q3ZGRkMmZAA0FwYWdlX3JhbmdlAAIBWQFlAW51bWJlcgAOIGNpZAAkOTA1ODM3ZDAtNTBmNS00M2Y1LTgzODMtNjVhODkyOTk2NjRiQANBcGFnZV9yYW5nZQACAWYBcgFudW1iZXIADyBjaWQAJGVmMjE1MTYxLWFhYTEtNDYzNC1hNzFiLTFiNWFmN2M4MWY2Y0ADQXBhZ2VfcmFuZ2UAAgFzAoAAAW51bWJlcgAQIGNpZAAkZjI0ZGMwNWItOWUxNC00MjRmLWFhYjktYTM0MjE0MTBjMjBhQANBcGFnZV9yYW5nZQACAoEAAo0AAW51bWJlcgARIGNpZAAkZWE5YTY0ZGItNTMyYS00NTk4LTkwMDUtNDY3NWJkMDQwNGVlQANBcGFnZV9yYW5nZQACAo4AAo4AAW51bWJlcgASIGNpZAAHNjUxNDY0YUADQXBhZ2VfcmFuZ2UAAgKPAAKPAAFudW1iZXIAEyBjaWQAJGNmZDI3YTk0LWViNTAtNDZiNy04MmMyLThiMDU5Y2M1NDQzNUADQXBhZ2VfcmFuZ2UAAgKQAAKbAAFudW1iZXIAFCBjaWQAJDI5NmYxMGJlLTg5NGYtNDg2MC1hMDNmLTJhZGM2ZTE2OGU4ZkADQXBhZ2VfcmFuZ2UAAgKcAAKoAAFudW1iZXIAFSBjaWQAJDlmODNkZTg3LTBhOWEtNGMyNS04NjAyLWViYzMxMjJhMDY2YUADQXBhZ2VfcmFuZ2UAAgKpAAK1AAFudW1iZXIAFiBjaWQAJDE0OWM5N2UxLTRlOGYtNDhmYi04YzAxLWIwN2U2NjUyYTU5YUADQXBhZ2VfcmFuZ2UAAgK2AALCAAFudW1iZXIAFyBjaWQAJDFlMGJmMGVmLTQ1NjMtNDE2ZC05Nzk0LTM5NzVmNTIwZGYyNUADQXBhZ2VfcmFuZ2UAAgLDAALPAAFudW1iZXIAGCBjaWQAJDk1NjViNzgyLTc4NTEtNGRkMy1hZDA1LTYwYmNiNjJmMzY0OUADQXBhZ2VfcmFuZ2UAAgLQAALcAAFudW1iZXIAGSBjaWQAJDcwYjUwMjM4LWI0MDUtNDMyOS1hYWVhLTMxNzIzNGEwZmM2ZUADQXBhZ2VfcmFuZ2UAAgLdAALdAAFudW1iZXIAGiBjaWQABzJlZjNlNmNAA0FwYWdlX3JhbmdlAAIC3gAC3gABbnVtYmVyABsgY2lkACRmODZhNmQwYS0wMzI3LTQwMmUtYTc3NS0zMjIyZTVjY2Y1MGJAA0FwYWdlX3JhbmdlAAIC3wAC7AABbnVtYmVyABwgY2lkACRkMTFjZGU5Ny1kOGE5LTRkM2UtOTM3My00Y2IwNjA5ZDRhNGFAA0FwYWdlX3JhbmdlAAIC7QAC/AABbnVtYmVyAB0gY2lkACQ4MjQzMGM3NC01OWFlLTRkOTQtYjQ1Mi1kMGM0NWVmZDdhNmRAA0FwYWdlX3JhbmdlAAIC/QACCwEBbnVtYmVyAB4gY2lkACQwYmYyZWUxZi1kZGVjLTQ4ZDItYjg2My1iZDc3YzQxM2JmNjZAA0FwYWdlX3JhbmdlAAICDAECGAEBbnVtYmVyAB8gY2lkACQ2MTZiYTVhNi01N2U3LTQwMDktYjM0Zi0yNWIwZmQ0MTQ1MDNAA0FwYWdlX3JhbmdlAAICGQECGQEBbnVtYmVyACAgY2lkAAcxMWE2ZjZjQANBcGFnZV9yYW5nZQACAhoBAhoBAW51bWJlcgAhIGNpZAAkZDhlY2E5Y2ItNjAwNi00MmI5LTg0ZDItYjQwZDMzMDM1NTM1QANBcGFnZV9yYW5nZQACAhsBAikBAW51bWJlcgAiIGNpZAAkYTM5YmQ0MjItYjVkZi00NmRlLTk3Y2UtNTkzZTA2ZTBiYjZhQANBcGFnZV9yYW5nZQACAioBAjkBAW51bWJlcgAjIGNpZAAkZDNhOGE5MDMtNDA0MS00ZjAxLWExOTYtYTZhNGMzNzlhNWYwQANBcGFnZV9yYW5nZQACAjoBAkgBAW51bWJlcgAkIGNpZAAkM2Q1ZjU5NDAtMGRmMC00NmY2LTllYzItY2RlMjJkZDcyYzk5QANBcGFnZV9yYW5nZQACAkkBAlUBAW51bWJlcgAlIGNpZAAkNzc3NjExODMtNTY3NS00NjMwLTllMjgtZWNjYTI5Njg2NzJmQANBcGFnZV9yYW5nZQACAlYBAmQBAW51bWJlcgAmIGNpZAAkODlkODFkNTktZjVlNy00ZWZmLWEyNDMtZjA1YjExY2RiOWRjQANBcGFnZV9yYW5nZQACAmUBAmUBAW51bWJlcgAnIGNpZAAHNWExZjJlM0ADQXBhZ2VfcmFuZ2UAAgJmAQJmAQFudW1iZXIAKCBjaWQAJDE1NDMwODNlLTY3ZTctNDViMi1hYmU5LWMwZjcxYWJmOGIwOUADQXBhZ2VfcmFuZ2UAAgJnAQJ1AQFudW1iZXIAKSBjaWQAJGRmNjlmMDAwLThlZjEtNGQ0My1iMjgxLWY4ZjkwZDRmMmFlOUADQXBhZ2VfcmFuZ2UAAgJ2AQKCAQFudW1iZXIAKiBjaWQAJDBkOWQyNGExLWEwYzEtNDY1Ni05YWJiLTQ0MDQ5YzVjZTg0ZEADQXBhZ2VfcmFuZ2UAAgKDAQKRAQFudW1iZXIAKyBjaWQAJGI5NDUxMDg2LWNhYjAtNDM1NC1hMWViLTc4MDE3Y2VjOWNmZEADQXBhZ2VfcmFuZ2UAAgKSAQKfAQFudW1iZXIALCBjaWQAJGI3YTcyYTI0LTc2NGItNDlhNy05OTFjLTIwN2M4ZmMyNTc4OUADQXBhZ2VfcmFuZ2UAAgKgAQKsAQFudW1iZXIALSBjaWQAJGI2MTNlYzgxLTNiZjEtNGNlYy05ODgwLTBjZjEwYmI2MmM1NEADQXBhZ2VfcmFuZ2UAAgKtAQK5AQFudW1iZXIALiBjaWQAJGRmZTdiMDcxLWIwMzYtNDY0Ni1iM2IwLWMxY2MyZGFhYjIwZEADQXBhZ2VfcmFuZ2UAAgK6AQK6AQFudW1iZXIALyBjaWQAJDdkNDk0MWVjLTY5OGYtNGYwOC04ZmVkLTJmMzlhNTdmNmQ1NCBib29rX2lkACA1ZTVjOWQ5MDJmYmQ0OWM4YWUxODYwYTE2MWE4MzI0MkB0b2MAA0FwYWdlX3JhbmdlAAIBAQEBQWNoaWxkcmVuAAlAA0FwYWdlX3JhbmdlAAIBAgECMWNoaWxkcmVuAAAgbmFtZQAM54mI5p2D5L+h5oGvQANBcGFnZV9yYW5nZQACAQMBBTFjaGlsZHJlbgAAIG5hbWUABuWJjeiogEADQXBhZ2VfcmFuZ2UAAgEGAQZBY2hpbGRyZW4ABUADQXBhZ2VfcmFuZ2UAAgEIARQxY2hpbGRyZW4AACBuYW1lACowMeOAgFlvdSdyZSBvdmVyIG1lPyDkvaDkuI3lnKjkuY7miJHkuobvvJ9AA0FwYWdlX3Jhbmdl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moTlkI3lrZflj6vnvZfmi4nvvIzlpbnmmK/kuKroiJ7lpbPigKbigKbigJ1AA0FwYWdlX3JhbmdlAAICoAECrAExY2hpbGRyZW4AACBuYW1lAFQzMOOAgEknZCBiZSBzYWQgc3VyZSwgYnV0IEkgd291bGRuJ3QgY3J5LiDmiJHlvZPnhLbkvJrkvKTlv4PnmoTvvIzkvYbmiJHkuI3kvJrlk63jgIJAA0FwYWdlX3JhbmdlAAICrQECuQExY2hpbGRyZW4AACBuYW1lAE4zMeOAgEkgaGF2ZSBmb3VuZCBteSBpZGVudGljYWwgaGFuZCB0d2luISDmiJHmib7liLDkuobmiJHnmoTlj4zog57og47miYvkuobvvIEgbmFtZQAlQ2hhcHRlciA2IOWQhOWFt+eJueiJsueahOe7j+WFuOiogOiuukADQXBhZ2VfcmFuZ2UAAgK6AQK6ATFjaGlsZHJlbgAAIG5hbWUAC+mZhOW9lUNE6aG1IG5hbWUAAEFwYWdlcwA8QAMBcGFnZV9udW1iZXIAASBwYWdlX2lkACdpc3NfMDA2MDAwOF8talVDdVQzRjBSTlJ6OGY1ck04ZnVpSHFBcG9BcG9zaXRpb24ABAEAAQABAAJSAUADAXBhZ2VfbnVtYmVyAAIgcGFnZV9pZAAnaXNzXzAwNjAwMDhfRXc2MTFlVk5ubThqSFh1eFNxR1pfVnN6dk1zQXBvc2l0aW9uAAQBAQEAAQACdAFAAwFwYWdlX251bWJlcgADIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X0VlYjRTcHh1ZzdNdFNsNi02ZDhnS3RGWGRUOEFwb3NpdGlvbgAEAQIBAAEAAnACQAMBcGFnZV9udW1iZXIABCBwYWdlX2lkACdpc3NfMDA2MDAwOF9sd0JZcVRxb080Y3BSRzNmSmpPWjBhTERDaWdBcG9zaXRpb24ABAECAQUBNgKVCEADAXBhZ2VfbnVtYmVyAAUgcGFnZV9pZAAnaXNzXzAwNjAwMDhfUkdsdkNuUXRtelZxVjRRb1VwN1lISTRoUHVBQXBvc2l0aW9uAAQBAgENAQACdg9AAwFwYWdlX251bWJlcgAGIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X01yMEpRNVlkcmJzX1YxbjVRbTFOZG85ZlpjZ0Fwb3NpdGlvbgAEAQMBAAEAAk0BQAMBcGFnZV9udW1iZXIAByBwYWdlX2lkACdpc3NfMDA2MDAwOF9mY1NjQUNoVGNycHhtTkR6ZERLVGFIMjI3WmtBcG9zaXRpb24ABAEEAQABAAKIAkADAXBhZ2VfbnVtYmVyAAggcGFnZV9pZAAnaXNzXzAwNjAwMDhfVGhFMXVfMjlZbHZWc0dhSmFaamd0SDlkeFVVQXBvc2l0aW9uAAQBBQEAAQACkAJAAwFwYWdlX251bWJlcgAJIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X3ZNMExJdmpyUjg2ZnlzQVU5ODZLN0VaVmlwb0Fwb3NpdGlvbgAEAQUBDAEAAq8LQAMBcGFnZV9udW1iZXIACiBwYWdlX2lkACdpc3NfMDA2MDAwOF8tcHVkaFBMN3J0WEtLcXN1Vnl1TFVIZ0xkRm9BcG9zaXRpb24ABAEFARkBAAKPFEADAXBhZ2VfbnVtYmVyAAsgcGFnZV9pZAAnaXNzXzAwNjAwMDhfNzNyaHUwWENMSHlfNVY0aThVOGlVY25RMjFZQXBvc2l0aW9uAAQBBQEmAQACjRpAAwFwYWdlX251bWJlcgAMIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X1VpUV96T0dVVC1vTTh3Nk5mYzVNYjhRYnlmQUFwb3NpdGlvbgAEAQUBLwE6Am0kQAMBcGFnZV9udW1iZXIADSBwYWdlX2lkACdpc3NfMDA2MDAwOF9CTWFiUWFhUXRaVldiSHpPVDFIYUUybHcxNVVBcG9zaXRpb24ABAEFAToBGgJbLEADAXBhZ2VfbnVtYmVyAA4gcGFnZV9pZAAnaXNzXzAwNjAwMDhfTjFxZTFiMm04cmtrWXVmNTJvcXNZTzk2M0ZRQXBvc2l0aW9uAAQBBQFEAWkCLTRAAwFwYWdlX251bWJlcgAPIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X3lLc1hoekExMmNxTmJxdTRwX2Ewa1RmWUJLY0Fwb3NpdGlvbgAEAQUBTQEAAho9QAMBcGFnZV9udW1iZXIAECBwYWdlX2lkACdpc3NfMDA2MDAwOF9DeTJDQXQ1emtrZU94Zk1lakJqRGRZczdyY2tBcG9zaXRpb24ABAEFAVcBAAIuRkADAXBhZ2VfbnVtYmVyABEgcGFnZV9pZAAnaXNzXzAwNjAwMDhfd2FSUUtVT2FSTkx6ZXo2N0hVR3BNOEEtWDhrQXBvc2l0aW9uAAQBBQFlAQAC7kxAAwFwYWdlX251bWJlcgASIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X2hlTGppVFhBdGNTcldGcnl5MXZBZ2ZybnBSc0Fwb3NpdGlvbgAEAQUBbQFSAm9WQAMBcGFnZV9udW1iZXIAEyBwYWdlX2lkACdpc3NfMDA2MDAwOF9Gem5aYWk1VWVFUFJoYlFyb1RZeFFIQkZuLUlBcG9zaXRpb24ABAEFAXMBAALHX0ADAXBhZ2VfbnVtYmVyABQgcGFnZV9pZAAnaXNzXzAwNjAwMDhfS0lYSnhZOWNJeV9rU3poeFlURUlfakEzbnEwQXBvc2l0aW9uAAQBBQF1AQACTmZAAwFwYWdlX251bWJlcgAVIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X215YUNoMHVzRXp0Q2c2Rk1WcGdKaEdGYW44RUFwb3NpdGlvbgAEAQYBAAEAApACQAMBcGFnZV9udW1iZXIAFiBwYWdlX2lkACdpc3NfMDA2MDAwOF92VXVoMmVFdjNGTHp3bE9ZckJGWEw5bi1sbDhBcG9zaXRpb24ABAEGAQ0BAAISDUADAXBhZ2VfbnVtYmVyABcgcGFnZV9pZAAnaXNzXzAwNjAwMDhfUkVmWjY5VnloU1ZRSWphMU81bDN1d0F2RHRFQXBvc2l0aW9uAAQBBgEaARoCthRAAwFwYWdlX251bWJlcgAYIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X29KZVBWQVlrQ3U4V0kwS2tQREhSTGRVb1JqTUFwb3NpdGlvbgAEAQYBJwEAAhMcQAMBcGFnZV9udW1iZXIAGSBwYWdlX2lkACdpc3NfMDA2MDAwOF8xUDY2T0I5TEhhUndtOXUzNnNYOWt0V3dyalFBcG9zaXRpb24ABAEGAS8BAAK/JUADAXBhZ2VfbnVtYmVyABogcGFnZV9pZAAnaXNzXzAwNjAwMDhfd3lUQnZfTEJSYmVqdWQxS0NkTEpUUWcxNjNnQXBvc2l0aW9uAAQBBgE6AQACNS5AAwFwYWdlX251bWJlcgAbIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X3d0UTBvaV9qTVJrVU1nSzI5c3dBQXluejA1d0Fwb3NpdGlvbgAEAQYBRwEAAp40QAMBcGFnZV9udW1iZXIAHCBwYWdlX2lkACdpc3NfMDA2MDAwOF96ZHFRTTVzMnZLY3BfYW1tZWNYT3VVbVBrRkVBcG9zaXRpb24ABAEGAVABAAINPUADAXBhZ2VfbnVtYmVyAB0gcGFnZV9pZAAnaXNzXzAwNjAwMDhfWUNBeTBiYW9DMnN0ZFpEMGQza0ZaVk5WdlNRQXBvc2l0aW9uAAQBBgFdAQAC2UlAAwFwYWdlX251bWJlcgAeIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X0lEUUxHSHBDYl9wMy0yRkNMUzVTdy13SEY0UUFwb3NpdGlvbgAEAQYBagEAApRTQAMBcGFnZV9udW1iZXIAHyBwYWdlX2lkACdpc3NfMDA2MDAwOF92b2NEOFh4UVJ6Q1ZRbDR0cVpGNjFFNGFaR1FBcG9zaXRpb24ABAEGAXQBAAK9W0ADAXBhZ2VfbnVtYmVyACAgcGFnZV9pZAAnaXNzXzAwNjAwMDhfYTFHT3VKM3hnN3RVYllvTS0wbVJBamp1SlprQXBvc2l0aW9uAAQBBgF6ARsCbWRAAwFwYWdlX251bWJlcgAhIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X0lBUExyS1BTT3lmZGpTWVRqRldTd0d2cUlxOEFwb3NpdGlvbgAEAQYBfAFAAo9xQAMBcGFnZV9udW1iZXIAIiBwYWdlX2lkACdpc3NfMDA2MDAwOF8wYnMwSzRNTlN4TGZYNjJzYjdsM3NVenhEaThBcG9zaXRpb24ABAEHAQABAAKQAkADAXBhZ2VfbnVtYmVyACMgcGFnZV9pZAAnaXNzXzAwNjAwMDhfUlRWcE9CSkpTZS1yUnZDVDNoVURER004Wl84QXBvc2l0aW9uAAQBBwEJAQACJAxAAwFwYWdlX251bWJlcgAkIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X3ZTTVB2aXZrWXdrcmVndzNxN0tuVHdCcDFZQUFwb3NpdGlvbgAEAQcBFAEAApEUQAMBcGFnZV9udW1iZXIAJSBwYWdlX2lkACdpc3NfMDA2MDAwOF9uRzJhd1JNQks3UFFzUDFVNXBTYml1Wl90WUFBcG9zaXRpb24ABAEHAR8BXAKxG0ADAXBhZ2VfbnVtYmVyACYgcGFnZV9pZAAnaXNzXzAwNjAwMDhfaVR5UjR3ejFEVDhUY0dCYXg4aW1rVGVPZHdFQXBvc2l0aW9uAAQBBwEoAX0CnSRAAwFwYWdlX251bWJlcgAnIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X1NVX1ZfWkJpX3FEUUZWMFFYa1hENEktc1lZNEFwb3NpdGlvbgAEAQcBMwEbAvEtQAMBcGFnZV9udW1iZXIAKCBwYWdlX2lkACdpc3NfMDA2MDAwOF9WV1ljTUszejVRTzhXVHFITldoaFRmRnhfSFVBcG9zaXRpb24ABAEHAT4BOwIkNUADAXBhZ2VfbnVtYmVyACkgcGFnZV9pZAAnaXNzXzAwNjAwMDhfWHhxeWttdFNqang3YlhJZlJPdkhFeWxwU09RQXBvc2l0aW9uAAQBBwFFAukAAts9QAMBcGFnZV9udW1iZXIAKiBwYWdlX2lkACdpc3NfMDA2MDAwOF83bGNUMGliRWM4Z090U0w4c2ZNRDBndHAyMjhBcG9zaXRpb24ABAEHAVEBOAJ3R0ADAXBhZ2VfbnVtYmVyACsgcGFnZV9pZAAnaXNzXzAwNjAwMDhfTHY2WS1HUE53N0s1b3NNdWZxcnN1elhJeU5VQXBvc2l0aW9uAAQBBwFdAQACHU5AAwFwYWdlX251bWJlcgAsIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X29naHBtQU1qVEhKa3kydnBxZ1ppR29hNUFFc0Fwb3NpdGlvbgAEAQcBaAEAAoFWQAMBcGFnZV9udW1iZXIALSBwYWdlX2lkACdpc3NfMDA2MDAwOF92MXVGejBGNVJBYWtPY2dpN2p1c1M2b0FObFVBcG9zaXRpb24ABAEHAWwBAAJ9YEADAXBhZ2VfbnVtYmVyAC4gcGFnZV9pZAAnaXNzXzAwNjAwMDhfaWlGOWN3Yno2R1ZPMmZ2aXpLZU5OVmZJbDJzQXBvc2l0aW9uAAQBBwFyAQACgWhAAwFwYWdlX251bWJlcgAvIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X1N4WEl4RVJRT1hyQVhlaFNFNE5SMlVud2p5SUFwb3NpdGlvbgAEAQgBAAEAApACQAMBcGFnZV9udW1iZXIAMCBwYWdlX2lkACdpc3NfMDA2MDAwOF95MV9JY3lIdjJKVVUzaXZDczNwZ2xWTzdXZUVBcG9zaXRpb24ABAEIAQoBAAJ9CkADAXBhZ2VfbnVtYmVyADEgcGFnZV9pZAAnaXNzXzAwNjAwMDhfbnQ0c0RNZ1g2dXc0UDJjZUQweTdnSVBZTUxrQXBvc2l0aW9uAAQBCAEXAQACrxNAAwFwYWdlX251bWJlcgAyIHBhZ2VfaWQAJ2lzc18wMDYwMDA4X09TaW5YanExTFJPdEp5bElCbWgzcUtxSXBza0Fwb3NpdGlvbgAEAQgBJAEAArgaQAMBcGFnZV9udW1iZXIAMyBwYWdlX2lkACdpc3NfMDA2MDAwOF9BV2Nha0h4LVV1SmZDMFJENlU4QVdISUFESHdBcG9zaXRpb24ABAEIAS0BAAIuJEADAXBhZ2VfbnVtYmVyADQgcGFnZV9pZAAnaXNzXzAwNjAwMDhfTC1sQmFIYVBhYXdFODFxVUVlakktTGZ2WDhJQXBvc2l0aW9uAAQBCAE4AQACAy9AAwFwYWdlX251bWJlcgA1IHBhZ2VfaWQAJ2lzc18wMDYwMDA4X2J5ZDNHcnpEUzlMSDhJeTFZbkFkOHAzMDZNRUFwb3NpdGlvbgAEAQgBQwEAAks1QAMBcGFnZV9udW1iZXIANiBwYWdlX2lkACdpc3NfMDA2MDAwOF9tWG9NLTd1OWlPV3FQTFNwR29aRGNfeG56TWdBcG9zaXRpb24ABAEIAUsBSgJ3PUADAXBhZ2VfbnVtYmVyADcgcGFnZV9pZAAnaXNzXzAwNjAwMDhfaFlfZjdQcDBQejlYNm4zd1VyT0RtQTVyUEVvQXBvc2l0aW9uAAQBCAFUAXkCfkdAAwFwYWdlX251bWJlcgA4IHBhZ2VfaWQAJ2lzc18wMDYwMDA4X2p4cHJhZWxOUFktbkJlM3FlaGN1eXVTVDBONEFwb3NpdGlvbgAEAQgBXwEAAgZOQAMBcGFnZV9udW1iZXIAOSBwYWdlX2lkACdpc3NfMDA2MDAwOF9PYmFrdWx4VlctM2xKMHB1aG5Da0hYaUZkRk1BcG9zaXRpb24ABAEIAWgBWgLUVUADAXBhZ2VfbnVtYmVyADogcGFnZV9pZAAnaXNzXzAwNjAwMDhfWmdLY0VqY3VJaXdvOTJKMlJKTGp0TkdVS0tnQXBvc2l0aW9uAAQBCAFuATcCy15AAwFwYWdlX251bWJlcgA7IHBhZ2VfaWQAJ2lzc18wMDYwMDA4X0JydVk2ZXlvNDgycDRGdDhEYk5NZlJBc0NyUUFwb3NpdGlvbgAEAQgBcQEAAgpiQAMBcGFnZV9udW1iZXIAPCBwYWdlX2lkACdpc3NfMDA2MDAwOF9GVzF0WWlXaHVkcUQtcU9iczVIOTd2bnNZN29BcG9zaXRpb24ABAEJAQABAAKQAiByZXZpc2lvbgAKMjAxNTExMTAuMQ=="), nil
	}
	return []byte(""), nil
}
