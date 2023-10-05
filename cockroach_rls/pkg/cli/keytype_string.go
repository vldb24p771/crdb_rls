// Code generated by "stringer"; DO NOT EDIT.

package cli

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[raw-0]
	_ = x[human-1]
	_ = x[rangeID-2]
	_ = x[hex-3]
}

func (i keyType) String() string {
	switch i {
	case raw:
		return "raw"
	case human:
		return "human"
	case rangeID:
		return "rangeID"
	case hex:
		return "hex"
	default:
		return "keyType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

var _keyTypes = map[string]keyType{
	"raw":     0,
	"human":   1,
	"rangeID": 2,
	"hex":     3,
}
