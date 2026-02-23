package internal

import (
	"encoding/binary"
	"io"
	"math/rand/v2"
)

// var Random = bufio.NewReaderSize(rand.Reader, 1<<15)

type p struct {
	r   *rand.ChaCha8
	off int
	ui  [8]byte
}

func (p *p) Read(b []byte) (n int, err error) {
	idx := 0
	for idx < len(b) {
		if p.off == 8 {
			binary.BigEndian.PutUint64(p.ui[:], p.r.Uint64())
			p.off = 0
		}

		b[idx] = p.ui[p.off]
		p.off++
		idx++
	}
	return len(b), nil
}

var Random io.Reader = &p{r: rand.NewChaCha8([32]byte{})}
