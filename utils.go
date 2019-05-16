package main

import "strconv"

// SafeStringToInt converts a string to an integer and defaults to `def`
func SafeStringToInt(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// IntMin returns the smallest integer
func IntMin(i1 int, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

// IntMax returns the largest integer
func IntMax(i1 int, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}
