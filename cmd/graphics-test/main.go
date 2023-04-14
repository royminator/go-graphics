package main

import (
	"go-graphics/pkg/gfx"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
    "go-graphics/pkg/input"
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
    inputCtx = input.Context{
        Actions: map[input.Action][]input.ButtonState{
            input.MOVE_NORTH: {
                {Id: input.W, Mode: input.PRESSED},
            },
        },
    }
    appState = input.AppState{
        InputContext: inputCtx,
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

	prog, err := gfx.LoadVertFragFromFile("shaders/", "triangle2d")
	if err != nil {
		panic(err)
	}

	vdata := gfx.FromVec3(verts)
	cdata := gfx.FromVec4(colors)
	vBuf := gfx.CreateVertexBufferG(vdata)
	cBuf := gfx.CreateVertexBufferG(cdata)
	gfx.CreateVertexArray(vBuf, cBuf)

	gl.UseProgram(prog.Handle)

	// Uniform setup
	mvp := mgl.Ident4()
	mvpLoc := gl.GetUniformLocation(prog.Handle, gl.Str("mvp\x00"))

	// Input processing
	shouldRun := true
	for shouldRun {
        input.ReadAndExecInputs(&appState)
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
