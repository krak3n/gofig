package fsnotify

import (
	"github.com/fsnotify/fsnotify"
)

// Notifier watches files for changes.
type Notifier struct {
	path string
	ch   chan error
}

// New constructs a new file Notifier.
func New(path string) *Notifier {
	return &Notifier{
		path: path,
	}
}

// Path returns the path to the file being watched.
func (n *Notifier) Path() string {
	return n.path
}

// Notify pushes a error value onto the channel when a file change occurs. This error could be nil
// or an actual error.
func (n *Notifier) Notify() <-chan error {
	ch := make(chan error, 1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ch <- err
		return ch
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if ok && event.Op&fsnotify.Write == fsnotify.Write {
					ch <- nil
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					ch <- err
				}
			}
		}
	}()

	if err := watcher.Add(n.Path()); err != nil {
		ch <- err
	}

	return ch
}

// Close stops listing for fsnotify events on the file.
func (n *Notifier) Close() {
	close(n.ch)
}
