package nasa

import (
	"math"
)

type World struct {
	Radius float64
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

/*
Calculates the distance between two locations
using the Spherical Law of Cosines.
*/
func (w World) distanceBetweenTwoLocations(point1, point2 Location) float64 {
	s1, c1 := math.Sincos(degreesToRadians(point1.Latitude))
	s2, c2 := math.Sincos(degreesToRadians(point2.Latitude))

	clong := math.Cos(degreesToRadians(point1.Longitude - point2.Longitude))
	return w.Radius * math.Acos(s1*s2+c1*c2*clong)
}

/*
Calculates the distance between two locations passed in
as string values.
*/
func (w World) DistanceBetweenTwoLocationStrings(from, to string) (float64, error) {
	_from, err := parseLocation(from)
	if err != nil {
		return 0, err
	}

	_to, err := parseLocation(to)
	if err != nil {
		return 0, err
	}

	distance := w.distanceBetweenTwoLocations(*_from, *_to)
	return distance, nil
}
