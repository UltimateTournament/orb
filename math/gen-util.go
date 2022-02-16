package math

import (
	"math"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Min[T Number](a, b T) T {

	if a < b {
		return a
	} else {
		return b
	}
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func Sqrt[T Number](num T) T {
	return T(math.Sqrt(float64(num)))
}

func MinOf[T Number]() (r T) {
	switch x := any(&r).(type) {
	case *int:
		*x = math.MinInt
	case *int8:
		*x = math.MinInt8
	case *int16:
		*x = math.MinInt16
	case *int32:
		*x = math.MinInt32
	case *int64:
		*x = math.MinInt64
	case *uint:
		*x = 0
	case *uint8:
		*x = 0
	case *uint16:
		*x = 0
	case *uint32:
		*x = 0
	case *uint64:
		*x = 0
	case *float32:
		*x = -math.MaxFloat32
	case *float64:
		*x = -math.MaxFloat64
	default:
		panic("unreachable")
	}
	return
}

func MaxOf[T Number]() (r T) {
	switch x := any(&r).(type) {
	case *int:
		*x = math.MaxInt
	case *int8:
		*x = math.MaxInt8
	case *int16:
		*x = math.MaxInt16
	case *int32:
		*x = math.MaxInt32
	case *int64:
		*x = math.MaxInt64
	case *uint:
		*x = math.MaxUint
	case *uint8:
		*x = math.MaxUint8
	case *uint16:
		*x = math.MaxUint16
	case *uint32:
		*x = math.MaxUint32
	case *uint64:
		*x = math.MaxUint64
	case *float32:
		*x = math.MaxFloat32
	case *float64:
		*x = math.MaxFloat64
	default:
		panic("unreachable")
	}
	return
}
