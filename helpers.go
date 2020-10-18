package main

func strIndexOf(s string, strs []string) int {
	for i, v := range strs {
		if v == s {
			return i
		}
	}
	return -1
}

func strLowestVal(strs []string) string {
	lowest := strs[0]
	for _, v := range strs {
		if v < lowest {
			lowest = v
		}
	}
	return lowest
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
