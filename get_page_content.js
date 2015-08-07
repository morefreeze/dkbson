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
    res_arr.push({'options': line, 'line':line,});
    (function do_req(k){
        var get_content = function(str){
            // duokan_page('bson_content'), remove single quote outside
            str = str.substr(str.indexOf("'")+1);
            str = str.substr(0, str.lastIndexOf("'"));
            res = dk.dkbson.decode(str);
            if (res.status == 'error'){
                console.error(res);
                process.exit(2);
            }
            // try to restore text, if current word y != last_y then add newline
            s='    ';
            last_y=0;
            items = res.items;
            for(var j in items){
                 if(items[j].type=='word'){
                    if(last_y!=items[j].y){
                        //s+="\n";
                        if(items[j].x >= 100) s+="\n    ";
                    }
                    s+=items[j].char;
                    last_y=items[j].y;
                 }
            }
            res_arr[k].text = s;
            finished += 1;
            if (finished == res_arr.length){
                for (var kk in res_arr){
                    if (undefined === res_arr[kk].text){
                        console.error('Page number content '+kk+' is missing ' + res_arr[kk].line);
                    }
                    else{
                        console.log(res_arr[kk].text);
                    }
                }
            }
        };
        dk.req(line, get_content);
    })(k);
    k += 1;
});

// todo: setTimeout if some request timeout then output those request have done
