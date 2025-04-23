package vehicles

import (
	"api-search/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	vs models.VehicleStore
}

func NewHandler(vs models.VehicleStore) *Handler {
	return &Handler{vs: vs}
}

func (h *Handler) RegiterRoutes(e *echo.Echo) {
	e.GET("/vehicles", h.List)
	e.GET("/vehicles/:id", h.Get)
	e.POST("/vehicles", h.Create)
	e.PUT("/vehicles/:id", h.Update)
	e.DELETE("/vehicles/:id", h.Delete)
}

func (h *Handler) List(c echo.Context) error {
	search := c.QueryParam("search")
	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil {
		size = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	order := c.QueryParam("order")
	if order == "" {
		order = "created_at:desc"
	}
	filter := c.QueryParam("filter")

	res, err := h.vs.List(size, page, search, order, filter)
	if err != nil {
		slog.Error("failed to retrieve vehicles", "error", err)
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	v, err := h.vs.Get(id)
	if err != nil {
		slog.Error("failed to retrieve a single vehicle", "error", err)
		return err
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) Create(c echo.Context) error {
	vw := &models.VehicleWrite{}
	if err := c.Bind(vw); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	if err := vw.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	v, err := h.vs.Create(vw)
	if err != nil {
		slog.Error("failed to create a single vehicle", "error", err)
		return err
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	vw := &models.VehicleWrite{}
	if err := c.Bind(vw); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	if err := vw.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	v, err := h.vs.Update(id, vw)
	if err != nil {
		slog.Error("failed to update vehicle", "error", err)
		return err
	}
	return c.JSON(http.StatusOK, v)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	_, err := h.vs.Delete(id)
	if err != nil {
		slog.Error("failed to delete vehicle", "error", err)
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}
