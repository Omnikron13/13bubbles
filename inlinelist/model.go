package inlinelist

import (
   "fmt"
   "strings"

   bt "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/bubbles/key"
)


// CachedItem is a struct that holds _unstyled_ rendered strings for the prefix, item, and suffix of an item,
// and the _styled_ rendered string.
type CachedItem[T any] struct {
   item *T
   suffix, main, prefix, listFocussedItemUnfocussed, listUnfocussedItemUnfocussed, listFocussedItemFocussed, listUnfocussedItemFocussed string
}


// Model is a widget that displays a list of items that flow horizontally as a paragraph, joined by a
// separator, with optional prefix and/or suffix for each item. Additionally it supports focussing items with a
// 'cursor'.
// TODO: add interface requirement for generic type T
type Model[T any] struct {
   Items []item[T]

   // These functions control how an item of type T should be rendered, with optional prefix and suffix which can have
   // their own styles, and excluded from any filtering, sorting, etc. of the list.
   RenderItem func(T) string
   RenderPrefix func(T) string
   RenderSuffix func(T) string

   // Separator that will be rendered between items
   separator string

   // Whether the list can take focus to allow selecting items, and a pointer to the focussed item (or nil if no item
   // is currently focussed).
   Focussable bool
   focussedItem *item[T]

   // Whether the list has focus, which will enable or disable keybindings, possibly change styles, etc.
   focussed bool

   // Customisable keybindings
   KeyBindings KeyMap

   // Customisable styling using lipgloss
   Styles StyleStates
}


// findIndex returns the index in the (displayed) list of the given item, or -1 if not found.
// This is intended to keep track of the focussed item if/when the list is changed, e.g. by filtering, sorting, etc.
func (m *Model[T]) findIndex(item *item[T]) int {
   for i := range m.Items {
      if &m.Items[i] == item {
         return i
      }
   }
   return -1
}


// New creates a new Model with the given items and options.
func New[T any](items ...T) (m *Model[T]) {
   m = &Model[T] {
      separator: ", ",
      RenderItem: func (item T) string { return fmt.Sprintf("%v", item) },
      KeyBindings: defaultKeyMap(),
      Styles: DefaultStyles(),
   }
   for _, i := range items {
      m.Items = append(m.Items, *newitem(i, m))
   }
   return
}


// Focus flags the list as having focus; enabling keybindings, possibly changing styles, etc.
func (m *Model[T]) Focus() {
   m.focussed = true
}


// GetFocussed returns a pointer to the focussed item, or nil if nothing is focussed (yet?), the focussed item is now
// gone (e.g. if the list has been filtered), or the list is not flagged as focusable in the first place.
func (m *Model[T]) GetFocussed() *item[T] {
   if !m.Focussable || m.findIndex(m.focussedItem) < 0 {
      return nil
   }
   return m.focussedItem
}


// Init initializes the Model; part of the bubbletea Model interface.
func (m *Model[T]) Init() (cmd bt.Cmd) {
   for _, i := range m.Items {
      i.renderCache()
   }
   return
}


// IsFocussed returns whether the Model has focus.
func (m *Model[T]) IsFocussed() bool {
   return m.focussed
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
         if m.Focussable {
            // TODO: deal with defocusing more elegantly?
            if m.focussedItem != nil { m.focussedItem.focussed = false }

            i := m.findIndex(m.focussedItem)
            switch {
               case key.Matches(msg, m.KeyBindings.Next):
                  if i < len(m.Items)-1 {
                     m.focussedItem = &m.Items[i+1]
                  }
               case key.Matches(msg, m.KeyBindings.Prev):
                  if i > 0 {
                     m.focussedItem = &m.Items[i-1]
                  } else if i < 0 {
                     m.focussedItem = &m.Items[len(m.Items)-1]
                  }
               case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
                  m.focussed = !m.focussed
            }

            if m.focussedItem != nil { m.focussedItem.focussed = true }
         }
   }

   model = m
   return
}


// View renders the Model; part of the bubbletea Model interface.
func (m *Model[T]) View() string {
   var sb strings.Builder

   var styles Style
   if m.focussed {
      styles = m.Styles.Focussed
   } else {
      styles = m.Styles.Unfocussed
   }

   for n, i := range m.Items {
      sb.WriteString(i.render())

      if n < len(m.Items)-1 {
         sb.WriteString(styles.Seperator.Render(m.separator))
      }
   }

   return m.Styles.Focussed.List.Render(sb.String())
}

