/*
 * Streams tweets from the vafan server using websockets.
 * Saul <saul@saulhoward.com>
 */
if ('undefined' === typeof vafan) {
    vafan = {};
}

vafan.twitter = {
    streamTweets: function() 
    {
        console.log("init");
        if (ws != null) {
            ws.close();
            ws = null;
        }
        var div = document.getElementById("tweets");
        //ws = new WebSocket("ws://localhost:8888/tweets/stream");
        var ws;
        ws = new WebSocket("ws://dev.convictfilms.com:8888/tweets/stream");
        ws.onopen = function () {
            div.innerText = "opened\n" + div.innerText;
        };
        ws.onmessage = function (e) {
            div.innerText = "msg:" + e.data + "\n" + div.innerText;
            if (e.data instanceof ArrayBuffer) {
                s = "ArrayBuffer: " + e.data.byteLength + "[";
                var view = new Uint8Array(e.data);
                for (var i = 0; i < view.length; ++i) {
                    s += " " + view[i];
                }
                s += "]";
                div.innerText = s + "\n" + div.innerText;
            }
        };
        ws.onclose = function (e) {
            div.innerText = "closed\n" + div.innerText;
        };
        console.log("init");
        div.innerText = "init\n" + div.innerText;
    }
}
