package head

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var VAO, VBO, EBO uint32

func Init(shaderProgram uint32) {

	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	// Create Vertex Buffer Object (VBO) for the rectangle
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(pyramidVertices), gl.Ptr(pyramidVertices), gl.STATIC_DRAW)

	// Create Element Buffer Object (EBO) for the rectangle
	gl.GenBuffers(1, &EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(pyramidIndices), gl.Ptr(pyramidIndices), gl.STATIC_DRAW)

	positionAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 3*4, nil)
}

func DrawHead(x, y, z float32, shaderProgram uint32, camera mgl32.Mat4) {
	gl.UseProgram(shaderProgram)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(x, y, z))
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	// Translate model to x, y, z location
	gl.BindVertexArray(VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(pyramidIndices)), gl.UNSIGNED_INT, nil)
}

var (
	pyramidVertices = []float32{
		0.0, 0.0, 0.0,
		0.1, 0.1, 0.1,
		0.1, -0.1, 0.1,
		-0.1, 0.1, 0.1,
		-0.1, -0.1, 0.1,
	}

	// Indices for vertices order
	pyramidIndices = []uint32{

		0, 1, 2,
		0, 2, 4,
		0, 3, 4,
		0, 1, 3,
		1, 2, 4,
		4, 1, 3,
	}
)
