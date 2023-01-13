package rdiff

import (
	"errors"
	"fmt"
	"io"

	"github.com/cespare/xxhash"
	"github.com/restic/chunker"
)

type SignatureFile struct {
	Chunk     chunker.Chunk
	Signature uint64
}

// Signature of the input file. If two signatures are different from two
// versions of the same file, use Delta to create a patch
func Signature(input io.Reader) ([]*SignatureFile, error) {
	b := make([]byte, chunker.MaxSize)
	c := chunker.New(input, 0x3dea92648f6e83)
	result := []*SignatureFile{}
	for {
		ch, err := c.Next(b)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("read: %w", err)
		}
		if err != nil {
			break
		}

		result = append(result, &SignatureFile{
			Chunk:     ch,
			Signature: xxhash.Sum64(ch.Data),
		})
	}
	return result, nil
}
