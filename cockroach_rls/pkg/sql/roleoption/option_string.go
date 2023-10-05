// Code generated by "stringer"; DO NOT EDIT.

package roleoption

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CREATEROLE-1]
	_ = x[NOCREATEROLE-2]
	_ = x[PASSWORD-3]
	_ = x[LOGIN-4]
	_ = x[NOLOGIN-5]
	_ = x[VALIDUNTIL-6]
	_ = x[CONTROLJOB-7]
	_ = x[NOCONTROLJOB-8]
	_ = x[CONTROLCHANGEFEED-9]
	_ = x[NOCONTROLCHANGEFEED-10]
	_ = x[CREATEDB-11]
	_ = x[NOCREATEDB-12]
	_ = x[CREATELOGIN-13]
	_ = x[NOCREATELOGIN-14]
	_ = x[VIEWACTIVITY-15]
	_ = x[NOVIEWACTIVITY-16]
	_ = x[CANCELQUERY-17]
	_ = x[NOCANCELQUERY-18]
	_ = x[MODIFYCLUSTERSETTING-19]
	_ = x[NOMODIFYCLUSTERSETTING-20]
	_ = x[VIEWACTIVITYREDACTED-21]
	_ = x[NOVIEWACTIVITYREDACTED-22]
	_ = x[REPLICATION-23]
	_ = x[NOREPLICATION-24]
	_ = x[SQLLOGIN-25]
	_ = x[NOSQLLOGIN-26]
	_ = x[VIEWCLUSTERSETTING-27]
	_ = x[NOVIEWCLUSTERSETTING-28]
}

func (i Option) String() string {
	switch i {
	case CREATEROLE:
		return "CREATEROLE"
	case NOCREATEROLE:
		return "NOCREATEROLE"
	case PASSWORD:
		return "PASSWORD"
	case LOGIN:
		return "LOGIN"
	case NOLOGIN:
		return "NOLOGIN"
	case VALIDUNTIL:
		return "VALID UNTIL"
	case CONTROLJOB:
		return "CONTROLJOB"
	case NOCONTROLJOB:
		return "NOCONTROLJOB"
	case CONTROLCHANGEFEED:
		return "CONTROLCHANGEFEED"
	case NOCONTROLCHANGEFEED:
		return "NOCONTROLCHANGEFEED"
	case CREATEDB:
		return "CREATEDB"
	case NOCREATEDB:
		return "NOCREATEDB"
	case CREATELOGIN:
		return "CREATELOGIN"
	case NOCREATELOGIN:
		return "NOCREATELOGIN"
	case VIEWACTIVITY:
		return "VIEWACTIVITY"
	case NOVIEWACTIVITY:
		return "NOVIEWACTIVITY"
	case CANCELQUERY:
		return "CANCELQUERY"
	case NOCANCELQUERY:
		return "NOCANCELQUERY"
	case MODIFYCLUSTERSETTING:
		return "MODIFYCLUSTERSETTING"
	case NOMODIFYCLUSTERSETTING:
		return "NOMODIFYCLUSTERSETTING"
	case VIEWACTIVITYREDACTED:
		return "VIEWACTIVITYREDACTED"
	case NOVIEWACTIVITYREDACTED:
		return "NOVIEWACTIVITYREDACTED"
	case REPLICATION:
		return "REPLICATION"
	case NOREPLICATION:
		return "NOREPLICATION"
	case SQLLOGIN:
		return "SQLLOGIN"
	case NOSQLLOGIN:
		return "NOSQLLOGIN"
	case VIEWCLUSTERSETTING:
		return "VIEWCLUSTERSETTING"
	case NOVIEWCLUSTERSETTING:
		return "NOVIEWCLUSTERSETTING"
	default:
		return "Option(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
