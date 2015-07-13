s = "duoka('dsfdasfjkljfs')";
s = s.substr(s.indexOf("'"));
s = s.substr(0, s.lastIndexOf("'")+1);
console.log(s);
var url = require('url');
var http = require("http");
var dk = require("./dkbson");
var options = {
    host: 'www.duokan.com',
    path: '/reader/book_info/'+'639c6fa4e78f11e1bbc800163e0123ac'+'/medium',
};
http.request(options, function(response){
    var str = '';
    response.on('data', function(chunk){
        str += chunk;
    });
    response.on('end', function(){
        res = dk.dkbson.decode(str);
        console.log(res.pages[0].position);
    });
}).end();
