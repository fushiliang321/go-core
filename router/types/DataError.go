package types

type DataError struct {
	text string
	data any
}

func (err *DataError) Error() string {
	return err.text
}
func (err *DataError) Data() any {
	return err.data
}

func NewDataError(text string, data any) *DataError {
	return &DataError{
		text: text,
		data: data,
	}
}
