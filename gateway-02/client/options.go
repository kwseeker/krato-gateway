package client

import "github.com/go-kratos/kratos/v2/selector"

// callInfo contains all related configuration and information about an RPC.
type callInfo struct {
	filters []selector.NodeFilter
}

// csAttempt implements a single transport stream attempt within a
// clientStream.
type csAttempt struct{}

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*callInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*callInfo, *csAttempt)
}

// FilterOption .
type FilterOption struct {
	Filters []selector.NodeFilter
}

func (o FilterOption) before(c *callInfo) error {
	c.filters = o.Filters
	return nil
}

func (o FilterOption) after(*callInfo, *csAttempt) {}

// WithFilter .
func WithFilter(filters []selector.NodeFilter) CallOption {
	return FilterOption{Filters: filters}
}

func defaultCallInfo() callInfo {
	return callInfo{}
}
