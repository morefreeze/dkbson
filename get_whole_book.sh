# /usr/bin/env bash
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
title=$(node get_book_info.js $1 title)
title="${1:0:4}${title// /_}"
count=0;
lines=""
mv ${title}.txt ${title}.txt.bak 2>/dev/null
while read line; do
    count=$((count+1))
    lines="$lines$line\n"
    if (( count % 100 == 0 ));then
        echo -en "$lines" | node get_page_content.js >> ${title}.txt
        lines=""
    fi
done << EOF
$(cat $1.jsurl)
EOF
echo -en "$lines" | node get_page_content.js >> ${title}.txt
