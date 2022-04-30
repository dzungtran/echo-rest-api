package defines

const (
	DefaultPerPage uint64 = 100
	MaximumPerPage uint64 = 250
)

type (
	CommonParamsForFetch struct {
		Page    uint64 `json:"page"`
		Limit   uint64 `json:"limit"`
		NoLimit bool   `json:"no_limit"`
	}
)
