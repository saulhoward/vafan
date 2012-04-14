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
    },

    linkifyTweets: function()
    {
        _.each($('.tweet-box span.text'), function(t) {
            $(t).html(vafan.twitter.linkify(t.innerHTML));
        });
    },

    linkify: function(inputText) 
    {
        var replaceText, replacePattern1, replacePattern2, replacePattern3;

        //URLs starting with http://, https://, or ftp://
        replacePattern1 = /(\b(https?|ftp):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/gim;
        replacedText = inputText.replace(replacePattern1, '<a href="$1" target="_blank">$1</a>');

        //URLs starting with "www." (without // before it, or it'd re-link the ones done above).
        replacePattern2 = /(^|[^\/])(www\.[\S]+(\b|$))/gim;
        replacedText = replacedText.replace(replacePattern2, '$1<a href="http://$2" target="_blank">$2</a>');

        //Change email addresses to mailto:: links.
        replacePattern3 = /(\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,6})/gim;
        replacedText = replacedText.replace(replacePattern3, '<a href="mailto:$1">$1</a>');

        return replacedText
    }
}
