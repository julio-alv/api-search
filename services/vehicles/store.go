package vehicles

import (
	"api-search/models"
	"context"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store struct {
	m          *mongo.Collection
	sm         meilisearch.ServiceManager
	updateChan chan models.UpdateMessage
}

func NewStore(m *mongo.Collection, sm meilisearch.ServiceManager, updateChan chan models.UpdateMessage) *Store {
	return &Store{
		m:          m,
		sm:         sm,
		updateChan: updateChan,
	}
}

func (s *Store) List(size, page int, search, order, filter string) (*meilisearch.SearchResponse, error) {
	return s.sm.Index("vehicles").Search(search, &meilisearch.SearchRequest{
		Limit:  int64(size),
		Offset: int64((page - 1) * size),
		Sort:   []string{order},
		Filter: filter,
	})
}

func (s *Store) Get(id string) (*models.Vehicle, error) {
	v := &models.Vehicle{}

	objID, _ := bson.ObjectIDFromHex(id)
	ctx := context.Background()

	return v, s.m.FindOne(ctx, bson.M{"_id": objID}).Decode(v)
}

func (s *Store) Create(vw *models.VehicleWrite) (*models.Vehicle, error) {
	v := &models.Vehicle{}

	vw.CreatedAt = time.Now().UTC().Unix()
	vw.UpdatedAt = time.Now().UTC().Unix()

	ctx := context.Background()
	res, err := s.m.InsertOne(ctx, vw)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID
	if err := s.m.FindOne(ctx, bson.M{"_id": id}).Decode(v); err != nil {
		return nil, err
	}

	s.updateChan <- models.UpdateMessage{
		Type:    "create",
		Payload: v,
		ID:      nil,
	}

	return v, nil
}

func (s *Store) Update(id string, vw *models.VehicleWrite) (*models.Vehicle, error) {
	v := &models.Vehicle{}

	vw.UpdatedAt = time.Now().UTC().Unix()

	objID, _ := bson.ObjectIDFromHex(id)
	ctx := context.Background()

	if _, err := s.m.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": vw}); err != nil {
		return nil, err
	}

	if err := s.m.FindOne(ctx, bson.M{"_id": objID}).Decode(v); err != nil {
		return nil, err
	}

	s.updateChan <- models.UpdateMessage{
		Type:    "update",
		Payload: v,
		ID:      &id,
	}

	return v, nil
}

func (s *Store) Delete(id string) (string, error) {
	objID, _ := bson.ObjectIDFromHex(id)
	ctx := context.Background()
	if _, err := s.m.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		return "", err
	}
	s.updateChan <- models.UpdateMessage{
		Type:    "delete",
		Payload: nil,
		ID:      &id,
	}
	return id, nil
}
