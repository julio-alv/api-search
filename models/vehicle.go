package models

import (
	"errors"

	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VehicleStore interface {
	List(size, page int, search, order, filter string) (*meilisearch.SearchResponse, error)
	Get(id string) (*Vehicle, error)
	Create(vw *VehicleWrite) (*Vehicle, error)
	Update(id string, vw *VehicleWrite) (*Vehicle, error)
	Delete(id string) (string, error)
}

type Vehicle struct {
	ID         bson.ObjectID `json:"id" bson:"_id"`
	CreatedAt  int64         `json:"created_at" bson:"created_at"`
	UpdatedAt  int64         `json:"updated_at" bson:"updated_at"`
	Name       string        `json:"name" bson:"name"`
	Color      string        `json:"color" bson:"color"`
	Seats      int           `json:"seats" bson:"seats"`
	Features   []string      `json:"features" bson:"features"`
	HP         int           `json:"horse_power" bson:"horse_power"`
	Torque     int           `json:"torque" bson:"torque"`
	DriveTrain string        `json:"drive_train" bson:"drive_train"`
}

type VehicleWrite struct {
	Name       string   `json:"name" bson:"name"`
	Color      string   `json:"color" bson:"color"`
	Seats      int      `json:"seats" bson:"seats"`
	Features   []string `json:"features,omitempty" bson:"features,omitempty"`
	HP         int      `json:"horse_power" bson:"horse_power"`
	Torque     int      `json:"torque" bson:"torque"`
	DriveTrain string   `json:"drive_train" bson:"drive_train"`
	CreatedAt  int64    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  int64    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (vw *VehicleWrite) Validate() error {
	if vw.Name == "" {
		return errors.New("name is required")
	}
	if vw.Color == "" {
		return errors.New("color is required")
	}
	if vw.Seats <= 0 {
		return errors.New("seats must be greater than 0")
	}
	if vw.HP <= 0 {
		return errors.New("horse power must be greater than 0")
	}
	if vw.Torque <= 0 {
		return errors.New("torque must be greater than 0")
	}
	if vw.DriveTrain == "" {
		return errors.New("drive train is required")
	}
	return nil
}
