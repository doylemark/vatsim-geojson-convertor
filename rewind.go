package main

import "reflect"

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Convert geojson polygons to clockwise direction per spec
func rewindBoundary(boundary *Boundary) *Boundary {
	coordinates := boundary.Geometry.Coordinates[0]

	var area float64 = 0
	len := len(coordinates)
	for i, j := 0, len-1; i < len; i, j = i+1, i+1 {
		area += (coordinates[i][0] - coordinates[j][0]) * (coordinates[j][1] + coordinates[i][1])
	}

	if area >= 0 {
		reverse(coordinates)
	}

	return boundary
}
