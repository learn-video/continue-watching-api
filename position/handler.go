package position

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type (
	Handler struct {
		r *redis.Client
	}
)

func NewHandler(r *redis.Client) *Handler {
	return &Handler{r: r}
}

func (h *Handler) Record(c echo.Context) error {
	validate := validator.New()
	pos := new(Position)
	if err := c.Bind(pos); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(pos); err != nil {
		return err
	}

	if err := Record(h.r, pos.UserID, pos.VideoID, pos.Position); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Fetch(c echo.Context) error {
	userID := c.QueryParam("user_id")
	videoID := c.QueryParam("video_id")
	pos, err := Fetch(h.r, userID, videoID)
	if err == ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, &Position{VideoID: videoID, Position: pos})
}
