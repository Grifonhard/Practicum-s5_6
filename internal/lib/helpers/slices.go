package helpers

func GroupBy[T any, C comparable, List ~[]T](list List, iterFunc func(elem T) C) map[C]List {
	res := map[C]List{}

	for i := range list {
		key := iterFunc(list[i])

		res[key] = append(res[key], list[i])
	}

	return res
}

func Filter[T any, List ~[]T](list List, predicateFunc func(elem T, index int) bool) List {
	res := make(List, 0, len(list))

	for i := range list {
		if predicateFunc(list[i], i) {
			res = append(res, list[i])
		}
	}

	return res
}
