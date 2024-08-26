package main

import (
   "os"

   bt "github.com/charmbracelet/bubbletea"

   "github.com/omnikron13/13bubbles/inlinelist"
)

func main() {
   l := inlinelist.New("one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten")

   l.Focussable = true
   l.Focus()

   l.RenderPrefix = func(s string) string {
      switch s {
         case "one":
            return "㊀"
         case "two":
            return "㊁"
         case "three":
            return "㊂"
         case "four":
            return "㊃"
         case "five":
            return "㊄"
         case "six":
            return "㊅"
         case "seven":
            return "㊆"
         case "eight":
            return "㊇"
         case "nine":
            return "㊈"
         case "ten":
            return "㊉"
         default:
            return ""
      }
   }

   l.RenderSuffix = func(s string) string {
      switch s {
         case "one":
            return " Ⅰ"
         case "two":
            return " Ⅱ"
         case "three":
            return " Ⅲ"
         case "four":
            return " Ⅳ"
         case "five":
            return " Ⅴ"
         case "six":
            return " Ⅵ"
         case "seven":
            return " Ⅶ"
         case "eight":
            return " Ⅷ"
         case "nine":
            return " Ⅸ"
         case "ten":
            return " Ⅹ"
         default:
            return ""
      }
   }

   p := bt.NewProgram(&l, bt.WithAltScreen())

   if m, err := p.Run(); err != nil { bt.ErrProgramKilled.Error() } else {
      _ = m
      os.Exit(0)
   }
}

