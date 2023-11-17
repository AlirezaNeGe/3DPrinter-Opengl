package utils

import (
	"gogl/collision"
	"gogl/interpreter"
	"gogl/unit"
	"math"
)

var unitSize = collision.UnitSize

func CalculateUnitNumber(start, end interpreter.GCode) int {
	dx := end.X - start.X
	dy := end.Y - start.Y
	dz := end.Z - start.Z

	distance := math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
	unitNumber := int(distance / unitSize)

	return unitNumber
}

func MakeUnits(values *[]interpreter.GCode, cubeShaderProgram uint32) []*unit.Unit {
	var units []*unit.Unit
	lastGcode := interpreter.GCode{
		G: 1,
		X: 0,
		Y: 0,
		Z: 0,
	}

	for _, gcode := range *values {
		steps := CalculateUnitNumber(lastGcode, gcode)
		if gcode.G == 1 {
			xStep := (gcode.X - lastGcode.X) / float32(steps)
			yStep := (gcode.Y - lastGcode.Y) / float32(steps)
			zStep := (gcode.Z - lastGcode.Z) / float32(steps)
			for i := 0; i < steps; i++ {
				units = append(units, unit.NewUnit(
					lastGcode.X+xStep*float32(i),
					lastGcode.Y+yStep*float32(i),
					lastGcode.Z+zStep*float32(i), 0.1, cubeShaderProgram))
			}
			units = append(units, unit.NewUnit(gcode.X, gcode.Y, gcode.Z, 0.1, cubeShaderProgram))
		}
		lastGcode = gcode
	}
	return units
}
