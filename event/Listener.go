package event

type (
	EventName = string
	Listen    struct {
		EventNames []EventName
		Process    func(Registered)
	}
)

var eventListeners = map[EventName][]func(Registered){}

func Listener(listen Listen) {
	for _, name := range listen.EventNames {
		if _, ok := eventListeners[name]; !ok {
			eventListeners[name] = []func(Registered){}
		}
		eventListeners[name] = append(eventListeners[name], listen.Process)
	}
}
