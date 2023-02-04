package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"log"
	"image"
	"image/color"
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

func initOpengl() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version: ", version)

}



func main() {
	// Init GLFW: Start /////////////////////////////////////////////////////////////////////////////////////////////////
	window := initGlfw()
	defer glfw.Terminate() // Terminate when we're done.
	// Init GLFW: Done //////////////////////////////////////////////////////////////////////////////////////////////////

	// Init OpenGL: Start ///////////////////////////////////////////////////////////////////////////////////////////////
	initOpengl()
	// Init OpenGL: Done ////////////////////////////////////////////////////////////////////////////////////////////////

	var texture uint32
	{
		gl.GenTextures(1, &texture)

		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		gl.BindImageTexture(0, texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
	}

	var framebuffer uint32
	{
			gl.GenFramebuffers(1, &framebuffer)
			gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)

			gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
			gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	}
	
	var w, h = window.GetSize()
	var img = image.NewRGBA(image.Rect(0, 0, w, h))
	var red = color.RGBA{255, 0, 0, 0}

	for i := 0; i < h/2; i++ {
		for j := 0; j < w; j++ {
			img.Set(j, i, red)
		}
	}

	for !window.ShouldClose() {


		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

		gl.BlitFramebuffer(0, 0, int32(w), int32(h), 0, 0, int32(w), int32(h), gl.COLOR_BUFFER_BIT, gl.LINEAR)
		// time := glfw.GetTime()
		// printTime(time)
		
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
