package event

type Registered struct {
	name EventName
	data any
}

func NewRegistered(name EventName, data any) *Registered {
	return &Registered{
		name: name,
		data: data,
	}
}

func (reg *Registered) GetName() string {
	return reg.name
}

func (reg *Registered) GetData() any {
	return reg.data
}
