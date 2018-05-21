package geo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToRadians(t *testing.T) {
	assert.Equal(t, 2*math.Pi, toRadians(360.0)) // 2pi radians
	assert.Equal(t, 0.0, toRadians(0))
}
func Test_Square(t *testing.T) {
	assert.Equal(t, 4.0, square(2.0))
}

func Test_Distance(t *testing.T) {
	assert.Equal(t, 0.0, DistanceInMetresBetweenPoints(Point{0.0, 0.0, 0.0}, Point{360.0, 360.0, 0.0}))
	assert.Equal(t, 0.0, DistanceInMetresBetweenPoints(Point{360.0, 360.0, 0.0}, Point{0.0, 0.0, 0.0}))

	// Eden Gardens Diameter
	assert.Equal(t, 173.7846, DistanceInMetresIncludingElevation(Point{22.565129, 88.342979, 0.0}, Point{22.564000, 88.343628, 100.0}))
}

func Test_Distance3D(t *testing.T) {
	assert.Equal(t, 100.0, DistanceInMetresIncludingElevation(Point{0.0, 0.0, 0.0}, Point{0.0, 0.0, 100.0}))
	assert.Equal(t, 494.9857, DistanceInMetresIncludingElevation(Point{28.663784, 77.086628, 100.0}, Point{28.663153, 77.091650, 100.0}))
}
