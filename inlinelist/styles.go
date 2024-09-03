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
   Unfocussed ItemStyle
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
   unfocussedListUnfocussedItem := ItemStyle {
      Main: lg.NewStyle().
         Foreground(lg.Color("#bac2de")),
      Prefix: lg.NewStyle(),
      Suffix: lg.NewStyle(),
   }
   unfocussedListFocussedItem := ItemStyle {
      Main: unfocussedListUnfocussedItem.Main.
         Foreground(lg.Color("#f5e0dc")),
      Prefix: unfocussedListUnfocussedItem.Prefix.
         Foreground(lg.Color("#f9e2af")),
      Suffix: unfocussedListUnfocussedItem.Suffix.
         Foreground(lg.Color("#99d1db")),
   }
   focussedListUnfocussedItem := ItemStyle {
      Main: unfocussedListUnfocussedItem.Main.
         Foreground(lg.Color("#cdd6f4")),
      Prefix: unfocussedListUnfocussedItem.Prefix.
         Foreground(lg.Color("#B9957F")),
      Suffix: unfocussedListUnfocussedItem.Suffix.
         Foreground(lg.Color("#99d1db")),
   }
   focussedListFocussedItem := ItemStyle {
      Main: focussedListUnfocussedItem.Main.
         Foreground(lg.Color("#f9e2af")),
      Prefix: focussedListUnfocussedItem.Prefix.
         Foreground(lg.Color("#fab387")),
      Suffix: focussedListUnfocussedItem.Suffix.
         Foreground(lg.Color("#89dceb")),
   }
   unfocussedList := Style {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Unfocussed: unfocussedListUnfocussedItem,
         Focussed: unfocussedListFocussedItem,
      },
   }
   focussedList := Style {
      List: lg.NewStyle(),
      Item: ItemStyleStates {
         Unfocussed: focussedListUnfocussedItem,
         Focussed: focussedListFocussedItem,
      },
   }
   return StyleStates{
      Unfocussed: unfocussedList,
      Focussed:   focussedList,
   }
}

// DefaultItemStyle returns a new ItemStyleStates struct with the given styles.

