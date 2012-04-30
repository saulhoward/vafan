/*
 * Loads the webfonts using Google's API
 * Saul <saul@saulhoward.com>
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.fonts = Backbone.View.extend({
    initialize: function ()
    {
        WebFontConfig = {
            google: { families: [ 
                        'Acme::latin',
                        'Bangers::latin',
                        'Ultra::latin', 
                        'Ubuntu::latin'
                            ] }
        };
        (function() {
            var wf = document.createElement('script');
            wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
            '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
        wf.type = 'text/javascript';
        wf.async = 'true';
        var s = document.getElementsByTagName('script')[0];
        s.parentNode.insertBefore(wf, s);
        })(); 
    }
});
