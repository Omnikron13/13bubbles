package inlinelist

import (
   "fmt"
   "strings"
   "unicode/utf8"

   bt "github.com/charmbracelet/bubbletea"
   lg "github.com/charmbracelet/lipgloss"
   "github.com/charmbracelet/bubbles/key"
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



// CachedItem is a struct that holds _unstyled_ rendered strings for the prefix, item, and suffix of an item,
// and the _styled_ rendered string.
type CachedItem[T any] struct {
   item *T
   suffix, main, prefix, focussedNormal, unfocussedNormal, focussedSelected, unfocussedSelected string
}


// Model is a widget that displays a list of items that flow horizontally as a paragraph, joined by a
// separator, with optional prefix and/or suffix for each item. Additionally if supports selecting items with a
// 'cursor'.
// TODO: add interface requirement for generic type T
type Model[T any] struct {
   Items []T

   // itemRenderCache is a cache of the strings the render functions return, operating on the assumption that they are
   // supposed to always return the same string for a given item (i.e. pure functions).
   // TODO: add a method to explicitly rebuild the cache for situations where the render functions are not pure.
   itemRenderCache []CachedItem[T]
   itemRenderCacheChannel chan *CachedItem[T]

   // These functions control how an item of type T should be rendered, with optional prefix and suffix which can have
   // their own styles, and excluded from any filtering, sorting, etc. of the list.
   RenderItem func(T) string
   RenderPrefix func(T) string
   RenderSuffix func(T) string

   // Separator that will be rendered between items
   separator string

   // Whether the list can take focus to allow selecting items, and a pointer to the selected item (or nil if no item
   // is currently selected).
   Selectable bool
   selected *T

   // Whether the list has focus, which will enable or disable keybindings, possibly change styles, etc.
   Focussed bool

   // Customisable keybindings
   KeyBindings KeyMap

   // Customisable styling using lipgloss
   Styles struct {
      Focussed Styles
      Unfocussed Styles
   }
}


// findIndex returns the index in the (displayed) list of the given item, or -1 if not found.
// This is intended to keep track of the selected item if/when the list is changed, e.g. by filtering, sorting, etc.
func (m *Model[T]) findIndex(item *T) int {
   for i := range m.Items {
      if &m.Items[i] == item {
         return i
      }
   }
   return -1
}


// New creates a new Model with the given items and options.
func New[T any](items ...T) (m Model[T]) {
   m = Model[T] {
      Items: items,
      separator: ", ",
      RenderItem: func (item T) string { return fmt.Sprintf("%v", item) },
      KeyBindings: defaultKeyMap(),
   }
   return
}


// GetSelected returns a pointer to the selected item, or nil if nothing is selected (yet?), the selected item is now
// gone (e.g. if the list has been filtered), or the list is not flagged as selectable in the first place.
func (m *Model[T]) GetSelected() *T {
   if !m.Selectable || m.findIndex(m.selected) < 0 {
      return nil
   }
   return m.selected
}


// Init initializes the Model; part of the bubbletea Model interface.
func (m *Model[T]) Init() (cmd bt.Cmd) {
   m.itemRenderCache = make([]CachedItem[T], 0, len(m.Items))
   m.itemRenderCacheChannel = make(chan *CachedItem[T], len(m.Items))
   go func() {
      for i := range m.Items {
         item := &m.Items[i]
         c := CachedItem[T] {
            item: item,
            main: m.RenderItem(*item),
         }
         if m.RenderPrefix != nil {
            c.prefix = m.RenderPrefix(*item)
         }
         if m.RenderSuffix != nil {
            c.suffix = m.RenderSuffix(*item)
         }
         c.unfocussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Normal.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.focussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Normal.Main.Render(c.main),
            m.Styles.Focussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.unfocussedSelected = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Selected.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Selected.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Selected.Suffix.Render(c.suffix),
         )
         c.focussedSelected = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Selected.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Selected.Main.Render(c.main),
            m.Styles.Focussed.Item.Selected.Suffix.Render(c.suffix),
         )
         m.itemRenderCacheChannel <- &c
      }
      close(m.itemRenderCacheChannel)
   }()

   return
}


// itemToString converts an item of type T to a string, using RenderPrefix(), RenderItem(), and RenderSuffix()
// functions (if set), returning a styled string and the unstyled rune length.
func (m *Model[T]) itemToString(item *T, style ItemStyleStates) (string, int) {
   var sb strings.Builder
   var n int

   if m.RenderPrefix != nil {
      s := m.RenderPrefix(*item)
      n += utf8.RuneCountInString(s)
      sb.WriteString(style.Normal.Prefix.Render(s))
   }

   s := m.RenderItem(*item)
   n += utf8.RuneCountInString(s)
   sb.WriteString(style.Normal.Main.Render(s))

   if m.RenderSuffix != nil {
      s := m.RenderSuffix(*item)
      n += utf8.RuneCountInString(s)
      sb.WriteString(style.Normal.Suffix.Render(s))
   }

   return sb.String(), n
}


// Update updates the Model; part of the bubbletea Model interface.
func (m *Model[T]) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   m.KeyBindings.Next.SetEnabled(m.Focussed)
   m.KeyBindings.Prev.SetEnabled(m.Focussed)

   switch msg := msg.(type) {
      case bt.KeyMsg:
         switch msg.String() {
            case "ctrl+c":
               cmd = bt.Quit
               return
         }
         if m.Selectable {
            i := m.findIndex(m.selected)
            switch {
               case key.Matches(msg, m.KeyBindings.Next):
                  if i < len(m.Items)-1 {
                     m.selected = &m.Items[i+1]
                  }
               case key.Matches(msg, m.KeyBindings.Prev):
                  if i > 0 {
                     m.selected = &m.Items[i-1]
                  } else if i < 0 {
                     m.selected = &m.Items[len(m.Items)-1]
                  }
               case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
                  m.Focussed = !m.Focussed
            }
         }
   }

   model = m
   return
}


// View renders the Model; part of the bubbletea Model interface.
func (m *Model[T]) View() string {
   for c := range m.itemRenderCacheChannel {
      m.itemRenderCache = append(m.itemRenderCache, *c)
   }

   var sb strings.Builder

   var styles Styles
   if m.Focussed {
      styles = m.Styles.Focussed
   } else {
      styles = m.Styles.Unfocussed
   }

   for i, c := range m.itemRenderCache {
      if m.Selectable && c.item == m.selected {
         if m.Focussed {
            sb.WriteString(c.focussedSelected)
         } else {
            sb.WriteString(c.unfocussedSelected)
         }
      } else {
         if m.Focussed {
            sb.WriteString(c.focussedNormal)
         } else {
            sb.WriteString(c.unfocussedNormal)
         }
      }

      if i < len(m.Items)-1 {
         sb.WriteString(styles.Seperator.Render(m.separator))
      }
   }

   return styles.List.Render(sb.String())
}

