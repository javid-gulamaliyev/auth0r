package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Application struct {
	ClientID        string   `json:"client_id"`
	Name            string   `json:"name"`
	ApplicationType string   `json:"app_type"`
	Description     string   `json:"description"`
	IsFirstParty    bool     `json:"first_party"`
	Callbacks       []string `json:"callbacks"`
	// TODO: add more fields
}

func ListApplications(ctx context.Context) ([]Application, error) {
	cmd := exec.Command("auth0", "apps", "list", "--json")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	var apps []Application
	if err := json.Unmarshal(output, &apps); err != nil {
		return nil, fmt.Errorf("failed to parse applications: %w", err)
	}

	return apps, nil
}
