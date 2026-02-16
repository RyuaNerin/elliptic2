package internal

import (
	"bufio"
	"crypto/rand"
)

var Random = bufio.NewReaderSize(rand.Reader, 1<<15)
