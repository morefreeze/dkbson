var request = require('request');
var sleep = require('sleep');
var readline = require('linebyline');
rl = readline(process.stdin);
var res_arr = [];
function get_url(url){
    return new Promise(function(resolve, reject) {
        var l_options = {'url':url};
        l_options.timeout = 3 * 1000;
        //console.error(l_options);
        var st_time = Date.now();
        console.error(2);
        request(l_options, function(error, response, body) {
            console.error('cost '+(Date.now()-st_time)/1000);
            if (error) {
                console.error(' ' +l_options.url + ' ' + error);
                return ;
            }
            //console.warn("succ "+body.length);
            resolve(body);
        });
    });
}
rl.on('line', function(line, lineCount, byteCount){
    var k = lineCount - 1;
        get_url(line).then(function(str){
            console.log("body "+str.length);
        },
        function(str){
            console.error("failed "+str);
        });
        //sleep.usleep(50000);
});

