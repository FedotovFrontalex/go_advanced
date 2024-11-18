package stat

import (
	"log"
	"server/pkg/event"
)

type StatServiceDeps struct {
	EventBus       *event.EventBys
	StatRepository *StatRepository
}

type StatService struct {
	EventBus       *event.EventBys
	StatRepository *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			id, ok := msg.Data.(uint)

			if !ok {
				log.Fatalln("Bsd EventLinkVisited Data:  ", msg.Data)
				continue
			}
			s.StatRepository.AddClick(id)
		}
	}
}
