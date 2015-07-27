var readline = require('readline');
var process = require('process');
var http = require("http");
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
rl.on('line', function(line){
    if ('' === line) return;
    var options = {
        host: 'www.duokan.com',
        path: '/reader/page/'+line,
    };
    res_arr.push({'options':options,'line':line,});
    (function do_req(k){
        var get_js = function(str){
            var res = dk.dkbson.decode(str);
            if (res.status == 'error'){
                console.log(res);
                process.exit(2);
            }
            res_arr[k].js_url = res.url;
            finished += 1;
            if (finished == res_arr.length){
                for (var kk in res_arr){
                    // retry 3 times every 300ms, or alert an error
                    if (undefined === res_arr[kk].js_url){
                        console.log('Page number js_url '+kk+' is missing ' + res_arr[kk].line);
                    }
                    else{
                        console.log(res_arr[kk].js_url);
                    }
                }
            }
        };
        dk.req(res_arr[k].options, get_js);
    })(k);
    k += 1;
});

// todo: setTimeout if some request timeout then output those request have done
