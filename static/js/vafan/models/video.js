/**
 * Vafan video model.
 * Saul <saul@saulhoward.com>
 */
if ('undefined' === typeof vafan){vafan={};}
if ('undefined' === typeof vafan.model){vafan.model={};}

vafan.model.video = Backbone.Model.extend({
    initialize: function(props) 
    { 
        this.url = props.url;
    }
});

