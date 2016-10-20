package chunker

import (
	"bytes"
	"io"
)

type Chunker interface {
	Chunk(io.Reader) ([]io.Reader, error)
}

func NewChunker(size int64) Chunker {
	return chunker{
		size: size,
	}
}

type chunker struct {
	size int64
}

func (c chunker) Chunk(r io.Reader) ([]io.Reader, error) {
	out := make([]io.Reader, 0)
	for {
		var buf []byte
		dst := bytes.NewBuffer(buf)
		_, err := io.CopyN(dst, r, int64(c.size))
		switch err {
		case io.EOF:
			// we are at the end of the reader, but haven't read any chunks
			if dst.Len() == 0 {
				return nil, io.EOF
			}

			out = append(out, dst)
			return out, nil
		case nil:
			out = append(out, dst)
		default:
			// something else went wrong reading from our reader - all is lost
			return nil, err
		}
	}
}
