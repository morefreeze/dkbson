var readline = require('readline');
var process = require('process');
var http = require("http");
var dk = require("./dkbson");
// read book_id/iss_xxxx
var rl = readline.createInterface({
    input: process.stdin,
    terminal: false,
});
var get_js = function(str){
    res = dk.dkbson.decode(str);
    if (res.status == 'error'){
        console.log(res);
        process.exit(2);
    }
    console.log(res.url);
};
rl.on('line', function(line){
    var options = {
        host: 'www.duokan.com',
        path: '/reader/page/'+line,
    };
    dk.req(options, get_js);
});
