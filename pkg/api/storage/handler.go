package storage

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"github.com/gavrilaf/wardrobe/pkg/utils/httpx"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

func Assemble(root *echo.Group, m Manager) {
	h := &handler{
		manager: m,
	}

	g := root.Group("/info_objects")
	{
		g.GET("/:id", h.getObject)

		g.POST("", h.createObject)
		g.POST("/:id/files", h.addFile)
		g.PUT("/:id/finilize", h.finilizeObject)

	}
}

type handler struct {
	manager Manager
}

func (h *handler) createObject(c echo.Context) error {
	ctx := c.Request().Context()

	var obj dto.InfoObject
	if err := c.Bind(&obj); err != nil {
		return httpx.BindingError(err)
	}

	id, err := h.manager.CreateInfoObject(ctx, obj)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": id})
}

func (h *handler) addFile(c echo.Context) error {
	ctx := c.Request().Context()

	infoObjectID, err := strconv.Atoi(c.Param("id"))
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

	fileMeta := dto.File{
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
	}

	fileID, err := h.manager.AddFile(ctx, infoObjectID, fileMeta, file)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": fileID})
}

func (h *handler) finilizeObject(c echo.Context) error {
	ctx := c.Request().Context()

	infoObjectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.ParameterError("id", err)
	}

	if err := h.manager.FinilizeInfoObject(ctx, infoObjectID); err != nil {
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

	obj, err := h.manager.GetInfoObject(ctx, id)
	if err != nil {
		return httpx.LogicError(err)
	}

	return c.JSON(http.StatusOK, obj)
}
