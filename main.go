package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type exercise struct {
	name  string
	sets  string
	tempo string
	notes string
}

type block struct {
	name      string
	duration  string
	exercises []exercise
}

var pullDay = []block{
	{
		name:     "Morning Activation",
		duration: "5-7 min",
		exercises: []exercise{
			{name: "Dead hangs from pull-up bar", sets: "2 x 20-30 sec"},
			{name: "Band or light DB face pulls", sets: "2 x 15"},
			{name: "Thoracic rotations", sets: "1 min each side"},
		},
	},
	{
		name:     "Strength Block A - Back Focus",
		duration: "10-15 min",
		exercises: []exercise{
			{name: "DB Bent Over Rows", sets: "3 x 10-12 each", tempo: "2 sec pull, 3 sec lower", notes: "Go heavy (35-50 lbs)"},
			{name: "Pull-ups or Assisted Pull-ups", sets: "3 x 6-10", tempo: "Controlled", notes: "Use band for assist if needed"},
			{name: "DB Pullover", sets: "3 x 12", tempo: "Slow stretch at bottom", notes: "Great for lats + serratus"},
		},
	},
	{
		name:     "Strength Block B - Biceps & Rear Delts",
		duration: "10-15 min",
		exercises: []exercise{
			{name: "Incline DB Curls", sets: "3 x 10-12", tempo: "3 sec down", notes: "Maximum stretch at bottom"},
			{name: "Hammer Curls", sets: "3 x 12", tempo: "Slow and strict", notes: "Forearms + brachialis"},
			{name: "Reverse Flyes (bent over)", sets: "3 x 15", tempo: "Squeeze at top", notes: "Light weight, rear delts"},
			{name: "Concentration Curls", sets: "2 x 10 each", tempo: "Squeeze at top", notes: "Finish with a pump"},
		},
	},
	{
		name:     "End of Day Flexibility",
		duration: "5-10 min",
		exercises: []exercise{
			{name: "Lat stretch on pull-up bar", sets: "30 sec each side"},
			{name: "Supine twist", sets: "1 min each side"},
			{name: "Forearm stretches", sets: "30 sec each"},
			{name: "Neck stretches", sets: "30 sec each direction"},
		},
	},
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1)

	blockStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	exerciseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

type model struct {
	blockIndex    int
	exerciseIndex int
	done          bool
}

func initialModel() model {
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
		case "enter", " ", "n":
			return m.nextExercise(), nil
		case "p", "b":
			return m.prevExercise(), nil
		case "s":
			return m.skipBlock(), nil
		}
	}
	return m, nil
}

func (m model) nextExercise() model {
	block := pullDay[m.blockIndex]
	if m.exerciseIndex < len(block.exercises)-1 {
		m.exerciseIndex++
	} else if m.blockIndex < len(pullDay)-1 {
		m.blockIndex++
		m.exerciseIndex = 0
	} else {
		m.done = true
	}
	return m
}

func (m model) prevExercise() model {
	if m.exerciseIndex > 0 {
		m.exerciseIndex--
	} else if m.blockIndex > 0 {
		m.blockIndex--
		m.exerciseIndex = len(pullDay[m.blockIndex].exercises) - 1
	}
	return m
}

func (m model) skipBlock() model {
	if m.blockIndex < len(pullDay)-1 {
		m.blockIndex++
		m.exerciseIndex = 0
	} else {
		m.done = true
	}
	return m
}

func (m model) View() string {
	if m.done {
		return titleStyle.Render("PULL DAY Complete!") + "\n\n" +
			exerciseStyle.Render("Great workout! Your back and biceps thank you.") + "\n\n" +
			helpStyle.Render("Press q to quit")
	}

	block := pullDay[m.blockIndex]
	ex := block.exercises[m.exerciseIndex]

	s := titleStyle.Render("PULL DAY") + "\n"
	s += blockStyle.Render(fmt.Sprintf("%s (%s)", block.name, block.duration)) + "\n"
	s += dimStyle.Render(fmt.Sprintf("Block %d/%d • Exercise %d/%d",
		m.blockIndex+1, len(pullDay),
		m.exerciseIndex+1, len(block.exercises))) + "\n\n"

	s += highlightStyle.Render(ex.name) + "\n"
	s += exerciseStyle.Render("Sets: "+ex.sets) + "\n"

	if ex.tempo != "" {
		s += exerciseStyle.Render("Tempo: "+ex.tempo) + "\n"
	}
	if ex.notes != "" {
		s += dimStyle.Render(ex.notes) + "\n"
	}

	s += helpStyle.Render("\n[enter/n] next • [p] previous • [s] skip block • [q] quit")

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
