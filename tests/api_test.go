package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/murielsilveira/gofus/internal/db/sqlc"
)

func TestDemoRoutes(t *testing.T) {
	app := newTestApp(t)

	t.Run("GET /", func(t *testing.T) {
		resp := request(t, app, http.MethodGet, "/", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Hello, world!", readBody(t, resp))
	})

	t.Run("GET /app.html", func(t *testing.T) {
		resp := request(t, app, http.MethodGet, "/app.html", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})

	t.Run("GET /db", func(t *testing.T) {
		resp := request(t, app, http.MethodGet, "/db", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "WORKED!!", readBody(t, resp))
	})
}

func TestBoards(t *testing.T) {
	t.Run("POST /api/v1/boards creates board", func(t *testing.T) {
		app := setup(t)
		resp := request(t, app, http.MethodPost, "/api/v1/boards", map[string]string{
			"name": "Sprint Board",
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var board sqlc.Board
		decodeJSON(t, resp, &board)
		require.NotEmpty(t, board.ID)
		require.Equal(t, "Sprint Board", board.Name)
		require.False(t, board.CreatedAt.IsZero())
		require.False(t, board.UpdatedAt.IsZero())
	})

	t.Run("POST /api/v1/boards rejects empty name", func(t *testing.T) {
		app := setup(t)
		resp := request(t, app, http.MethodPost, "/api/v1/boards", map[string]string{
			"name": "",
		})
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("POST /api/v1/boards rejects invalid JSON", func(t *testing.T) {
		app := setup(t)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/boards", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("GET /api/v1/boards lists boards", func(t *testing.T) {
		app := setup(t)
		created := createBoard(t, app, "List Board")

		resp := request(t, app, http.MethodGet, "/api/v1/boards", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var boards []sqlc.Board
		decodeJSON(t, resp, &boards)
		require.Len(t, boards, 1)
		require.Equal(t, created.ID, boards[0].ID)
	})

	t.Run("GET /api/v1/boards/:id returns board", func(t *testing.T) {
		app := setup(t)
		created := createBoard(t, app, "Get Board")

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/boards/%s", created.ID), nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var board sqlc.Board
		decodeJSON(t, resp, &board)
		require.Equal(t, created.ID, board.ID)
		require.Equal(t, "Get Board", board.Name)
	})

	t.Run("GET /api/v1/boards/:id invalid UUID", func(t *testing.T) {
		app := setup(t)
		resp := request(t, app, http.MethodGet, "/api/v1/boards/not-a-uuid", nil)
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("GET /api/v1/boards/:id not found", func(t *testing.T) {
		app := setup(t)
		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/boards/%s", uuid.New()), nil)
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("PATCH /api/v1/boards/:id updates board", func(t *testing.T) {
		app := setup(t)
		created := createBoard(t, app, "Old Name")

		resp := request(t, app, http.MethodPatch, fmt.Sprintf("/api/v1/boards/%s", created.ID), map[string]string{
			"name": "New Name",
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var board sqlc.Board
		decodeJSON(t, resp, &board)
		require.Equal(t, "New Name", board.Name)
	})

	t.Run("PATCH /api/v1/boards/:id rejects empty body", func(t *testing.T) {
		app := setup(t)
		created := createBoard(t, app, "Patch Board")

		resp := request(t, app, http.MethodPatch, fmt.Sprintf("/api/v1/boards/%s", created.ID), map[string]any{})
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("DELETE /api/v1/boards/:id deletes board", func(t *testing.T) {
		app := setup(t)
		created := createBoard(t, app, "Delete Board")

		resp := request(t, app, http.MethodDelete, fmt.Sprintf("/api/v1/boards/%s", created.ID), nil)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()

		getResp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/boards/%s", created.ID), nil)
		assertError(t, getResp, http.StatusNotFound, "not found")
	})
}

func TestColumns(t *testing.T) {
	t.Run("POST /api/v1/boards/:boardID/columns creates column", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")

		position := int32(1)
		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/boards/%s/columns", board.ID), map[string]any{
			"name":     "Todo",
			"position": position,
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var column sqlc.Column
		decodeJSON(t, resp, &column)
		require.NotEmpty(t, column.ID)
		require.Equal(t, board.ID, column.BoardID)
		require.Equal(t, "Todo", column.Name)
		require.Equal(t, int32(1), column.Position)
	})

	t.Run("POST /api/v1/boards/:boardID/columns missing board", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/boards/%s/columns", uuid.New()), map[string]string{
			"name": "Todo",
		})
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("POST /api/v1/boards/:boardID/columns rejects empty name", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")

		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/boards/%s/columns", board.ID), map[string]string{
			"name": "",
		})
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("GET /api/v1/boards/:boardID/columns lists columns", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")
		created := createColumn(t, app, board.ID, "In Progress")

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/boards/%s/columns", board.ID), nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var columns []sqlc.Column
		decodeJSON(t, resp, &columns)
		require.Len(t, columns, 1)
		require.Equal(t, created.ID, columns[0].ID)
	})

	t.Run("GET /api/v1/boards/:boardID/columns missing board", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/boards/%s/columns", uuid.New()), nil)
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("GET /api/v1/columns/:id returns column", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")
		created := createColumn(t, app, board.ID, "Done")

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s", created.ID), nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var column sqlc.Column
		decodeJSON(t, resp, &column)
		require.Equal(t, created.ID, column.ID)
	})

	t.Run("GET /api/v1/columns/:id not found", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s", uuid.New()), nil)
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("PATCH /api/v1/columns/:id updates column", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")
		created := createColumn(t, app, board.ID, "Review")

		position := int32(3)
		resp := request(t, app, http.MethodPatch, fmt.Sprintf("/api/v1/columns/%s", created.ID), map[string]any{
			"name":     "QA",
			"position": position,
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var column sqlc.Column
		decodeJSON(t, resp, &column)
		require.Equal(t, "QA", column.Name)
		require.Equal(t, int32(3), column.Position)
	})

	t.Run("DELETE /api/v1/columns/:id deletes column", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Column Board")
		created := createColumn(t, app, board.ID, "Archive")

		resp := request(t, app, http.MethodDelete, fmt.Sprintf("/api/v1/columns/%s", created.ID), nil)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()

		getResp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s", created.ID), nil)
		assertError(t, getResp, http.StatusNotFound, "not found")
	})
}

func TestTasks(t *testing.T) {
	t.Run("POST /api/v1/columns/:columnID/tasks creates task", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")

		position := int32(2)
		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/columns/%s/tasks", column.ID), map[string]any{
			"title":       "Write tests",
			"description": "Cover all endpoints",
			"position":    position,
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var task sqlc.Task
		decodeJSON(t, resp, &task)
		require.NotEmpty(t, task.ID)
		require.Equal(t, column.ID, task.ColumnID)
		require.Equal(t, "Write tests", task.Title)
		require.Equal(t, "Cover all endpoints", task.Description)
		require.Equal(t, int32(2), task.Position)
	})

	t.Run("POST /api/v1/columns/:columnID/tasks missing column", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/columns/%s/tasks", uuid.New()), map[string]string{
			"title": "Orphan task",
		})
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("POST /api/v1/columns/:columnID/tasks rejects empty title", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")

		resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/columns/%s/tasks", column.ID), map[string]string{
			"title": "",
		})
		assertError(t, resp, http.StatusBadRequest, "bad request")
	})

	t.Run("GET /api/v1/columns/:columnID/tasks lists tasks", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")
		created := createTask(t, app, column.ID, "List task")

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s/tasks", column.ID), nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var tasks []sqlc.Task
		decodeJSON(t, resp, &tasks)
		require.Len(t, tasks, 1)
		require.Equal(t, created.ID, tasks[0].ID)
	})

	t.Run("GET /api/v1/columns/:columnID/tasks missing column", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s/tasks", uuid.New()), nil)
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("GET /api/v1/tasks/:id returns task", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")
		created := createTask(t, app, column.ID, "Get task")

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/tasks/%s", created.ID), nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var task sqlc.Task
		decodeJSON(t, resp, &task)
		require.Equal(t, created.ID, task.ID)
	})

	t.Run("GET /api/v1/tasks/:id not found", func(t *testing.T) {
		app := setup(t)

		resp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/tasks/%s", uuid.New()), nil)
		assertError(t, resp, http.StatusNotFound, "not found")
	})

	t.Run("PATCH /api/v1/tasks/:id updates task and moves column", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")
		otherColumn := createColumn(t, app, board.ID, "Done")
		created := createTask(t, app, column.ID, "Move me")

		resp := request(t, app, http.MethodPatch, fmt.Sprintf("/api/v1/tasks/%s", created.ID), map[string]any{
			"column_id": otherColumn.ID,
			"title":     "Moved task",
		})
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var task sqlc.Task
		decodeJSON(t, resp, &task)
		require.Equal(t, otherColumn.ID, task.ColumnID)
		require.Equal(t, "Moved task", task.Title)
	})

	t.Run("DELETE /api/v1/tasks/:id deletes task", func(t *testing.T) {
		app := setup(t)
		board := createBoard(t, app, "Task Board")
		column := createColumn(t, app, board.ID, "Todo")
		created := createTask(t, app, column.ID, "Delete me")

		resp := request(t, app, http.MethodDelete, fmt.Sprintf("/api/v1/tasks/%s", created.ID), nil)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp.Body.Close()

		getResp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/tasks/%s", created.ID), nil)
		assertError(t, getResp, http.StatusNotFound, "not found")
	})
}

func TestAPI_FullKanbanFlow(t *testing.T) {
	app := setup(t)

	board := createBoard(t, app, "Flow Board")
	todo := createColumn(t, app, board.ID, "Todo")
	done := createColumn(t, app, board.ID, "Done")

	task := createTask(t, app, todo.ID, "Ship feature")

	moveResp := request(t, app, http.MethodPatch, fmt.Sprintf("/api/v1/tasks/%s", task.ID), map[string]any{
		"column_id": done.ID,
	})
	require.Equal(t, http.StatusOK, moveResp.StatusCode)

	var moved sqlc.Task
	decodeJSON(t, moveResp, &moved)
	require.Equal(t, done.ID, moved.ColumnID)

	listResp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s/tasks", done.ID), nil)
	var doneTasks []sqlc.Task
	decodeJSON(t, listResp, &doneTasks)
	require.Len(t, doneTasks, 1)
	require.Equal(t, task.ID, doneTasks[0].ID)

	deleteResp := request(t, app, http.MethodDelete, fmt.Sprintf("/api/v1/boards/%s", board.ID), nil)
	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
	deleteResp.Body.Close()

	columnResp := request(t, app, http.MethodGet, fmt.Sprintf("/api/v1/columns/%s", todo.ID), nil)
	assertError(t, columnResp, http.StatusNotFound, "not found")
}

func createBoard(t *testing.T, app *fiber.App, name string) sqlc.Board {
	t.Helper()

	resp := request(t, app, http.MethodPost, "/api/v1/boards", map[string]string{"name": name})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var board sqlc.Board
	decodeJSON(t, resp, &board)
	return board
}

func createColumn(t *testing.T, app *fiber.App, boardID uuid.UUID, name string) sqlc.Column {
	t.Helper()

	resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/boards/%s/columns", boardID), map[string]string{
		"name": name,
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var column sqlc.Column
	decodeJSON(t, resp, &column)
	return column
}

func createTask(t *testing.T, app *fiber.App, columnID uuid.UUID, title string) sqlc.Task {
	t.Helper()

	resp := request(t, app, http.MethodPost, fmt.Sprintf("/api/v1/columns/%s/tasks", columnID), map[string]string{
		"title": title,
	})
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var task sqlc.Task
	decodeJSON(t, resp, &task)
	return task
}
