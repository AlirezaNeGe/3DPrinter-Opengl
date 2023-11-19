package collision

import (
	"gogl/unit"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// Minimum distance for collision is unit size
	UnitSize = 0.02
	// Angle between distance vector and z-axis
	theta = 40
)

func UnitsCollide(unitA, unitB *unit.Unit) bool {

	// Calculate the distance between unit centers
	distanceVector := mgl32.Vec3{unitA.X, unitA.Y, unitA.Z}.Sub(mgl32.Vec3{unitB.X, unitB.Y, unitB.Z})
	distance := distanceVector.Len()
	if distance < UnitSize {
		// Define the z-axis direction
		zAxis := mgl32.Vec3{0, 0, 1}

		// Calculate the angle between the distance vector and the z-axis
		angleRadians := math.Acos(float64(distanceVector.Normalize().Dot(zAxis.Normalize())))
		angle := mgl32.RadToDeg(float32(angleRadians))
		return angle < theta
	}

	return false
}
