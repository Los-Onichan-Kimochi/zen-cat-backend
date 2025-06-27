package plan_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetPlanSuccessfully(t *testing.T) {
	/*
		GIVEN: A plan exists in the database
		WHEN:  GetPostgresqlPlan is called with the plan ID
		THEN:  The plan is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN
	result, err := adapter.GetPostgresqlPlan(plan.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, plan.Id, result.Id)
	assert.Equal(t, plan.Fee, result.Fee)
	assert.Equal(t, plan.Type, result.Type)
	assert.Equal(t, plan.ReservationLimit, result.ReservationLimit)
}

func TestGetPlanNotFound(t *testing.T) {
	/*
		GIVEN: No plan exists with the given ID
		WHEN:  GetPostgresqlPlan is called with non-existent ID
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	plan, err := adapter.GetPostgresqlPlan(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, errors.ObjectNotFoundError.PlanNotFound, *err)
}

func TestFetchPlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple plans exist in the database
		WHEN:  FetchPostgresqlPlans is called with specific IDs
		THEN:  All matching plans are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan3 := factories.NewPlanModel(db, factories.PlanModelF{})

	planIds := []uuid.UUID{plan1.Id, plan2.Id}

	// WHEN
	plans, err := adapter.FetchPostgresqlPlans(planIds)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.Equal(t, 2, len(plans))

	// Verify the correct plans were returned
	foundPlan1 := false
	foundPlan2 := false
	foundPlan3 := false
	for _, plan := range plans {
		if plan.Id == plan1.Id {
			foundPlan1 = true
		}
		if plan.Id == plan2.Id {
			foundPlan2 = true
		}
		if plan.Id == plan3.Id {
			foundPlan3 = true
		}
	}
	assert.True(t, foundPlan1)
	assert.True(t, foundPlan2)
	assert.False(t, foundPlan3) // Should not be included
}

func TestFetchPlansWithEmptyIds(t *testing.T) {
	/*
		GIVEN: Plans exist in the database
		WHEN:  FetchPostgresqlPlans is called with empty ID list
		THEN:  All plans are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewPlanAdapterTestWrapper(t)

	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN
	plans, err := adapter.FetchPostgresqlPlans([]uuid.UUID{})

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.GreaterOrEqual(t, len(plans), 2)

	// Find our created plans
	foundPlan1 := false
	foundPlan2 := false
	for _, plan := range plans {
		if plan.Id == plan1.Id {
			foundPlan1 = true
		}
		if plan.Id == plan2.Id {
			foundPlan2 = true
		}
	}
	assert.True(t, foundPlan1)
	assert.True(t, foundPlan2)
}

func TestFetchPlansNotFound(t *testing.T) {
	/*
		GIVEN: No plans exist with the given IDs
		WHEN:  FetchPostgresqlPlans is called with non-existent IDs
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	nonExistentIds := []uuid.UUID{uuid.New(), uuid.New()}

	// WHEN
	plans, err := adapter.FetchPostgresqlPlans(nonExistentIds)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.Equal(t, 0, len(plans))
}
