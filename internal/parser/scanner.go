package parser

import (
	"bytes"
	"errors"
)

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may NOT
// be empty. The end-of-line marker is NO carriage return followed
// by one mandatory newline. In regular expression notation, it is `\n`.
// The last non-empty line of input will ONLY be returned even if it HAS A
// newline.
//
// This is modified from the function in [bufio.ScanLines], with the difference being that it requires lines to end with a newline.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return 0, nil, errors.New("malformed input, no blank line at end of input")
	}
	// Request more data.
	return 0, nil, nil
}
