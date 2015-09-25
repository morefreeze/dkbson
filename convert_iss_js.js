var readline = require('readline');
var process = require('process');
var sleep = require('sleep');
var dk = require("./dkbson");
// read book_id/iss_xxxx
var rl = readline.createInterface({
    input: process.stdin,
    terminal: false,
});
// store final order res
var res_arr = [];
var k = 0;
var finished = 0;
var all_promises = [];
rl.on('line', function(line){
    if ('' === line) return;
    var url = 'http://www.duokan.com' + '/reader/page/'+line;
    res_arr.push({'url':url,'line':line,});
    all_promises.push(null);
    (function do_req(k){
        console.error(k);
        all_promises[k] = dk.req(res_arr[k].url, function(str){
            var res = dk.dkbson.decode(str);
            if (res.status == 'error'){
                console.log(res);
                process.exit(2);
            }
            res_arr[k].js_url = res.url;
            console.error(k + " " + res_arr[k].url + " " + res.url);
            /*
            finished += 1;
            if (finished == res_arr.length){
                for (var kk in res_arr){
                    if (undefined === res_arr[kk].js_url){
                        console.error('Page number js_url '+kk+' is missing ' + res_arr[kk].line);
                    }
                    else{
                        console.log(res_arr[kk].js_url);
                    }
                }
            }
            */
            });
        sleep.usleep(20000);
    })(k);
    k += 1;
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
