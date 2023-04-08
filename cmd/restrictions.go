package main

func minimalist(_table Table, _column Column) bool {
	return _column.Key == ""
}

func standard(_table Table, _column Column) bool {
	return false
}
