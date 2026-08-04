package protocol

type ValueType string

const (
	SimpleString ValueType = "simple_string"
	Error        ValueType = "error"
	Integer      ValueType = "integer"
	BulkString   ValueType = "bulk_string"
	Array        ValueType = "array"
	Null         ValueType = "null"
)

type Value struct {
	Type ValueType

	Str string
	Int int64

	Array []Value
}

func NewSimpleString(s string) Value {
	return Value{
		Type: SimpleString,
		Str:  s,
	}
}

func NewError(s string) Value {
	return Value{
		Type: Error,
		Str:  s,
	}
}

func NewInteger(i int64) Value {
	return Value{
		Type: Integer,
		Int:  i,
	}
}

func NewBulkString(s string) Value {
	return Value{
		Type: BulkString,
		Str:  s,
	}
}

func NewNullBulkString() Value {
	return Value{
		Type: Null,
	}
}

func NewArray(values []Value) Value {
	return Value{
		Type:  Array,
		Array: values,
	}
}
