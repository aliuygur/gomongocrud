package gomongocrud

import (
	"context"
	"time"
)

type tasksRepo interface {
	Store(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	FetchAll(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

type idGenerator interface {
	Generate(ctx context.Context) string
}

type TasksService struct {
	Tasks tasksRepo
	IDs   idGenerator
}

func (s *TasksService) Store(ctx context.Context, task *Task) error {
	now := time.Now()
	id := s.IDs.Generate(ctx)

	task.ID = id
	task.UpdatedAt = now
	task.CreatedAt = now

	return s.Tasks.Store(ctx, task)
}

func (s *TasksService) GetByID(ctx context.Context, id string) (*Task, error) {
	return s.Tasks.GetByID(ctx, id)
}

func (s *TasksService) FetchAll(ctx context.Context) ([]Task, error) {
	return s.Tasks.FetchAll(ctx)
}

func (s *TasksService) Update(ctx context.Context, task *Task) error {
	return s.Tasks.Update(ctx, task)
}

func (s *TasksService) Delete(ctx context.Context, id string) error {
	return s.Tasks.Delete(ctx, id)
}
