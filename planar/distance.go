package planar

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/math"
)

// Distance returns the distance between two points in 2d euclidean geometry.
func Distance[T math.Number](p1, p2 orb.PointOf[T]) T {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

// DistanceSquared returns the square of the distance between two points in 2d euclidean geometry.
func DistanceSquared[T math.Number](p1, p2 orb.PointOf[T]) T {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return d0*d0 + d1*d1
}
