package plan_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestUpdatePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  UpdatePostgresqlPlan is called with new data
		THEN:  The plan is updated and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	newFee := 99.99
	newReservationLimit := 15
	updatedBy := "test-admin"

	// WHEN
	updatedPlan, err := adapter.UpdatePostgresqlPlan(
		plan.Id,
		&newFee,
		nil, // Don't update type
		&newReservationLimit,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedPlan)
	assert.Equal(t, plan.Id, updatedPlan.Id)
	assert.Equal(t, newFee, updatedPlan.Fee)
	assert.Equal(t, &newReservationLimit, updatedPlan.ReservationLimit)
}

func TestUpdatePlanWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  UpdatePostgresqlPlan is called with empty updatedBy
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	newFee := 99.99
	updatedBy := ""

	// WHEN
	updatedPlan, err := adapter.UpdatePostgresqlPlan(
		plan.Id,
		&newFee,
		nil, nil,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, updatedPlan)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestUpdatePlanWithCompleteData(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  UpdatePostgresqlPlan is called with all fields
		THEN:  The plan is completely updated
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	newFee := 149.99
	newType := model.PlanTypeAnual
	newReservationLimit := 25
	updatedBy := "test-admin"

	// WHEN
	updatedPlan, err := adapter.UpdatePostgresqlPlan(
		plan.Id,
		&newFee,
		&newType,
		&newReservationLimit,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedPlan)
	assert.Equal(t, plan.Id, updatedPlan.Id)
	assert.Equal(t, newFee, updatedPlan.Fee)
	assert.Equal(t, newType, updatedPlan.Type)
	assert.Equal(t, &newReservationLimit, updatedPlan.ReservationLimit)
}

func TestUpdatePlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the given ID
		WHEN:  UpdatePostgresqlPlan is called
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	nonExistentId := uuid.New()
	newFee := 99.99
	updatedBy := "test-admin"

	// WHEN
	updatedPlan, err := adapter.UpdatePostgresqlPlan(
		nonExistentId,
		&newFee,
		nil, nil,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, updatedPlan)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *err)
}

func TestDeletePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  DeletePostgresqlPlan is called
		THEN:  The plan is soft deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN
	err := adapter.DeletePostgresqlPlan(plan.Id)

	// THEN
	assert.Nil(t, err)

	// Verify plan is deleted by trying to get it
	deletedPlan, getErr := adapter.GetPostgresqlPlan(plan.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedPlan)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *getErr)
}

func TestDeletePlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the given ID
		WHEN:  DeletePostgresqlPlan is called
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	err := adapter.DeletePostgresqlPlan(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *err)
}

func TestBulkDeletePlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple plans exist in the database
		WHEN:  BulkDeletePostgresqlPlans is called
		THEN:  All specified plans are soft deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan3 := factories.NewPlanModel(db, factories.PlanModelF{})

	planIds := []uuid.UUID{plan1.Id, plan2.Id}

	// WHEN
	err := adapter.BulkDeletePostgresqlPlans(planIds)

	// THEN
	assert.Nil(t, err)

	// Verify deleted plans cannot be found
	_, getErr1 := adapter.GetPostgresqlPlan(plan1.Id)
	assert.NotNil(t, getErr1)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *getErr1)

	_, getErr2 := adapter.GetPostgresqlPlan(plan2.Id)
	assert.NotNil(t, getErr2)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *getErr2)

	// Verify non-deleted plan still exists
	plan3Result, getErr3 := adapter.GetPostgresqlPlan(plan3.Id)
	assert.Nil(t, getErr3)
	assert.NotNil(t, plan3Result)
	assert.Equal(t, plan3.Id, plan3Result.Id)
}

func TestBulkDeletePlansEmpty(t *testing.T) {
	/*
		GIVEN: An empty list of plan IDs
		WHEN:  BulkDeletePostgresqlPlans is called
		THEN:  No error occurs
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	emptyIds := []uuid.UUID{}

	// WHEN
	err := adapter.BulkDeletePostgresqlPlans(emptyIds)

	// THEN
	assert.Nil(t, err)
}
