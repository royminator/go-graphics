package main

import (
	"go-graphics/pkg/gfx"
	"io/ioutil"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	verts = []mgl.Vec3{
		{-0.5, -0.5, 0},
		{0.5, -0.5, 0},
		{0.0, 0.5, 0},
	}
	colors = []mgl.Vec4{
        {1.0, 0.0, 0.0, 1.0},
        {0.0, 1.0, 0.0, 1.0},
        {0.0, 0.0, 1.0, 1.0},
    }
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

	vBuf := gfx.CreateVertexBuffer(verts)
	cBuf := gfx.CreateColorBuffer(colors)
	gfx.CreateVertexArray(vBuf, cBuf)

	gl.UseProgram(prog)

	// Uniform setup
	mvp := mgl.Ident4()
	mvpLoc := gl.GetUniformLocation(prog, gl.Str("mvp\x00"))

	// Input processing
	shouldRun := true
	for shouldRun {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				shouldRun = false
				break
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.GetKeyFromName("w"):
						mvp = mgl.Translate3D(mvp.At(0, 3), 0.01+mvp.At(1, 3), mvp.At(2, 3))
					case sdl.GetKeyFromName("a"):
						mvp = mgl.Translate3D(mvp.At(0, 3)-0.01, mvp.At(1, 3), mvp.At(2, 3))
					case sdl.GetKeyFromName("r"):
						mvp = mgl.Translate3D(mvp.At(0, 3), mvp.At(1, 3)-0.01, mvp.At(2, 3))
					case sdl.GetKeyFromName("s"):
						mvp = mgl.Translate3D(0.01+mvp.At(0, 3), mvp.At(1, 3), mvp.At(2, 3))
					}
				}
				break
			}
		}
		draw(mvpLoc, mvp)
		window.GLSwap()
	}
}

func draw(mvpLoc int32, mvp mgl.Mat4) {
	gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func panicIfErr(err error) {

	if err != nil {
		panic(err)
	}
}
