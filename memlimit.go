package mm

import "context"

type (
	// Instance is NOT thread safe, do not attempt to call Close() at the same time Release() or Request() are called.
	Instance interface {
		// Request 'x' number of bytes of memory. Can be called multiple times as memory requirements of the instance
		// increase. If no bytes are available, will block until the requested number of bytes is available. If the
		// block takes too long, the context can be canceled and the call will return immediately.
		Request(ctx context.Context, x int) error
		// Release 'x' number of bytes of memory. Can be called multiple times as memory requirements of the instance
		// decrease.
		Release(x int)

		// ID returns the instance id
		ID() int

		// Close releases all bytes requested by this instance, No more calls to Request() or Release() should be called
		// after Close is called
		Close()
	}
)

type instance struct {
	m     *manager
	total int
	id    int
}

type req struct {
	resp <-chan struct{}
	x    int
}

func (i *instance) Request(ctx context.Context, x int) error {
	// FAST PATH
	// If the requested amount is available, a mutex request is sufficient.
	ch, err := i.m.request(i.id, x)
	if err != nil {
		return err
	}

	// SLOW PATH
	// If the manager returned a non-nil channel, then the manager wants us to queue our request
	if ch != nil {
		r := req{
			resp: make(<-chan struct{}),
			x:    x,
		}

		select {
		case ch <- r:
		case <-ctx.Done():
			return ctx.Err()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-r.resp:
		}
	}
	i.total += x
	return nil
}

func (i *instance) Release(x int) {
	i.m.release(i.id, x)
}

func (i *instance) ID() int {
	return i.id
}

func (i *instance) Close() {
	i.m.release(i.id, i.total)
}

// NewManager returns a new manager with a hard and soft limit defined.
func NewManager(soft int, hard int, deadLockCount int) Manager {
	return nil
}

// Manager coordinates all the instances created. It manages the total memory pool available and
// sets priority for newer instances.
type Manager interface {
	// NewInstance creates a new instance. The 'id' is an optional parameter used to determine priority when the
	// soft limit is reached, It is the callers responsibility to ensure the id is not already in use.
	//
	// Callers should ensure newer instances have a lower 'id' than older instances such that when
	// a soft limit is reached, the older instances get priority over newer instances.
	NewInstance(id int) Instance

	// Wait for all the instances to close. Does not notify the instances it's waiting, it is up to the implementation
	// to signal to the instances that the manager is waiting for them to release the bytes requested.
	Wait(ctx context.Context) error
}

type manager struct{}

func (m *manager) NewInstance(id int) Instance {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Wait(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) request(id, x int) (chan<- req, error) {
	//TODO implement me
	panic("implement me")
	return nil, nil
}

func (m *manager) release(id, x int) {
	//TODO implement me
	panic("implement me")
	return
}
