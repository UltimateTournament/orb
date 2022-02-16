package orb

import "github.com/paulmach/orb/math"

// Ring represents a set of ring on the earth.
type RingOf[T math.Number] LineStringOf[T]
type Ring = RingOf[float64]

// GeoJSONType returns the GeoJSON type for the object.
func (r RingOf[T]) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Ring is a 2d object.
func (r RingOf[T]) Dimensions() int {
	return 2
}

// Closed will return true if the ring is a real ring.
// ie. 4+ points and the first and last points match.
// NOTE: this will not check for self-intersection.
func (r RingOf[T]) Closed() bool {
	return (len(r) >= 4) && (r[0] == r[len(r)-1])
}

// Reverse changes the direction of the ring.
// This is done inplace, ie. it modifies the original data.
func (r RingOf[T]) Reverse() {
	LineStringOf[T](r).Reverse()
}

// Bound returns a rect around the ring. Uses rectangular coordinates.
func (r RingOf[T]) Bound() BoundOf[T] {
	return MultiPointOf[T](r).Bound()
}

// Orientation returns 1 if the the ring is in couter-clockwise order,
// return -1 if the ring is the clockwise order and 0 if the ring is
// degenerate and had no area.
func (r RingOf[T]) Orientation() Orientation {
	var area T = 0

	// This is a fast planar area computation, which is okay for this use.
	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		area += (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
	}

	if area > 0 {
		return CCW
	}

	if area < 0 {
		return CW
	}

	// degenerate case, no area
	return 0
}

// Equal compares two rings. Returns true if lengths are the same
// and all points are Equal.
func (r RingOf[T]) Equal(ring RingOf[T]) bool {
	return MultiPointOf[T](r).Equal(MultiPointOf[T](ring))
}

// Clone returns a new copy of the ring.
func (r RingOf[T]) Clone() RingOf[T] {
	if r == nil {
		return nil
	}

	ps := MultiPointOf[T](r)
	return RingOf[T](ps.Clone())
}
