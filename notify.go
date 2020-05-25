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

// FileNotifyParser parses and watches for notifications from a notifier.
type FileNotifyParser struct {
	*FileParser

	notifier FileNotifier
}

// NewFileNotifyParser constructs a new FileNotifyParser.
func NewFileNotifyParser(parser ParseReadCloser, notifier FileNotifier) *FileNotifyParser {
	return &FileNotifyParser{
		FileParser: NewFileParser(parser, notifier.Path()),

		notifier: notifier,
	}
}

// Notify calls the wrapped function returning the values from the returned Notifier Notify method.
func (p *FileNotifyParser) Notify() <-chan error {
	return p.notifier.Notify()
}

// Close calls the wrapped function returning the values from the returned Notifier Close method.
func (p *FileNotifyParser) Close() error {
	return p.notifier.Close()
}

// FromFileAndNotify reads a file anf notifies of changes.
func FromFileAndNotify(parser ParseReadCloser, notifier FileNotifier) NotifyParser {
	return NewFileNotifyParser(parser, notifier)
}
