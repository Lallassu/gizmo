//=============================================================
// fragmentshaders.go
//-------------------------------------------------------------
// Various fragment shaders used.
//=============================================================
package main

//=============================================================
// Fragment shader for lights
//=============================================================
var fragmentShaderLight = `
#version 330 core

in vec2  vTexCoords;
in vec2  vPosition;
in vec4  vColor;

out vec4 fragColor;

uniform float uPosX;
uniform float uPosY;
uniform float uRadius;

void main() {
   vec4 c = vColor;

   // Normalized distance where min(0), max(radius)
   float dist = abs(distance(vec2(uPosX, uPosY), vPosition))/uRadius;

   fragColor = vec4(c.r-(dist/c.r), c.g-(dist/c.g), c.b-(dist/c.b), c.a-dist);
}

`

//=============================================================
// Full screen fragment shader
//=============================================================
var fragmentShaderFullScreen = `
#version 330 core

in vec2  vTexCoords;
in vec4  vColor;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
   vec4 c = vColor;
   vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
   vec4 tx = texture(uTexture, t);
   if (c.r == 1) {
      fragColor = vec4(tx.r, tx.g, tx.b, tx.a);
   } else {
       if (c.a == 0.1111) {
          c *= 2;
          vec3 fc = vec3(1.0, 0.3, 0.1);
          vec2 borderSize = vec2(0.1); 

          vec2 rectangleSize = vec2(1.0) - borderSize; 

          float distanceField = length(max(abs(c.x)-rectangleSize,0.0) / borderSize);

          float alpha = 1.0 - distanceField;
          fc *= abs(0.8 / (sin( c.x + sin(c.y)+ 1.3 ) * 5.0) );
          fragColor = vec4(fc, alpha*5);
       } else {
          fragColor = vColor;
       }
   }
}

`
