package api

import (
	"api-search/models"
	"api-search/services/vehicles"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type API struct {
	m  *mongo.Database
	sm meilisearch.ServiceManager
}

func NewAPI(m *mongo.Database, sm meilisearch.ServiceManager) *API {
	return &API{
		m:  m,
		sm: sm,
	}
}

func (a *API) Start(address string) error {
	e := echo.New()

	updateChan := make(chan models.UpdateMessage, 100)

	vehicleStore := vehicles.NewStore(a.m.Collection("vehicles"), a.sm, updateChan)
	vehicleHandler := vehicles.NewHandler(vehicleStore)
	vehicleHandler.RegiterRoutes(e)

	// This is how you would update the "sortable/filterable" attributes
	attributes := []string{
		"name",
		"color",
		"seats",
		"torque",
		"horse_power",
		"drive_train",
		"created_at",
		"updated_at",
	}
	if _, err := a.sm.Index("vehicles").UpdateSortableAttributes(&attributes); err != nil {
		return err
	}
	if _, err := a.sm.Index("vehicles").UpdateFilterableAttributes(&attributes); err != nil {
		return err
	}

	// setup another thread to listen for updates
	go func() {
		for {
			msg := <-updateChan

			switch msg.Type {
			case "create":
				if _, err := a.sm.Index("vehicles").AddDocuments(msg.Payload, "id"); err != nil {
					slog.Error(err.Error())
				}
			case "update":
				if _, err := a.sm.Index("vehicles").UpdateDocuments(msg.Payload, "id"); err != nil {
					slog.Error(err.Error())
				}
			case "delete":
				if _, err := a.sm.Index("vehicles").DeleteDocument(*msg.ID); err != nil {
					slog.Error(err.Error())
				}
			}
		}
	}()

	return e.Start(address)
}
