package unit

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var unitVAO, unitVBO, unitEBO uint32

type Unit struct {
	ShaderProgram uint32
	Velocity      float32
	Still         bool
	X             float32
	Y             float32
	Z             float32
}

func Init(shaderProgram uint32) {
	// VAO
	gl.GenVertexArrays(1, &unitVAO)
	gl.BindVertexArray(unitVAO)

	// Create Vertex Buffer Object (VBO) for the cube
	gl.GenBuffers(1, &unitVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, unitVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(unitVertices), gl.Ptr(unitVertices), gl.STATIC_DRAW)

	// Create Element Buffer Object (EBO) for the cube
	gl.GenBuffers(1, &unitEBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, unitEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(unitIndices), gl.Ptr(unitIndices), gl.STATIC_DRAW)

	cubepositionAttrib := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
	gl.VertexAttribPointer(cubepositionAttrib, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(cubepositionAttrib)
}

func NewUnit(x, y, z float32, velocity float32, shaderProgram uint32) *Unit {
	return &Unit{
		ShaderProgram: shaderProgram,
		Velocity:      velocity,
		X:             x,
		Y:             y,
		Z:             z,
	}
}

func (u *Unit) Draw(camera mgl32.Mat4) {
	gl.UseProgram(u.ShaderProgram)

	// Translate model to x, y, z location
	model := mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(u.X, u.Y, u.Z))
	modelUniform := gl.GetUniformLocation(u.ShaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	unitCameraUniform := gl.GetUniformLocation(u.ShaderProgram, gl.Str("camera\x00"))

	gl.UniformMatrix4fv(unitCameraUniform, 1, false, &camera[0])
	gl.BindVertexArray(unitVAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(unitIndices)), gl.UNSIGNED_INT, nil)
}

var (
	unitVertices = []float32{
		// Vertices (x, y, z)
		-0.02, -0.02, -0.02,
		0.02, -0.02, -0.02,
		0.02, 0.02, -0.02,
		-0.02, 0.02, -0.02,
		-0.02, -0.02, 0.02,
		0.02, -0.02, 0.02,
		0.02, 0.02, 0.02,
		-0.02, 0.02, 0.02,
	}

	unitIndices = []uint32{
		0, 1, 2,
		2, 3, 0,
		4, 5, 6,
		6, 7, 4,
		0, 4, 7,
		7, 3, 0,
		1, 5, 6,
		6, 2, 1,
		0, 1, 5,
		5, 4, 0,
		2, 3, 7,
		7, 6, 2,
	}
)
