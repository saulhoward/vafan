/**
 * Vafan video view.
 * A video as seen on a 'video page' (eg, /videos/brighton-wok-trailer)
 * Saul <saul@saulhoward.com>
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.video = Backbone.View.extend({

    el: '#video',

    initialize: function(props)
    {
        var v = this;
        this.$el = $(this.el);
        this.appView = props.appView;
        this.model.on('change', this.render);

        // Edit features
        if (this.model != null) {
            if (this.model.get('video').isEditable == true) {
                this.appView.on('vafan:startEdit', function() {
                    v.startEdit();
                });
                this.appView.on('vafan:stopEdit', function() {
                    v.stopEdit();
                });
                this.appView.on('vafan:saveEdits', function(){
                    v.saveVideo();
                });
            }
        }

        this.render();
    },

    render: function ()
    {
       // var v = this;
       // return this;
    },

    // Open the interface for editing.
    startEdit: function()
    {
        var v, $desc, $descTextarea, $title, $shortDesc;
        v = this;

        // Description
        $desc = $('.description', v.el);
        this.prevDescVal = $desc.html();
        $descTextarea = $('<textarea/>', {
            val: v.model.get('video').description
        });
        $desc.html($descTextarea);

        //Title
        $title = $('.title', v.el);
        $title.attr('contenteditable', 'true');
        $title.html(v.model.get('video').title);

        // ShortDescription
        $shortDesc = $('.shortDescription', v.el);
        $shortDesc.attr('contenteditable', 'true');
        $shortDesc.html(v.model.get('video').shortDescription);
    },

    stopEdit: function()
    {
        var v, $desc, $descTextarea, $title, $shortDesc;
        v = this;

        // Description
        $desc = $('.description', v.el);
        $desc.html(this.prevDescVal);

        //Title
        $title = $('.title', v.el);
        $title.attr('contenteditable', 'false');
        $title.html(v.model.get('video').title);

        // ShortDescription
        $shortDesc = $('.shortDescription', v.el);
        $shortDesc.attr('contenteditable', 'false');
        $shortDesc.html(v.model.get('video').shortDescription);
    },

    saveVideo: function()
    {
        var v, desc, shortDesc, title;

        v = this;
        v.appView.trigger('vafan:saving');

        title = $('.title', v.el).html();
        title = v.stripHTML(title);
        $('.title', v.el).html(title);
        v.model.set('title', title);

        shortDesc = $('.shortDescription', v.el).html();
        shortDesc = v.stripHTML(shortDesc);
        $('.shortDescription', v.el).html(shortDesc);
        v.model.set('shortDescription', shortDesc);

        desc = $('.description textarea', v.el).val();
        v.model.set('description', desc);

        // Save model
        v.model.save({
            // empty properties
        }, {
            success: function() 
            {
                v.appView.trigger('vafan:saved');
            }
        });

    },

    stripHTML: function(html) 
    {
        var tempDiv = document.createElement("DIV");
        tempDiv.innerHTML = html;
        return tempDiv.textContent;
    }
});

