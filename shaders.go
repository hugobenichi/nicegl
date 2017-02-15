package nicegl

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func LoadShader(path string, shaderType uint32) (uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	buffer := make([]byte, info.Size()+1)
	if _, err := f.Read(buffer); err != nil {
		return 0, err
	}
	buffer[len(buffer)-1] = 0

	return Compile(unsafeString(buffer), shaderType)
}

func unsafeString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func LoadVertexShader(path string) (uint32, error) {
	return LoadShader(path, gl.VERTEX_SHADER)
}

func LoadFragmentShader(path string) (uint32, error) {
	return LoadShader(path, gl.FRAGMENT_SHADER)
}

func Compile(shader string, shaderType uint32) (uint32, error) {
	id := gl.CreateShader(shaderType)

	cstring, free := gl.Strs(shader)
	defer free()

	gl.ShaderSource(id, 1, cstring, nil)
	gl.CompileShader(id)

	if err := checkShader(id); err != nil {
		gl.DeleteShader(id)
		return 0, err
	}
	return id, nil
}

func NewProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := Compile(vertexShaderSource+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := Compile(fragmentShaderSource+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	return LinkProgram(vertexShader, fragmentShader)
}

func LinkProgram(vertexShader, fragmentShader uint32) (uint32, error) {
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	if err := checkProgram(program); err != nil {
		return 0, err
	}

	return program, nil
}

func VBO(data []float32) uint32 {
	var id uint32
	gl.GenBuffers(1, &id)
	gl.BindBuffer(gl.ARRAY_BUFFER, id)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	return id
}

type GlGetter struct {
	Getter func(uint32, uint32, *int32)
	InfoId uint32
}

func (c GlGetter) Ok(id uint32) bool {
	var s int32
	c.Getter(id, c.InfoId, &s)
	return s == gl.TRUE
}

var (
	StatusProgramLink   = GlGetter{gl.GetProgramiv, gl.LINK_STATUS}
	StatusShaderCompile = GlGetter{gl.GetShaderiv, gl.COMPILE_STATUS}
)

func checkShader(id uint32) error {
	if StatusShaderCompile.Ok(id) {
		return nil
	}

	// TODO factor this into LogGetter if this repeat again
	var length int32
	gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)
	log := strings.Repeat("\x00", int(length+1))
	gl.GetShaderInfoLog(id, length, nil, gl.Str(log))

	return fmt.Errorf("failed to compile shader id %d: %s", id, log)
}

func checkProgram(id uint32) error {
	if StatusProgramLink.Ok(id) {
		return nil
	}

	var length int32
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &length)
	log := strings.Repeat("\x00", int(length+1))
	gl.GetProgramInfoLog(id, length, nil, gl.Str(log))

	return fmt.Errorf("failed to link program: %v", log)
}
