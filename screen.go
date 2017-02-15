package nicegl

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func panic_if(err error) {
	if err != nil {
		panic(err)
	}
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

func (cfg *ScreenCfg) InitAndMakeWindow() (*glfw.Window, func()) {
	err := glfw.Init()
	panic_if(err)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, cfg.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, cfg.Minor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w, err := glfw.CreateWindow(cfg.Width, cfg.Height, cfg.Title, nil, nil)
	panic_if(err)
	w.MakeContextCurrent()

	err = gl.Init()
	panic_if(err)

	return w, glfw.Terminate
}

func (cfg *ScreenCfg) AspectRatio() float32 {
	return float32(cfg.Width) / float32(cfg.Height)
}

const (
	fps_weight = 0.1
	fps_denom  = 1.0 / (fps_weight + 1)
)

func UpdateFPS(fps *float64, t *time.Time) int64 {
	next := time.Now()
	df := next.Sub(*t) / 1000 // usec
	*fps = (1000*1000*fps_weight/float64(df) + *fps) * fps_denom
	*t = next
	return int64(df)
}
