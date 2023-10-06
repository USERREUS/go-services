package mongostore

import (
	"notification-service/internal/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	table = "inventory"
)

type Notification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MsgType     string             `bson:"message_type"`
	Description string             `bson:"description"`
	Data        string             `bson:"data,omitempty"`
}

type Repository struct {
	store *Store
}

func (r *Repository) Create(m *model.Model) error {
	temp := Notification{
		MsgType:     m.MsgType,
		Data:        time.Now().String(),
		Description: m.Description,
	}

	_, err := r.store.collection.InsertOne(r.store.context, temp)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindOne(id string) (*model.Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var result Notification
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &model.Model{
		MsgType:     result.MsgType,
		Data:        result.Data,
		Description: result.Description,
	}, nil
}

func (r *Repository) FindAll() (map[string]*model.Model, error) {
	records := make(map[string]*model.Model)

	filter := bson.M{}
	cursor, err := r.store.collection.Find(r.store.context, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.store.context)

	var results []Notification
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records[data.ID.Hex()] = &model.Model{
			MsgType:     data.MsgType,
			Description: data.Description,
			Data:        data.Data,
		}
	}

	return records, nil
}
