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

    // red point light shining from the front
    //var pointLight = new THREE.PointLight( 0xff0000 );
    //pointLight.position.set( 0, 0, -10 );
    //scene.add( pointLight );

    // Create the materials
    var materialClass = WebGLSupported ? THREE.MeshPhongMaterial : THREE.MeshBasicMaterial;
    var darkGrey =  new materialClass( { color: 0x333333 } );
    var dvdCover = new materialClass( {
        color: 0xffffff,
        shininess: 100,
        specular:  0x333333,
        map: THREE.ImageUtils.loadTexture( '/img/brighton-wok/dvd/front.png' ) 
    } );
    var dvdSpine = new materialClass( { color: 0x000000 } );
    var dvdRight = new materialClass( { color: 0x000000 } );
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
    document.getElementById('dvd').classList.add('three-dee');
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
        camera.position.x = base;
        camera.position.z = base + 13;

        // animate the light to bob slightly
        var rand = Math.random();
        light.position.x = (base * 50);
        light.position.z = (base * 50) + 50;

        // animate the pointLight's intensity
        light.intensity = base + 0.9;

        frame += 0.03; // speed
        renderer.render( scene, camera );
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
