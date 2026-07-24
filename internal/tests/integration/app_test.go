package integration

import (
	"techzone/internal/app"
	"testing"
)

func mustNewTestApp(t *testing.T) *app.App {
	t.Helper()
	requireIntegration(t)

	application, err := app.NewServer(true)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(application.Close)

	return application
}
