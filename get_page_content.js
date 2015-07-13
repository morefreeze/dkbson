var readline = require('readline');
var process = require('process');
var http = require("http");
var dk = require("./dkbson");
// read book_id/iss_xxxx
var rl = readline.createInterface({
    input: process.stdin,
    terminal: false,
});
var get_content = function(str){
    // duokan_page('bson_content'), remove single quote outside
    str = str.substr(str.indexOf("'")+1);
    str = str.substr(0, str.lastIndexOf("'"));
    //console.log(str);
    res = dk.dkbson.decode(str);
    if (res.status == 'error'){
        console.log(res);
        process.exit(2);
    }
    // try to restore text, if current word y != last_y then add newline
    s='    ';
    last_y=0;
    items = res.items;
    for(var j in items){
         if(items[j].type=='word'){
            if(last_y!=items[j].y){
                s+="\n";
                if(items[j].x >= 100) s+="    ";
            }
            s+=items[j].char;
            last_y=items[j].y;
         }
    }
    console.log(s);
};
rl.on('line', function(line){
    dk.req(line, get_content);
});
