# dkbson
Extract the dkbson in duokan.com for decode duokan data

### get_whole_book.sh
`sh get_whole_book.sh book_md5`

**DO RUN WITH `sh` INSTEAD OF `bash`**

It will save all middle result including iss list and js url list, finally it
will save ${title}.txt

If you meet this error:
```
events.js:85
      throw er; // Unhandled 'error' event
```
use `bash get_whole_book.sh [book_md5]` explicitly

Also, you can only know this one-key shell, following is separate js script

### get_book_info.js
`Usage: node get_book_info.js [title|iss]`

If omit the third argument, it will output the result json, but some object
only display `[Object]` instead of expanding it

### convert_iss_js.js
`Usage: cat foo.iss | node convert_iss_js.js` foo.iss is output of
`node get_book_info.js iss`

It will output url of js file each line, and it has the same order of the iss file

### get_page_content.js
`Usage: cat foo.jsurl | node get_page_content.js` foo.jsurl is output of
`convert_iss_js.js`

It will output the text of the js file decoded. It will check `y` of position
to make new line, but it is useless when the book including lots of formula
or mathematical symbols.

## TODO
I am trying to let the duokan web page load local iss file to display rich text
of the book.

