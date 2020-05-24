package fsnotify

import (
	"crypto/md5" //nolint: gosec
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

// Notifier watches files for changes.
type Notifier struct {
	path    string
	hash    string
	watcher *fsnotify.Watcher
	ch      chan error
	doneCh  chan struct{}
}

// New constructs a new file Notifier.
func New(path string) *Notifier {
	return &Notifier{
		path:   path,
		ch:     make(chan error),
		doneCh: make(chan struct{}),
	}
}

// Path returns the path to the file being watched.
func (n *Notifier) Path() string {
	return n.path
}

// Notify pushes a error value onto the channel when a file change occurs. This error could be nil
// or an actual error.
func (n *Notifier) Notify() <-chan error {
	// Generate a hash of the files contents. We will use this to compare contents of the file to
	// guard against duplicate events. See issue:
	// https://github.com/fsnotify/fsnotify/issues/324
	hash, err := n.fileHash()
	if err != nil {
		n.ch <- err
		return n.ch
	}

	n.hash = hash

	// Create the watcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		n.ch <- err
		return n.ch
	}

	n.watcher = w

	go n.notify()

	if err := n.watcher.Add(n.Path()); err != nil {
		n.ch <- err
	}

	return n.ch
}

// Close stops listing for fsnotify events on the file.
func (n *Notifier) Close() error {
	defer close(n.ch)

	if n.watcher == nil {
		return nil
	}

	err := n.watcher.Close()

	<-n.doneCh

	return err
}

func (n *Notifier) notify() {
	defer close(n.doneCh)

	for {
		select {
		case event, ok := <-n.watcher.Events:
			if !ok {
				return // The channel has closed
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				// Get a new hash of the files contents.
				hash, err := n.fileHash()
				if err != nil {
					n.ch <- err
					continue
				}

				// If the hashes do not match send to notifcation channel
				if hash != n.hash {
					n.ch <- nil
				}

				n.hash = hash // Update the hash value
			}
		case err, ok := <-n.watcher.Errors:
			if !ok {
				return // The channel has closed
			}

			n.ch <- err
		}
	}
}

func (n *Notifier) fileHash() (string, error) {
	f, err := os.Open(n.Path())
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New() //nolint: gosec
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
