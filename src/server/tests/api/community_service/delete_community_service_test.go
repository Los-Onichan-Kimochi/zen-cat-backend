package community_service_test

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

func TestDeleteCommunityServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A community-service association exists
		WHEN:  DELETE /community-service/{communityId}/{serviceId}/ is called
		THEN:  The association is deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	url := "/community-service/" + communityService.CommunityId.String() + "/" + communityService.ServiceId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.CommunityService{}).Where("community_id = ? AND service_id = ?", communityService.CommunityId, communityService.ServiceId).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteCommunityServiceNotFound(t *testing.T) {
	/*
		GIVEN: A community-service association does not exist
		WHEN:  DELETE /community-service/{communityId}/{serviceId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	communityId := uuid.New()
	serviceId := uuid.New()

	// WHEN
	url := "/community-service/" + communityId.String() + "/" + serviceId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCommunityServiceInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided for communityId or serviceId
		WHEN:  DELETE /community-service/{communityId}/{serviceId}/ is called
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// WHEN
	url := "/community-service/invalid-uuid/" + communityService.ServiceId.String() + "/"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
