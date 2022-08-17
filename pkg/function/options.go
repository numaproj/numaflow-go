package function

import "time"

type handleOptions struct {
	eventTime time.Time
}

type HandleOption interface {
	apply(*handleOptions)
}

type eventTimeOption time.Time

func (c eventTimeOption) apply(opts *handleOptions) {
	opts.eventTime = time.Time(c)
}

func WithEventTime(c time.Time) HandleOption {
	return eventTimeOption(c)
}
