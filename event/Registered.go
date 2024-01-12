package event

type Registered struct {
	name Name
	data any
}

func NewRegistered(name Name, datas ...any) *Registered {
	var data any
	if len(datas) > 0 {
		data = datas[0]
	}
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
