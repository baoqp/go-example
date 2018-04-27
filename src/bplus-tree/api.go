package bplus_tree

func Open(filename string) (*Tree, error) {
	return open(filename, false)
}

func (t *Tree) Close() error {
	t.rmLock.Lock()
	t.file.Close()
	if t.header.page != nil {
		t.header.page.destroy()
		t.header.page = nil
	}
	return nil
}

func (t *Tree) Get(key []byte) ([]byte, error) {
	Key := NewKey(key)
	var Value Value
	err := t.get(Key, &Value)
	if err != nil {
		return nil, err
	}
	return Value.value, nil
}

func (t *Tree) Update(key []byte, value []byte, updateCb UpdateCallback, arg []byte) error {
	Key := NewKey(key)
	Value := NewValue(value)
	return t.update(Key, Value, updateCb, arg)
}

func (t *Tree) BulkUpdate(count uint64, keys [][]byte, values [][]byte, updateCb UpdateCallback, arg []byte) error {
	var Keys []*Key
	for _, key := range keys {
		Keys = append(Keys, NewKey(key))
	}

	var Values []*Value
	for _, key := range values {
		Values = append(Values, NewValue(key))
	}

	return t.bulkUpdate(count, Keys, Values, updateCb, arg)
}

func (t *Tree) Set(key []byte, value []byte) error {
	return t.Update(key, value, nil, nil)
}

func (t *Tree) BulkSet(count uint64, keys [][]byte, values [][]byte, ) error {
	return t.BulkUpdate(count, keys, values, nil, nil)
}

func (t *Tree) Remove(key []byte, removeCb RemoveCallback, arg []byte) error {
	Key := NewKey(key)
	return t.remove(Key, removeCb, arg)
}

func (t *Tree) GetFilteredRange(start []byte, end []byte, callback FilterCallback,
	rangeCallback RangeCallback, arg []byte) error {

	Start := NewKey(start)
	End := NewKey(end)

	return t.getFilteredRange(Start, End, callback, rangeCallback, arg)
}

func (t *Tree) GetRange(start []byte, end []byte, rangeCallback RangeCallback,
	arg []byte) error {

	Start := NewKey(start)
	End := NewKey(end)

	return t.getRange(Start, End, rangeCallback, arg)
}
