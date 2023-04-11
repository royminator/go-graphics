package main

import (
	"go-graphics/pkg/gfx"
	"io/ioutil"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
    verts = []float32 {-0.5, -0.5, 0.5, -0.5, 0.0, 0.5}
    colors = []float32 {1.0,0.0,0.0,1.0,0.0,1.0,0.0,1.0,0.0,0.0,1.0,1.0}
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("go-graphics", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_OPENGL)
	defer window.Destroy()
	panicIfErr(err)

	ctx, err := window.GLCreateContext()
	defer sdl.GLDeleteContext(ctx)
	panicIfErr(err)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertSrc, err := ioutil.ReadFile("shaders/triangle2d.vert")
	panicIfErr(err)

	fragSrc, err := ioutil.ReadFile("shaders/triangle2d.frag")
	panicIfErr(err)

	vs, err := gfx.CreateShader(gl.VERTEX_SHADER, string(vertSrc))
	panicIfErr(err)

	fs, err := gfx.CreateShader(gl.FRAGMENT_SHADER, string(fragSrc))
	panicIfErr(err)

    prog := gfx.CreateShaderProgram([]uint32{vs, fs})

    vBuf := gfx.CreateBuffer(verts)
    cBuf := gfx.CreateBuffer(colors)
    gfx.CreateVertexArray(vBuf, cBuf)

    gl.UseProgram(prog)

    shouldRun := true
    for shouldRun {
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                shouldRun = false
            }
        }
        draw()
        window.GLSwap()
    }
}

func draw() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func panicIfErr(err error) {

	if err != nil {
		panic(err)
	}
}
