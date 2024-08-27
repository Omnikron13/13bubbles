package inlinelist

import (
   "strings"

   "golang.org/x/text/unicode/norm"
   "github.com/rivo/uniseg"
)

// Although the string utils here are very tiny pass-throughs, they are here mostly to provide a consistent interface
// in case the implementation changes in the future.


// stringOrBytes is a type that can be either a string or a byte slice, because Go string handling is a pain.
type stringOrBytes interface{string|[]byte}


// 'NO-BREAK SPACE' ('NBSP')
const nbsp = "\u00A0"


// spaceReplacer replaces all codepoints in category “Space Separator” with U+00A0; 'NO-BREAK SPACE' ('NBSP').
var spaceReplacer = strings.NewReplacer(
   "\u0020", nbsp, // 'SPACE' ('SP')
   "\u1680", nbsp, // 'OGHAM SPACE MARK'
   "\u2000", nbsp, // 'EN QUAD'
   "\u2001", nbsp, // 'EM QUAD'
   "\u2002", nbsp, // 'EN SPACE'
   "\u2003", nbsp, // 'EM SPACE'
   "\u2004", nbsp, // 'THREE-PER-EM SPACE'
   "\u2005", nbsp, // 'FOUR-PER-EM SPACE'
   "\u2006", nbsp, // 'SIX-PER-EM SPACE'
   "\u2007", nbsp, // 'FIGURE SPACE'
   "\u2008", nbsp, // 'PUNCTUATION SPACE'
   "\u2009", nbsp, // 'THIN SPACE'
   "\u200A", nbsp, // 'HAIR SPACE'
   "\u202F", nbsp, // 'NARROW NO-BREAK SPACE' ('NNBSP')
   "\u205F", nbsp, // 'MEDIUM MATHEMATICAL SPACE' ('MMSP')
   "\u3000", nbsp, // 'IDEOGRAPHIC SPACE'
)


// countGraphemes returns the number of graphemes in the given byte slice.
// Technically this is counting the number of terminal columns required to display the string rather than actual
// graphemes, as East Asian wide characters are counted as two columns.
func countGraphemes[T stringOrBytes](s T) int {
   return uniseg.StringWidth(string(s))
}


// noBreak returns a copy of the given string with all spaces replaced with non-breaking spaces.
func noBreak[T stringOrBytes](s T) T {
   return T(spaceReplacer.Replace(string(s)))
}


// normalise returns the Unicode Normalization Form C (NFC) of the given string.
func normalise[T stringOrBytes](s T) []byte {
   return norm.NFKC.Bytes([]byte(s))
}

