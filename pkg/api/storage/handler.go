package storage

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/gavrilaf/wardrobe/pkg/api/dto"
	"github.com/gavrilaf/wardrobe/pkg/utils/httpx"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

func Assemble(root *echo.Group, m Manager) {
	h := &handler{
		manager: m,
	}

	g := root.Group("/storage")
	{
		g.POST("", h.createObject)
		g.PUT("/:id", h.uploadContent)

		g.GET("/:id", h.getObject)
		g.GET("/:id/meta", h.getObjectMeta)
	}
}

type handler struct {
	manager Manager
}

func (h *handler) createObject(c echo.Context) error {
	ctx := c.Request().Context()

	var fo dto.FO
	if err := c.Bind(&fo); err != nil {
		return httpx.BindingError(err)
	}

	id, err := h.manager.CreateObject(ctx, fo)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": id})
}

func (h *handler) uploadContent(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.ParameterError("id", err)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return httpx.ParameterError("file", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return httpx.ParameterError("file", fmt.Errorf("failed to open file, %w", err))
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.WithError(log.FromContext(ctx), closeErr).Error("can't close file stream")
		}
	}()

	err = h.manager.UploadContent(ctx, id, file, fileHeader.Size)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) getObject(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.ParameterError("id", err)
	}

	fo, err := h.manager.GetObject(ctx, id)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.JSON(http.StatusOK, fo)
}

func (h *handler) getObjectMeta(c echo.Context) error {
	return nil
}
