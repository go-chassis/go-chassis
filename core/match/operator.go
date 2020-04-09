package match

import (
	"regexp"
	"strconv"
	"strings"
)

func exact(value, express string) bool {
	return value == express
}

func contains(value, express string) bool {
	return strings.Contains(value, express)
}

func regex(value, express string) bool {
	reg := regexp.MustCompilePOSIX(express)
	if !reg.Match([]byte(value)) {
		return false
	}
	return true
}

func noEqu(value, express string) bool {
	return !(value == express)
}

func noLess(value, express string) bool {
	return cmpInt(value, express, func(v, e int) bool {
		return v >= e
	})
}

func less(value, express string) bool {
	return cmpInt(value, express, func(v, e int) bool {
		return v < e
	})
}

func noGreater(value, express string) bool {
	return cmpInt(value, express, func(v, e int) bool {
		return v <= e
	})
}

func greater(value, express string) bool {
	return cmpInt(value, express, func(v, e int) bool {
		return v > e
	})
}

func cmpInt(value, express string, op func(v, e int) bool) bool {
	v, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	exp, err := strconv.Atoi(express)
	if err != nil {
		return false
	}
	return op(v, exp)
}
