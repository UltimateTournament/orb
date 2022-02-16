package orb

import "github.com/paulmach/orb/math"

// Polygon is a closed area. The first LineString is the outer ring.
// The others are the holes. Each LineString is expected to be closed
// ie. the first point matches the last.
type PolygonOf[T math.Number] []RingOf[T]
type Polygon = PolygonOf[float64]

// GeoJSONType returns the GeoJSON type for the object.
func (p PolygonOf[T]) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Polygon is a 2d object.
func (p PolygonOf[T]) Dimensions() int {
	return 2
}

// Bound returns a bound around the polygon.
func (p PolygonOf[T]) Bound() BoundOf[T] {
	if len(p) == 0 {
		return emptyBoundOf[T]()
	}
	return p[0].Bound()
}

// Equal compares two polygons. Returns true if lengths are the same
// and all points are Equal.
func (p PolygonOf[T]) Equal(polygon PolygonOf[T]) bool {
	if len(p) != len(polygon) {
		return false
	}

	for i := range p {
		if !p[i].Equal(polygon[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the polygon.
// All of the rings are also cloned.
func (p PolygonOf[T]) Clone() PolygonOf[T] {
	if p == nil {
		return p
	}

	np := make(PolygonOf[T], 0, len(p))
	for _, r := range p {
		np = append(np, r.Clone())
	}

	return np
}
