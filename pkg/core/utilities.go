package core

func in_set(_element string, _set []string) bool {
	for _, element := range _set {
		if _element == element {
			return true
		}
	}
	return false
}

func equal_set(_set1 []string, _set2 []string) bool {
	if len(_set1) != len(_set2) {
		return false
	}
	for _, element := range _set1 {
		if !in_set(element, _set2) {
			return false
		}
	}
	return true
}

// Wraps a table name in a group tag if it is in a group. Otherwise, returns the table name.
// This is used to ensure that links in the d2 graph are drawn correctly when table groups are used.
func wrap_name_in_group(_table_name string, _grouping []TableGroup) string {
	for _, group := range _grouping {
		if in_set(_table_name, group.Tables) {
			return group.Tag + "." + _table_name
		}
	}
	return _table_name
}
