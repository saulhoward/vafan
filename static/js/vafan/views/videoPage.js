/**
 * Vafan video page view.
 * Saul <saul@saulhoward.com
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.videoPage = Backbone.View.extend({

    el: '#video',

    initialize: function()
    {
       this.makeEditable($('.description', this.el));
    },

    makeEditable: function(div)
    {
        div.attr('contenteditable', 'true');
        $button = $('<button>Save</button>')
        div.after($button);

    // $('button').on('click', function() {
        // var siblings = $(this).siblings(),
        // content1 = $(siblings[1]).html(),
        // content2 = $(siblings[2]).html(),
        // dataString = 'content1=' + content1 + '&content2=' + content2;

    // $.ajax({
        // type: 'post',
        // url: 'update.php',
        // data: dataString
    // });

// });

}


});

