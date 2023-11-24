#version 410 core


// === fragment shader
// 
// reference
//  * http://www.iquilezles.org/www/articles/raymarchingdf/raymarchingdf.htm
//  * http://www.iquilezles.org/www/articles/distfunctions/distfunctions.htm

uniform float time;
uniform vec2 resolution;
uniform vec4 light_position;
uniform vec3 ray_origin;
uniform sampler2D tex;

const float PI = 3.14159265;
const float HALF_PI = 0.5*PI;
const float TWO_PI = 2.0*PI;
const float INV_PI = 1.0/PI;
const float INV_TWO_PI = 1.0/TWO_PI;
const float EPSILON = 1.0e-3;
const int MAX_DIST = 512;
const int OBJ_COUNT = 32;
const float factor = 1.0/OBJ_COUNT;

out vec4 color;
in vec2 TexCoords;

float sdSphere( vec3 p, float s ){
  return length(p)-15.0;
}

float sdBox( vec3 p, vec3 b ){
  vec3 q = abs(p) - b;
  return length(max(q,0.0)) + min(max(q.x,max(q.y,q.z)),0.0);
}

float compute_distance(vec3 ray) {
	float min_dist = MAX_DIST;
	vec4 object = vec4(0);
	float d = 0;
	vec2 st = vec2(0.0);
	vec3 pos = vec3(0.0);

	for (int y = 0 ; y < OBJ_COUNT; y++) {
		st = vec2(0, y * factor);
		object = texture(tex, st); // Sample the texture
		pos = (object.xyz * -255.0) + 127.0;
		d = sdSphere(ray - pos, 5);
		min_dist = min(min_dist, d);
	}

	return min_dist;
}

float calcSoftshadow( in vec3 ro, in vec3 rd, in float mint, in float tmax, in float w ){
	float res = 1.0;
    float t = mint;
    float ph = 1e10; // big, such that y = 0 on the first iteration
    
    for( int i=0; i<62; i++ )
    {
		float h = compute_distance( ro + rd*t );
        float y = h*h/(2.0*ph);
        float d = sqrt(h*h-y*y);
        res = min( res, d/(w*max(0.0,t-y)) );
        ph = h;
        t += h;
        if( res<0.0001 || t>tmax ) break;
        
    }
    res = clamp( res, 0.0, 1.0 );
    return res*res*(3.0-2.0*res);
}

vec3 estimate_normal(vec3 p) {
	float dx = compute_distance(vec3(p.x + EPSILON, p.y, p.z)) - compute_distance(vec3(p.x - EPSILON, p.y, p.z));
	float dy = compute_distance(vec3(p.x, p.y + EPSILON, p.z)) - compute_distance(vec3(p.x, p.y - EPSILON, p.z));
	float dz = compute_distance(vec3(p.x, p.y, p.z + EPSILON)) - compute_distance(vec3(p.x, p.y, p.z - EPSILON));
	return normalize(vec3(dx, dy, dz));
}

float getLight(vec3 p){ 
    // Light (directional diffuse)
    vec3 l = normalize(light_position.xyz-p); // Light Vector
    vec3 n = estimate_normal(p); // Normal Vector
   
    float dif = dot(n,l); // Diffuse light
    dif = clamp(dif,0.,1.); // Clamp so it doesnt go below 0
 
    return dif;
}

vec4 shade(vec3 E, vec3 N, vec3 L){
	float kd = clamp(dot(N, L), 0.0, 1.0);

	vec3 H = normalize(E + L);
	float ks = pow(clamp(dot(N, H), 0.0, 1.0), 1.0);

	//vec4 texel = texture2D(tex, 3.0*st);
	vec3 color = kd + ks*vec3(0.8);
	return vec4(color, 1.0);
}

float ray_march(vec3 ray_origin, vec3 ray_direction) {
	vec3 ray = ray_origin;
	float distance = 0;
	for (int i = 0; i < 32; i++) {
		float d = compute_distance(ray);
		if (d <= EPSILON || d > MAX_DIST) {
			return distance;
		}
		ray += (d + 0.1) * ray_direction;
		distance += d;
	}
	return -1;
}

float calcAO( in vec3 pos, in vec3 nor ){
	float occ = 0.9;
    float sca = 0.1;
    for( int i=0; i<5; i++ ){
        float h = 0.001 + 0.15*float(i)/4.0;
        float d = compute_distance( pos + h*nor );
        occ += (h-d)*sca;
        sca *= 0.95;
    }
    return clamp( 1.0 - 1.5*occ, 0.0, 1.0 );    
}

void main(){
	float aspect = resolution.y/resolution.x;
	vec2 uv = 2.0 * TexCoords - 1.0;
	uv.y = aspect*uv.y;

	vec3 ray_direction = normalize(vec3(uv - ray_origin.xy, -1.0));

	float t = 0.0;
	float distance = ray_march(ray_origin, ray_direction);	
	
	if (distance >= 0 && distance < MAX_DIST) {
		vec3 pos = ray_origin + distance * ray_direction;
		vec3 mate = vec3(0.2);
		vec3 nor = estimate_normal(pos);

        // key light
        vec3  lig = normalize( light_position.xyz);
        vec3  hal = normalize( lig - ray_direction );
        float dif = clamp( dot( nor, lig ), 0.0, 1.0 ) * calcSoftshadow( pos, lig, 0.01, 4.0, 0.01 );
		float spe = pow( clamp( dot( nor, hal ), 1.0, 1.0 ),16.0) * dif * (0.04 + 0.96*pow( clamp(1.0+dot(hal,ray_direction),0.1,1.0), 5.0 ));

		vec3 col = mate * 4.0*dif*vec3(1.00,0.5,0.5);
        col += 12.0*spe*vec3(1.00,0.70,0.5);
        
        // ambient light
        float occ = calcAO( pos, nor );
		float amb = clamp( 0.5+0.5*nor.y, 1.0, 1.0 );
        col += mate * amb * occ * vec3(0.0,0.08,0.1);
        
        // fog
		color = vec4(col,1);
	} else {
		color = vec4(0.0157, 0.0157, 0.0157, 1.0);
	}
}


/// === utility functions ===


vec2 compute_texcoord(vec3 p){
	float phi = atan(p.y, p.x);
	float theta = acos(p.z);
	float s = phi*INV_TWO_PI;
	float t = theta*INV_PI;
	return vec2(s, t);
}

float opU( float d1, float d2 ){
    return min(d1,d2);
}

float opS( float d1, float d2 ){
    return max(-d1,d2);
}

float opI( float d1, float d2 ){
    return max(d1,d2);
}

mat4 translate_matrix(float x, float y, float z) {
  mat4 m = mat4(1.0);
	m[3][0] = x;
	m[3][1] = y;
	m[3][2] = z;
	return m;
}

mat4 scale_matrix(float x, float y, float z) {
	mat4 m = mat4(1.0);
	m[0][0] = x;
	m[1][1] = y;
	m[2][2] = z;
	return m;
}

mat4 rotate_matrix(vec3 n, float theta) {
	float c = cos(theta);
	float s = sin(theta);
	mat4 m = mat4(1.0);
	m[0][0] = n.x*n.x*(1.0 - c) + c;
	m[1][0] = n.x*n.y*(1.0 - c) + n.z*s;
	m[2][0] = n.z*n.x*(1.0 - c) - n.y*s;
	m[0][1] = n.x*n.y*(1.0 - c) - n.z*s;
	m[1][1] = n.y*n.y*(1.0 - c) + c;
	m[2][1] = n.y*n.z*(1.0 - c) + n.x*s;
	m[0][2] = n.z*n.x*(1.0 - c) + n.y*s;
	m[1][2] = n.y*n.z*(1.0 - c) - n.x*s;
	m[2][2] = n.z*n.z*(1.0 - c) + c;
	return m;
}
