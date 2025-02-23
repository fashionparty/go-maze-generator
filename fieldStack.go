package main

type FieldStack struct {
	stack []*Field
}

func (fs *FieldStack) Add(field *Field) {
	fs.stack = append(fs.stack, field)
}

func (fs *FieldStack) Pop() (*Field, bool) {
	if len(fs.stack) == 0 {
		return nil, false
	}
	lastIndex := len(fs.stack) - 1
	field := fs.stack[lastIndex]
	fs.stack = fs.stack[:lastIndex]
	return field, true
}
