/**
 * Vafan main onload script.
 * This script is run on every page.
 * Saul <saul@saulhoward.com>
 */

$(function () {
    var fonts, dvd, bwT,
        tweetBox, tweetBubble,
        video, videoView;

    // Webfonts.
    fonts = new vafan.view.fonts();

    // 3D DVD Case - only if webgl.
    if ((Modernizr.webgl) && $('#movie .dvd').length > 0) {
        dvd = new vafan.view.threeDeeDVD({
            el: "#movie .dvd"
        });
    }

    // Modal videos - trailer only for now.
    if ($('body.index #video').length > 0) {
        bwT = new vafan.view.brightonWokTrailer();
    }

    // Tweets.
    if ($('#featured-tweets').length > 0) {
        tweetBox = new vafan.view.twitter({
            el:      '#featured-tweets .tweet-box',
            tweetEl: 'span.text'
        });
    }
    if ($('#latest-tweet').length > 0) {
        tweetBubble = new vafan.view.twitter({
            el:      "#latest-tweet",
            tweetEl: ".tweet-bubble"
        });
    }

    // 'carousel' is from bootstrap main lib.
    if ($('.carousel').length > 0) {
        $('.carousel').carousel();
    }

    // 'datepicker' comes courtesy of bootstrap_datepicker.
    if ($('.datepicker').length > 0) {
        $('.datepicker').datepicker();
    }

    // Video page view
    if ($('body.video').length > 0) {
        // Create a video model, with the JSON from this page's URL.
        video = new vafan.model.video({
            url: window.jsonURL
        });
        // Fetch the video, and start the view on success.
        video.fetch({
            success: function() 
            {
                videoView = new vafan.view.video({
                    model: video
                });
            }
        });
    }
});

