package eventbus

func newError[Payload any](err error, p Payload) Error[Payload] {
	return Error[Payload]{
		Origin:  err,
		Payload: p,
	}
}

type Error[Payload any] struct {
	Origin  error
	Payload Payload
}

func (err Error[Payload]) Error() string {
	return err.Origin.Error()
}
