/* Author: 
 * http://www.elated.com/articles/rotatable-3d-product-boxshot-threejs/
 */
window.onload = function() {
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
    document.getElementById('dvd').appendChild( renderer.domElement );

    // Create the scene to hold the object
    var scene = new THREE.Scene();
    // Create the camera
    var camera = new THREE.Camera(
        42,                       // Field of view
        width / height,           // Aspect ratio
        .1,                       // Near plane distance
        10000                     // Far plane distance
    );

    // Position the camera
    camera.position.set( -5, -2, 12 );

    // Add the lights
    var light = new THREE.PointLight( 0xffffff, .4 );
    light.position.set( 10, 10, 10 );
    scene.addLight( light );
    ambientLight = new THREE.AmbientLight( 0xbbbbbb );
    scene.addLight( ambientLight );

    // Create the materials
    var materialClass = WebGLSupported ? THREE.MeshLambertMaterial : THREE.MeshBasicMaterial;
    var darkGrey =  new materialClass( { color: 0x333333 } );
    var bookCover = new materialClass( { color: 0xffffff, map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/front.png' ) } );
    var bookSpine = new materialClass( { color: 0x000000 } );
    var bookPages = new materialClass( { color: 0x000000 } );
    var bookPagesTop = new materialClass( { color: 0x333333 } );
    var bookPagesBottom = new materialClass( { color: 0x333333 } );

    var materials = [
        bookSpine,          // Left side
        bookPages,          // Right side
        bookPagesTop,       // Top side
        bookPagesBottom,    // Bottom side
        bookCover,          // Front side
        darkGrey            // Back side
    ];

    // Create the dvd and add it to the scene
    document.getElementById('dvd').classList.add('three-dee');
    var dvd =  new THREE.Mesh( new THREE.CubeGeometry( 6, 8.55, 0.5, 4, 4, 1, materials ), new THREE.MeshFaceMaterial() );
    scene.addChild( dvd );

    // Begin the animation
    animate();


    /*
       Animate a frame
       */

    function animate() {
        // Rotate the book based on the current mouse position
        dvd.rotation.y = mouseX * 0.0005;
        dvd.rotation.x = mouseY * 0.0005;
        // Render the frame
        renderer.render( scene, camera );
        // Keep the animation going
        requestAnimFrame( animate );
    }
}

/*
 * Check if the browser supports WebGL
 * Adapted from http://doesmybrowsersupportwebgl.com/
 **/
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
