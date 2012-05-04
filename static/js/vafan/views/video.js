/**
 * Vafan video view.
 * Saul <saul@saulhoward.com
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.video = Backbone.View.extend({

    el: '#video',

    $saveButton: null,

    initialize: function()
    {
        this.$el = $(this.el);
        log(this.model);
        this.model.bind('change', this.render);
        this.render();
    },

    render: function ()
    {
        // Edit features
        if (this.model != null) {
            if (this.model.get('video').isEditable == true) {
                console.log("rendering editable interface...");
                this.makeEditable();
            }
        }
        return this;
    },

    // Open the interface for editing.
    makeEditable: function()
    {
        var v, $desc, $descTextarea, $title, $shortDesc;
        v = this;

        // Description
        $desc = $('.description', v.el);
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

        // Save button
        $('header.navbar .page-actions.btn-group').prepend(v.getSaveButton());

        //v.$el.before(v.getSaveButton());
    },

    saveVideo: function()
    {
        var v, desc, shortDesc, title;

        v = this;
        v.disableSave();

        title = $('.title', v.el).html();
        title = v.stripHTML(title);
        $('.title', v.el).html(title);

        shortDesc = $('.shortDescription', v.el).html();
        shortDesc = v.stripHTML(shortDesc);
        $('.shortDescription', v.el).html(shortDesc);

        desc = $('.description textarea', v.el).val();

        // Save model
        v.model.save({
            title: title,
            description: desc,
            shortDescription: shortDesc
        }, {
            success: function() 
            {
                v.enableSave();
            }
        });

    },

    disableSave: function() 
    {
        this.getSaveButton().attr({"disabled": "disabled"});    
    },

    enableSave: function() 
    {
        this.getSaveButton().removeAttr('disabled');
    },

    getSaveButton: function()
    {
        var v = this;
        if (v.$saveButton == null) {
            v.$saveButton = $('<button/>', {
                id: 'saveEdits',
                html: 'Save Edits',
                "class": 'btn btn-success pull-right',
                click: function(e) {
                    log("Saving...");
                    v.saveVideo();
                }
            });
            v.$saveButton.prepend('<i class="icon-ok"></i> ');
        }
        return v.$saveButton;
    },

    stripHTML: function(html) 
    {
        var tempDiv = document.createElement("DIV");
        tempDiv.innerHTML = html;
        return tempDiv.textContent;
    }
});

