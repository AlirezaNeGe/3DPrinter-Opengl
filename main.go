package main

import (
	"bufio"
	"fmt"
	"gogl/camera"
	"gogl/collision"
	"gogl/head"
	"gogl/interpreter"
	"gogl/scene"
	"gogl/shaders"
	"gogl/unit"
	"gogl/utils"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
	fps    = 500
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	// Extract Gcode values
	gcodeFilePath := "./gcode.example"
	file, err := os.Open(gcodeFilePath) // Replace with your file name
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := interpreter.ExtractValues(scanner)

	if err := glfw.Init(); err != nil {
		fmt.Println("failed to initialize glfw:", err)
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, "OpenGL Separated unit and Rectangle", nil, nil)
	if err != nil {
		fmt.Println("failed to create window:", err)
		return
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("failed to initialize OpenGL:", err)
		return
	}

	gl.Viewport(0, 0, width, height)

	// Create and compile the shaders
	vertexShader, err := shaders.CompileShader(shaders.VertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("failed to compile vertex shader:", err)
		return
	}
	fragmentShader, err := shaders.CompileShader(shaders.FragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("failed to compile fragment shader:", err)
		return
	}

	// unit
	unitVertexShader, err := shaders.CompileShader(shaders.UnitVertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("failed to compile vertex shader:", err)
		return
	}
	unitFragmentShader, err := shaders.CompileShader(shaders.UnitFragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("failed to compile fragment shader:", err)
		return
	}

	// unit
	HeadVertexShader, err := shaders.CompileShader(shaders.HeadVertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		fmt.Println("failed to compile vertex shader:", err)
		return
	}
	HeadFragmentShader, err := shaders.CompileShader(shaders.HeadFragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("failed to compile fragment shader:", err)
		return
	}

	// Create the shader program scene
	shaderProgram := shaders.NewProgram(vertexShader, fragmentShader)
	scene.Init(shaderProgram)

	// Create the shader program units
	unitShaderProgram := shaders.NewProgram(unitVertexShader, unitFragmentShader)
	unit.Init(unitShaderProgram)

	// Create the shader program head
	headShaderProgram := shaders.NewProgram(HeadVertexShader, HeadFragmentShader)
	head.Init(headShaderProgram)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)

	gl.UseProgram(shaderProgram)
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// unit coordinates
	gl.UseProgram(unitShaderProgram)
	unitProjectionUniform := gl.GetUniformLocation(unitShaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(unitProjectionUniform, 1, false, &projection[0])

	// head coordinates
	gl.UseProgram(headShaderProgram)
	headUniform := gl.GetUniformLocation(headShaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(headUniform, 1, false, &model[0])

	headProjectionUniform := gl.GetUniformLocation(headShaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(headProjectionUniform, 1, false, &projection[0])

	var acceleration float32 = -9.8
	previousTime := glfw.GetTime()

	// units := utils.MakeUnits(&values, unitShaderProgram)
	path := utils.MakePath(&values)

	var pathToDraw []*utils.Path
	pathToDraw = append(pathToDraw, path[0])

	var unitsToDraw []*unit.Unit

	pl := 0
	lenPath := len(path)

	gl.Enable(gl.DEPTH_TEST)
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	window.SetCursorPosCallback(camera.CursorPosCallback)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		camera.ProcessInput(window)
		// Update
		t := glfw.GetTime()
		elapsed := t - previousTime
		previousTime = t

		camera := mgl32.LookAtV(camera.CameraPos, camera.CameraFront.Add(camera.CameraFront), camera.CameraUp)

		// Render the scene
		scene.DrawScene(shaderProgram, camera)

		// Previous units
		for i, u := range unitsToDraw {
			// Check for collision
			for j, uB := range unitsToDraw {
				if i != j {
					if collision.UnitsCollide(u, uB) {
						if u.Z < uB.Z {
							uB.Still = true
						} else {
							u.Still = true
						}
					}

				}
			}
			if u.Z > 0 && u.Still == false {
				u.Velocity += acceleration * float32(elapsed)
				u.Z += u.Velocity * float32(elapsed)
				if u.Z <= 0 {
					u.Z = 0
				}
			}
			u.Draw(camera)
		}

		head.DrawHead(path[pl].X, path[pl].Y, path[pl].Z, headShaderProgram, camera)
		if pl < lenPath-1 {
			if path[pl].Feed == true {
				u := unit.NewUnit(path[pl].X, path[pl].Y, path[pl].Z, 0, unitShaderProgram)
				u.Draw(camera)
				unitsToDraw = append(unitsToDraw, u)
			}
			pl++
		}

		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Second / time.Duration(fps))
	}
}
