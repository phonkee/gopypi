/*
receiver is interface for signal receiver that will be called.

*/
package signals

import "context"

/*
Receiver is interface for signal handling
*/
type Receiver interface {
	Dispatch(context.Context)
}

type ReceiverFunc func(context.Context)

func (r ReceiverFunc) Dispatch(c context.Context) {
	r(c)
}
