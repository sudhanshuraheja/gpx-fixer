package geo

import "math"

const radiusOfEarthInMetres = 6371 * 1000

// Point is a location
type Point struct {
	Latitude  float64
	Longitude float64
	Elevation float64
}

func toRadians(x float64) float64 {
	return ((x / 180) * math.Pi)
}

func square(number float64) float64 {
	return number * number
}

// DistanceInMetresBetweenPoints calculates the distance in metres between two points
// http://www.movable-type.co.uk/scripts/latlong.html
func DistanceInMetresBetweenPoints(point1, point2 Point) float64 {
	latitudeDelta := toRadians(point1.Latitude - point2.Latitude)
	longitudeDelta := toRadians(point1.Longitude - point2.Longitude)
	currentLatitude := toRadians(point1.Latitude)
	newLatitude := toRadians(point2.Latitude)

	a := square(math.Sin(latitudeDelta/2)) + (square(math.Sin(longitudeDelta/2)) * math.Cos(currentLatitude) * math.Cos(newLatitude))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := radiusOfEarthInMetres * c
	return math.Round(distance)
}

// DistanceInMetresIncludingElevation calculates distance and uses elevation
// Uses pythagoras theorem
func DistanceInMetresIncludingElevation(point1, point2 Point) float64 {
	elevationDelta := point1.Elevation - point2.Elevation
	groundDistance := DistanceInMetresBetweenPoints(point1, point2)

	distance := math.Sqrt(square(elevationDelta) + square(groundDistance))
	return math.Round(distance)
}
