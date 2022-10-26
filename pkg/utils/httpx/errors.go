package httpx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type ForbiddenError struct {
}

func (e ForbiddenError) Error() string {
	return "access denied"
}

//

func ParameterError(name string, err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid parameter: %s, %v", name, err))
}

func BindingError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid json: %v", err))
}

func LogicError(err error) *echo.HTTPError {
	switch {
	case errors.As(err, &ForbiddenError{}):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	return echo.NewHTTPError(http.StatusConflict, err.Error())
}
