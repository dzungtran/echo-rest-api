package contexts

type (
	CommonParamsForFetch struct {
		Page    uint64 `json:"page"`
		Limit   uint64 `json:"limit"`
		NoLimit bool   `json:"no_limit"`
	}
)
