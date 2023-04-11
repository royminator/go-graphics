#version 150
in vec2 position;
in vec4 color;

out vec4 color_vsout;

void main() {
    gl_Position = vec4(position, 0.0, 1.0);
    color_vsout = color;
}
