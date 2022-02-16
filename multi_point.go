package orb

import "github.com/paulmach/orb/math"

// A MultiPoint represents a set of points in the 2D Eucledian or Cartesian plane.
type MultiPointOf[T math.Number] []PointOf[T]
type MultiPoint = MultiPointOf[float64]

// GeoJSONType returns the GeoJSON type for the object.
func (mp MultiPointOf[T]) GeoJSONType() string {
	return "MultiPoint"
}

// Dimensions returns 0 because a MultiPoint is a 0d object.
func (mp MultiPointOf[T]) Dimensions() int {
	return 0
}

// Clone returns a new copy of the points.
func (mp MultiPointOf[T]) Clone() MultiPointOf[T] {
	if mp == nil {
		return nil
	}

	points := make([]PointOf[T], len(mp))
	copy(points, mp)

	return MultiPointOf[T](points)
}

// Bound returns a bound around the points. Uses rectangular coordinates.
func (mp MultiPointOf[T]) Bound() BoundOf[T] {
	if len(mp) == 0 {
		return emptyBoundOf[T]()
	}

	b := BoundOf[T]{mp[0], mp[0]}
	for _, p := range mp {
		b = b.Extend(p)
	}

	return b
}

// Equal compares two MultiPoint objects. Returns true if lengths are the same
// and all points are Equal, and in the same order.
func (mp MultiPointOf[T]) Equal(multiPoint MultiPointOf[T]) bool {
	if len(mp) != len(multiPoint) {
		return false
	}

	for i := range mp {
		if !mp[i].Equal(multiPoint[i]) {
			return false
		}
	}

	return true
}
