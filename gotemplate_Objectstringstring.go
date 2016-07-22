package contenttype

// template type Object(First, Second, FirstArray, SecondArray)

//Object has js-like functions
type Objectstringstring map[string]string

/*
Keys return array of keys of the Obecjt
*/
func (obj Objectstringstring) Keys() Arraystring {
	keys := make(Arraystring, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	return Arraystring(keys)
}
