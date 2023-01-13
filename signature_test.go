package rdiff_test

import (
	"io"
	"rdiff"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	alphaOld, _, alphaNew, _ := LoadFile(t, "alpha")
	duplicationOld, _, duplicationNew, _ := LoadFile(t, "duplication")
	prependOld, _, prependNew, _ := LoadFile(t, "prepend")
	loremOld, _, loremNew, _ := LoadFile(t, "insertion")
	removalOld, _, removalNew, _ := LoadFile(t, "removal")
	type args struct {
		input io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"alpha_old", args{alphaOld}, false},
		{"alpha_new", args{alphaNew}, false},
		{"duplication_old", args{duplicationOld}, false},
		{"duplication_new", args{duplicationNew}, false},
		{"prepend_old", args{prependOld}, false},
		{"prepend_new", args{prependNew}, false},
		{"lorem_old", args{loremOld}, false},
		{"lorem_new", args{loremNew}, false},
		{"removal_old", args{removalOld}, false},
		{"removal_new", args{removalNew}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdiff.Signature(tt.args.input)
			assert.Nil(t, err)

			for _, el := range got {
				t.Logf("pos: %04d len: %04d digest: %04d signature: %x\n",
					el.Chunk.Start,
					el.Chunk.Length,
					el.Chunk.Cut,
					el.Signature,
				)
			}
		})
	}
}
