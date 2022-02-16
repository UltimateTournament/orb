// Package quadtree implements a quadtree using rectangular partitions.
// Each point exists in a unique node in the tree or as leaf nodes.
// This implementation is based off of the d3 implementation:
// https://github.com/mbostock/d3/wiki/Quadtree-Geom
package quadtree

import (
	"errors"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/math"
	"github.com/paulmach/orb/planar"
)

var (
	// ErrPointOutsideOfBounds is returned when trying to add a point
	// to a quadtree and the point is outside the bounds used to create the tree.
	ErrPointOutsideOfBounds = errors.New("quadtree: point outside of bounds")
)

type Quadtree = QuadtreeOf[float64]

// Quadtree implements a two-dimensional recursive spatial subdivision
// of orb.Pointers. This implementation uses rectangular partitions.
type QuadtreeOf[T math.Number] struct {
	bound orb.BoundOf[T]
	root  *node[T]
}

// A FilterFunc is a function that filters the points to search for.
type FilterFunc[T math.Number] func(p orb.PointerOf[T]) bool

// node represents a node of the quad tree. Each node stores a Value
// and has links to its 4 children
type node[T math.Number] struct {
	Value    orb.PointerOf[T]
	Children [4]*node[T]
}

// New creates a new quadtree for the given bound. Added points
// must be within this bound.
func New[T math.Number](bound orb.BoundOf[T]) *QuadtreeOf[T] {
	return &QuadtreeOf[T]{bound: bound}
}

// Bound returns the bounds used for the quad tree.
func (q *QuadtreeOf[T]) Bound() orb.BoundOf[T] {
	return q.bound
}

// Add puts an object into the quad tree, must be within the quadtree bounds.
// This function is not thread-safe, ie. multiple goroutines cannot insert into
// a single quadtree.
func (q *QuadtreeOf[T]) Add(p orb.PointerOf[T]) error {
	if p == nil {
		return nil
	}

	point := p.Point()
	if !q.bound.Contains(point) {
		return ErrPointOutsideOfBounds
	}

	if q.root == nil {
		q.root = &node[T]{
			Value: p,
		}
		return nil
	}

	q.add(q.root, p, p.Point(),
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min[0], q.bound.Max[0],
		q.bound.Min[1], q.bound.Max[1],
	)

	return nil
}

// add is the recursive search to find a place to add the point
func (q *QuadtreeOf[T]) add(n *node[T], p orb.PointerOf[T], point orb.PointOf[T], left, right, bottom, top T) {
	i := 0

	// figure which child of this internal node the point is in.
	if cy := (bottom + top) / 2.0; point[1] <= cy {
		top = cy
		i = 2
	} else {
		bottom = cy
	}

	if cx := (left + right) / 2.0; point[0] >= cx {
		left = cx
		i++
	} else {
		right = cx
	}

	if n.Children[i] == nil {
		n.Children[i] = &node[T]{Value: p}
		return
	}

	// proceed down to the child to see if it's a leaf yet and we can add the pointer there.
	q.add(n.Children[i], p, point, left, right, bottom, top)
}

// Remove will remove the pointer from the quadtree. By default it'll match
// using the points, but a FilterFunc can be provided for a more specific test
// if there are elements with the same point value in the tree. For example:
//	func(pointer orb.Pointer) {
//		return pointer.(*MyType).ID == lookingFor.ID
//	}
func (q *QuadtreeOf[T]) Remove(p orb.PointerOf[T], eq FilterFunc[T]) bool {
	if eq == nil {
		point := p.Point()
		eq = func(pointer orb.PointerOf[T]) bool {
			return point.Equal(pointer.Point())
		}
	}

	b := q.bound
	v := &findVisitor[T]{
		point:          p.Point(),
		filter:         eq,
		closestBound:   &b,
		minDistSquared: math.MaxOf[T](),
	}

	newVisit[T](v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min[0], q.bound.Max[0],
		q.bound.Min[1], q.bound.Max[1],
	)

	if v.closest == nil {
		return false
	}

	removeNode(v.closest)
	return true
}

// removeNode is the recursive fixing up of the tree when we remove a node.
func removeNode[T math.Number](n *node[T]) {
	var i int
	for {
		i = -1
		if n.Children[0] != nil {
			i = 0
		} else if n.Children[1] != nil {
			i = 1
		} else if n.Children[2] != nil {
			i = 2
		} else if n.Children[3] != nil {
			i = 3
		}

		if i == -1 {
			n.Value = nil
			return
		}

		if n.Children[i].Value == nil {
			n.Children[i] = nil
			continue
		}

		break
	}

	n.Value = n.Children[i].Value
	removeNode(n.Children[i])
}

// Find returns the closest Value/Pointer in the quadtree.
// This function is thread safe. Multiple goroutines can read from
// a pre-created tree.
func (q *QuadtreeOf[T]) Find(p orb.PointOf[T]) orb.PointerOf[T] {
	return q.Matching(p, nil)
}

// Matching returns the closest Value/Pointer in the quadtree for which
// the given filter function returns true. This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *QuadtreeOf[T]) Matching(p orb.PointOf[T], f FilterFunc[T]) orb.PointerOf[T] {
	if q.root == nil {
		return nil
	}

	b := q.bound
	v := &findVisitor[T]{
		point:          p,
		filter:         f,
		closestBound:   &b,
		minDistSquared: math.MaxOf[T](),
	}

	newVisit[T](v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min[0], q.bound.Max[0],
		q.bound.Min[1], q.bound.Max[1],
	)

	if v.closest == nil {
		return nil
	}
	return v.closest.Value
}

// KNearest returns k closest Value/Pointer in the quadtree.
// This function is thread safe. Multiple goroutines can read from a pre-created tree.
// An optional buffer parameter is provided to allow for the reuse of result slice memory.
// The points are returned in a sorted order, nearest first.
// This function allows defining a maximum distance in order to reduce search iterations.
func (q *QuadtreeOf[T]) KNearest(buf []orb.PointerOf[T], p orb.PointOf[T], k int, maxDistance ...T) []orb.PointerOf[T] {
	return q.KNearestMatching(buf, p, k, nil, maxDistance...)
}

// KNearestMatching returns k closest Value/Pointer in the quadtree for which
// the given filter function returns true. This function is thread safe.
// Multiple goroutines can read from a pre-created tree. An optional buffer
// parameter is provided to allow for the reuse of result slice memory.
// The points are returned in a sorted order, nearest first.
// This function allows defining a maximum distance in order to reduce search iterations.
func (q *QuadtreeOf[T]) KNearestMatching(buf []orb.PointerOf[T], p orb.PointOf[T], k int, f FilterFunc[T], maxDistance ...T) []orb.PointerOf[T] {
	if q.root == nil {
		return nil
	}

	b := q.bound
	v := &nearestVisitor[T]{
		point:          p,
		filter:         f,
		k:              k,
		maxHeap:        make(maxHeap[T], 0, k+1),
		closestBound:   &b,
		maxDistSquared: math.MaxOf[T](),
	}

	if len(maxDistance) > 0 {
		v.maxDistSquared = maxDistance[0] * maxDistance[0]
	}

	newVisit[T](v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min[0], q.bound.Max[0],
		q.bound.Min[1], q.bound.Max[1],
	)

	//repack result
	if cap(buf) < len(v.maxHeap) {
		buf = make([]orb.PointerOf[T], len(v.maxHeap))
	} else {
		buf = buf[:len(v.maxHeap)]
	}

	for i := len(v.maxHeap) - 1; i >= 0; i-- {
		buf[i] = v.maxHeap.Pop().point
	}

	return buf
}

// InBound returns a slice with all the pointers in the quadtree that are
// within the given bound. An optional buffer parameter is provided to allow
// for the reuse of result slice memory. This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *QuadtreeOf[T]) InBound(buf []orb.PointerOf[T], b orb.BoundOf[T]) []orb.PointerOf[T] {
	return q.InBoundMatching(buf, b, nil)
}

// InBoundMatching returns a slice with all the pointers in the quadtree that are
// within the given bound and matching the give filter function. An optional buffer
// parameter is provided to allow for the reuse of result slice memory. This function
// is thread safe.  Multiple goroutines can read from a pre-created tree.
func (q *QuadtreeOf[T]) InBoundMatching(buf []orb.PointerOf[T], b orb.BoundOf[T], f FilterFunc[T]) []orb.PointerOf[T] {
	if q.root == nil {
		return nil
	}

	var p []orb.PointerOf[T]
	if len(buf) > 0 {
		p = buf[:0]
	}
	v := &inBoundVisitor[T]{
		bound:    &b,
		pointers: p,
		filter:   f,
	}

	newVisit[T](v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min[0], q.bound.Max[0],
		q.bound.Min[1], q.bound.Max[1],
	)

	return v.pointers
}

// The visit stuff is a more go like (hopefully) implementation of the
// d3.quadtree.visit function. It is not exported, but if there is a
// good use case, it could be.

type visitor[T math.Number] interface {
	// Bound returns the current relevant bound so we can prune irrelevant nodes
	// from the search. Using a pointer was benchmarked to be 5% faster than
	// having to copy the bound on return. go1.9
	Bound() *orb.BoundOf[T]
	Visit(n *node[T])

	// Point should return the specific point being search for, or null if there
	// isn't one (ie. searching by bound). This helps guide the search to the
	// best child node first.
	Point() orb.PointOf[T]
}

// visit provides a framework for walking the quad tree.
// Currently used by the `Find` and `InBound` functions.
type visit[T math.Number] struct {
	visitor visitor[T]
}

func newVisit[T math.Number](v visitor[T]) *visit[T] {
	return &visit[T]{
		visitor: v,
	}
}

func (v *visit[T]) Visit(n *node[T], left, right, bottom, top T) {
	b := v.visitor.Bound()
	// if left > b.Right() || right < b.Left() ||
	// 	bottom > b.Top() || top < b.Bottom() {
	// 	return
	// }
	if left > b.Max[0] || right < b.Min[0] ||
		bottom > b.Max[1] || top < b.Min[1] {
		return
	}

	if n.Value != nil {
		v.visitor.Visit(n)
	}

	if n.Children[0] == nil && n.Children[1] == nil &&
		n.Children[2] == nil && n.Children[3] == nil {
		// no children check
		return
	}

	cx := (left + right) / 2.0
	cy := (bottom + top) / 2.0

	i := childIndex(cx, cy, v.visitor.Point())
	for j := i; j < i+4; j++ {
		if n.Children[j%4] == nil {
			continue
		}

		if k := j % 4; k == 0 {
			v.Visit(n.Children[0], left, cx, cy, top)
		} else if k == 1 {
			v.Visit(n.Children[1], cx, right, cy, top)
		} else if k == 2 {
			v.Visit(n.Children[2], left, cx, bottom, cy)
		} else if k == 3 {
			v.Visit(n.Children[3], cx, right, bottom, cy)
		}
	}
}

type findVisitor[T math.Number] struct {
	point          orb.PointOf[T]
	filter         FilterFunc[T]
	closest        *node[T]
	closestBound   *orb.BoundOf[T]
	minDistSquared T
}

func (v *findVisitor[T]) Bound() *orb.BoundOf[T] {
	return v.closestBound
}

func (v *findVisitor[T]) Point() orb.PointOf[T] {
	return v.point
}

func (v *findVisitor[T]) Visit(n *node[T]) {
	// skip this pointer if we have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Point()
	if d := planar.DistanceSquared(point, v.point); d < v.minDistSquared {
		v.minDistSquared = d
		v.closest = n

		d = math.Sqrt(d)
		v.closestBound.Min[0] = v.point[0] - d
		v.closestBound.Max[0] = v.point[0] + d
		v.closestBound.Min[1] = v.point[1] - d
		v.closestBound.Max[1] = v.point[1] + d
	}
}

// type pointsQueueItem struct {
// 	point    orb.Pointer
// 	distance float64 // distance to point and priority inside the queue
// 	index    int     // point index in queue
// }

// type pointsQueue []pointsQueueItem

// func newPointsQueue(capacity int) pointsQueue {
// 	// We make capacity+1 because we need additional place for the greatest element
// 	return make([]pointsQueueItem, 0, capacity+1)
// }

// func (pq pointsQueue) Len() int { return len(pq) }

// func (pq pointsQueue) Less(i, j int) bool {
// 	// We want pop longest distances so Less was inverted
// 	return pq[i].distance > pq[j].distance
// }

// func (pq pointsQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }

// func (pq *pointsQueue) Push(x interface{}) {
// 	n := len(*pq)
// 	item := x.(pointsQueueItem)
// 	item.index = n
// 	*pq = append(*pq, item)
// }

// func (pq *pointsQueue) Pop() interface{} {
// 	old := *pq
// 	n := len(old)
// 	item := old[n-1]
// 	item.index = -1
// 	*pq = old[0 : n-1]
// 	return item
// }

type nearestVisitor[T math.Number] struct {
	point          orb.PointOf[T]
	filter         FilterFunc[T]
	k              int
	maxHeap        maxHeap[T]
	closestBound   *orb.BoundOf[T]
	maxDistSquared T
}

func (v *nearestVisitor[T]) Bound() *orb.BoundOf[T] {
	return v.closestBound
}

func (v *nearestVisitor[T]) Point() orb.PointOf[T] {
	return v.point
}

func (v *nearestVisitor[T]) Visit(n *node[T]) {
	// skip this pointer if we have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Point()
	if d := planar.DistanceSquared(point, v.point); d < v.maxDistSquared {
		v.maxHeap.Push(n.Value, d)
		if len(v.maxHeap) > v.k {

			v.maxHeap.Pop()

			// Actually this is a hack. We know how heap works and obtain
			// top element without function call
			top := v.maxHeap[0]

			v.maxDistSquared = top.distance

			// We have filled queue, so we start to restrict searching range
			d = math.Sqrt(top.distance)
			v.closestBound.Min[0] = v.point[0] - d
			v.closestBound.Max[0] = v.point[0] + d
			v.closestBound.Min[1] = v.point[1] - d
			v.closestBound.Max[1] = v.point[1] + d
		}
	}
}

type inBoundVisitor[T math.Number] struct {
	bound    *orb.BoundOf[T]
	pointers []orb.PointerOf[T]
	filter   FilterFunc[T]
}

func (v *inBoundVisitor[T]) Bound() *orb.BoundOf[T] {
	return v.bound
}

func (v *inBoundVisitor[T]) Point() (p orb.PointOf[T]) {
	return
}

func (v *inBoundVisitor[T]) Visit(n *node[T]) {
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	p := n.Value.Point()
	if v.bound.Min[0] > p[0] || v.bound.Max[0] < p[0] ||
		v.bound.Min[1] > p[1] || v.bound.Max[1] < p[1] {
		return

	}
	v.pointers = append(v.pointers, n.Value)
}

func childIndex[T math.Number](cx, cy T, point orb.PointOf[T]) int {
	i := 0
	if point[1] <= cy {
		i = 2
	}

	if point[0] >= cx {
		i++
	}

	return i
}
