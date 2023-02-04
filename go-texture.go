package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"log"
	"image"
	"image/color"
)

const (
	WINDOW_TITLE = "GLFW Template"
	WINDOW_WIDTH = 1980 
	WINDOW_HEIGHT = 1080
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// initGlfw: Initialize GLFW and return a window.
func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, WINDOW_TITLE, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)	// Enable vsync
	
	return window
}

// initOpengl: Function to initialize the OpenGL component.
func initOpengl() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version: ", version)

}

func initTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.BindImageTexture(0, texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)

	return texture
}

func initFramebuffer(texture uint32) uint32 {
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return framebuffer
}


func Draw() []uint8 {

	w := WINDOW_WIDTH
	h := WINDOW_HEIGHT

	var img = image.NewRGBA(image.Rect(0, 0, w/2, h/2))
	var red = color.RGBA{255, 0, 0, 0}
	var blue = color.RGBA{0, 0, 255, 0}
	for i := 0; i < h/2; i++ {
		for j := 0; j < w/2; j++ {

			if (j % 2 > 0) {
				img.Set(j, i, red)

			} else {
				img.Set(j, i, blue)

			}
		}
	}

	return img.Pix
}

func main() {
	window := initGlfw()
	defer glfw.Terminate() // Terminate when we're done.

	initOpengl()

	texture := initTexture()
	initFramebuffer(texture)

	var w, h = window.GetSize()

	for !window.ShouldClose() {

		imgData := Draw()

		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(w/2), int32(h/2), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgData))
		gl.BlitFramebuffer(0, 0, int32(w/2), int32(h/2), 0, 0, int32(w/2), int32(h/2), gl.COLOR_BUFFER_BIT, gl.LINEAR)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
