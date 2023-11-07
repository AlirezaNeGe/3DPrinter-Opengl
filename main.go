package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
)

var vertices = []float32{
	-1.0, -1.0,
	1.0, -1.0,
	1.0, 1.0,
	-1.0, 1.0,
}

var indices = []uint32{
	0, 1, 2,
	2, 3, 0,
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, "OpenGL Rectangle", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL:", err)
	}

	// Create Vertex Array Object (VAO)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Create Vertex Buffer Object (VBO) for vertices
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Create Element Buffer Object (EBO) for indices
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// Vertex shader source code
	vertexShaderSource := `
    #version 330 core
    in vec2 position;
    uniform mat4 model;
    uniform mat4 camera;
    uniform mat4 projection;
    void main()
    {
        gl_Position = projection * camera * model * vec4(position, 0.0, 1.0);
    }
    `

	// Fragment shader source code
	fragmentShaderSource := `
    #version 330 core
    out vec4 color;
    void main()
    {
        color = vec4(0.5, 0.5, 0.5, 1.0);
    }
    `

	// Create and compile the shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalln("failed to compile vertex shader:", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatalln("failed to compile fragment shader:", err)
	}

	// Create the shader program
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.UseProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// Specify the layout of the vertex data
	positionAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 2*4, nil)

	// Set the projection, view, and model matrices
	projection := mgl32.Perspective(mgl32.DegToRad(60.0), float32(width)/float32(height), 0.1, 10.0)
	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, -3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model := mgl32.HomogRotate3D(mgl32.DegToRad(45), mgl32.Vec3{1, 0, 0})

	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))

	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 10.0)
	// projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	// gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	//
	// camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	// gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	//
	// model := mgl32.Ident4()
	// modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		return 0, fmt.Errorf("failed to compile %v: %v", source, string(log))
	}
	return shader, nil
}
