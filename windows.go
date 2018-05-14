package main

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func createWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

func setupWindowFrame() *windowFrame {
	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	frame := &windowFrame{
		program:    program,
		projection: mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0),
		camera:     mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}),
		model:      mgl32.Ident4(),

		projectionUniform: gl.GetUniformLocation(program, gl.Str("projection\x00")),
		cameraUniform:     gl.GetUniformLocation(program, gl.Str("camera\x00")),
		modelUniform:      gl.GetUniformLocation(program, gl.Str("model\x00")),
		textureUniform:    gl.GetUniformLocation(program, gl.Str("tex\x00")),

		vertexAttrib: uint32(gl.GetAttribLocation(program, gl.Str("vert\x00"))),
		uvAttrib:     uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00"))),
	}

	gl.UseProgram(program)

	gl.UniformMatrix4fv(frame.projectionUniform, 1, false, &frame.projection[0])
	gl.UniformMatrix4fv(frame.cameraUniform, 1, false, &frame.camera[0])
	gl.UniformMatrix4fv(frame.modelUniform, 1, false, &frame.model[0])
	gl.Uniform1i(frame.textureUniform, 0)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	frame.texture, err = newTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	gl.GenVertexArrays(1, &frame.vao)
	gl.BindVertexArray(frame.vao)

	gl.GenBuffers(1, &frame.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, frame.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(frame.vertexAttrib)
	gl.VertexAttribPointer(frame.vertexAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(frame.uvAttrib)
	gl.VertexAttribPointer(frame.uvAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return frame
}
