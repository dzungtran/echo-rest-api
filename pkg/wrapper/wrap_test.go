package wrapper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTranslate_SuccessResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data:   map[string]string{"key": "value"},
		Status: http.StatusOK,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, map[string]interface{}{"key": "value"}, resp["data"])
	assert.Nil(t, resp["message"])
	assert.Nil(t, resp["metadata"])
}

func TestTranslate_SuccessWithMetadata(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data:         []string{"item1", "item2"},
		Status:       http.StatusOK,
		Total:        100,
		IncludeTotal: true,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, []interface{}{"item1", "item2"}, resp["data"])
	assert.NotNil(t, resp["metadata"])
	metadata := resp["metadata"].(map[string]interface{})
	assert.Equal(t, float64(100), metadata["total"])
}

func TestTranslate_ErrorResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Error:  errors.New("something went wrong"),
		Status: http.StatusInternalServerError,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "something went wrong", resp["message"])
	assert.Nil(t, resp["data"])
}

func TestTranslate_NotFoundResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Error:  errors.New("resource not found"),
		Status: http.StatusNotFound,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "resource not found", resp["message"])
}

func TestTranslate_ForbiddenResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Error:  errors.New("access denied"),
		Status: http.StatusForbidden,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "access denied", resp["message"])
}

func TestTranslate_CreatedResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data:   map[string]int{"id": 123},
		Status: http.StatusCreated,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, map[string]interface{}{"id": float64(123)}, resp["data"])
}

func TestTranslate_NilData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data:   nil,
		Status: http.StatusOK,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Nil(t, resp["data"])
}

func TestTranslate_IncludeTotalFalse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data:         []int{1, 2, 3},
		Status:       http.StatusOK,
		Total:        50,
		IncludeTotal: false,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, []interface{}{float64(1), float64(2), float64(3)}, resp["data"])
	assert.Nil(t, resp["metadata"])
}

func TestTranslate_DefaultStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Data: "some data",
		// Status not set
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTranslate_ValidationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	res := Response{
		Error:  errors.New("invalid input: name is required"),
		Status: http.StatusBadRequest,
	}

	err := Translate(c, res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "invalid input: name is required", resp["message"])
}
