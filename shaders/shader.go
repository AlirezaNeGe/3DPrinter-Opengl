package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Vertex shader source code
var VertexShaderSource = `
    #version 330 core
    in vec3 position;
    uniform mat4 model;
    uniform mat4 projection;
    uniform mat4 camera;

    void main()
    {
        gl_Position = projection * camera * model * vec4(position, 1.0);
    }
    `

// Fragment shader source code
var FragmentShaderSource = `
    #version 330 core
    out vec4 color;
    void main()
    {
        color = vec4(0.2, 0.2, 0.2, 1.0);
    }
    `

// Create a separate shader program for the unit
var UnitVertexShaderSource = `
    #version 330 core
    in vec3 position;
    uniform mat4 model;
    uniform mat4 camera;
    uniform mat4 projection;
    void main()
    {
        gl_Position = projection * camera * model * vec4(position, 1.0);
    }
    `

var UnitFragmentShaderSource = `
    #version 330 core
    out vec4 color;
    void main()
    {
        color = vec4(0.0, 1.0, 0.0, 1.0);
    }
    `

var HeadVertexShaderSource = `
    #version 330 core
    in vec3 position;
    uniform mat4 model;
    uniform mat4 camera;
    uniform mat4 projection;
    void main()
    {
        gl_Position = projection * camera * model * vec4(position, 1.0);
    }
    `

var HeadFragmentShaderSource = `
    #version 330 core
    out vec4 color;
    void main()
    {
        color = vec4(0.0, 0.0, 1.0, 1.0);
    }
    `

func CompileShader(source string, shaderType uint32) (uint32, error) {
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

func NewProgram(vertexShader, fragmentShader uint32) uint32 {
	// Create the shader program for the rectangle
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram
}
