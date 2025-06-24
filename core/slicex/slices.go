package slicex

type Key interface {
	Key() string
}

func Union[T Key](first, second []T) []T {
	set := make(map[string]T)

	for index, each := range first {
		set[each.Key()] = first[index]
	}
	for index, each := range second {
		set[each.Key()] = second[index]
	}

	merged := make([]T, 0, len(set))
	for k := range set {
		merged = append(merged, set[k])
	}

	return merged
}

func Intersect[T Key](first, second []T) []T {
	set := make(map[string]T)

	for _, each := range first {
		set[each.Key()] = each
	}

	intersect := make([]T, 0, len(set))
	for _, each := range second {
		if _, ok := set[each.Key()]; ok {
			intersect = append(intersect, each)
		}
	}

	return intersect
}

func IntersectMap[T Key](first, second []T) map[string]T {
	set := make(map[string]T)

	for _, each := range first {
		set[each.Key()] = each
	}

	intersect := make(map[string]T)
	for _, each := range second {
		if _, ok := set[each.Key()]; ok {
			intersect[each.Key()] = each
		}
	}

	return intersect
}

func Contains[T Key](list []T, item T) int {
	for index, each := range list {
		if each.Key() == item.Key() {
			return index
		}
	}

	return -1
}
