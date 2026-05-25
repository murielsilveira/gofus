package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/murielsilveira/gofus/internal/platform/db"
	"github.com/murielsilveira/gofus/internal/platform/server"
)

var (
	testPool        *pgxpool.Pool
	skipIntegration bool
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}

	ctx := context.Background()

	testURL := db.TestDatabaseURL()

	pool, err := db.Connect(ctx, testURL)
	if err != nil {
		skipIntegration = true
		os.Exit(m.Run())
	}
	testPool = pool

	if err := runMigrations(testURL); err != nil {
		pool.Close()
		panic(err)
	}

	code := m.Run()

	pool.Close()
	os.Exit(code)
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func requireDB(t *testing.T) {
	t.Helper()
	if skipIntegration {
		t.Skip("database not available")
	}
}

func resetDB(t *testing.T) {
	t.Helper()
	requireDB(t)

	_, err := testPool.Exec(context.Background(), "TRUNCATE boards RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func newTestApp(t *testing.T) *fiber.App {
	t.Helper()
	requireDB(t)
	return server.NewWithPool(testPool)
}

func setup(t *testing.T) *fiber.App {
	t.Helper()
	resetDB(t)
	return newTestApp(t)
}

func request(t *testing.T, app *fiber.App, method, path string, body any) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req)
	require.NoError(t, err)

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(dst))
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(data)
}

type errorResponse struct {
	Error string `json:"error"`
}

func assertError(t *testing.T, resp *http.Response, status int, message string) {
	t.Helper()
	require.Equal(t, status, resp.StatusCode)

	var body errorResponse
	decodeJSON(t, resp, &body)
	require.Equal(t, message, body.Error)
}
