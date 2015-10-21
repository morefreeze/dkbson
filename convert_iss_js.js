var readline = require('linebyline');
var process = require('process');
var sleep = require('sleep');
var dk = require("./dkbson");
// read book_id/iss_xxxx
rl = readline(process.stdin);
// store final order res
var res_arr = [];
var all_promises = [];
function decode_retry(url){
    var MAX_DEPTH = 3;
    return new Promise(function(resolve, reject) {
        function retry(depth) {
            if (depth > MAX_DEPTH) return reject("decode max retry times");
            dk.req(url)
                .then(function(str) {
                    var res = dk.dkbson.decode(str);
                    if (res.status == 'error'){
                        console.error('================'+res);
                        process.exit(2);
                    }
                    //console.error(k + " " + res_arr[k].url + " " + res.url);
                    resolve(res.url);
                    //return res.url;
                })
                .catch(function(error) {
                    console.error('decode error');
                    console.error(error);
                    retry(depth+1);
                });
        }
        retry(1);
    });
}
rl.on('line', function(line, lineCount, byteCount){
    if ('' === line) return;
    var url = 'http://www.duokan.com' + '/reader/page/'+line;
    res_arr.push({'url':url,'line':line,});
    all_promises.push(null);
    var k = lineCount - 1;
    all_promises[k] = decode_retry(res_arr[k].url);
    sleep.usleep(2000000);
})
.on('close', function() {
    var p = Promise.all(all_promises);
    p.then(function(body_arr){
        for (var k in body_arr){
            console.log(body_arr[k]);
        }
        console.error("all done");
    }, function(reason){
        console.error("some one error " + reason);
    });
});
