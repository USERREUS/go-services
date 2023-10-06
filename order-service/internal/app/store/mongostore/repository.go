package mongostore

import (
	"order-service/internal/app/model"
	"order-service/internal/app/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Date      string             `bson:"data,omitempty"`
	OrderItem *OrderItem         `bson:"order_item"`
}

type OrderItem struct {
	ProductCode string `bson:"product_code"`
	Name        string `bson:"name"`
	Count       int    `bson:"count"`
	Cost        int    `bson:"cost"`
}

type Repository struct {
	store *Store
}

func (r *Repository) Create(m *model.Model) error {
	item := &OrderItem{
		ProductCode: m.ProductCode,
		Name:        m.Name,
		Count:       m.Count,
		Cost:        m.Cost,
	}

	//FIX
	order := Order{
		Date:      time.Now().String(),
		OrderItem: item,
	}

	_, err := r.store.collection.InsertOne(r.store.context, order)
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

	var result Order
	err = r.store.collection.FindOne(r.store.context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &model.Model{
		ProductCode: result.OrderItem.ProductCode,
		Name:        result.OrderItem.Name,
		Count:       result.OrderItem.Count,
		Cost:        result.OrderItem.Cost,
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

	var results []Order
	if err := cursor.All(r.store.context, &results); err != nil {
		return nil, err
	}

	for _, data := range results {
		records[data.ID.Hex()] = &model.Model{
			ProductCode: data.OrderItem.ProductCode,
			Name:        data.OrderItem.Name,
			Count:       data.OrderItem.Count,
			Cost:        data.OrderItem.Cost,
		}
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
