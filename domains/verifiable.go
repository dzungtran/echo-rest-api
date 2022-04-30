package domains

type Verifiable interface {
	Verify() error
}

func VerifyObject(o interface{}) error {
	obj, ok := o.(Verifiable)
	if !ok {
		return nil
	}
	return obj.Verify()
}
