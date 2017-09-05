var process = require('process');
var fs = require('fs');
var dk = require("./dkbson");
if (process.argv.length <= 1){
    console.log("Usage: node "+process.argv[1]+" bson_file");
    process.exit(1);
}
var decode_bson = function(str){
    res = dk.dkbson.decode(str.trim());
    if (res.status == 'error'){
        console.log(res);
        process.exit(2);
    }
    console.log("%j", res);
};

file_name = process.argv[2];

fs.readFile(file_name, 'utf8', function(err, data) {
    if (err) {
        return console.log(err);
    }
    decode_bson(data);
});
