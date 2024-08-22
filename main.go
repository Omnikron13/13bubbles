package main

import (
   "os"

   bt "github.com/charmbracelet/bubbletea"

   "github.com/omnikron13/13bubbles/inlinelist"
)

func main() {
   l := inlinelist.New([]string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"})

   l.Selectable = true
   l.Focussed = true

   p := bt.NewProgram(&l, bt.WithAltScreen())

   if m, err := p.Run(); err != nil { bt.ErrProgramKilled.Error() } else {
      _ = m
      os.Exit(0)
   }
}

