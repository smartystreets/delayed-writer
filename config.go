package delayed

import (
	"fmt"
	"io"
)

func New(options ...option) Writer {
	var config configuration
	Options.apply(options...)(&config)
	return newWriter(config)
}

func (singleton) Source(value func() fmt.Stringer) option {
	return func(this *configuration) { this.Source = value }
}
func (singleton) Target(value io.WriteCloser) option {
	return func(this *configuration) { this.Target = value }
}
func (singleton) PoolSize(value int) option {
	return func(this *configuration) { this.PoolSize = value }
}
func (singleton) ChannelSize(value int) option {
	return func(this *configuration) { this.ChannelSize = value }
}
func (singleton) Monitor(value Monitor) option {
	return func(this *configuration) { this.Monitor = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, item := range Options.defaults(options...) {
			item(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	empty := &nop{}

	return append([]option{
		Options.Source(func() fmt.Stringer { return empty }),
		Options.Target(empty),
		Options.PoolSize(1024),
		Options.ChannelSize(128),
		Options.Monitor(empty),
	}, options...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type configuration struct {
	Source      func() fmt.Stringer
	Target      io.WriteCloser
	PoolSize    int
	ChannelSize int
	Monitor     Monitor
}
type option func(*configuration)
type singleton struct{}

var Options singleton

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type nop struct{}

func (*nop) Buffered()                       {}
func (*nop) Discarded(fmt.Stringer)          {}
func (*nop) Written()                        {}
func (*nop) WriteFailed(fmt.Stringer, error) {}

func (*nop) String() string { return "" }

func (*nop) Write([]byte) (int, error) { return 0, nil }
func (*nop) Close() error              { return nil }
