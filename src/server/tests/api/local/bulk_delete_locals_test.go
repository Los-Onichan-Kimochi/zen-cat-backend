package local_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple locals exist in the database
		WHEN:  DELETE /local/bulk-delete/ is called with a list of local IDs
		THEN:  The locals are deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})
	request := schemas.BulkDeleteLocalRequest{
		Locals: []string{local1.Id.String(), local2.Id.String()},
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/local/bulk-delete/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.Local{}).Where("id IN ?", []string{local1.Id.String(), local2.Id.String()}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestBulkDeleteLocalsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /local/bulk-delete/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/local/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
