package event

const (
	EventLinkVisited = "link.visited"
)

type Event struct {
	Type string
	Data any
}

type EventBys struct {
	bus chan Event
}

func NewEventBus() *EventBys {
	return &EventBys{
		bus: make(chan Event),
	}
}

func (e *EventBys) Publish(event Event) {
	e.bus <- event
}

func (e *EventBys) Subscribe() <-chan Event {
	return e.bus
}
