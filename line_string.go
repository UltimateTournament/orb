package orb

import "github.com/paulmach/orb/math"

// LineString represents a set of points to be thought of as a polyline.
type LineStringOf[T math.Number] []PointOf[T]
type LineString = LineStringOf[float64]

// GeoJSONType returns the GeoJSON type for the object.
func (ls LineStringOf[T]) GeoJSONType() string {
	return "LineString"
}

// Dimensions returns 1 because a LineString is a 1d object.
func (ls LineStringOf[T]) Dimensions() int {
	return 1
}

// Reverse will reverse the line string.
// This is done inplace, ie. it modifies the original data.
func (ls LineStringOf[T]) Reverse() {
	l := len(ls) - 1
	for i := 0; i <= l/2; i++ {
		ls[i], ls[l-i] = ls[l-i], ls[i]
	}
}

// Bound returns a rect around the line string. Uses rectangular coordinates.
func (ls LineStringOf[T]) Bound() BoundOf[T] {
	return MultiPointOf[T](ls).Bound()
}

// Equal compares two line strings. Returns true if lengths are the same
// and all points are Equal.
func (ls LineStringOf[T]) Equal(lineString LineStringOf[T]) bool {
	return MultiPointOf[T](ls).Equal(MultiPointOf[T](lineString))
}

// Clone returns a new copy of the line string.
func (ls LineStringOf[T]) Clone() LineStringOf[T] {
	ps := MultiPointOf[T](ls)
	return LineStringOf[T](ps.Clone())
}
