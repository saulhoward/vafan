/**
 * Vafan video view.
 * Saul <saul@saulhoward.com
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.video = Backbone.View.extend({

    el: '#video',

    initialize: function()
    {
        var v = this;
        console.log(this.model);
        v.$el = $(v.el);
        if (v.model.get('video').isEditable == true) {
            v.makeEditable();
        }
    },

    // Open the interface for editing.
    makeEditable: function()
    {
        var v, $desc;
        v = this;
        // Description
        $desc = $('.description', v.el);
        $desc.attr('contenteditable', 'true');

        // Save button
        v.$el.before(v.getSaveButton());
    },

    saveVideo: function()
    {
        var v;
        v = this;
        //v.stopEditable();

        v.model.save({
            description: v.$('.description', v.el).html()
        });

    },

    getSaveButton: function()
    {
        var v = this;
        return $('<button/>', {
            id: 'saveEdits',
            html: 'Save',
            click: function(e) {
                log("Saving...");
                v.saveVideo();
            }
        });
    }
});

