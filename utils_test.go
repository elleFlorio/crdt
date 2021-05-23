package main

import (
	"testing"
)

func TestArrayDiff(t *testing.T) {
	empty := [0]GSetElement{}

	first := [3]GSetElement{
		{Id: 1, Value: "a"},
		{Id: 2, Value: "b"},
		{Id: 3, Value: "c"},
	}

	second := [3]GSetElement{
		{Id: 1, Value: "a"},
		{Id: 2, Value: "b"},
		{Id: 4, Value: "d"},
	}

	test := arrayDif(first[:], second[:])
	if len(test) != 1 {
		t.Fatalf("Array diff has wrong number of element. Expexted %d found %d", 1, len(test))
	}

	if test[0].Id != 4 {
		t.Fatalf("Array diff has wrong elements. Expexted Id %d found Id %d", 4, test[0].Id)
	}

	test2 := arrayDif(second[:], first[:])
	if len(test2) != 1 {
		t.Fatalf("Array diff has wrong number of element. Expexted %d found %d", 1, len(test2))
	}

	if test2[0].Id != 3 {
		t.Fatalf("Array diff has wrong elements. Expexted Id %d found Id %d", 4, test2[0].Id)
	}

	testEmpty := arrayDif(empty[:], first[:])
	if len(testEmpty) != 3 {
		t.Fatalf("Array diff has wrong number of element. Expexted %d found %d", 3, len(testEmpty))
	}

	for i, v := range testEmpty {
		if first[i].Id != v.Id {
			t.Fatalf("Array diff has wrong elements. Expexted %v found %v", first, testEmpty)
		}
	}

	testEmpty2 := arrayDif(first[:], empty[:])
	if len(testEmpty2) != 0 {
		t.Fatalf("Array diff has wrong number of element. Expexted %d found %d", 0, len(testEmpty2))
	}
}
