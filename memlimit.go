package mm

import "context"

// Instance is NOT thread safe, do not attempt to call Close() at the same time Release() or Request() are called.
type Instance interface {
	// Request 'x' number of bytes of memory. Can be called multiple times as memory requirements of the instance
	// increase. If no bytes are available, will block until the requested number of bytes is available. If the
	// block takes too long, the context can be canceled and the call will return immediately.
	Request(ctx context.Context, x int) error
	// Release 'x' number of bytes of memory. Can be called multiple times as memory requirements of the instance
	// decrease.
	Release(x int)
	// Close releases all bytes requested by this instance, No more calls to Request() or Release() should be called
	// after Close is called
	Close()
}

// NewManager returns a new manager with a hard and soft limit defined.
func NewManager(soft int64, hard int64, deadLockCount int) Manager {
	return nil
}

// Manager is the coordinator for all the instances created. It manages the total memory pool available and
// sets priority for newer instances.
type Manager interface {
	// NewInstance creates a new instance, the 'id' is used to determine priority when the soft limit is reached
	// In general callers should ensure newer instances have a lower 'id' than older instances such that when
	// a soft limit is reached, the older instances get priority over newer instances.
	NewInstance(id int) Instance
}
