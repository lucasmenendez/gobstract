package gobstract

func min(nums ...int) (min int) {
	min = nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
	}

	return min
}

func levenshtain(w1, w2 string) (diff int) {
	var l1, l2 int = len(w1), len(w2)
	var costs [][]int = make([][]int, l1 + 1)
	for i := 0; i < l1 + 1; i++ {
		costs[i] = make([]int, l2 + 1)
	}

	for i := 0; i < l1 + 1; i++ {
		costs[i][0] = i
	}

	for j := 0; j < l2 + 1; j++ {
		costs[0][j] = j
	}

	for i := 1; i < l1 + 1; i++ {
		for j := 1; j < l2 + 1; j++ {
			var o int = 1
			if w1[i - 1] == w2[j - 1] {
				o = 0
			}

			costs[i][j] = min(
				costs[i][j - 1] + 1,
				costs[i - 1][j] + 1,
				costs[i - 1][j - 1] + o,
			)
		}
	}

	return costs[l1][l2]
}
