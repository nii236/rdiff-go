package rdiff_test

import (
	"fmt"
	"io"
	"os"
	"rdiff"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
)

func LoadFile(t *testing.T, name string) (io.Reader, int64, io.Reader, int64) {
	oldFile := fmt.Sprintf("./testdata/%s.old", name)
	newFile := fmt.Sprintf("./testdata/%s.new", name)
	fOld, err := os.Open(oldFile)
	if err != nil {
		t.Fatal(err)
	}
	fNew, err := os.Open(newFile)
	if err != nil {
		t.Fatal(err)
	}
	fOldStat, err := os.Stat(oldFile)
	if err != nil {
		t.Fatal(err)
	}
	fNewStat, err := os.Stat(newFile)
	if err != nil {
		t.Fatal(err)
	}
	return fOld, fOldStat.Size(), fNew, fNewStat.Size()
}

func TestDelta(t *testing.T) {
	alphaOld, alphaOldSize, alphaNew, alphaNewSize := LoadFile(t, "alpha")
	duplicationOld, duplicationOldSize, duplicationNew, duplicationNewSize := LoadFile(t, "duplication")
	prependOld, prependOldSize, prependNew, prependNewSize := LoadFile(t, "prepend")
	loremOld, loremOldSize, loremNew, loremNewSize := LoadFile(t, "insertion")
	identicalOld, identicalOldSize, identicalNew, identicalNewSize := LoadFile(t, "identical")
	removalOld, removalOldSize, removalNew, removalNewSize := LoadFile(t, "removal")
	type args struct {
		oldBytesTotal int64
		newBytesTotal int64
		oldSig        []*rdiff.SignatureFile
		newFile       io.Reader
	}
	fOldSignature, err := rdiff.Signature(alphaOld)
	if err != nil {
		t.Fatal(err)
	}
	duplicationOldSignature, err := rdiff.Signature(duplicationOld)
	if err != nil {
		t.Fatal(err)
	}
	prependOldSignature, err := rdiff.Signature(prependOld)
	if err != nil {
		t.Fatal(err)
	}
	loremOldSignature, err := rdiff.Signature(loremOld)
	if err != nil {
		t.Fatal(err)
	}
	identicalOldSignature, err := rdiff.Signature(identicalOld)
	if err != nil {
		t.Fatal(err)
	}
	removalOldSignature, err := rdiff.Signature(removalOld)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name           string
		args           args
		wantDiffLen    int
		wantByteMapLen int
		wantErr        bool
	}{
		{"alphabet_replace", args{oldBytesTotal: alphaOldSize, newBytesTotal: alphaNewSize, oldSig: fOldSignature, newFile: alphaNew}, 1, 1, false},
		{"appended_duplicate", args{oldBytesTotal: duplicationOldSize, newBytesTotal: duplicationNewSize, oldSig: duplicationOldSignature, newFile: duplicationNew}, 11, 11, false},
		{"prepend_string", args{oldBytesTotal: prependOldSize, newBytesTotal: prependNewSize, oldSig: prependOldSignature, newFile: prependNew}, 2, 2, false},
		{"insert_center", args{oldBytesTotal: loremOldSize, newBytesTotal: loremNewSize, oldSig: loremOldSignature, newFile: loremNew}, 1, 1, false},
		{"identical_files", args{oldBytesTotal: identicalOldSize, newBytesTotal: identicalNewSize, oldSig: identicalOldSignature, newFile: identicalNew}, 0, 0, false},
		{"removal", args{oldBytesTotal: removalOldSize, newBytesTotal: removalNewSize, oldSig: removalOldSignature, newFile: removalNew}, 1, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdiff.Delta(tt.args.oldSig, tt.args.newFile)
			assert.Nil(t, err)

			t.Logf("diffs: %d old: %s new: %s\n", len(got.Diffs), humanize.Bytes(uint64(tt.args.oldBytesTotal)), humanize.Bytes(uint64(tt.args.newBytesTotal)))
			for _, el := range got.Diffs {
				t.Logf("pos: %d len: %d\n", el.Start, el.Length)
			}

			assert.Equal(t, tt.wantDiffLen, len(got.Diffs))
			assert.Equal(t, tt.wantByteMapLen, len(got.ByteMap))
		})
	}
}
