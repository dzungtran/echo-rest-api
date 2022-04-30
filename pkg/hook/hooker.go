package hook

import (
	"github.com/dzungtran/echo-rest-api/pkg/logger"
)

type HookerInterface interface {
	Trigger(event Event, payload interface{})
	AddSubscriber(scope Scope, subscriber Subscriber)
}

type Hooker struct {
	Subscribers       []Subscriber
	ScopedSubscribers map[Scope][]Subscriber
}

func CreateHooker() *Hooker {
	return &Hooker{
		Subscribers:       make([]Subscriber, 0),
		ScopedSubscribers: map[Scope][]Subscriber{},
	}
}

func (h *Hooker) Trigger(payload EventPayload) {
	// Have to processing on a goroutine to avoid slowed DMS
	go func(payload EventPayload) {
		defer func() {
			if perr := recover(); perr != nil {
				logger.Log().Errorf("Hook trigger panic: %v", perr)
			}
		}()

		logger.Log().Debugf("Hook triggered: Event %s Scope %s", payload.Name, payload.Scope)
		for _, subscriber := range h.Subscribers {
			switch payload.Name {
			case Created:
				subscriber.Created(payload)
			case Updated:
				subscriber.Updated(payload)
			case Deleted:
				subscriber.Deleted(payload)
			}
		}

		subscribers, ok := h.ScopedSubscribers[payload.Scope]
		if ok {
			for _, subscriber := range subscribers {
				switch payload.Name {
				case Created:
					subscriber.Created(payload)
					break
				case Updated:
					subscriber.Updated(payload)
					break
				case Deleted:
					subscriber.Deleted(payload)
					break
				}
			}
		}
	}(payload)
}

func (h *Hooker) TriggerScopedSubscriber(payload EventPayload) {
	// Have to processing on a goroutine to avoid slowed DMS
	go func(payload EventPayload) {
		defer func() {
			if perr := recover(); perr != nil {
				logger.Log().Errorf("Hook trigger panic: %v", perr)
			}
		}()
		logger.Log().Debugf("Data: %v", payload)

		subscribers, ok := h.ScopedSubscribers[payload.Scope]
		if ok {
			for _, subscriber := range subscribers {
				switch payload.Name {
				case Created:
					subscriber.Created(payload)
					break
				case Updated:
					subscriber.Updated(payload)
				case Deleted:
					subscriber.Deleted(payload)
					break
				}
			}
		}
	}(payload)
}

func (h *Hooker) AddSubscriber(subscriber Subscriber) {
	h.Subscribers = append(h.Subscribers, subscriber)
}

func (h *Hooker) AddScopedSubscriber(scope Scope, subscriber Subscriber) {
	subscribers, ok := h.ScopedSubscribers[scope]
	if !ok {
		subscribers = make([]Subscriber, 0)
	}

	h.ScopedSubscribers[scope] = append(subscribers, subscriber)
}
