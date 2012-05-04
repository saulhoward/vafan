/**
 * Vafan app view. Sort of the main controller view.
 * Saul <saul@saulhoward.com
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.view){vafan.view={};}

vafan.view.app = Backbone.View.extend({

    el: 'body',

    isEditing: false,

    $saveButton: null,
    $editButton: null,
    $pageActions: $('header.navbar .page-actions.btn-group'),

    initialize: function()
    {
        var fonts, dvd, bwT, tweetBox, tweetBubble, appView;

        appView = this;
        this.$el = $(this.el);

        // Bindings
        //this.on('vafan:startEdit', this.startEdit);
        this.on('vafan:saving', this.disableSave);
        this.on('vafan:saved', this.enableSave);

        // General javascript features (every page).
        if ($('.carousel').length > 0) {
            _.each( $('.carousel'), function(c){
                $(c).carousel();
            });
        }
        if ($('.datepicker').length > 0) {
            _.each($('.datepicker'), function(d){
                $(d).datepicker();
            });
        }
        fonts = new vafan.view.fonts();

        // Index page
        if (appView.$el.hasClass('index')) {
            // 3D DVD Case - only if webgl.
            if ((Modernizr.webgl) && $('#movie .dvd').length > 0) {
                dvd = new vafan.view.threeDeeDVD({
                    el: "#movie .dvd"
                });
            }

            // BWok trailer
            if ($('#video').length > 0) {
                bwT = new vafan.view.brightonWokTrailer();
            }
        }

        // Video page
        else if (appView.$el.hasClass('video')) {
            // Create a video model, with the JSON from this page's URL.
            video = new vafan.model.video({
                url: window.jsonURL
            });
            // Fetch the video, and start the view on success.
            video.fetch({
                success: function() {
                            videoView = new vafan.view.video({
                                model: video,
                                      appView:   appView
                            });
                            if (video.get('video').isEditable == true) {
                                appView.addEditButton();
                            }
                        }
            });
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
    },

    saveEdits: function ()
    {
        this.trigger('vafan:saveEdits');
    },

    startEdit: function ()
    {
        this.isEditing = true;
        this.$pageActions.prepend(this.getSaveButton());
        this.trigger('vafan:startEdit');
    },

    stopEdit: function ()
    {
        this.isEditing = false;
        $('#saveEdits', this.$pageActions).remove();
        this.trigger('vafan:stopEdit');
    },

    disableSave: function() 
    {
        this.getSaveButton().button('loading');    
    },

    enableSave: function() 
    {
        this.getSaveButton().button('reset');    
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
                    v.saveEdits();
                }
            });
            v.$saveButton.attr("data-loading-text", "Saving...");
            v.$saveButton.prepend('<i class="icon-ok"></i> ');
        }
        return v.$saveButton;
    },

    addEditButton: function()
    {
        this.$pageActions.prepend(this.getEditButton());
    },

    getEditButton: function()
    {
        var v = this;
        if (v.$editButton == null) {
            v.$editButton = $('<button/>', {
                id: 'editPage',
                html: 'Edit Page',
                "class": 'btn btn-primary pull-right',
                click: function(e) {
                    if (v.isEditing) {
                        v.stopEdit();
                    } else {
                        v.startEdit();
                    }
                }
            });
            v.$editButton.attr("data-toggle", "button");
            v.$editButton.prepend('<i class="icon-pencil"></i> ');
        }
        return v.$editButton;
    }


});

