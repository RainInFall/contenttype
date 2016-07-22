package contenttype

func (array Arraystring) Len() int {
	return len(array)
}

func (array Arraystring) Less(i, j int) bool {
	return array[i] < array[j]
}

func (array Arraystring) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}
