package gobstract

const scalingFactor float64 = 0.1

func jaroSimilarity(t1, t2 string) (dj float64) {
	if len(t1) == 0 || len(t2) == 0 {
		return
	} else if len(t2) < len(t1) {
		t1, t2 = t2, t1
	}

	var (
		t1l     float64 = float64(len(t1))
		t2l     float64 = float64(len(t2))
		matches []bool  = make([]bool, len(t2))
		m       float64
		t       float64
	)

	for i, c1 := range t1 {
		for j, c2 := range t2 {
			if matches[j] {
				continue
			}
			if i <= j && c1 == c2 {
				if i != j {
					t++
				}
				matches[j] = true
				m++
				break
			}
		}
	}

	if m == 0 {
		return
	}

	t /= 2
	dj = ((m / t1l) + (m / t2l) + (m-t)/m) / 3
	return
}

func jaroWinklerSimilarity(t1, t2 string) (d float64) {
	var dj = jaroSimilarity(t1, t2)

	if dj == 0 {
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

func strSimilarity(t1, t2 string) float64 {
	return jaroWinklerSimilarity(t1, t2)
}