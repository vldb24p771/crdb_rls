// Code generated by "stringer"; DO NOT EDIT.

package base

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ReplicationAuto-0]
	_ = x[ReplicationManual-1]
}

func (i TestClusterReplicationMode) String() string {
	switch i {
	case ReplicationAuto:
		return "ReplicationAuto"
	case ReplicationManual:
		return "ReplicationManual"
	default:
		return "TestClusterReplicationMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
