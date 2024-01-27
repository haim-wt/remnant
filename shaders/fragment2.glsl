#version 150

in VertexData
{
    vec4 v_position;
    vec3 v_normal;
    vec2 v_texcoord;
} inData;

out vec4 fragColor;

uniform vec2 resolution;
uniform float time;
uniform vec4 mouse; 

const float PI          = 3.14159265359;
const float PI2         = 6.28318530718;

const float MAX_DIST    = 100.;
const float MIN_DIST    = .001;

float hash21(vec2 p) {  
    return fract(sin(dot(p, vec2(27.609, 57.583)))*43758.5453); 
}
mat2 rot(float a){
    return mat2(cos(a),sin(a),-sin(a),cos(a));
}

// gyroid function
float sdGry(vec3 p, float s, float t, float b) {
    p *=s;
    float g = abs(dot(sin(p*1.25), cos(p.zxy))-b)/(s*1.25)-t;
    return g;
}

//global vars cause its just a demo
float g2,g3,g4,g5;
vec3 hitPoint;
mat2 rx,ry;

vec2 map(vec3 p) {
    vec2 res=vec2(1000.,0.);
    
    p.zy*=rx;
    p.xz*=ry;
    
    float spc = .6;
    float hlf = spc/2.;
    vec3 q = p;
    
    // move scene
    p.z -= time*.2;
    p.y -= .3;

    float g1 = sdGry(p, 3.31, .03, 1.45);
    g2 = sdGry(p, 36.3, .003, 1.2);
    g3 = sdGry(p, 20.3, .03, 1.2);
    g4 = sdGry(p, 16.3, .03, .2);
    g5 = sdGry(p, 24.3, .03, .2);

    g1 -= (g2 *.15);
    g1 -= (g3 *.02);
    g1 -= (g4 *.32);
    g1 -= (g5 *.01);

    if(g1<res.x) res=vec2(g1*.5,1.);
    
    hitPoint = p;
    return res;
}

vec3 normal(vec3 p, float t) {
    t*=MIN_DIST;
    float d = map(p).x;
    vec2 e = vec2(t,0);
    vec3 n = d - vec3(
        map(p-e.xyy).x,
        map(p-e.yxy).x,
        map(p-e.yyx).x
        );
    return normalize(n);
}

vec2 marcher( in vec3 ro, in vec3 rd, int x) {
    float t = 0.;
    float m = 0.;
    for( int i=0; i<x; i++ ) {
        vec2 d = map(ro + rd * t);
        m = d.y;
        if(d.x<.0001*t||t>MAX_DIST) break;
        t += d.x*.6;
    }
    return vec2(t,m);
}

vec3 shade(vec3 p) {
    float cks1 = smoothstep(.01,.011, abs(abs(min(g5,g3)-.0035)-.0025)-.015);
    float cks2 = smoothstep(.01,.011, abs(abs(g4-.005)-.0025)-.015);
    vec3 ckr = mix(vec3(.8),vec3(.5,0,0),clamp(cks1,0.,1.)); 
    return mix(ckr,vec3(0),clamp(cks2,0.,1.));   
}

void main() {
    // Set Coordinates from -1 to 1
    vec2 F = gl_FragCoord.xy;
    vec2 uv = (2. * gl_FragCoord.xy - resolution.xy )/resolution.y;

    vec3 ro = vec3(0,0,.15);
    vec3 rd = normalize(vec3(uv,-1));
    
        
    // mouse //
    float x = -(mouse.y/resolution.y*2.-1.)*PI;
    float y = -(mouse.x/resolution.x*2.-1.)*PI;
    
    // rotation of ray order & direction
    rx =rot(x);
    ry =rot(y);
    
    vec2 ray = marcher(ro, rd, 128);

    vec3 C = vec3(0);
    float t = ray.x;
    float m = ray.y;
    if(t<MAX_DIST) {
        vec3 p = ro + t * rd;
        vec3 n = normal(p, t);
        
        float diff = n.x*.5+.5;
        
        if(ray.y == 1.){
            C += diff * shade(p);  
        }

        C = mix( C, vec3(.8), 1.-exp(-.075*t*t*t));
    } 

    C= pow(C, vec3(0.4545));
    fragColor = vec4(C,1.0);
}
