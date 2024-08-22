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
   Focussed ItemStyle
}

// Style is a type grouping lipgloss styles for a list and its items, to facilitate applying different styles
// to different states of the list (e.g. focussed, etc.)
type Style struct {
   List lg.Style
   Item ItemStyleStates
   Seperator lg.Style
}

// StyleStates is a type grouping Style structs for different states of a list
type StyleStates struct {
   Unfocussed Style
   Focussed Style
}

// DefaultStyles returns a new StyleStates struct with the given styles.
func DefaultStyles() StyleStates {
   unfocussedItemNormal := ItemStyle {
      Main: lg.NewStyle().
         Foreground(lg.Color("#bac2de")),
      Prefix: lg.NewStyle(),
      Suffix: lg.NewStyle(),
   }
   unfocussedItemFocussed := ItemStyle {
      Main: unfocussedItemNormal.Main.
         Foreground(lg.Color("#f5e0dc")),
      Prefix: unfocussedItemNormal.Prefix.
         Foreground(lg.Color("#f9e2af")),
      Suffix: unfocussedItemNormal.Suffix.
         Foreground(lg.Color("#99d1db")),
   }
   focussedItemNormal := ItemStyle {
      Main: unfocussedItemNormal.Main.
         Foreground(lg.Color("#cdd6f4")),
      Prefix: unfocussedItemNormal.Prefix.
         Foreground(lg.Color("#B9957F")),
      Suffix: unfocussedItemNormal.Suffix.
         Foreground(lg.Color("#99d1db")),
   }
   focussedItemFocussed := ItemStyle {
      Main: focussedItemNormal.Main.
         Foreground(lg.Color("#f9e2af")).
         Bold(true),
      Prefix: focussedItemNormal.Prefix.
         Foreground(lg.Color("#fab387")).
         Bold(true),
      Suffix: focussedItemNormal.Suffix.
         Foreground(lg.Color("#89dceb")),
   }
   unfocussed := Style {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Normal: unfocussedItemNormal,
         Focussed: unfocussedItemFocussed,
      },
   }
   focussed := Style {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Normal: focussedItemNormal,
         Focussed: focussedItemFocussed,
      },
   }
   return StyleStates{
      Unfocussed: unfocussed,
      Focussed:   focussed,
   }
}

// DefaultItemStyle returns a new ItemStyleStates struct with the given styles.

