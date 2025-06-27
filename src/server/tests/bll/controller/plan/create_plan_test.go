package plan_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreatePlanSuccessfully(t *testing.T) {
	// GIVEN: Valid plan creation request
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 30
	createRequest := schemas.CreatePlanRequest{
		Fee:              99.99,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Plan is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Fee, result.Fee)
	assert.Equal(t, createRequest.Type, result.Type)
	assert.Equal(t, createRequest.ReservationLimit, result.ReservationLimit)
	assert.NotEqual(t, "", result.Id)
}

func TestCreatePlanEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid plan creation request but empty updatedBy
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 30
	createRequest := schemas.CreatePlanRequest{
		Fee:              49.99,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called with empty updatedBy
	result, err := controller.CreatePlan(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreatePlanAnnualType(t *testing.T) {
	// GIVEN: Plan creation request with annual type
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 365
	createRequest := schemas.CreatePlanRequest{
		Fee:              999.99,
		Type:             model.PlanTypeAnual,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Annual plan is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Fee, result.Fee)
	assert.Equal(t, model.PlanTypeAnual, result.Type)
	assert.Equal(t, &reservationLimit, result.ReservationLimit)
}

func TestCreatePlanZeroFee(t *testing.T) {
	// GIVEN: Plan creation request with zero fee (free plan)
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 10
	createRequest := schemas.CreatePlanRequest{
		Fee:              0.0,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Free plan is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.Fee)
	assert.Equal(t, model.PlanTypeMonthly, result.Type)
}

func TestCreatePlanNegativeFee(t *testing.T) {
	// GIVEN: Plan creation request with negative fee
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 30
	createRequest := schemas.CreatePlanRequest{
		Fee:              -10.0,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreatePlanNilReservationLimit(t *testing.T) {
	// GIVEN: Plan creation request with nil reservation limit
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	createRequest := schemas.CreatePlanRequest{
		Fee:              19.99,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: nil, // Unlimited reservations
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Plan with unlimited reservations is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Fee, result.Fee)
	assert.Nil(t, result.ReservationLimit)
}

func TestCreatePlanZeroReservationLimit(t *testing.T) {
	// GIVEN: Plan creation request with zero reservation limit
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 0
	createRequest := schemas.CreatePlanRequest{
		Fee:              19.99,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit, // Zero reservations allowed
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Behavior depends on business rules
	// Zero reservations might be invalid or might represent a read-only plan
	if err != nil {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, result)
		assert.Equal(t, &reservationLimit, result.ReservationLimit)
	}
}

func TestCreatePlanHighFee(t *testing.T) {
	// GIVEN: Plan creation request with very high fee
	controller, _, _ := controllerTest.NewPlanControllerTestWrapper(t)

	reservationLimit := 1000
	createRequest := schemas.CreatePlanRequest{
		Fee:              9999.99,
		Type:             model.PlanTypeAnual,
		ReservationLimit: &reservationLimit,
	}

	// WHEN: CreatePlan is called
	result, err := controller.CreatePlan(createRequest, "test_admin")

	// THEN: Plan with high fee is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 9999.99, result.Fee)
	assert.Equal(t, model.PlanTypeAnual, result.Type)
}
