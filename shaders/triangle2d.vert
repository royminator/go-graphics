#version 150
in vec3 position;
in vec4 color;

uniform mat4 mvp;

out vec4 color_vsout;

void main() {
    gl_Position = mvp * vec4(position, 1.0);
    color_vsout = color;
}
