package gobstract

const scalingFactor float64 = 0.1

// jaroDistance functions implements Jaro String Similarity algorithm, that
// determines how similar are two strings. Jaro Similarity returns value between
// 0 and 1, that represents the distance between both strings. Check out here:
// https://goo.gl/se7e3g
func jaroDistance(t1, t2 string) (d float64) {
	if len(t1) == 0 || len(t2) == 0 {
		return
	} else if len(t2) < len(t1) {
		t1, t2 = t2, t1
	}

	var (
		t1l float64 = float64(len(t1))
		t2l float64 = float64(len(t2))
		m   float64
		t   float64
	)

	for i, c1 := range t1 {
		for j, c2 := range t2 {
			if i <= j && c1 == c2 {
				if i != j {
					t++
				}
				m++
				break
			}
		}
	}

	if m == 0 {
		return
	}

	t /= 2
	d = ((m / t1l) + (m / t2l) + (m-t)/m) / 3
	return
}

// jaroWinklerDistance function extends Jaro String Similarity algorithm,
// measuring common strings prefixes. Check out here: https://goo.gl/6fvpAM
func jaroWinklerDistance(t1, t2 string) (d float64) {
	var dj float64
	if dj = jaroDistance(t1, t2); dj == 0 {
		return
	}

	var l float64
	for _, c1 := range t1 {
		for _, c2 := range t2 {
			if c1 == c2 && l < 4 {
				l++
			}
		}
	}

	d = dj + (l * scalingFactor * (1 - dj))
	return
}

// strDistance functions is the generic function to call across the code. The
// functions call the chosen strings distance algorithm.
func strDistance(t1, t2 string) float64 {
	return jaroWinklerDistance(t1, t2)
}
