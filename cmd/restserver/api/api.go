package api

import (
	"errors"
	"net/http"

	"github.com/alioygur/gomongocrud"
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate types -o types.gen.go oas3.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate server,spec -o server.gen.go oas3.yaml

type Handler struct {
	Tasks *gomongocrud.TasksService
}

func (h *Handler) FindTasks(ctx echo.Context) error {
	tasks, err := h.Tasks.FetchAll(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// bind domain model to response model
	tt := make([]*Task, 0, len(tasks))
	for _, r := range tasks {
		tt = append(tt, newTask(&r))
	}

	return ctx.JSON(http.StatusOK, tt)
}

func (h *Handler) AddTask(ctx echo.Context) error {
	var newtask NewTask
	if err := ctx.Bind(&newtask); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var task gomongocrud.Task
	task.Name = newtask.Name
	task.Description = newtask.Description

	if err := h.Tasks.Store(ctx.Request().Context(), &task); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// send 201 Created http status code with created task object
	return ctx.JSON(http.StatusCreated, newTask(&task))
}

func (h *Handler) DeleteTask(ctx echo.Context, id string) error {
	if err := h.Tasks.Delete(ctx.Request().Context(), id); err != nil {
		// don't send 500 http status code on client error
		if errors.Is(err, gomongocrud.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// send 204 no content http status code
	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) FindTaskById(ctx echo.Context, id string) error {
	task, err := h.Tasks.GetByID(ctx.Request().Context(), id)
	if err != nil {
		// don't send 500 http status code on client error
		if errors.Is(err, gomongocrud.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, newTask(task))
}

// newTask creates view model from domain model
func newTask(t *gomongocrud.Task) *Task {
	var task Task
	task.Id = t.ID
	task.Name = t.Name
	task.Description = t.Description
	task.CreatedAt = t.CreatedAt
	task.UpdatedAt = t.UpdatedAt
	return &task
}
