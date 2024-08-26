package inlinelist

import (
   "golang.org/x/text/unicode/norm"
   "github.com/rivo/uniseg"
)

// Although the string utils here are very tiny pass-throughs, they are here mostly to provide a consistent interface
// in case the implementation changes in the future.


// stringOrBytes is a type that can be either a string or a byte slice, because Go string handling is a pain.
type stringOrBytes interface{string|[]byte}


// countGraphemes returns the number of graphemes in the given byte slice.
// Technically this is counting the number of terminal columns required to display the string rather than actual
// graphemes, as East Asian wide characters are counted as two columns.
func countGraphemes[T stringOrBytes](s T) int {
   return uniseg.StringWidth(string(s))
}


// normalise returns the Unicode Normalization Form C (NFC) of the given string.
func normalise[T stringOrBytes](s T) []byte {
   return norm.NFKC.Bytes([]byte(s))
}

