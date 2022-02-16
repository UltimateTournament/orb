package orb

import "github.com/paulmach/orb/math"

var emptyBound = Bound{Min: Point{1, 1}, Max: Point{-1, -1}}

func emptyBoundOf[T math.Number]() BoundOf[T] {
	return BoundOf[T]{Min: PointOf[T]{1, 1}, Max: PointOf[T]{0, 0}}
}

// A Bound represents a closed box or rectangle.
// To create a bound with two points you can do something like:
//	orb.MultiPoint{p1, p2}.Bound()
type BoundOf[T math.Number] struct {
	Min, Max PointOf[T]
}

// A Bound represents a closed box or rectangle.
// To create a bound with two points you can do something like:
//	orb.MultiPoint{p1, p2}.Bound()
type Bound = BoundOf[float64]

// GeoJSONType returns the GeoJSON type for the object.
func (b BoundOf[T]) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Bound is a 2d object.
func (b BoundOf[T]) Dimensions() int {
	return 2
}

// ToPolygon converts the bound into a Polygon object.
func (b BoundOf[T]) ToPolygon() PolygonOf[T] {
	return PolygonOf[T]{b.ToRing()}
}

// ToRing converts the bound into a loop defined
// by the boundary of the box.
func (b BoundOf[T]) ToRing() RingOf[T] {
	return RingOf[T]{
		b.Min,
		PointOf[T]{b.Max[0], b.Min[1]},
		b.Max,
		PointOf[T]{b.Min[0], b.Max[1]},
		b.Min,
	}
}

// Extend grows the bound to include the new point.
func (b BoundOf[T]) Extend(point PointOf[T]) BoundOf[T] {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	return BoundOf[T]{
		Min: PointOf[T]{
			math.Min(b.Min[0], point[0]),
			math.Min(b.Min[1], point[1]),
		},
		Max: PointOf[T]{
			math.Max(b.Max[0], point[0]),
			math.Max(b.Max[1], point[1]),
		},
	}
}

// Union extends this bound to contain the union of this and the given bound.
func (b BoundOf[T]) Union(other BoundOf[T]) BoundOf[T] {
	if other.IsEmpty() {
		return b
	}

	b = b.Extend(other.Min)
	b = b.Extend(other.Max)
	b = b.Extend(other.LeftTop())
	b = b.Extend(other.RightBottom())

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b BoundOf[T]) Contains(point PointOf[T]) bool {
	if point[1] < b.Min[1] || b.Max[1] < point[1] {
		return false
	}

	if point[0] < b.Min[0] || b.Max[0] < point[0] {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b BoundOf[T]) Intersects(bound BoundOf[T]) bool {
	if (b.Max[0] < bound.Min[0]) ||
		(b.Min[0] > bound.Max[0]) ||
		(b.Max[1] < bound.Min[1]) ||
		(b.Min[1] > bound.Max[1]) {
		return false
	}

	return true
}

// Pad extends the bound in all directions by the given value.
func (b BoundOf[T]) Pad(d T) BoundOf[T] {
	b.Min[0] -= d
	b.Min[1] -= d

	b.Max[0] += d
	b.Max[1] += d

	return b
}

// Center returns the center of the bounds by "averaging" the x and y coords.
func (b BoundOf[T]) Center() PointOf[T] {
	return PointOf[T]{
		(b.Min[0] + b.Max[0]) / 2.0,
		(b.Min[1] + b.Max[1]) / 2.0,
	}
}

// Top returns the top of the bound.
func (b BoundOf[T]) Top() T {
	return b.Max[1]
}

// Bottom returns the bottom of the bound.
func (b BoundOf[T]) Bottom() T {
	return b.Min[1]
}

// Right returns the right of the bound.
func (b BoundOf[T]) Right() T {
	return b.Max[0]
}

// Left returns the left of the bound.
func (b BoundOf[T]) Left() T {
	return b.Min[0]
}

// LeftTop returns the upper left point of the bound.
func (b BoundOf[T]) LeftTop() PointOf[T] {
	return PointOf[T]{b.Left(), b.Top()}
}

// RightBottom return the lower right point of the bound.
func (b BoundOf[T]) RightBottom() PointOf[T] {
	return PointOf[T]{b.Right(), b.Bottom()}
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b BoundOf[T]) IsEmpty() bool {
	return b.Min[0] > b.Max[0] || b.Min[1] > b.Max[1]
}

// IsZero return true if the bound just includes just null island.
func (b BoundOf[T]) IsZero() bool {
	return b.Max == PointOf[T]{} && b.Min == PointOf[T]{}
}

// Bound returns the the same bound.
func (b BoundOf[T]) Bound() BoundOf[T] {
	return b
}

// Equal returns if two bounds are equal.
func (b BoundOf[T]) Equal(c BoundOf[T]) bool {
	return b.Min == c.Min && b.Max == c.Max
}
