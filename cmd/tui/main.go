package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if _, err := tea.NewProgram(newModel(false)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
