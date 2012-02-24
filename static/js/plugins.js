
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

// Check if the browser supports WebGL
// Adapted from http://doesmybrowsersupportwebgl.com/
function isWebGLSupported() {
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

// Web fonts
function loadWebFonts() {
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

// Vafan specific functions
if ('undefined' === typeof vafan) {
    vafan = {};
}
// DVD 3D Box
// Uses Three.js
vafan.dvd = function () 
{
    // Set up some variables and add a mousemove handler to the page
    var mouseX = 0;                                               // Mouse X pos relative to window centre
    var mouseY = 0;                                               // Mouse Y pos relative to window centre
    var width = 380; // renderer w
    var height = 482; // renderer h
    var windowCentreX = window.innerWidth / 2;                    // Window centre (X pos)
    var windowCentreY = window.innerHeight / 2;                   // Window centre (Y pos)
    var WebGLSupported = isWebGLSupported();                      // Check for WebGL support

    document.addEventListener( 'mousemove', function( event ) {
        // Update mouseX and mouseY based on the new mouse X and Y positions
        mouseX = ( event.clientX - windowCentreX );
        mouseY = ( event.clientY - windowCentreY );
    }, false );

    // Create the renderer and add it to the page's body element
    var renderer = WebGLSupported ? new THREE.WebGLRenderer() : new THREE.CanvasRenderer();
    renderer.setSize( width, height );
    //document.getElementById('dvd').appendChild( renderer.domElement );
    $('#dvd').append(renderer.domElement);

    // Create the scene to hold the object
    var scene = new THREE.Scene();
    // Create the camera
    var camera = new THREE.PerspectiveCamera(
        42,                       // Field of view
        width / height,           // Aspect ratio
        .1,                       // Near plane distance
        10000                     // Far plane distance
    );

    // Position the camera
    //camera.position.set( -5, -2, 12 );
    camera.position.z = 13;

    // Add the lights
    var light = new THREE.PointLight( 0xffffff, 0.8 );
    light.position.z = 50;
    scene.add( light );

    var ambientLight = new THREE.AmbientLight( 0xbbbbbb );
    ambientLight.position.z = 50;
    scene.add( ambientLight );

    // red point light shining from the back
    var pointLight = new THREE.PointLight( 0xff0000 );
    pointLight.position.set( 0, 0, -10 );
    scene.add( pointLight );

    // Create the materials
    var materialClass = WebGLSupported ? THREE.MeshPhongMaterial : THREE.MeshBasicMaterial;
    var darkGrey =  new materialClass( { color: 0x333333 } );
    var dvdCover = new materialClass( {
        color: 0xffffff,
        shininess: 100,
        specular:  0x333333,
        map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/front.png' ) 
    } );
    var dvdSpine = new materialClass( { color: 0x151515 } );
    var dvdRight = new materialClass( { color: 0x151515 } );
    var dvdTop = new materialClass( { color: 0x333333 } );
    var dvdBottom = new materialClass( { color: 0x333333 } );

    var materials = [
        dvdSpine,          // Left side
        dvdRight,          // Right side
        dvdTop,       // Top side
        dvdBottom,    // Bottom side
        dvdCover,          // Front side
        darkGrey            // Back side
    ];

    // Create the dvd and add it to the scene
    $('#dvd').addClass('three-dee');
    var dvd =  new THREE.Mesh( new THREE.CubeGeometry( 6, 8.55, 0.5, 4, 4, 1, materials ), new THREE.MeshFaceMaterial() );
    scene.add( dvd );

    // Begin the animation
    animate();

    var frame = 0;
    function animate() {
        var base = Math.sin(frame) * 0.1;
        // Rotate the dvd based on the current mouse position
        dvd.rotation.y = mouseX * 0.0005;
        dvd.rotation.x = mouseY * 0.0005;

        // animate the camera to bob slightly
        var rand = Math.random();
        camera.position.x = (base * 2);
        camera.position.z = (base * 2) + 13;

        // animate the light to bob slightly
        var rand = Math.random();
        light.position.x = (base * 50);
        light.position.z = (base * 50) + 50;

        // animate the pointLight's intensity
        light.intensity = base + 0.9;

        frame += 0.03; // speed
        renderer.render( scene, camera );
        requestAnimationFrame( animate );
    }
}

