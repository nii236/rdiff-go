package rdiff

import (
	"io"
)

// Diff is a single diff operation for a chunk generated from the Rabin fingerprint.
type Diff struct {
	Signature uint64
	Start     uint
	Length    uint
	Cut       uint64
}

// PatchFile is the Delta, or Description. It contains the data needd to patch an old file to a new version.
type PatchFile struct {
	// Index of diffs to apply the patch
	Diffs []*Diff

	// Separate hashmap for literal bytes (signature as key) means we can
	// reuse the map value if a byte pattern occurs multiple times throughout the file
	ByteMap map[uint64][]byte
}

// Delta returns a patch file, used to update OldFile to NewFile, given the Signature of the OldFile
func Delta(oldSignature []*SignatureFile, newFile io.Reader) (*PatchFile, error) {
	newSignatureTable, err := Signature(newFile)
	if err != nil {
		return nil, err
	}

	result := &PatchFile{
		Diffs:   []*Diff{},
		ByteMap: map[uint64][]byte{},
	}

	for i := 0; i < len(newSignatureTable); i++ {
		newSig := newSignatureTable[i]

		// Handle case when new file is longer than old file
		if len(oldSignature) < i+1 {
			result.Diffs = append(result.Diffs, &Diff{
				Start:     newSig.Chunk.Start,
				Signature: newSig.Signature,
				Length:    newSig.Chunk.Length,
				Cut:       newSig.Chunk.Cut,
			},
			)
			_, ok := result.ByteMap[newSig.Signature]
			if !ok {
				result.ByteMap[newSig.Signature] = newSig.Chunk.Data
			}
			continue
		}

		// Standard case
		if oldSignature[i].Signature != newSignatureTable[i].Signature {
			result.Diffs = append(result.Diffs, &Diff{
				Start:     newSig.Chunk.Start,
				Signature: newSig.Signature,
				Length:    newSig.Chunk.Length,
				Cut:       newSig.Chunk.Cut,
			})
			_, ok := result.ByteMap[newSig.Signature]
			if !ok {
				result.ByteMap[newSig.Signature] = newSig.Chunk.Data
			}
		}
	}

	return result, nil
}
