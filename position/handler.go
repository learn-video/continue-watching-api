package position

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Position struct {
	VideoID  string `json:"video_id" validate:"required"`
	Position int    `json:"position" validate:"required"`
}

type PositionDetail struct {
	Position int `json:"position"`
}

type (
	Handler struct {
		r *redis.Client
	}
)

func NewHandler(r *redis.Client) *Handler {
	return &Handler{r: r}
}

func (h *Handler) Record(c echo.Context) error {
	userID, err := c.Cookie("user_id")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	validate := validator.New()
	pos := new(Position)
	if err := c.Bind(pos); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := validate.Struct(pos); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := Record(h.r, userID.Value, pos.VideoID, pos.Position); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Fetch(c echo.Context) error {
	userID, err := c.Cookie("user_id")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	videoID := c.QueryParam("video_id")
	pos, err := Fetch(h.r, userID.Value, videoID)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, &PositionDetail{Position: pos})
}
