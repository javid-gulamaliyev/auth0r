package main

import (
	"context"
	"fmt"

	"github.com/javid-gulamaliyev/auth0r/auth0"
)

func main() {
	println(" === scratch ====")
	apps, err := auth0.ListApplications(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, app := range apps {
		fmt.Printf("Name: %s, Client ID: %s\n", app.Name, app.ClientID)
	}
}
