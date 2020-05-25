package gofig

import "context"

// A Notifier is a Parser that notifies via a channel if changes to configuration have occurred.
// Remember to check the error on the channel.
type Notifier interface {
	Notify() <-chan error
	Close() error
}

// A FileNotifier is a Notifier that also returns the path to the file being watched.
type FileNotifier interface {
	Notifier

	Path() string
}

// A NotifyParser is a Notifier that also Parses configuration.
type NotifyParser interface {
	Parser
	Notifier
}

// Notify notifies when a change to configuration has occurred.
func (l *Loader) Notify(c chan<- error, notifiers ...NotifyParser) {
	l.NotifyWithContext(context.Background(), c, notifiers...)
}

// NotifyWithContext notifies when a change to configuration has occurred.
func (l *Loader) NotifyWithContext(ctx context.Context, c chan<- error, notifiers ...NotifyParser) {
	l.notifiers = append(l.notifiers, notifiers...)
	l.wg.Add(len(notifiers))

	for _, n := range notifiers {
		go func(n NotifyParser) {
			defer l.wg.Done()

			ch := n.Notify()
			for {
				select {
				case <-ctx.Done():
					return
				case err, ok := <-ch:
					if !ok {
						return // Channel is closed
					}

					if err == nil {
						err = l.Parse(n)
					}

					c <- err
				}
			}
		}(n)
	}
}

// Close stops listening for notification events. This only needs to be called if Notify or
// NotifyWithContext are being used.
func (l *Loader) Close() error {
	var err CloseError

	for _, n := range l.notifiers {
		if e := n.Close(); e != nil {
			err.Add(e)
		}
	}

	l.wg.Wait()

	return err.NilOrError()
}

// ParseNotifierFunc implements the Notifier and Parser interface.
type ParseNotifierFunc func() (Parser, Notifier)

// Keys consumes the keys but does nothing with them.
func (fn ParseNotifierFunc) Keys(c <-chan string) error {
	for {
		_, ok := <-c
		if !ok {
			return nil
		}
	}
}

// Values calls the wrapped function returning the values from the returned Parser Values method.
func (fn ParseNotifierFunc) Values() (<-chan func() (string, interface{}), error) {
	p, _ := fn()

	return p.Values()
}

// Notify calls the wrapped function returning the values from the returned Notifier Notify method.
func (fn ParseNotifierFunc) Notify() <-chan error {
	_, n := fn()

	return n.Notify()
}

// Close calls the wrapped function returning the values from the returned Notifier Close method.
func (fn ParseNotifierFunc) Close() error {
	_, n := fn()

	return n.Close()
}

// FromFileWithNotify reads a file anf notifies of changes.
func FromFileWithNotify(parser ReaderParser, notifier FileNotifier) NotifyParser {
	return ParseNotifierFunc(func() (Parser, Notifier) {
		return FromFile(parser, notifier.Path()), notifier
	})
}
