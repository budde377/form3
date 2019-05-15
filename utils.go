package main

import "strconv"

func SafeStringToInt(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func IntMin(i1 int, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

func IntMax(i1 int, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}