NOTE: This is just an idea, nothing has been implemented yet.

## A Memory Limiter Library
Because we don't know what requests in flight will take up a ton of memory, we could
OOM during normal operation. This library is an attempt at creating an upper limit on
the total amount of memory an application can use, so we don't go over our total
memory usage and OOM. Instead, we block some handlers until other handlers
can finish their work.

THIS LIBRARY DOES NOT MANAGE MEMORY, it only provides a way for handlers to record
and limit the total amount of memory allocated.

### Over the Soft limit
As long as we are below the soft limit, there is no restriction on
the number of bytes that can be requested by any handler. Once we are above the
soft limit, only a percentage of handlers will be able to Request() more
bytes. The closer we get to the hard limit, the lower the percentage of
handlers will be able to Request() more bytes.

Imagine we have a soft limit of 1000 bytes and a hard limit of 2000 bytes.

* Handler 1 Requests 200 bytes
* Handler 2 Requests 200 bytes
* Handler 3 Requests 200 bytes
* Handler 4 Requests 200 bytes
* Handler 5 Requests 200 bytes

We are now at the soft limit of 1000 bytes, any new handlers created by
calling `Manager.NewHandler() will not be allowed to request any bytes until
the total bytes requested drops below the soft limit. (Calls to Request() will
block). Only handlers that are currently in flight (have already requested bytes)
will be considered.

Soft limit rule is that only when 100% of the soft limit is available, the top
100% (ordered by handler number) of handlers are allowed to request more
bytes. So Handlers 1-5 can request bytes.

* Handler 1 Requests 200 bytes (Is allowed)

We have now used up 20% of our total soft limit; as a result, only 80% of the
top handlers can now request bytes. For example, Handler 5 will block if it
requests more bytes.

* Handler 2 Requests 200 bytes (Is allowed)

Now 40% of our total soft limit is used, now only 60% of the top handlers are
allowed to get new bytes.  This continues until 100% of the hard limit is used
up. Once all 100% of the soft limit is used up only the top `deadLockCount`
number of handlers can request new bytes. This ensures we don't deadlock
waiting for handlers to release memory. If the `deadLockCount` is set to zero,
then ALL pending Request() calls will deadlock until their respective contexts are
canceled.
