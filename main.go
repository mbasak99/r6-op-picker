package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type operator struct {
	Name   string `json:"name"`
	ImgURL string `json:"icon_url"`
	Side   string `json:"side"`
}

// Model of what the program is going to work with
type model struct {
	Cursor           int
	AvailableOptions []string
	Operators        []operator
	PlayerOpMap      map[string]string // Each player will have an operator assigned to them
	Side             string            // Team's current side, side is either Attacker or Defender
}

// Initialize all the ops to later choose from
func initModel() model {
	// Fetch all the operators from both sides
	resp, err := http.Get("https://r6-api.vercel.app/api/operators")
	if err != nil {
		log.Fatalf("Failed to retrieve operators from API: %+v\n", err)
	}
	defer resp.Body.Close()

	rawOperators, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read operator data from API: %+v", err)
	}

	// Store convert JSON to Go struct
	var operatorsJson []operator
	err = json.Unmarshal([]byte(rawOperators), &operatorsJson)
	if err != nil {
		log.Fatalf("Failed to parse operator JSON: %+v\n", err)
	}

	// Only have the operators initialized
	// Get the other info when the user starts interacting with the terminal
	return model{
		Operators:        operatorsJson,
		Cursor:           0,
		AvailableOptions: []string{"Enter players in the team", "Select current side", "Reset "},
	}
}

func (m model) basicInformationSet() bool {
	userFilledNeededInfo := false

	// check player map for at least one player
	for k, v := range m.PlayerOpMap {
		fmt.Println(k, v)
		if k != "" {
			userFilledNeededInfo = true
		}
	}

	return userFilledNeededInfo
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
		case "up":
			if m.Cursor > 0 {
				m.Cursor++
			}
		case "down":
			if m.Cursor < len(m.AvailableOptions)-1 {
				m.Cursor--
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Please choose one of the options below:\n\n"

	for i, v := range m.AvailableOptions {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %d) %s\n", cursor, i+1, v)
	}

	s += "\n\nPress q to quit.\n"
	return s + "\n"
}

func main() {
	program := tea.NewProgram(initModel())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Something went wrong with the program: %+v\n", err)
	}
}
