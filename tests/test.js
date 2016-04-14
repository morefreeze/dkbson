var url = require('url');
var dk = require("./dkbson");
var url = 'http://pages.read.duokan.com/mfsv2/download/s010/p01frvInFioU/TQ1KxPO3jFZA1LQ.js';
dk.req(url, function(body) {
    //console.log('body: ' + body);
    data = body.slice(13,-3);
    //console.log('data: ' + data);
    res = dkbson.decode(data);
    text = '';
    for (var i = 0;i < 100;++i){
        if (res.items[i].type == 'word')
            text += res.items[i].char;
    }
    console.log(text);
    //console.log(res.items[0].pos);
});
