package hook

type Event string

const (
	// Define events
	Created Event = "created"
	Updated Event = "updated"
	Deleted Event = "deleted"

	// Define scopes
	UserScope Scope = "user"

	// Event source
	SourceUserAPI = "user_api"
)

type Scope string

type EventPayload struct {
	Name       Event
	Scope      Scope
	Source     string
	PayloadOld interface{}
	Payload    interface{}
}

type RawEventPayload struct {
	Name    Event  `json:"name"`
	Scope   Scope  `json:"scope"`
	Source  string `json:"source"`
	Payload []byte `json:"payload"`
}
