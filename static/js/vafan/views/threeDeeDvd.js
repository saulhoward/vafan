/*
 * DVD 3D Box
 * Uses Three.js
 * Saul <saul@saulhoward.com>
 */

if ('undefined' === typeof vafan){vafan = {};}
if ('undefined' === typeof vafan.view){vafan.view = {};}

vafan.view.threeDeeDVD = Backbone.View.extend({

    el: "dvd",

    initialize: function() {
        this.$el = $(this.el);
        this.threeDee();
    },

    /*
     * Create the 3d box
     */
    threeDee: function () 
    {
        var dvdView = this;
        // Set up some variables and add a mousemove handler to the page
        var mouseX = 0;                                       // Mouse X pos relative to window centre
        var mouseY = 0;                                       // Mouse Y pos relative to window centre
        var width = 450; // renderer w
        var height = 594; // renderer h
        var windowCentreX = window.innerWidth / 2;            // Window centre (X pos)
        var windowCentreY = window.innerHeight / 2;           // Window centre (Y pos)
        var WebGLSupported = dvdView.isWebGLSupported();  // Check for WebGL support

        document.addEventListener( 'mousemove', function( event ) {
            // Update mouseX and mouseY based on the new mouse X and Y positions
            mouseX = ( event.clientX - windowCentreX );
            mouseY = ( event.clientY - windowCentreY );
        }, false );

        // Create the renderer and add it to the page's body element
        var renderer = WebGLSupported ? new THREE.WebGLRenderer() : new THREE.CanvasRenderer();
        renderer.setSize( width, height );
        //document.getElementById('dvd').appendChild( renderer.domElement );
        dvdView.$el.append(renderer.domElement);

        // Create the scene to hold the object
        var scene = new THREE.Scene();
        // Create the camera
        var camera = new THREE.PerspectiveCamera(
                32,                       // Field of view
                width / height,           // Aspect ratio
                .1,                       // Near plane distance
                10000                     // Far plane distance
                );

        // Position the camera
        //camera.position.set( -5, -2, 12 );
        camera.position.z = 13;

        // Add the lights
        scene.add(new THREE.AmbientLight(0x444444));
        var mlight = new THREE.PointLight( 0xffffff, 0.3 );
        mlight.position.set( 25, -10, 50 );
        scene.add( mlight );
        var llight = new THREE.PointLight( 0xffffff, 0.3 );
        llight.position.set( -25, 10, 70 );
        scene.add( llight );
        var rlight = new THREE.PointLight( 0xffffff, 0.2 );
        rlight.position.set( 25, 0, 70 );
        scene.add( rlight );

        // Create the materials
        var materialClass = WebGLSupported ? THREE.MeshPhongMaterial : THREE.MeshBasicMaterial;
        var darkGrey =  new materialClass( { color: 0x333333 } );
        var dvdCover = new materialClass( {
            color: 0xffffff,
            shininess: 100,
            specular:  0x666666,
            map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/cover.png' ) 
        } );
        var dvdSpine = new materialClass( {
            color: 0xffffff,
            shininess: 100,
            specular:  0x666666,
            map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/spine.png' ) 
        } );
        var dvdRight = new materialClass( {
            color: 0x151515,
            map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/front.png' ) 
        } );
        var dvdTop = new materialClass( {
            color: 0x333333,
            map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/top.png' ) 
        } );

        var materials = [
            dvdRight,     // Right side
            dvdSpine,     // Left side
            dvdTop,       // Top side
            dvdTop,    // Bottom side
            dvdCover,     // Front side
            darkGrey      // Back side
                ];

        // Create the dvd and add it to the scene
        dvdView.$el.addClass('three-dee');
        var dvd =  new THREE.Mesh( new THREE.CubeGeometry( 4.24, 6, 0.45, 4, 4, 1, materials ), new THREE.MeshFaceMaterial() );
        scene.add(dvd);

        // Begin the animation
        animate();

        var step = 0.002;
        var frame = 0;
        function animate() {
            // Rotate the dvd based on the current mouse position
            dvd.rotation.y =  dvdView.getRotation(mouseX * step, dvd.rotation.y, -1, 1);
            dvd.rotation.x =  dvdView.getRotation(mouseY * step, dvd.rotation.x, -0.4, 0.65);
            // animate the camera to bob slightly
            camera.position.x = (Math.sin(frame) * 0.1);
            camera.position.y = (Math.cos(frame) * 0.1);

            frame += 0.01; // speed 
            renderer.render( scene, camera );
            requestAnimationFrame( animate );
        }
    },

    // Calculate DVD bounded rotation
    getRotation: function(pos, rotation, posLow, posHigh) 
    {
        var difference = function (a, b) { return Math.abs(a - b) }
        var step = 0.002;
        if ((pos < posHigh) && (pos > posLow)) {
            if (difference(pos, rotation) < 0.1) {
                return pos;
            } else {
                if (pos > rotation) {
                    return rotation + 0.05;
                } else if (pos < rotation) {
                    return rotation - 0.05;
                }
            }
        } else if (rotation > 0.2) {
            return rotation - step;
        } else if (rotation < -0.1) {
            return rotation + step;
        }
        return rotation;
    },

    // Check if the browser supports WebGL
    // Adapted from http://doesmybrowsersupportwebgl.com/
    isWebGLSupported: function () 
    {
        var cvs = document.createElement('canvas');
        var contextNames = ["webgl","experimental-webgl","moz-webgl","webkit-3d"];
        var ctx;
        if ( navigator.userAgent.indexOf("MSIE") >= 0 ) {
            try {
                ctx = WebGLHelper.CreateGLContext(cvs, 'canvas');
            } catch(e) {}
        } else {
            for ( var i = 0; i < contextNames.length; i++ ) {
                try {
                    ctx = cvs.getContext(contextNames[i]);
                    if ( ctx ) break;
                } catch(e){}
            }
        }
        if ( ctx ) return true;
        return false;
    }


});

