/**
 * Vafan video controller.
 *
 * Uses Vimeo lib and bootstrap modals.
 *
 * Modal: http://twitter.github.com/bootstrap/javascript.html#modals
 *
 * Vimeo: http://vimeo.com/api/docs/player-js
 *        https://github.com/vimeo/player-api/tree/master/javascript
 *
 *        http://vimeo.com/api/oembed.json?url=http%3A//vimeo.com/1285896
 *        http://vimeo.com/api/v2/video/1285896.json
 *
 * Saul <saul@saulhoward.com
 */
if ('undefined' === typeof vafan) {
    vafan = {};
}

vafan.video = {
    /*
     * Set up any video links
     */
    start: function()
    {
        var modalHtml = $('#modal-template').html();
        $('#video .content').append(modalHtml);

        var videoHtml = $('#video-template').html();

        $('#brighton-wok-trailer').on('show', function () {
            var $modalBody = $('#brighton-wok-trailer .modal-body');
            if ($('iframe', $modalBody).length == 0) {
                $modalBody.html(videoHtml);
            }
            jQuery('iframe.vimeo', this).each(function(){
                Froogaloop(this).addEvent('ready', vafan.video.videoReady);
            });
        })
    },

    /* Set up a Vimeo Iframe
    */
    videoReady: function (playerID)
    {
        Froogaloop(playerID).api('play');
        $('#brighton-wok-trailer').on('hide', function () {
            Froogaloop(playerID).api('pause');
        })
        $('#brighton-wok-trailer').on('show', function () {
            Froogaloop(playerID).api('play');
        })
    }
}

