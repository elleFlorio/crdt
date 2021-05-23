package main

func arrayDif(first []GSetElement, second []GSetElement) []GSetElement {
	check := make(map[uint32]bool, len(first))
	diff := make([]GSetElement, 0, len(second))

	for _, element := range first {
		check[element.Id] = true
	}

	for _, element := range second {
		if !check[element.Id] {
			diff = append(diff, element)
		}
	}

	return diff
}
