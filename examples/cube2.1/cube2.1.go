package main

/*
	dependencies
		go get -u github.com/go-gl/gl/v2.1/gl
		go get -u github.com/go-gl/glfw/v3.2/glfw
*/

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Vec2 [2]float32
type Vec3 [3]float32
type Vec4 [4]float32

func (v *Vec2) Addr() *float32 {
	return &v[0]
}
func (v *Vec2) GlTextureCoord() {
	gl.TexCoord2fv(v.Addr())
}

func (v *Vec3) Addr() *float32 {
	return &v[0]
}

func (v *Vec3) GlVertex() {
	gl.Vertex3fv(v.Addr())
}

func (v *Vec4) Addr() *float32 {
	return &v[0]
}

type Quad struct {
	Color   Vec4
	Normal  Vec3
	Texture [4]Vec2
	Vertex  [4]Vec3
}

func (q *Quad) DrawImm() {
	gl.Color4fv(q.Color.Addr())
	gl.Normal3fv(q.Normal.Addr())
	q.Texture[0].GlTextureCoord()
	q.Vertex[0].GlVertex()
	q.Texture[1].GlTextureCoord()
	q.Vertex[1].GlVertex()
	q.Texture[2].GlTextureCoord()
	q.Vertex[2].GlVertex()
	q.Texture[3].GlTextureCoord()
	q.Vertex[3].GlVertex()
}

type Cube [6]Quad

func (c *Cube) DrawImm() {
	//gl.BindTexture(gl.TEXTURE_2D, texture)
	for i := 0; i < 6; i++ {
		(&c[i]).DrawImm()
	}
}

// vs[2]-----vs[3]
// |\         | \
// | vs[0]-------vs[1]
// | |        |  |
// | |        |  |
// vs|6]-----vs[7|
//  \|          \|
//	 vs[4]-------vs[5]
func (c *Cube) Make(vs *[8]Vec3, colors *[6]Vec4) {
	c[0].Color = colors[0]
	c[0].Normal = Vec3{0, 0, 1}
	c[0].Vertex[0] = vs[0]
	c[0].Vertex[1] = vs[1]
	c[0].Vertex[2] = vs[2]
	c[0].Vertex[3] = vs[3]

	c[1].Color = colors[1]
	c[0].Normal = Vec3{0, 0, -1}
	c[0].Vertex[0] = vs[4]
	c[0].Vertex[1] = vs[5]
	c[0].Vertex[2] = vs[6]
	c[0].Vertex[3] = vs[7]

	c[2].Color = colors[2]
	c[0].Normal = Vec3{0, -1, 0}
	c[0].Vertex[0] = vs[0] // TODO, and below
	c[0].Vertex[1] = vs[1]
	c[0].Vertex[2] = vs[2]
	c[0].Vertex[3] = vs[3]

	c[3].Color = colors[3]
	c[0].Normal = Vec3{0, 1, 0}
	c[0].Vertex[0] = vs[0]
	c[0].Vertex[1] = vs[1]
	c[0].Vertex[2] = vs[2]
	c[0].Vertex[3] = vs[3]

	c[4].Color = colors[4]
	c[0].Normal = Vec3{1, 0, 0}
	c[0].Vertex[0] = vs[0]
	c[0].Vertex[1] = vs[1]
	c[0].Vertex[2] = vs[2]
	c[0].Vertex[3] = vs[3]

	c[5].Color = colors[5]
	c[0].Normal = Vec3{-1, 0, 0}
	c[0].Vertex[0] = vs[0]
	c[0].Vertex[1] = vs[1]
	c[0].Vertex[2] = vs[2]
	c[0].Vertex[3] = vs[3]

}

func init() {
	// Goroutine running GLFW event handling must run on main thread
	runtime.LockOSThread()
}

func panic_if(err error) {
	if err != nil {
		panic(err)
	}
}

func glfw_must_init() func() {
	err := glfw.Init()
	panic_if(err)
	return glfw.Terminate
}

type ScreenCfg struct {
	Title  string
	Width  int
	Height int
	Major  int
	Minor  int
	// Monitor ??
	// shared Window ??
}

func (cfg *ScreenCfg) MakeWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, cfg.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, cfg.Minor)

	w, err := glfw.CreateWindow(cfg.Width, cfg.Height, cfg.Title, nil, nil)
	panic_if(err)
	win.MakeContextCurrent()

	err = gl.Init()
	panic_if(err)

	return w
}

type SceneCfg struct {
	Light      LightCfg
	ClearColor [4]float32
	Frustum    [6]float64
	Drawable   []func()
}

func (cfg *SceneCfg) Setup() {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	if cfg.Light != (LightCfg{}) {
		gl.Enable(gl.LIGHTING)
	}
}

func (cfg *SceneCfg) DrawFrustum() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(cfg.Frustum[0], cfg.Frustum[1], cfg.Frustum[2], cfg.Frustum[3], cfg.Frustum[4], cfg.Frustum[5])
}

func (cfg *SceneCfg) DrawImm() {
	// TODO: add fps
	gl.ClearColor(cfg.ClearColor[0], cfg.ClearColor[1], cfg.ClearColor[2], cfg.ClearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	cfg.Light.Set()
	cfg.DrawFrustum()

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)

	gl.Rotatef(rotationX, 1, 0, 0)
	gl.Rotatef(rotationY, 0, 1, 0)
	rotationX += 0.5
	rotationY += 0.5

	gl.Begin(gl.QUADS)
	for _, d := range cfg.Drawable {
		d()
	}
	gl.End()
}

type LightCfg struct {
	Ambient       [4]float32
	Diffuse       [4]float32
	LightPosition [4]float32
}

func (l *LightCfg) Set() {
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &l.Ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &l.Diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &l.LightPosition[0])
	gl.Enable(gl.LIGHT0)
}

var (

	// Geometry
	cube Cube = [6]Quad{
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{0, 0, 1},
			[4]Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			[4]Vec3{{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1}},
			// 0, 1, 3, 4
		},
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{0, 0, -1},
			[4]Vec2{{1, 0}, {1, 1}, {0, 1}, {0, 0}},
			[4]Vec3{{-1, -1, -1}, {-1, 1, -1}, {1, 1, -1}, {1, -1, -1}},
		},
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{0, 1, 0},
			[4]Vec2{{0, 1}, {0, 0}, {1, 0}, {1, 1}},
			[4]Vec3{{-1, 1, -1}, {-1, 1, 1}, {1, 1, 1}, {1, 1, -1}},
		},
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{0, -1, 0},
			[4]Vec2{{1, 1}, {0, 1}, {0, 0}, {1, 0}},
			[4]Vec3{{-1, -1, -1}, {1, -1, -1}, {1, -1, 1}, {-1, -1, 1}},
		},
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{1, 0, 0},
			[4]Vec2{{1, 0}, {1, 1}, {0, 1}, {0, 0}},
			[4]Vec3{{1, -1, -1}, {1, 1, -1}, {1, 1, 1}, {1, -1, 1}},
		},
		Quad{
			[4]float32{1, 1, 1, 1},
			[3]float32{-1, 0, 0},
			[4]Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			[4]Vec3{{-1, -1, -1}, {-1, -1, 1}, {-1, 1, 1}, {-1, 1, -1}},
		},
	}

	light = LightCfg{
		Ambient:       [4]float32{0.5, 0.5, 0.5, 1},
		Diffuse:       [4]float32{1, 1, 1, 1},
		LightPosition: [4]float32{-5, 5, 10, 0},
	}

	screen = ScreenCfg{
		Title:  "hw",
		Width:  800,
		Height: 600,
		Major:  2,
		Minor:  1,
	}

	scene = SceneCfg{
		Light:      light,
		ClearColor: [4]float32{0.5, 0.8, 0.5, 0},
		Frustum:    [6]float64{-1, 1, -1, 1, 1.0, 10.0},
		Drawable:   []func(){cube.DrawImm},
	}

	rotationX float32
	rotationY float32
)

func loop(w *glfw.Window, s *SceneCfg) {
	s.Setup()
	t := time.Now()
	fps := float64(0)
	for !w.ShouldClose() {
		s.DrawImm()
		w.SwapBuffers()
		glfw.PollEvents()
		nicegl.UpdateFps(&fps, &t)
		fmt.Printf("%.1f\n", fps)
	}
}

func main() {
	defer glfw_must_init()()

	win := screen.MakeWindow()

	loop(win, &scene)
	//texture = newTexture("square.png")
	//defer gl.DeleteTextures(1, &texture)
}
