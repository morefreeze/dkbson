if [[ $1 == "" ]];then
    echo "Usage: sh $0 book_id"
    exit 1
fi
if [ ! -f ${1}.iss ];then
    node get_book_info.js $1 iss > ${1}.iss
fi
if [ ! -f ${1}.jsurl ];then
    cat ${1}.iss | node convert_iss_js.js > ${1}.jsurl
fi
cat ${1}.jsurl| node get_page_content.js > ${1}.txt
