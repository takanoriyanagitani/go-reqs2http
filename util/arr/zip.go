package arr

import (
	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
)

func Zip[L, R any](l []L, r []R) []pair.Pair[L, R] {
	var li int = len(l)
	var ri int = len(r)
	var smaller int = min(li, ri)
	ret := make([]pair.Pair[L, R], 0, smaller)

	for j := 0; j < smaller; j++ {
		ret = append(ret, pair.Pair[L, R]{
			Left:  l[j],
			Right: r[j],
		})
	}
	return ret
}
