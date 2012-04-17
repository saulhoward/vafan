// -- GLOBAL FUNCTIONS

// Console Log
// usage: log('inside coolFunc', this, arguments);
// paulirish.com/2009/log-a-lightweight-wrapper-for-consolelog/
window.log = function(){
  log.history = log.history || [];   // store logs to an array for reference
  log.history.push(arguments);
  if(this.console) {
    arguments.callee = arguments.callee.caller;
    var newarr = [].slice.call(arguments);
    (typeof console.log === 'object' ? log.apply.call(console.log, console, newarr) : console.log.apply(console, newarr));
  }
};
// make it safe to use console.log always
(function(b){function c(){}for(var d="assert,count,debug,dir,dirxml,error,exception,group,groupCollapsed,groupEnd,info,log,timeStamp,profile,profileEnd,time,timeEnd,trace,warn".split(","),a;a=d.pop();){b[a]=b[a]||c}})((function(){try
{console.log();return window.console;}catch(err){return window.console={};}})());


// i have no idea what this is but i imagine it must be awesoWHAAT
var kkeys = [], konami = "38,38,40,40,37,39,37,39,66,65";
$(document).keydown(function(e) {
  kkeys.push( e.keyCode );
  if ( kkeys.toString().indexOf( konami ) >= 0 ){
    $(document).unbind('keydown',arguments.callee);
    $.fn.invert && $('body').invert();
  }
});

// jquery invert plugin
// by paul irish

// some (bad) code from this css color inverter
//    http://plugins.jquery.com/project/invert-color
// some better code via Opera to inverse images via canvas
//    http://dev.opera.com/articles/view/html-5-canvas-the-basics/#insertingimages
// and some imagesLoaded stuff from me
//    http://gist.github.com/268257

// the code is still pretty shit. 
// needs better hex -> rgb conversion..
// it could use color caching and image caching and also some perf-y skip elements type stuff.
(function ($) {
    $.fn.invert = function () {
        $(this).find('*').andSelf().each(function () {

            var t = $(this);
            change_part('color', t);
            change_part('backgroundColor', t);

            $.each(['right', 'bottom', 'left', 'top'], function (i, val) {
                change_part('border-' + val + '-color', t);
            });

            changeImage(t);

        });

        return this;
    };


    function change_part(prop, t) {

        var c = to_rgb(t.css(prop));

        if (c != 'transparent' && c.substring(0, 4) != 'rgba') {
            var rgb = /rgb\((.+),(.+),(.+)\)/.exec(c);
            t.css(prop, 'rgb(' + (255 - rgb[1]) + ',' + (255 - rgb[2]) + ',' + (255 - rgb[3]) + ')');
        }
    }



    function to_rgb(c) {
        if (c.substring(0, 3) == 'rgb' || c == 'transparent') return c;

        if (c.substring(0, 1) == '#') {
            if (c.length == 4) {
                c = '#' + c.substring(1, 2) + c.substring(1, 2) + c.substring(2, 3) + c.substring(2, 3) + c.substring(3, 4) + c.substring(3, 4)
            }

            return 'rgb(' + parseInt(c.substring(1, 3), 16) + ',' + parseInt(c.substring(3, 5), 16) + ',' + parseInt(c.substring(5, 7), 16) + ')';
        }

        var name = ['black', 'white', 'red', 'yellow', 'gray', 'blue', 'green', 'lime', 'fuchsia', 'aqua', 'maroon', 'navy', 'olive', 'purple', 'silver', 'teal'],
            replace = ['(0, 0, 0)', '(255, 255, 255)', '(255, 0, 0)', '(255, 255, 0)', '(128, 128, 128)', '(0, 0, 255)', '(0,128,0)', '(0, 255, 0)', '(255, 0, 255)', '(0, 255, 255)', '(128, 0, 0)', '(0, 0, 128)', '(128, 128, 0)', '(128, 0, 128)', '(192, 192, 192)', '(0, 128, 128)'];

        c = replace[$.inArray(c, name)] || c;

        return 'rgb' + c;
    }

    function changeImage(t) {
        // only operate on img tags and backgroundImages in dataurl or url form.
        if (!(t.is('img') || /^(data|url)/.test(t.css('backgroundImage')))) return;


        var img = new Image();

        $(img).bind('load', function () {

            var elem = document.createElement('canvas');
            elem.width = img.width;
            elem.height = img.height;

            var context = elem.getContext('2d'),
                x = 0,
                y = 0,
                result;

            try {
                    // Draw the image on canvas.
                    context.drawImage(img, x, y);
                    // Get the pixels.
                    var imgd = context.getImageData(x, y, img.width, img.height);
                    var pix = imgd.data;
                    // Loop over each pixel and invert the color.
                    for (var i = 0, n = pix.length; i < n; i += 4) {
                        pix[i] = 255 - pix[i]; // red
                        pix[i + 1] = 255 - pix[i + 1]; // green
                        pix[i + 2] = 255 - pix[i + 2]; // blue
                        // i+3 is alpha (the fourth element)
                    }
                    // Draw the ImageData object.
                    context.putImageData(imgd, x, y);
                    result = elem.toDataURL();
                }
            catch (e) {
                    // image is on a different domain.
                }

            if (t.is('img')) t.attr('src', result)
            else t.css('backgroundImage', 'url(' + result + ')');

        }).each(function () {
            // cached images don't fire load sometimes, so we reset src.
            if (this.complete || this.complete === undefined) {
                var src = this.src;
                this.src = '#';
                this.src = src;
            }
        });

        var match = t.css('backgroundImage').match(/url\((.*?)\)/),
            bg = match && match[1];
        // img src, url bg, datauri background.
        img.src = t[0].src || bg || t.css('backgroundImage');

    }

})(jQuery);
