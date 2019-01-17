//=============================================================
// fragmentshaders.go
//-------------------------------------------------------------
// Various fragment shaders used.
//=============================================================
package main

//=============================================================
// Fragment shader for menu items
//=============================================================
var fragmentShaderMenuItem = `
#version 330 core

in vec2  vTexCoords;
in vec2  vPosition;
in vec4  vColor;

out vec4 fragColor;

uniform float uTime;
uniform int uSelected;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() 
{
	
   vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
   vec4 c = texture(uTexture, t);
   if (uSelected != 0 && c.r != 0) {
   		 float add = clamp(sin(uTime*15), 0.2,0.8);
         fragColor = vec4(c.rg, tan(c.b+add), c.a-add);
   } else {
		fragColor = c;
   }
}
`

//=============================================================
// Fragment shader for doors
//=============================================================
var fragmentShaderDoor = `
#version 330 core

in vec2  vTexCoords;
in vec2  vPosition;
in vec4  vColor;

out vec4 fragColor;

uniform float uTime;
uniform vec2 uPos;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() 
{
	
    vec2 t = gl_FragCoord.xy / uTexBounds.zw;
	   vec4 c_ = texture(uTexture, t);
   if (c_.r == 0 && c_.g == 0 && c_.b == 1) {
       float dist = distance(uPos, vPosition)/uPos.y;
	   float o = sin(dist);
	   vec4 c = texture(uTexture, clamp(t*o, 0.5, 0.8));
	   fragColor = vec4(dist/clamp(cos(uTime),0.4, 0.6)*c.r*o, dist*c.g*o, c.b/2, c.a*clamp((1.0-dist), 0.2, 0.8));
   } else {
	   fragColor = c_;
   }
}
`

//=============================================================
// Fragment shader for minimap
//=============================================================
var fragmentShaderMinimap = `
#version 330 core

in vec2  vTexCoords;
in vec2  vPosition;
in vec4  vColor;

out vec4 fragColor;

uniform float uTime;
uniform vec2 uPos;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
   vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
   vec4 tx = texture(uTexture, t);
   vec4 c = vec4(tx.r, tx.g, tx.b, tx.a);
   float dist = distance(uPos, vPosition)/uPos.y;
   if (dist < 0.02) {
   		if (sin(uTime*15) < 0) {
			fragColor = vec4(1.0, 0, 0, 1.0);
		} else {
			fragColor = vec4(0, 0 ,0 ,0);
		}
 //  } else if (dist < 0.5) {
 //       fragColor = vec4(c.r/dist/2, c.g, c.b/dist/2, c.a/dist/2);
   } else if (dist < 0.6) {
		fragColor = c;
   } else {
		fragColor = vec4(0, 0, 0, 0);
   }
}
`

//=============================================================
// Fragment shader for portals
//=============================================================
var fragmentShaderPortal = `
#version 330 core

in vec2  vTexCoords;
in vec2  vPosition;
in vec4  vColor;

out vec4 fragColor;

uniform float uTime;
uniform vec2 uPos;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
   vec4 c = vColor;
   
   if (c.r == 0 && c.g == 0 && c.b == 1) {
	   float py = uPos.y+cos(uTime*5)*3;
	   float px = uPos.x+sin(uTime*5)*5;
       float dist = distance(vec2(px, py), vPosition)/uPos.y;
	   float o = clamp(sin(uTime/dist), 0.5, 1.0);
	   if (dist < 0.2) {
	   	   float t = dist*2;
	       fragColor = vec4(t+sin(uTime), t, t, dist);
	   } else {
	       fragColor = vec4(dist*o*(1.0-dist), c.g*o*(1.0-dist), (c.b/2)*o*(1.0-dist), c.a*clamp((1.0-dist), 0.2, 0.8));
	   }
   } else {
	   fragColor = c;
   }
}
`

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
   float dist = distance(vec2(uPosX, uPosY), vPosition)/uRadius;
   fragColor = vec4(c.r*(1.0-dist), c.g*(1.0-dist), c.b*(1.0-dist), c.a*(1.0-dist));
   
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
	      // particleFire
          c *= 2;
          vec3 fc = vec3(1.0, 0.3, 0.1);
          vec2 borderSize = vec2(0.1); 

          vec2 rectangleSize = vec2(1.0) - borderSize; 

          float distanceField = length(max(abs(c.x)-rectangleSize,0.0) / borderSize);

          float alpha = 1.0 - distanceField;
          fc *= abs(0.8 / (sin( c.x + sin(c.y)+ 1.3 ) * 5.0) );
          fragColor = vec4(fc, alpha*5);
	   } else if (c.a == 0.2222) {
	      // particleEvaporate
          vec3 fc = vec3(c.r, c.g, c.b);
          vec2 borderSize = vec2(0.1); 

          vec2 rectangleSize = vec2(1.0) - borderSize; 

          float distanceField = length(max(abs(c.x)-rectangleSize,0.0) / borderSize);

          float alpha = 1.0 - distanceField;
          fc *= abs(0.5 / (sin( c.x + sin(c.y)+ 1.3 ) * 2.0) );
          fragColor = vec4(fc.r, fc.g, fc.b, alpha/2);
       } else {
          fragColor = vColor;
       }
   }
}
`
