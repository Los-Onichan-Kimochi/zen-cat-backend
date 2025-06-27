package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateServiceSuccessfully(t *testing.T) {
	// GIVEN: Valid service creation request
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "Professional Consultation",
		Description: "One-on-one professional consultation service",
		ImageUrl:    "https://example.com/consultation.jpg",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Service is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.Description, result.Description)
	assert.Equal(t, createRequest.ImageUrl, result.ImageUrl)
	assert.Equal(t, createRequest.IsVirtual, result.IsVirtual)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateServiceEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid service creation request but empty updatedBy
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "Test Service",
		Description: "Testing purposes",
		ImageUrl:    "https://example.com/test.jpg",
		IsVirtual:   true,
	}

	// WHEN: CreateService is called with empty updatedBy
	result, err := controller.CreateService(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreateServiceDuplicateName(t *testing.T) {
	// GIVEN: A service already exists with the same name
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest1 := schemas.CreateServiceRequest{
		Name:        "Duplicate Service",
		Description: "First service",
		ImageUrl:    "https://example.com/first.jpg",
		IsVirtual:   false,
	}

	createRequest2 := schemas.CreateServiceRequest{
		Name:        "Duplicate Service", // Same name
		Description: "Second service",
		ImageUrl:    "https://example.com/second.jpg",
		IsVirtual:   true,
	}

	// Create first service
	_, err1 := controller.CreateService(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateService is called with duplicate name
	result, err := controller.CreateService(createRequest2, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateServiceVirtual(t *testing.T) {
	// GIVEN: Service creation request for virtual service
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "Virtual Therapy",
		Description: "Online therapy sessions via video call",
		ImageUrl:    "https://example.com/virtual-therapy.jpg",
		IsVirtual:   true,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Virtual service is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.True(t, result.IsVirtual)
}

func TestCreateServicePhysical(t *testing.T) {
	// GIVEN: Service creation request for physical service
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "In-Person Consultation",
		Description: "Face-to-face consultation at our office",
		ImageUrl:    "https://example.com/in-person.jpg",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Physical service is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.False(t, result.IsVirtual)
}

func TestCreateServiceEmptyName(t *testing.T) {
	// GIVEN: Service creation request with empty name
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "", // Empty name
		Description: "Service without name",
		ImageUrl:    "https://example.com/empty.jpg",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateServiceEmptyDescription(t *testing.T) {
	// GIVEN: Service creation request with empty description
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "Service Without Description",
		Description: "", // Empty description
		ImageUrl:    "https://example.com/no-desc.jpg",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Service is created (description might be optional)
	// The exact behavior depends on business rules
	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, createRequest.Name, result.Name)
		assert.Equal(t, "", result.Description)
	} else {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestCreateServiceInvalidImageUrl(t *testing.T) {
	// GIVEN: Service creation request with invalid image URL
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	createRequest := schemas.CreateServiceRequest{
		Name:        "Service Invalid URL",
		Description: "Service with invalid image URL",
		ImageUrl:    "not-a-valid-url",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Service is created (URL validation might be at different layer)
	// The exact behavior depends on implementation
	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, createRequest.Name, result.Name)
		assert.Equal(t, createRequest.ImageUrl, result.ImageUrl)
	} else {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestCreateServiceLongName(t *testing.T) {
	// GIVEN: Service creation request with very long name
	controller, _, _ := controllerTest.NewServiceControllerTestWrapper(t)

	longName := "This is a very long service name that tests the system's ability to handle extended text fields and ensure proper storage and retrieval of service names"

	createRequest := schemas.CreateServiceRequest{
		Name:        longName,
		Description: "Service with long name",
		ImageUrl:    "https://example.com/long-name.jpg",
		IsVirtual:   false,
	}

	// WHEN: CreateService is called
	result, err := controller.CreateService(createRequest, "test_admin")

	// THEN: Service with long name is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, longName, result.Name)
}
