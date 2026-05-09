package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
)

func main() {
	if _, err := tea.NewProgram(newModel(false)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
