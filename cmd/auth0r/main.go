package main

import (
	"context"
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/javid-gulamaliyev/auth0r/auth0"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	applicationsTable   table.Model
	applicationsSpinner spinner.Model
	viewportWidth       int
	isLoading           bool
}

type appLoadedMsg struct {
	apps []auth0.Application
	err  error
}

func fetchApps() tea.Cmd {
	return func() tea.Msg {
		apps, err := auth0.ListApplications(context.Background())
		return appLoadedMsg{
			apps: apps,
			err:  err,
		}
	}
}

func makeDynamicColumns(width int) []table.Column {
	w := width - 14 // subtract border (2) + per-column padding (6 cols × 2)
	return []table.Column{
		{Title: "Name", Width: w * 20 / 100},
		{Title: "Client ID", Width: w * 20 / 100},
		{Title: "Description", Width: w * 25 / 100},
		{Title: "First Party", Width: w * 10 / 100},
		{Title: "Type", Width: w * 15 / 100},
		{Title: "Callbacks", Width: w * 10 / 100},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.applicationsSpinner.Tick,
		fetchApps(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var spinnerCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.applicationsTable.Focused() {
				m.applicationsTable.Blur()
			} else {
				m.applicationsTable.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.applicationsTable.SelectedRow()[1]),
			)
		}

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.applicationsTable.SetWidth(msg.Width - 2)
		m.applicationsTable.SetColumns(makeDynamicColumns(msg.Width))

	case appLoadedMsg:
		if msg.err != nil {
			return m, tea.Quit
		}
		rows := make([]table.Row, len(msg.apps))
		for i, app := range msg.apps {
			rows[i] = []string{
				app.Name,
				app.ClientID,
				app.Description,
				fmt.Sprintf("%t", app.IsFirstParty),
				app.ApplicationType,
				fmt.Sprintf("%d", len(app.Callbacks)),
			}
		}
		m.isLoading = false
		m.applicationsTable.SetRows(rows)
	}

	m.applicationsTable, cmd = m.applicationsTable.Update(msg)
	m.applicationsSpinner, spinnerCmd = m.applicationsSpinner.Update(msg)
	return m, tea.Batch(cmd, spinnerCmd)
}

func (m model) View() tea.View {
	if m.isLoading {
		return tea.NewView(baseStyle.Render(m.applicationsSpinner.View()))
	}
	return tea.NewView(baseStyle.Render(m.applicationsTable.View()) + "\n  " + m.applicationsTable.HelpView() + "\n")
}

func main() {
	columns := makeDynamicColumns(80)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(80),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	spinner := spinner.New()
	spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := model{applicationsTable: t, isLoading: true, applicationsSpinner: spinner}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
