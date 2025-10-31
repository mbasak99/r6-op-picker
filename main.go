package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// global var because I don't wanna keep hitting the API endpoint and get rate-limited
// var operators map[string][]string

// Model of what the program is going to work with
type model struct {
	operators   map[string][]string
	playerOpMap map[string]string // Each player will have an operator assigned to them
	side        string            // Side is either Attacker or Defender
}

// Initialize all the ops to later choose from
func initModel() model {
	// Fetch all the operators from both sides
	resp, err := http.Get("https://r6-api.vercel.app/api/operators")
	if err != nil {
		log.Fatalf("Failed to retrieve operators from API: %+v\n", err)
	}
	defer resp.Body.Close()

	ops, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read operator data from API: %+v", err)
	}

	fmt.Printf("%+s", ops)
	// Sort operator by side

	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "TODO\n"
}

func main() {
	program := tea.NewProgram(initModel())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Something went wrong with the program: %+v\n", err)
	}
}
