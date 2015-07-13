var process = require('process');
var http = require("http");
var dk = require("./dkbson");
if (process.argv.length <= 2){
    console.log("Usage: node "+process.argv[1]+" book_md5 [title/iss]");
    process.exit(1);
}
var options = {
    host: 'www.duokan.com',
    path: '/reader/book_info/'+process.argv[2]+'/medium',
};
var parse_iss = function(str){
    res = dk.dkbson.decode(str);
    if (res.status == 'error'){
        console.log(res);
        process.exit(2);
    }
    var field = process.argv[3];
    if ('title' == field){
        console.log(res.title);
    }
    else if('iss' == field){
        for (var i in res.pages){
            console.log(res.book_id+'/'+res.pages[i].page_id);
        }
    }
    else{
        console.log(res);
    }
};
dk.req(options, parse_iss);
