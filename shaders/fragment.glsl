#version 410 core

uniform float time;
uniform vec2 resolution;
uniform sampler2D tex;

uniform vec3 camera_position;
uniform vec3 camera_direction;
uniform vec3 camera_up;
uniform float camera_fov;

uniform vec3 light;

const float PI = 3.14159265;
const float RADIAN = PI / 180.0;
const float EPSILON = 1.0e-4;
const int MAX_DIST = 1024;
const int OBJ_COUNT = 1;
const int MAX_STEPS = 128;
const float factor = 1.0 / float(OBJ_COUNT);

out vec4 color;
in vec2 TexCoords;

float rand(vec2 c){
	return fract(sin(dot(c.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

vec2 hash( vec2 p )      // this hash is not production ready, please
{                        // replace this by something better
	p = vec2( dot(p,vec2(127.1,311.7)),
			  dot(p,vec2(269.5,183.3)));

	return -1.0 + 2.0*fract(sin(p)*43758.5453123);
}

vec3 noise(vec2 x )
{
    vec2 i = floor( x );
    vec2 f = fract( x );

    vec2 u = f*f*f*(f*(f*6.0-15.0)+10.0);
    vec2 du = 30.0*f*f*(f*(f-2.0)+1.0);
    
    vec2 ga = hash( i + vec2(0.0,0.0) );
    vec2 gb = hash( i + vec2(1.0,0.0) );
    vec2 gc = hash( i + vec2(0.0,1.0) );
    vec2 gd = hash( i + vec2(1.0,1.0) );
    
    float va = dot( ga, f - vec2(0.0,0.0) );
    float vb = dot( gb, f - vec2(1.0,0.0) );
    float vc = dot( gc, f - vec2(0.0,1.0) );
    float vd = dot( gd, f - vec2(1.0,1.0) );

    return vec3( va + u.x*(vb-va) + u.y*(vc-va) + u.x*u.y*(va-vb-vc+vd),   // value
                 ga + u.x*(gb-ga) + u.y*(gc-ga) + u.x*u.y*(ga-gb-gc+gd) +  // derivatives
                 du * (u.yx*(va-vb-vc+vd) + vec2(vb,vc) - va));
}

float fbm(vec2 x, float H )
{    
    float G = exp2(-H);
    float f = 0.5;
    float a = 0.5;
    float t = 0.2;
    for( int i=0; i<8; i++ )
    {
        t += a*noise(f*x).x;
        f *= 2.0;
        a *= G;
    }
    return t;
}

float sdSphere(vec3 p, float s) {
    vec3 n = normalize(vec3(0,1,0));
    //return dot(p,n)+5 - fbm(p.xz, 1);
    return length(p) - (8 + fbm(p.xy, 1));
}

float sdBox(vec3 p, vec3 b) {
    vec3 q = abs(p) - b;
    return length(max(q, 0.0)) + min(max(q.x, max(q.y, q.z)), 0.0);
}

float rand(vec2 uv)
{
    return fract(sin(dot(uv, vec2(12.9898, 78.233))) * 43758.5453);
}

float noise(vec2 p)
{
    vec2 ip = floor(p);
    float r00 = rand(ip);
    float r01 = rand(ip + vec2(0, 1));
    float r10 = rand(ip + vec2(1, 0));
    float r11 = rand(ip + vec2(1, 1));
    vec2 fp = smoothstep(0., 1., p - ip);
    return mix(mix(r00, r01, fp.y), mix(r10, r11, fp.y), fp.x);
}

float compute_distance(vec3 ray) {
    float min_dist = MAX_DIST;
    vec4 object = vec4(0);
    float d = 0;
    vec2 st = vec2(0.0);
    vec3 pos = camera_position;

    for (int y = 0; y < OBJ_COUNT; y++) {
        st = vec2(0.0, y * factor);
        object = texture(tex, st); // Sample the texture
		pos = (object.xyz * -8.0) + 4.0;
		d = sdSphere(ray - pos, y % 8);
        min_dist = min(min_dist, d);
    }

    return min_dist;
}

float calcSoftshadow(vec3 ro, vec3 rd, float mint, float tmax, float w) {
    float res = 1.0;
    float t = mint;
    float ph = 1e10;

    for (int i = 0; i < 32; i++) {
        float h = compute_distance(ro + rd * t);
        float y = h * h / (2.0 * ph);
        float d = sqrt(h * h - y * y);
        res = min(res, d / (w * max(0.0, t - y)));
        ph = h;
        t += h;
        if (res < 0.0001 || t > tmax) break;
    }
    return clamp(res, 0.0, 1.0);
}

vec3 estimate_normal(vec3 p) {
    float dx = compute_distance(vec3(p.x + EPSILON, p.y, p.z)) - compute_distance(vec3(p.x - EPSILON, p.y, p.z));
    float dy = compute_distance(vec3(p.x, p.y + EPSILON, p.z)) - compute_distance(vec3(p.x, p.y - EPSILON, p.z));
    float dz = compute_distance(vec3(p.x, p.y, p.z + EPSILON)) - compute_distance(vec3(p.x, p.y, p.z - EPSILON));
    return normalize(vec3(dx, dy, dz));
}

float ray_march(vec3 ray_origin, vec3 ray_direction) {
    vec3 ray = ray_origin;
    float totalDistance = 0.0;
    for (int i = 0; i < MAX_STEPS; i++) {
        float distanceToSurface = compute_distance(ray);
        if (distanceToSurface < EPSILON) {
            return totalDistance;
        }
        ray += distanceToSurface * ray_direction;
        totalDistance += distanceToSurface;
        if (totalDistance > MAX_DIST) {
            return -1.0;
        }
    }
    return -1.0;
}

void main() {
    float aspect = resolution.y / resolution.x;
    vec2 uv = 2.0 * TexCoords - 1.0;
    uv.y *= aspect;

    float fovFactor = tan(camera_fov * 0.5 * RADIAN);
    vec3 forward = normalize(camera_direction);
    vec3 right = normalize(cross(camera_up, forward));
    vec3 up = normalize(cross(forward, right));


    vec3 ray_direction = normalize(forward + fovFactor * uv.x * right + fovFactor * uv.y * up);
	//vec3 ray_direction = normalize(vec3(uv - camera_position.xy, -1.0));

    float distance = ray_march(camera_position, ray_direction);
    if (distance >= 0.0) {
        vec3 pos = camera_position + distance * ray_direction;
        vec3 nor = estimate_normal(pos);

        vec3 lig = normalize(light - pos);
        float dif = clamp(dot(nor, lig), 0.0, 1.0) * calcSoftshadow(pos, lig, 0.1, 2.0, 0);
        vec3 col = vec3(0.8549, 0.5843, 0.5843) * dif * vec3(0.7); // Simple color multiplication for demonstration

        // fog
        float fogFactor = 1.0 - exp(-0.0001 * distance );
        col = mix(col, vec3(0.1), fogFactor); // Simple linear interpolation for fog effect
        color = vec4(col, 1.0);
    } else {
        color = vec4(0.1); // Default color when no hit is detected
    }
}
