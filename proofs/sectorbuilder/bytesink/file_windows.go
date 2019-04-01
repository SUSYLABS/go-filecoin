package bytesink

import (
	"os"

	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
)

type winFileByteSink struct {
	file *os.File
	path string
}

// Write writes the provided buffer to the underlying file.
func (s *winFileByteSink) Write(buf []byte) (n int, err error) {
	return s.file.Write(buf)
}

// Close calls fsync on the underlying temp file and then closes and removes it.
func (s *winFileByteSink) Close() (err error) {
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

// NewByteSink creates a ByteSink which writes to the file at the given path.
func NewByteSink(path string) (ByteSink, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open pipe")
	}

	return &winFileByteSink{
		file: file,
		path: path,
	}, nil
}
