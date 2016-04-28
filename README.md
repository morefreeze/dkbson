# dkbson
Extract the dkbson in duokan.com for decoding duokan data.

DO NOT use for business!!

`npm install`

`go run main.go`

You can feel free to modify `main.go` to any book you want, even you can set
your own cookie(download it you can use Chrome extension like `cookie.txt export`).

## TODO
1. I am trying to generate more pretty page(alternative) and load image even
book internal link:
    1. let the duokan web page load local iss file to display rich text
of the book(perfect but hard).
    1. write file to html with simple format(acceptable).
1. Add queue when downloading.
1. Save meta data into durability store(like MongoDB, Redis).
1. Dependent on last, support breakpoint download.
