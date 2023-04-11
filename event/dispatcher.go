package event

func Dispatch(reg *Registered) {
	funs, ok := eventListeners[reg.name]
	if !ok {
		return
	}
	for _, fun := range funs {
		fun(*reg)
	}
}
