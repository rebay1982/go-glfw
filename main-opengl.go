package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"log"
	"fmt"
	"strings"
)

const (
	vertexShaderSource = `
		#version 460
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 460
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(0, 0, 1, 1.0);
		}
	` + "\x00"
)

var (
	square = []float32 {
		-1.0, 1.0, 0.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,

		-1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0,
	}
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func printTime(time float64) {
	log.Printf("Time elapsed %f\n", time)
}

func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(640, 480, "Hellow world.", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Set VSYCH on.
	glfw.SwapInterval(1)
	
	return window
}

func initOpengl() uint32 {
	err := gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version: ", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program
}


func makeVao(vertices []float32) uint32 {
	var vbo uint32	// Vertex Buffer Object
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) << 2, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32	// Vertex Array Object
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

// compileShaer: No idea how this sorcery works, but apparently I don't need to yet....
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}


func main() {
	// Init GLFW: Start /////////////////////////////////////////////////////////////////////////////////////////////////
	window := initGlfw()
	defer glfw.Terminate() // Terminate when we're done.
	// Init GLFW: Done //////////////////////////////////////////////////////////////////////////////////////////////////

	// Init OpenGL: Start ///////////////////////////////////////////////////////////////////////////////////////////////
	program := initOpengl()
	// Init OpenGL: Done ////////////////////////////////////////////////////////////////////////////////////////////////


	vao := makeVao(square)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
		// time := glfw.GetTime()
		// printTime(time)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
