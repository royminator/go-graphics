package main

import (
	"go-graphics/pkg/gfx"

	"go-graphics/pkg/input"

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

type Camera struct {
	Pose mgl.Mat4
}

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
	mvpLoc := gl.GetUniformLocation(prog.Handle, gl.Str("mvp\x00"))
	camera := Camera{mgl.Ident4()}
	inputMapper := input.MakeMapper()
	inputMapper.RegisterEntity(input.MOVE_NORTH, &camera)

	// Input processing
	shouldRun := true
	for shouldRun {
		input.ReadAndExecInputs(inputMapper)
		draw(mvpLoc, camera)
		window.GLSwap()
	}
}

func draw(mvpLoc int32, camera Camera) {
	gl.UniformMatrix4fv(mvpLoc, 1, false, &camera.Pose[0])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func panicIfErr(err error) {

	if err != nil {
		panic(err)
	}
}

func (c *Camera) ExecInput(action input.Action) {
	switch action {
	case input.MOVE_NORTH:
		c.Pose = c.Pose.Mul4(mgl.Translate3D(0.0, 0.01, 0.0))
		break
	}
}

func (c *Camera) Id() int {
	return 1
}
