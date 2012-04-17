/**
 * Vafan main onload script
 * Saul <saul@saulhoward.com
 */
$(function () {

    vafan.fonts.load();

    if ($('.datepicker').length > 0) {
        $('.datepicker').datepicker();
    }

    // 3D DVD Case
    if ($('#dvd').length > 0) {
        //vafan.threeDeeDvd.start();
    }

    // Modal videos
    if ($('#video').length > 0) {
        vafan.video.start();
    }

    // Tweets
    if ($('.tweet-box').length > 0) {
        vafan.twitter.linkifyTweets()
       //vafan.twitter.streamTweets()
    }

    $('.carousel').carousel();

});

