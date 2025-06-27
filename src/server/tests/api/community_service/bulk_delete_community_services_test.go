package community_service_test

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

func TestBulkDeleteCommunityServicesSuccessfully(t *testing.T) {
	/*
		GIVEN: Community-service associations exist
		WHEN:  DELETE /community-service/bulk-delete/ is called with a list of associations
		THEN:  The associations are deleted and a HTTP_204_NO_CONTENT status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	communityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	communityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	request := schemas.BulkDeleteCommunityServiceRequest{
		CommunityServices: []*schemas.DeleteCommunityServiceRequest{
			{CommunityId: communityService1.CommunityId, ServiceId: communityService1.ServiceId},
			{CommunityId: communityService2.CommunityId, ServiceId: communityService2.ServiceId},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community-service/bulk-delete/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	var count int64
	db.Model(&model.CommunityService{}).Where("community_id IN ? AND service_id IN ?",
		[]string{communityService1.CommunityId.String(), communityService2.CommunityId.String()},
		[]string{communityService1.ServiceId.String(), communityService2.ServiceId.String()},
	).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestBulkDeleteCommunityServicesInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /community-service/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": "json"}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/community-service/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
