package hook

type Subscriber interface {
	Created(payload EventPayload)
	Updated(payload EventPayload)
	Deleted(payload EventPayload)
}
