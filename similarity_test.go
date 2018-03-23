package gobstract

import (
	"testing"
	//"fmt"
)

func TestJaroDistance(t *testing.T) {
	var (
		c1 string = "MOON"
		c2 string = "EARTH"
		c3 string = "EARTH PLANET"
	)

	var (
		d1 float64 = jaroDistance(c1, c2)
		d2 float64 = jaroDistance(c2, c2)
		d3 float64 = jaroDistance(c2, c3)
	)

	if d1 != 0 {
		t.Errorf("Expected 0, got %f", d1)
	}

	if d2 != 1 {
		t.Errorf("Expected 1, got %f", d2)
	}

	if d3 != 0.8055555555555557 {
		t.Errorf("Expected 0.8055555555555557, got %f", d3)
	}
}

func TestJaroWinklerDistance(t *testing.T) {
	var (
		c1 string = "MOON"
		c2 string = "EARTH"
		c3 string = "EARTH PLANET"
	)

	var (
		d1 float64 = jaroWinklerDistance(c1, c2)
		d2 float64 = jaroWinklerDistance(c2, c2)
		d3 float64 = jaroWinklerDistance(c2, c3)
	)

	if d1 != 0 {
		t.Errorf("Expected 0, got %f", d1)
	}

	if d2 != 1 {
		t.Errorf("Expected 1, got %f", d2)
	}

	if d3 != 0.8833333333333334 {
		t.Errorf("Expected 0.8055555555555557, got %f", d3)
	}
}
