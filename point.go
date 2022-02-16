package orb

import "github.com/paulmach/orb/math"

// A Point is a Lon/Lat 2d point.
type PointOf[T math.Number] [2]T
type Point = PointOf[float64]

var _ Pointer = Point{}

// GeoJSONType returns the GeoJSON type for the object.
func (p PointOf[T]) GeoJSONType() string {
	return "Point"
}

// Dimensions returns 0 because a point is a 0d object.
func (p PointOf[T]) Dimensions() int {
	return 0
}

// Bound returns a single point bound of the point.
func (p PointOf[T]) Bound() BoundOf[T] {
	return BoundOf[T]{p, p}
}

// Point returns itself so it implements the Pointer interface.
func (p PointOf[T]) Point() PointOf[T] {
	return p
}

// Y returns the vertical coordinate of the point.
func (p PointOf[T]) Y() T {
	return p[1]
}

// X returns the horizontal coordinate of the point.
func (p PointOf[T]) X() T {
	return p[0]
}

// Lat returns the vertical, latitude coordinate of the point.
func (p PointOf[T]) Lat() T {
	return p[1]
}

// Lon returns the horizontal, longitude coordinate of the point.
func (p PointOf[T]) Lon() T {
	return p[0]
}

// Equal checks if the point represents the same point or vector.
func (p PointOf[T]) Equal(point PointOf[T]) bool {
	return p[0] == point[0] && p[1] == point[1]
}
