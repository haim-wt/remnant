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

const float PI = 3.14159265;
const float HALF_PI = 0.5*PI;
const float TWO_PI = 2.0*PI;
const float INV_PI = 1.0/PI;
const float INV_TWO_PI = 1.0/TWO_PI;
const float EPSILON = 1.0e-4;
const int MAX_DIST = 100;
const int OBJ_COUNT = 256;

out vec4 color;
in vec2 TexCoords;

vec3[OBJ_COUNT] objects;
float[OBJ_COUNT] distances;

float sdSphere( vec3 p, float s )
{
  return length(p)-5.0;
}

float sdBox( vec3 p, vec3 b )
{
  vec3 q = abs(p) - b;
  return length(max(q,0.0)) + min(max(q.x,max(q.y,q.z)),0.0);
}

float opU( float d1, float d2 )
{
    return min(d1,d2);
}

float opS( float d1, float d2 )
{
    return max(-d1,d2);
}

float opI( float d1, float d2 )
{
    return max(d1,d2);
}

mat4 translate_matrix(float x, float y, float z) 
{
  mat4 m = mat4(1.0);
	m[3][0] = x;
	m[3][1] = y;
	m[3][2] = z;
	return m;
}

mat4 scale_matrix(float x, float y, float z) 
{
	mat4 m = mat4(1.0);
	m[0][0] = x;
	m[1][1] = y;
	m[2][2] = z;
	return m;
}

mat4 rotate_matrix(vec3 n, float theta) 
{
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

void getObjects() {
	for (int i=0;i<OBJ_COUNT;i++){
		objects[i] = vec3(-10 + 8*i, -10 + 8*i,  -50 - 8*i);
	}
}

float compute_distances(vec3 ray) 
{
	float min_dist = MAX_DIST;
	for (int obj = 0 ; obj < OBJ_COUNT; obj++) {
		distances[obj] = sdSphere(ray - objects[obj], 1);
		min_dist = min(min_dist, distances[obj]);
	}

	return min_dist;
}

int compute_near(vec3 ray) 
{
	float max_delta = 0;
	float delta = 0;
	int near = 0;

	for (int i = 0 ; i < 50; i++) {
		delta = sdSphere(ray - objects[i], 2);
		delta = distances[i] - delta;
		if (delta > max_delta) {
			max_delta = delta;
			near = i;
		}
	}

	return near;
}

vec3 estimate_normal(vec3 p) 
{
	// float dx = compute_distance(vec3(p.x + EPSILON, p.y, p.z)) - compute_distance(vec3(p.x - EPSILON, p.y, p.z));
	// float dy = compute_distance(vec3(p.x, p.y + EPSILON, p.z)) - compute_distance(vec3(p.x, p.y - EPSILON, p.z));
	// float dz = compute_distance(vec3(p.x, p.y, p.z + EPSILON)) - compute_distance(vec3(p.x, p.y, p.z - EPSILON));
	// return normaliz

	return vec3(0);
}

// 'p' must be normalized
vec2 compute_texcoord(vec3 p)
{
	float phi = atan(p.y, p.x);
	float theta = acos(p.z);
	float s = phi*INV_TWO_PI;
	float t = theta*INV_PI;
	return vec2(s, t);
}

vec4 shade(vec3 E, vec3 N, vec3 L, vec2 st)
{
	float kd = clamp(dot(N, L), 0.0, 1.0);

	vec3 H = normalize(E + L);
	float ks = pow(clamp(dot(N, H), 0.0, 1.0), 1.0);

	//vec4 texel = texture2D(tex, 3.0*st);
	vec3 color = kd + ks*vec3(0.8);
	return vec4(color, 1.0);
}

float ray_march(vec3 ray_origin, vec3 ray_direction) {
	vec3 ray = ray_origin * ray_direction;
	for (int i = 0; i < MAX_DIST; i++) {
		float d = compute_distances(ray);
		if (d < EPSILON) {
			return 1;
		}
		ray += d * ray_direction;
	}
	return -1;
}

void main()
{
	float aspect = resolution.y/resolution.x;
	vec2 uv = 2.0 * TexCoords - 1.0;
	uv.y = aspect*uv.y;

	getObjects();

	vec3 ray_direction = normalize(vec3(uv - ray_origin.xy, -1.0));

	vec3 P = vec3(0.0);
	vec3 E = vec3(0.0);
	vec3 L = vec3(0.0);
	vec3 N = vec3(0.0);
	vec2 st = vec2(0.0);

	float t = 0.0;
	float distance = ray_march(ray_origin, ray_direction);	

	// N = estimate_normal( P );
	// L = normalize((light_position).xyz - P);
	// E = normalize( P );
	//st = compute_texcoord(E);

	if (distance >= 0) {
		color = vec4(1, 1, 1, 1.0);
	} else {
		color = vec4(0.05, 0.0, 0.0, 1.0);
	}
}