package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScanLines(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		atEOF   bool
		wantAdv int
		wantTok []byte
		wantErr bool
	}{
		{
			name:    "normal Line",
			input:   []byte("hello\n"),
			atEOF:   false,
			wantAdv: 6,
			wantTok: []byte("hello"),
			wantErr: false,
		},
		{
			name:    "end Of File",
			input:   []byte("hello"),
			atEOF:   true,
			wantAdv: 0,
			wantTok: nil,
			wantErr: true,
		},
		{
			name:    "empty at EOF",
			input:   nil,
			atEOF:   true,
			wantAdv: 0,
			wantTok: nil,
			wantErr: false,
		},
		{
			name:    "multiple Lines",
			input:   []byte("line1\nline2\n"),
			atEOF:   false,
			wantAdv: 6,
			wantTok: []byte("line1"),
			wantErr: false,
		},
		{
			name:    "no Newline",
			input:   []byte("line1 line2"),
			atEOF:   false,
			wantAdv: 0,
			wantTok: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adv, tok, err := ScanLines(tt.input, tt.atEOF)
			require.Equal(t, tt.wantAdv, adv)
			require.Equal(t, tt.wantTok, tok)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
