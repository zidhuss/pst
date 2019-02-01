package db

import (
	"testing"
)

func TestIllegalID(t *testing.T) {
	// TODO:
}

func TestIntToIDAndBack(t *testing.T) {
	tt := []struct {
		input  uint64
		output string
	}{
		{50, "P"},
		{81734, "mXU"},
		{764547281, "UeleZ"},
	}

	for _, tc := range tt {
		// out := IntToID(tc.input)
		id := intToID(tc.input)
		if tc.output != id {
			t.Errorf("%d to id expected to return %s but got %s", tc.input,
				tc.output, id)
		}
		integer, _ := idToInt(tc.output)
		if tc.input != integer {
			t.Errorf("%s to int expected to return %d but got %d", tc.output,
				tc.input, integer)
		}
	}
}
