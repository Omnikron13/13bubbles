package inlinelist

import (
   "fmt"
   "strings"
)

// cachedItem is a struct that stores a styled rendered string representation of every possible item state along with
// each string width in graphemes/columns, and likely in future a normalised form of just the unstyled 'item' string.
type cachedItem struct {
   focussedList struct {
      focussedItem struct {
         render string
         width int
      }
      unfocussedItem struct {
         render string
         width int
      }
   }
   unfocussedList struct {
      focussedItem struct {
         render string
         width int
      }
      unfocussedItem struct {
         render string
         width int
      }
   }
}


// item is a struct that wraps the type T objects that are stored in the list, encapsulating them with additional
// information that is used for rendering/caching, flagging, etc.
type item[T any] struct {
   Object T

   focussed bool
   cache *cachedItem
   cacheChannel chan *cachedItem
   list *Model[T]
}


// newitem creates a new item with the given object.
// The item spawn off a goroutine to render the item in its various possible states into the cache on creation. An
// important implication of this is that the item should be able to expect there to either be a populated cache or a
// channel to extract and populate the cache from at any given time so that it can render without blocking.
func newitem[T any](object T, list *Model[T]) *item[T] {
   i := &item[T] {
      Object: object,
      list: list,
      //cacheChannel: make(chan *cachedItem, 1),
   }

   return i
}


// render returns the styled string repreentation of the item using RenderItem(), as well as RenderPrefix() and/or
// RenderSuffix() if they are set.
// It wiil save a reference to the rendered string and its length in graphemes/columns internally, integrating the
// previously external cache structures.
func (i *item[T]) render() string {
   if i.cache == nil {
      select {
         case i.cache = <-i.cacheChannel:
            _ = 0
         default:
            i.renderCache()
      }
   }

   switch {
      case i.list.Focussable && i.list.focussed:
         switch {
            case i.focussed:
               return i.cache.focussedList.focussedItem.render
            case !i.focussed:
               return i.cache.focussedList.unfocussedItem.render
         }
      case !i.list.Focussable || !i.list.focussed:
         switch {
            case i.focussed:
               return i.cache.unfocussedList.focussedItem.render
            case !i.focussed:
               return i.cache.unfocussedList.unfocussedItem.render
         }
   }

   panic(fmt.Errorf("Error: unreachable codepath reached"))
}


// renderCache renders out a fresh cache of styled strings for each possiible item state.
func (i *item[T]) renderCache() {
   i.cache = nil

   // TODO: tidy this mess up with additional types, methods on cache types, etc.
   c := &cachedItem{}
   c.focussedList.focussedItem.render, c.focussedList.focussedItem.width = i.renderStyle(&i.list.Styles.Focussed.Item.Focussed)
   c.focussedList.unfocussedItem.render, c.focussedList.unfocussedItem.width = i.renderStyle(&i.list.Styles.Focussed.Item.Unfocussed)
   c.unfocussedList.focussedItem.render, c.unfocussedList.focussedItem.width = i.renderStyle(&i.list.Styles.Unfocussed.Item.Focussed)
   c.unfocussedList.unfocussedItem.render, c.unfocussedList.unfocussedItem.width = i.renderStyle(&i.list.Styles.Unfocussed.Item.Unfocussed)

   i.cache = c
}


// renderStyle is the function that actually renders the styled string representation of the item using RenderItem(),
// as well as RenderPrefix() and/or RenderSuffix() if they are set.
// It is meant to be deferred to by more specific render functions (which may be being deferred to be more general
// render functions. Such fun!)
// Along with the rendered string, it also returns the width of that string in graphemes/columns.
func (i *item[T]) renderStyle(style *ItemStyle) (str string, n int) {
   var sb strings.Builder

   if i.list.RenderPrefix != nil {
      s := noBreak(i.list.RenderPrefix(i.Object))
      n += countGraphemes(s)
      sb.WriteString(style.Prefix.Render(s))
   }

   s := i.list.RenderItem(i.Object)
   n += countGraphemes(s)
   sb.WriteString(style.Main.Render(s))

   if i.list.RenderSuffix != nil {
      s := i.list.RenderSuffix(i.Object)
      n += countGraphemes(s)
      sb.WriteString(style.Suffix.Render(s))
   }

   str = sb.String()

   return
}

