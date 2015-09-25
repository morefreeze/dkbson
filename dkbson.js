BinaryParser = function(t, e) {
    this.bigEndian = t;
    this.allowExceptions = e;
};
p = BinaryParser.prototype;
with (p.encodeFloat = function(t, e, i) {
    var n, o, s, a, l, c = Math.pow(2, i - 1) - 1, u = -c + 1, h = c, d = u - e, f = isNaN(v = parseFloat(t)) || v == -1 / 0 || v == +1 / 0 ? v : 0, p = 0, g = 2 * c + 1 + e + 3, _ = Array(g), m = 0 > (v = 0 !== f ? 0 : v), v = Math.abs(v), y = Math.floor(v), b = v - y;
    for (n = g; n; _[--n] = 0)
        ;
    for (n = c + 2; y && n; _[--n] = y % 2, y = Math.floor(y / 2))
        ;
    for (n = c + 1; b > 0 && n; (_[++n] = ((b *= 2) >= 1) - 0) && --b)
        ;
    for (n = -1; g > ++n && !_[n]; )
        ;
    if (_[(o = e - 1 + (n = (p = c + 1 - n) >= u && h >= p ? n + 1 : c + 1 - (p = u - 1))) + 1]) {
        if (!(s = _[o]))
            for (a = o + 2; !s && g > a; s = _[a++])
                ;
        for (a = o + 1; s && --a >= 0; (_[a] = !_[a] - 0) && (s = 0))
            ;
    }
    for (n = 0 > n - 2 ? -1 : n - 3; g > ++n && !_[n]; )
        ;
    for ((p = c + 1 - n) >= u && h >= p ? ++n : u > p && (p != c + 1 - g && d > p && this.warn("encodeFloat::float underflow"), n = c + 1 - (p = u - 1)), (y || 0 !== f) && (this.warn(y ? "encodeFloat::float overflow" : "encodeFloat::" + f), p = h + 1, n = c + 2, f == -1 / 0 ? m = 1 : isNaN(f) && (_[n] = 1)), v = Math.abs(p + c), a = i + 1, l = ""; --a; l = v % 2 + l, v = v >>= 1)
        ;
    for (v = 0, a = 0, n = (l = (m ? "1" : "0") + l + _.slice(n, n + e).join("")).length, r = []; n; v += (1 << a) * l.charAt(--n), 7 == a && (r[r.length] = String.fromCharCode(v), v = 0), a = (a + 1) % 8)
        ;
    return r[r.length] = v ? String.fromCharCode(v) : "", (this.bigEndian ? r.reverse() : r).join("");
}, p.encodeInt = function(t, e) {
    var i = Math.pow(2, e), n = [];
    for ((t >= i || -(i >> 1) > t) && this.warn("encodeInt::overflow") && (t = 0), 0 > t && (t += i); t; n[n.length] = String.fromCharCode(t % 256), t = Math.floor(t / 256))
        ;
    for (e = -(-e >> 3) - n.length; e--; n[n.length] = "\0")
        ;
    return (this.bigEndian ? n.reverse() : n).join("");
}, p.decodeFloat = function(t, e, i) {
    var n, r, o, s = ((s = new this.Buffer(this.bigEndian, t)).checkBuffer(e + i + 1), s), a = Math.pow(2, i - 1) - 1, l = s.readBits(e + i, 1), c = s.readBits(e, i), u = 0, h = 2, d = s.buffer.length + (-e >> 3) - 1;
    do{
        for (n = s.buffer[++d], r = e % 8 || 8, o = 1 << r; o >>= 1; n & o && (u += 1 / h), h *= 2)
            ;
        e -= r;
    }while (e);
    return c == (a << 1) + 1 ? u ? 0 / 0 : l ? -1 / 0 : +1 / 0 : (1 + -2 * l) * (c || u ? c ? Math.pow(2, c - a) * (1 + u) : Math.pow(2, -a + 1) * u : 0);
}, p.decodeInt = function(t, e, i) {
    var n = new this.Buffer(this.bigEndian, t), r = n.readBits(0, e), o = Math.pow(2, e);
    return i && r >= o / 2 ? r - o : r;
}, {p: (p.Buffer = function(t, e) {
        this.bigEndian = t || 0, this.buffer = [], this.setBuffer(e);
    }).prototype})
    p.readBits = function(t, e) {
        function i(t, e) {
            for (++e; --e; t = 1073741824 == (1073741824 & (t %= 2147483648)) ? 2 * t : 2 * (t - 1073741824) + 2147483647 + 1)
                ;
            return t;
        }
        if (0 > t || 0 >= e)
            return 0;
        this.checkBuffer(t + e);
        for (var n, r = t % 8, o = this.buffer.length - (t >> 3) - 1, s = this.buffer.length + (-(t + e) >> 3), a = o - s, l = (this.buffer[o] >> r & (1 << (a ? 8 - r : e)) - 1) + (a && (n = (t + e) % 8) ? (this.buffer[s++] & (1 << n) - 1) << (a-- << 3) - r : 0); a; l += i(this.buffer[s++], (a-- << 3) - r))
            ;
        return l;
    }, p.setBuffer = function(t) {
        if (t) {
            for (var e, i = e = t.length, n = this.buffer = Array(e); i; n[e - i] = t.charCodeAt(--i))
                ;
            this.bigEndian && n.reverse();
        }
    }, p.hasNeededBits = function(t) {
        return this.buffer.length >= -(-t >> 3);
    }, p.checkBuffer = function(t) {
        if (!this.hasNeededBits(t))
            throw Error("checkBuffer::missing bytes");
    };
p.warn = function(t) {
    if (this.allowExceptions)
        throw Error(t);
    return 1;
}, p.toSmall = function(t) {
    return this.decodeInt(t, 8, !0);
}, p.fromSmall = function(t) {
    return this.encodeInt(t, 8, !0);
}, p.toByte = function(t) {
    return this.decodeInt(t, 8, !1);
}, p.fromByte = function(t) {
    return this.encodeInt(t, 8, !1);
}, p.toShort = function(t) {
    return this.decodeInt(t, 16, !0);
}, p.fromShort = function(t) {
    return this.encodeInt(t, 16, !0);
}, p.toWord = function(t) {
    return this.decodeInt(t, 16, !1);
}, p.fromWord = function(t) {
    return this.encodeInt(t, 16, !1);
}, p.toInt = function(t) {
    return this.decodeInt(t, 32, !0);
}, p.fromInt = function(t) {
    return this.encodeInt(t, 32, !0);
}, p.toDWord = function(t) {
    return this.decodeInt(t, 32, !1);
}, p.fromDWord = function(t) {
    return this.encodeInt(t, 32, !1);
}, p.toFloat = function(t) {
    return this.decodeFloat(t, 23, 8);
}, p.fromFloat = function(t) {
    return this.encodeFloat(t, 23, 8);
}, p.toDouble = function(t) {
    return this.decodeFloat(t, 52, 11);
}, p.fromDouble = function(t) {
    return this.encodeFloat(t, 52, 11);
};
base64 = function() {
        function t(t, e) {
            var i = o.indexOf(t.charAt(e));
            if (-1 === i)
                throw "Cannot decode base64";
            return i;
        }
        function e(e) {
            var i, n, o = 0, s = e.length, a = [];
            if (e += "", 0 === s)
                return e;
            if (0 !== s % 4)
                throw "Cannot decode base64";
            for (e.charAt(s - 1) === r && (o = 1, e.charAt(s - 2) === r && (o = 2), s -= 4), i = 0; s > i; i += 4)
                n = t(e, i) << 18 | t(e, i + 1) << 12 | t(e, i + 2) << 6 | t(e, i + 3), a.push(String.fromCharCode(n >> 16, 255 & n >> 8, 255 & n));
            switch (o) {
                case 1:
                    n = t(e, i) << 18 | t(e, i + 1) << 12 | t(e, i + 2) << 6, a.push(String.fromCharCode(n >> 16, 255 & n >> 8));
                    break;
                case 2:
                    n = t(e, i) << 18 | t(e, i + 1) << 12, a.push(String.fromCharCode(n >> 16));
            }
            return a.join("");
        }
        function i(t, e) {
            var i = t.charCodeAt(e);
            if (i > 255)
                throw "INVALID_CHARACTER_ERR: DOM Exception 5";
            return i;
        }
        function n(t) {
            if (1 !== arguments.length)
                throw "SyntaxError: exactly one argument required";
            t += "";
            var e, n, s = [], a = t.length - t.length % 3;
            if (0 === t.length)
                return t;
            for (e = 0; a > e; e += 3)
                n = i(t, e) << 16 | i(t, e + 1) << 8 | i(t, e + 2), s.push(o.charAt(n >> 18)), s.push(o.charAt(63 & n >> 12)), s.push(o.charAt(63 & n >> 6)), s.push(o.charAt(63 & n));
            switch (t.length - a) {
                case 1:
                    n = i(t, e) << 16, s.push(o.charAt(n >> 18) + o.charAt(63 & n >> 12) + r + r);
                    break;
                case 2:
                    n = i(t, e) << 16 | i(t, e + 1) << 8, s.push(o.charAt(n >> 18) + o.charAt(63 & n >> 12) + o.charAt(63 & n >> 6) + r);
            }
            return s.join("");
        }
        var r = "=", o = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", s = "1.0";
        return {decode: e,encode: n,VERSION: s};
    }();

dkbson = function() {
        function t(t) {
            return decodeURIComponent(escape(t));
        }
        function e(t, e) {
            for (var i = 1, n = t.charCodeAt(e), r = n, o = 128; 4 >= i && 0 != (o & n); )
                o >>= 1, i++, r = (r << 1) % 256;
            for (var s = 0; i - 1 > s; s++)
                r >>= 1;
            var a = r;
            for (s = 1; i > s; s++)
                a = (a << 8) + t.charCodeAt(e + s);
            return [e + i, a];
        }
        function i(i, n, r) {
            var o, s = new BinaryParser();
            if (i == l.TYPE_INT8)
                o = s.toSmall(n.substr(r, 1)), r++;
            else if (i == l.TYPE_INT16)
                o = s.toShort(n.substr(r, 2)), r += 2;
            else if (i == l.TYPE_INT32)
                o = s.toInt(n.substr(r, 4)), r += 4;
            else if (i == l.TYPE_INT64)
                o = s.decodeInt(n.substr(r, 8), 64, !0), r += 8;
            else if (i == l.TYPE_DOUBLE)
                o = s.toDouble(n.substr(r, 8)), r += 8;
            else if (i == l.TYPE_FLOAT){
                o = s.toFloat(n.substr(r, 4)), r += 4;
            }
            else if (i == l.TYPE_STRING) {
                var a = e(n, r);
                r = a[0];
                var c = a[1];
                o = n.substr(r, c), o = t(o), r += c;
            } else
                i == l.TYPE_NULL ? (o = null, r++) : i == l.TYPE_BOOL ? (o = 0 === n.charCodeAt(r) ? !1 : !0, r++) : i == l.TYPE_REAL16 ? (o = s.toShort(n.substr(r, 2)), o /= 100, r += 2) : i == l.TYPE_REAL24 ? (o = s.toInt(n.substr(r, 3) + "\0"), o /= 1e3, r += 3) : i == l.TYPE_REAL32 ? (o = s.toInt(n.substr(r, 4)), o /= 1e4, r += 4) : console.log("error: unsupported type:" + i.charCodeAt(0));
            return [r, o];
        }
        function n(t, o) {
            var s = [], a = e(t, o);
            o = a[0];
            for (var c = a[1], u = c; u-- > 0; ) {
                var h, d, f = t.charCodeAt(o), p = o + 1;
                f == l.TYPE_OBJECT ? (h = r(t, p), o = h[0], d = h[1]) : f == l.TYPE_ARRAY ? (h = n(t, p), o = h[0], d = h[1]) : (h = i(f, t, p), o = h[0], d = h[1]), s.push(d);
            }
            return [o, s];
        }
        function r(t, o) {
            var s = e(t, o);
            o = s[0];
            for (var a = s[1], c = {}, u = a; u-- > 0; ) {
                var h, d = t.charCodeAt(o), f = "", p = "";
                for (h = o + 1; "\0" != t.charAt(h); h++)
                    f += t.charAt(h);
                h++;
                var g;
                d == l.TYPE_OBJECT ? (g = r(t, h), o = g[0], p = g[1]) : d == l.TYPE_ARRAY ? (g = n(t, h), o = g[0], p = g[1]) : (g = i(d, t, h), o = g[0], p = g[1]), c[f] = p
            }
            return [o, c];
        }
        var a = base64, l = {TYPE_INT8: 1,TYPE_INT16: 2,TYPE_INT32: 3,TYPE_INT64: 4,TYPE_FLOAT: 16,TYPE_DOUBLE: 17,TYPE_REAL16: 18,TYPE_REAL24: 19,TYPE_REAL32: 20,TYPE_STRING: 32,TYPE_BOOL: 48,TYPE_NULL: 49,TYPE_OBJECT: 64,TYPE_ARRAY: 65};
        return {decode: function(t) {
                var e = a.decode(t), i = r(e, 0);
                return i[1];
            }};
    }();

exports.dkbson = dkbson;

var request = require('request');
var sleep = require('sleep');
// DO NOT ADD console.log, use console.error instead or it will make other program confused
var req = function(options_or_url, end_cb){
    var MAX_RETRY = 3;
    var TIMEOUT_MS = 3000;
    var SLEEP_MS = 500;
    var options = {};
    if ('string' == typeof(options_or_url)){
        var url = options_or_url;
        options.url = url;
    }
    else{
        options = options_or_url;
    }
    function get_url(options){
        return new Promise(function(resolve, reject) {
            function retry(timeout) {
                if (timeout > 3) return Promise.reject(new Error("max retry times"));
                var l_options = {'url': url, 'timeout': timeout*1000};
                request(l_options, function(error, response, body) {
                    if (error){
                        //console.error(timeout + ' ' +options.url + ' ' + error);
                        retry(timeout+1);
                        return ;
                    }
                    resolve(body);
                });
            }
            retry(1);
        });
    }
    return get_url(options);
};
exports.req = req;
