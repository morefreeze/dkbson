if [[ $1 == "" ]];then
    echo "Usage: sh $0 book_id"
    exit 1
fi
node get_book_info.js $1 iss | head | node convert_iss_js.js | node get_page_content.js
