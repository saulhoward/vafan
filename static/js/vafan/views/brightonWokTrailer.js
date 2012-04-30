/**
 * Brighton Wok trailer specific video view.
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
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.brightonWokTrailer = Backbone.View.extend({
    /*
     * Set up any video links
     */
    initialize: function()
    {
        var vidView, modalHtml, videoHtml;

        vidView = this;
        modalHtml = $('#modal-template').html();
        $('#video .video-wrapper').append(modalHtml);

        videoHtml = $('#video-template').html();

        $('#brighton-wok-trailer').on('show', function () {
            var $modalBody = $(
                '#brighton-wok-trailer .modal-body');
            if ($('iframe', $modalBody).length == 0) {
                $modalBody.html(videoHtml);
            }
            jQuery('iframe.vimeo', this).each(function(){
                Froogaloop(this).addEvent(
                    'ready', vidView.videoReady);
            });
            $('#video-selector').addClass('playing');
        })
        $('#brighton-wok-trailer').on('hide', function () {
            $('#video-selector').removeClass('playing');
        })
 
        /* Needed as the 'video-wrapper' div is relative 
         * and doesn't close the modal */
        $('#video-selector .video-wrapper, #video-selector .video').live('click', function(){
            if ($('#brighton-wok-trailer').is(":visible")) {
                $('#brighton-wok-trailer').modal('hide');
            } else {
                $('#brighton-wok-trailer').modal('show');
            }
            return true;
        });
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
});

