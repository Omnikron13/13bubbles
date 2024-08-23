package inlinelist

import (
   "fmt"
   "strings"
   "unicode/utf8"

   bt "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/bubbles/key"
)


// CachedItem is a struct that holds _unstyled_ rendered strings for the prefix, item, and suffix of an item,
// and the _styled_ rendered string.
type CachedItem[T any] struct {
   item *T
   suffix, main, prefix, listFocussedNormal, listUnfocussedNormal, listFocussedItemFocussed, listUnfocussedItemFocussed string
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
   focussed bool

   // Customisable keybindings
   KeyBindings KeyMap

   // Customisable styling using lipgloss
   Styles StyleStates
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
      Styles: DefaultStyles(),
   }
   return
}


// Focus flags the list as having focus; enabling keybindings, possibly changing styles, etc.
func (m *Model[T]) Focus() {
   m.focussed = true
}


// GetFocussed returns a pointer to the selected item, or nil if nothing is selected (yet?), the selected item is now
// gone (e.g. if the list has been filtered), or the list is not flagged as selectable in the first place.
func (m *Model[T]) GetFocussed() *T {
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
         c.listUnfocussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Normal.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.listFocussedNormal = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Normal.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Normal.Main.Render(c.main),
            m.Styles.Focussed.Item.Normal.Suffix.Render(c.suffix),
         )
         c.listUnfocussedItemFocussed = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Unfocussed.Item.Focussed.Prefix.Render(c.prefix),
            m.Styles.Unfocussed.Item.Focussed.Main.Render(c.main),
            m.Styles.Unfocussed.Item.Focussed.Suffix.Render(c.suffix),
         )
         c.listFocussedItemFocussed = fmt.Sprintf(
            "%s%s%s",
            m.Styles.Focussed.Item.Focussed.Prefix.Render(c.prefix),
            m.Styles.Focussed.Item.Focussed.Main.Render(c.main),
            m.Styles.Focussed.Item.Focussed.Suffix.Render(c.suffix),
         )
         m.itemRenderCacheChannel <- &c
      }
      close(m.itemRenderCacheChannel)
   }()

   return
}


// IsFocussed returns whether the Model has focus.
func (m *Model[T]) IsFocussed() bool {
   return m.focussed
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


// Unfocus flags the list as not having focus; disabling keybindings, possibly changing styles, etc.
func (m *Model[T]) Unfocus() {
   m.focussed = false
}


// Update updates the Model; part of the bubbletea Model interface.
func (m *Model[T]) Update(msg bt.Msg) (model bt.Model, cmd bt.Cmd) {
   m.KeyBindings.Next.SetEnabled(m.focussed)
   m.KeyBindings.Prev.SetEnabled(m.focussed)

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
                  m.focussed = !m.focussed
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

   var styles Style
   if m.focussed {
      styles = m.Styles.Focussed
   } else {
      styles = m.Styles.Unfocussed
   }

   for i, c := range m.itemRenderCache {
      if m.Selectable && c.item == m.selected {
         if m.focussed {
            sb.WriteString(c.listFocussedItemFocussed)
         } else {
            sb.WriteString(c.listUnfocussedItemFocussed)
         }
      } else {
         if m.focussed {
            sb.WriteString(c.listFocussedNormal)
         } else {
            sb.WriteString(c.listUnfocussedNormal)
         }
      }

      if i < len(m.Items)-1 {
         sb.WriteString(styles.Seperator.Render(m.separator))
      }
   }

   return styles.List.Render(sb.String())
}

