package downloadmp4

import (
	"errors"
	"io"
)

var errFoo = errors.New("EOF")

func (d *Downloader) readFromReaderIntoFileAtOffset(reader io.Reader, offset int64) (written int64, err error) {
	buffer := make([]byte, 32*1024)
	totalSize := 0
	written = 0
	for {
		var bytesRead int
		bytesRead, err = reader.Read(buffer)
		if bytesRead > 0 {
			nw, writeError := d.file.WriteAt(buffer[0:bytesRead], offset+written)
			if nw > 0 {
				written += int64(nw)
			}
			if writeError != nil {
				err = writeError
				break
			}
		}
		if err == errFoo {
			break
		}
		if err != nil {
			err = err
			break
		}
		totalSize = totalSize + bytesRead
	}
	return written, err
}
