package gomongocrud

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	Client *mongo.Client
}

func (repo *MongoRepo) Store(ctx context.Context, task *Task) error {
	if _, err := tasksCollection(repo.Client).InsertOne(ctx, task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (repo *MongoRepo) FetchAll(ctx context.Context) ([]Task, error) {
	cursor, err := tasksCollection(repo.Client).Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var tasks []Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, errors.WithStack(err)
	}
	return tasks, nil
}

func (repo *MongoRepo) GetByID(ctx context.Context, id string) (*Task, error) {
	var task Task
	if err := tasksCollection(repo.Client).FindOne(ctx, bson.M{"_id": id}).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrRecordNotFound
		}
		return nil, errors.WithStack(err)
	}
	return &task, nil
}

func (repo *MongoRepo) Update(ctx context.Context, task *Task) error {
	res, err := tasksCollection(repo.Client).ReplaceOne(ctx, bson.M{"_id": task.ID}, task)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.ModifiedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (repo *MongoRepo) Delete(ctx context.Context, id string) error {
	res, err := tasksCollection(repo.Client).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.WithStack(err)
	}

	if res.DeletedCount == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func tasksCollection(c *mongo.Client) *mongo.Collection {
	return c.Database("tasks").Collection("tasks")
}
