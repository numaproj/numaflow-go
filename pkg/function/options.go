package function

import "time"

// TODO: this options file is currently only a placeholder

type mapOptions struct {
	eventTime time.Time
}

type MapOption interface {
	apply(*mapOptions)
}

type eventTimeOption time.Time

func (c eventTimeOption) apply(opts *mapOptions) {
	opts.eventTime = time.Time(c)
}

func WithEventTime(c time.Time) MapOption {
	return eventTimeOption(c)
}
