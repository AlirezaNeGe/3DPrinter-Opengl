package camera

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	width       = 800
	height      = 600
	CameraPos   = mgl32.Vec3{0, -4, 3}
	CameraFront = mgl32.Vec3{0, 0, -0.5}
	CameraUp    = mgl32.Vec3{0, 1, 0}
	cameraRight mgl32.Vec3
	yaw                 = float32(-90.0)
	pitch               = float32(0.0)
	lastX               = float32(width) / 2.0
	lastY               = float32(height) / 2.0
	firstRun            = true
	cameraSpeed float32 = 0.2
)

func ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		CameraPos = CameraPos.Add(CameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		CameraPos = CameraPos.Sub(CameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		CameraPos = CameraPos.Sub(CameraFront.Cross(CameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		CameraPos = CameraPos.Add(CameraFront.Cross(CameraUp).Normalize().Mul(cameraSpeed))
	}
}

func CursorPosCallback(window *glfw.Window, xpos float64, ypos float64) {
	if firstRun {
		lastX, lastY = float32(xpos), float32(ypos)
		firstRun = false
	}

	xoffset := float32(xpos) - lastX
	yoffset := lastY - float32(ypos)

	lastX, lastY = float32(xpos), float32(ypos)

	sensitivity := 0.1
	xoffset *= float32(sensitivity)
	yoffset *= float32(sensitivity)

	yaw += xoffset
	pitch += yoffset

	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
	}
	CameraFront = front.Normalize()
}
