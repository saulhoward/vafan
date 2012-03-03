/*
 * Loads the webfonts for vafan's use
 * Saul <saul@saulhoward.com>
 */
if ('undefined' === typeof vafan) {
    vafan = {};
}

vafan.fonts = {
    // Web fonts
    load: function ()
    {
        WebFontConfig = {
            google: { families: [ 
                        'Acme::latin',
                        'Bangers::latin',
                        'Ultra::latin' 
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
}
