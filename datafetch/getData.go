package datafetch

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func GetData(c echo.Context) error {
	result := c.Get("result")
	data, ok := result.(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	storedData, _ := data["data"].(string)

	return c.JSON(
		http.StatusFound, map[string]string{
			"data": storedData,
		},
	)
}
