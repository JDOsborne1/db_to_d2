package main

func in_set(_element string, _set []string) bool {
	for _, element := range _set {
		if _element == element {
			return true
		}
	}
	return false
}

func wrap_name_in_group(_table_name string, _grouping []TableGroup) string {
	for _, group := range _grouping {
		if in_set(_table_name, group.Tables) {
			return group.Tag + "." + _table_name
		}
	}
	return _table_name
}
