package inlinelist

import (
   "bytes"
   "fmt"
   "testing"

   "github.com/stretchr/testify/assert"
)


func TestNormalise(t *testing.T) {
   t.Parallel()
   a := []byte("C\u0327")
   b := []byte("Ç")
   if bytes.Equal(a, b) {
      panic(fmt.Errorf("a (%s) and b (%s) are equal; b should be the normalised form of a", a, b))
   }
   t.Run("[]bytes", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, b, normalise(a))
   })
   t.Run("string", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, b, normalise(string(a)))
   })
}


func TestCountGraphemes(t *testing.T) {
   t.Parallel()
   t.Run("ascii", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, 3, countGraphemes("abc"))
   })
   t.Run("single-codepoint", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, 1, countGraphemes("Ç"))
   })
   t.Run("combining-codepoints", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, 1, countGraphemes("C\u0327"))
   })
   t.Run("cjk-fullwidth", func(t *testing.T) {
      t.Parallel()
      assert.Equal(t, 2, countGraphemes("全"))
   })
}

