// +build !windows

package bytesink

import (
	"os"
	"syscall"

	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
)

type fifoByteSink struct {
	file *os.File
	path string
}

// Write writes the provided buffer to the underlying file.
func (s *fifoByteSink) Write(buf []byte) (n int, err error) {
	return s.file.Write(buf)
}

// Close calls fsync on the underlying temp file and then closes and removes it.
func (s *fifoByteSink) Close() (err error) {
	err = s.file.Sync()
	if err != nil {
		return err
	}

	defer func() {
		rerr := os.Remove(s.path)
		if err == nil {
			err = rerr
		}
	}()

	defer func() {
		cerr := s.file.Close()
		if err == nil {
			err = cerr
		}
	}()

	return nil
}

// NewByteSink creates a FIFO pipe and returns the address of a fifoByteSink,
// which satisfies the ByteSink interface. The FIFO pipe is used to stream bytes
// to rust-fil-proofs from Go during the piece-adding flow. Writes to the pipe
// are buffered automatically by the OS; the size of the buffer varies.
func NewByteSink(path string) (ByteSink, error) {
	err := syscall.Mkfifo(path, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "mkfifo failed")
	}

	file, err := os.OpenFile(path, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open pipe")
	}

	return &fifoByteSink{
		file: file,
		path: path,
	}, nil
}
