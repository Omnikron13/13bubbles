package inlinelist

import (
   lg "github.com/charmbracelet/lipgloss"
)

// ItemStyle is a type grouping lipgloss styles for a list item and its prefix and suffix, to facilitate
// applying a variety of different styles for different states of the list item (e.g. selected, focussed, etc.)
type ItemStyle struct {
   Main lg.Style
   Prefix lg.Style
   Suffix lg.Style
}

// ItemStyleStates is a type grouping ItemStyle structs for different states of a list item
type ItemStyleStates struct {
   Normal ItemStyle
   Selected ItemStyle
}

// Styles is a type grouping lipgloss styles for a list and its items, to facilitate applying different styles
// to different states of the list (e.g. focussed, etc.)
type Styles struct {
   List lg.Style
   Item ItemStyleStates
   Seperator lg.Style
}

// DefaultStyle returns a new Styles struct with the given styles.
func DefaultStyles() (unfocussed Styles, focussed Styles) {
   unfocussedItemNormal := ItemStyle {
      Main: lg.NewStyle().
         Foreground(lg.Color("#bac2de")),
      Prefix: lg.NewStyle(),
      Suffix: lg.NewStyle(),
   }
   unfocussedItemSelected := ItemStyle {
      Main: unfocussedItemNormal.Main.
         Bold(true).
         Foreground(lg.Color("#f5e0dc")),
      Prefix: unfocussedItemNormal.Prefix.
         Foreground(lg.Color("#f9e2af")),
      Suffix: unfocussedItemNormal.Suffix,
   }
   focussedItemNormal := ItemStyle {
      Main: unfocussedItemNormal.Main.
         Foreground(lg.Color("#cdd6f4")),
      Prefix: unfocussedItemNormal.Prefix,
      Suffix: unfocussedItemNormal.Suffix,
   }
   focussedItemSelected := ItemStyle {
      Main: focussedItemNormal.Main.
         Foreground(lg.Color("#f9e2af")),
      Prefix: focussedItemNormal.Prefix.
         Foreground(lg.Color("#fab387")),
      Suffix: focussedItemNormal.Suffix,
   }
   unfocussed = Styles {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Normal: unfocussedItemNormal,
         Selected: unfocussedItemSelected,
      },
   }
   focussed = Styles {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Normal: focussedItemNormal,
         Selected: focussedItemSelected,
      },
   }
   return
}

// DefaultItemStyle returns a new ItemStyleStates struct with the given styles.

