package local_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A local exists in the database
		WHEN:  DELETE /local/{localId}/ is called
		THEN:  The local is deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/local/"+local.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.Local{}).Where("id = ?", local.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteLocalNotFound(t *testing.T) {
	/*
		GIVEN: A local with the given ID does not exist
		WHEN:  DELETE /local/{localId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/local/"+nonExistentId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteLocalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  DELETE /local/{localId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/local/invalid-uuid/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
