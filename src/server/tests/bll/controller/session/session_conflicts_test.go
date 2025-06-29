package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

// TestCreateSessionWithProfessionalConflict verifica que no se pueda crear una sesión
// cuando el profesional ya tiene otra sesión en el mismo horario
func TestCreateSessionWithProfessionalConflict(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión existente para el profesional
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Intentar crear una sesión conflictiva (mismo profesional, mismo horario)
	conflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC)
	conflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Conflicting Session",
		Date:           sessionDate,
		StartTime:      conflictStartTime,
		EndTime:        conflictEndTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession es llamado con conflicto de profesional
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionWithLocalConflict verifica que no se pueda crear una sesión
// cuando el local ya tiene otra sesión en el mismo horario
func TestCreateSessionWithLocalConflict(t *testing.T) {
	// GIVEN: Un local con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente en el local
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional1.Id,
		LocalId:            &testLocal.Id,
		CommunityServiceId: &testCommunityService.Id, // Mismo tipo de servicio
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Intentar crear una sesión conflictiva (mismo local, mismo horario, mismo servicio, diferente profesional)
	conflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 30, 0, 0, time.UTC)
	conflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 30, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:              "Conflicting Local Session",
		Date:               sessionDate,
		StartTime:          conflictStartTime,
		EndTime:            conflictEndTime,
		Capacity:           10,
		ProfessionalId:     testProfessional2.Id,     // Diferente profesional
		LocalId:            &testLocal.Id,            // Mismo local
		CommunityServiceId: &testCommunityService.Id, // Mismo tipo de servicio
	}

	// WHEN: CreateSession es llamado con conflicto de local
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionNoConflict verifica que se puedan crear sesiones
// cuando no hay conflictos
func TestCreateSessionNoConflict(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Crear una sesión sin conflictos (horario diferente)
	noConflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	noConflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "No Conflict Session",
		Date:           sessionDate,
		StartTime:      noConflictStartTime,
		EndTime:        noConflictEndTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession es llamado sin conflictos
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
	assert.Equal(t, createRequest.ProfessionalId, result.ProfessionalId)
	assert.Equal(t, createRequest.LocalId, result.LocalId)
}

// TestCreateSessionConsecutiveTimes verifica que se puedan crear sesiones
// consecutivas (sin superposición)
func TestCreateSessionConsecutiveTimes(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Crear una sesión consecutiva (empieza cuando termina la anterior)
	consecutiveStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	consecutiveEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 12, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Consecutive Session",
		Date:           sessionDate,
		StartTime:      consecutiveStartTime,
		EndTime:        consecutiveEndTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession es llamado con horarios consecutivos
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
}

// TestCreateSessionDifferentDays verifica que se puedan crear sesiones
// en días diferentes sin conflictos
func TestCreateSessionDifferentDays(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión existente
	sessionDate1 := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate1.Year(), sessionDate1.Month(), sessionDate1.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate1.Year(), sessionDate1.Month(), sessionDate1.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate1,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Crear una sesión en un día diferente
	sessionDate2 := time.Now().Add(48 * time.Hour) // Día siguiente
	sameTimeStart := time.Date(sessionDate2.Year(), sessionDate2.Month(), sessionDate2.Day(), 10, 0, 0, 0, time.UTC)
	sameTimeEnd := time.Date(sessionDate2.Year(), sessionDate2.Month(), sessionDate2.Day(), 11, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Different Day Session",
		Date:           sessionDate2,
		StartTime:      sameTimeStart,
		EndTime:        sameTimeEnd,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession es llamado en día diferente
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
}

// TestCreateSessionWithCancelledSession verifica que las sesiones canceladas
// no bloqueen la creación de nuevas sesiones
func TestCreateSessionWithCancelledSession(t *testing.T) {
	// GIVEN: Un profesional con una sesión cancelada
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión cancelada
	sessionDate := time.Now().Add(24 * time.Hour)
	cancelledStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	cancelledEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	cancelledState := model.SessionStateCancelled

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &cancelledStartTime,
		EndTime:        &cancelledEndTime,
		State:          &cancelledState,
	})

	// Intentar crear una sesión en el mismo horario (debería permitirse)
	createRequest := schemas.CreateSessionRequest{
		Title:          "New Session After Cancelled",
		Date:           sessionDate,
		StartTime:      cancelledStartTime,
		EndTime:        cancelledEndTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession es llamado con sesión cancelada existente
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente (las canceladas no bloquean)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
}

// TestCreateSessionVirtualNoLocalConflict verifica que las sesiones virtuales
// no tengan conflictos de local
func TestCreateSessionVirtualNoLocalConflict(t *testing.T) {
	// GIVEN: Un local con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión presencial existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional1.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Crear una sesión virtual en el mismo horario (debería permitirse)
	sessionLink := "https://meet.example.com/virtual"
	createRequest := schemas.CreateSessionRequest{
		Title:          "Virtual Session",
		Date:           sessionDate,
		StartTime:      existingStartTime,
		EndTime:        existingEndTime,
		Capacity:       20,
		SessionLink:    &sessionLink,
		ProfessionalId: testProfessional2.Id,
		LocalId:        nil, // Sesión virtual
	}

	// WHEN: CreateSession es llamado para sesión virtual
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión virtual se debe crear exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
	assert.Nil(t, result.LocalId) // Debe ser nil para sesión virtual
	assert.Equal(t, sessionLink, *result.SessionLink)
}

// TestUpdateSessionWithConflict verifica que no se pueda actualizar una sesión
// a un horario conflictivo
func TestUpdateSessionWithConflict(t *testing.T) {
	// GIVEN: Dos sesiones existentes
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear primera sesión
	sessionDate := time.Now().Add(24 * time.Hour)
	session1StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	session1EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session1StartTime,
		EndTime:        &session1EndTime,
	})

	// Crear segunda sesión
	session2StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	session2EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session2StartTime,
		EndTime:        &session2EndTime,
	})

	// Intentar actualizar la segunda sesión para que conflictúe con la primera
	conflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC)
	conflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	updateRequest := schemas.UpdateSessionRequest{
		StartTime: &conflictStartTime,
		EndTime:   &conflictEndTime,
	}

	// WHEN: UpdateSession es llamado con conflicto
	result, err := controller.UpdateSession(session2.Id, updateRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestUpdateSessionNoConflict verifica que se pueda actualizar una sesión
// a un horario sin conflictos
func TestUpdateSessionNoConflict(t *testing.T) {
	// GIVEN: Una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear sesión
	sessionDate := time.Now().Add(24 * time.Hour)
	originalStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	originalEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	session := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &originalStartTime,
		EndTime:        &originalEndTime,
	})

	// Actualizar a un horario sin conflictos
	newStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	updateRequest := schemas.UpdateSessionRequest{
		StartTime: &newStartTime,
		EndTime:   &newEndTime,
	}

	// WHEN: UpdateSession es llamado sin conflictos
	result, err := controller.UpdateSession(session.Id, updateRequest, "test_admin")

	// THEN: La sesión se debe actualizar exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	// Comparar en UTC para evitar errores de zona horaria
	assert.Equal(t, newStartTime.UTC(), result.StartTime.UTC())
	assert.Equal(t, newEndTime.UTC(), result.EndTime.UTC())
}

// TestBulkCreateSessionsWithConflict verifica que la creación masiva
// detecte conflictos
func TestBulkCreateSessionsWithConflict(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Crear solicitudes de creación masiva con conflicto
	conflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC)
	conflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	bulkRequest := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{
			{
				Title:          "Bulk Session 1",
				Date:           sessionDate,
				StartTime:      conflictStartTime,
				EndTime:        conflictEndTime,
				Capacity:       10,
				ProfessionalId: testProfessional.Id,
				LocalId:        &testLocal.Id,
			},
			{
				Title:          "Bulk Session 2",
				Date:           sessionDate,
				StartTime:      conflictStartTime,
				EndTime:        conflictEndTime,
				Capacity:       10,
				ProfessionalId: testProfessional.Id,
				LocalId:        &testLocal.Id,
			},
		},
	}

	// WHEN: BulkCreateSessions es llamado con conflictos
	result, err := controller.BulkCreateSessions(bulkRequest.Sessions, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestBulkCreateSessionsNoConflict verifica que la creación masiva
// funcione sin conflictos
func TestBulkCreateSessionsNoConflict(t *testing.T) {
	// GIVEN: Un profesional sin sesiones existentes
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Crear solicitudes de creación masiva sin conflictos
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime1 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime1 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	startTime2 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	endTime2 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	bulkRequest := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{
			{
				Title:          "Bulk Session 1",
				Date:           sessionDate,
				StartTime:      startTime1,
				EndTime:        endTime1,
				Capacity:       10,
				ProfessionalId: testProfessional.Id,
				LocalId:        &testLocal.Id,
			},
			{
				Title:          "Bulk Session 2",
				Date:           sessionDate,
				StartTime:      startTime2,
				EndTime:        endTime2,
				Capacity:       10,
				ProfessionalId: testProfessional.Id,
				LocalId:        &testLocal.Id,
			},
		},
	}

	// WHEN: BulkCreateSessions es llamado sin conflictos
	result, err := controller.BulkCreateSessions(bulkRequest.Sessions, "test_admin")

	// THEN: Las sesiones se deben crear exitosamente
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Sessions, 2)
	assert.Equal(t, "Bulk Session 1", result.Sessions[0].Title)
	assert.Equal(t, "Bulk Session 2", result.Sessions[1].Title)
}

// TestCreateSessionSameLocalDifferentActivities verifica que NO se pueda crear una sesión
// con el mismo local pero diferente actividad al mismo tiempo
func TestCreateSessionSameLocalDifferentActivities(t *testing.T) {
	// GIVEN: Un local con una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente en el local
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional1.Id,
		LocalId:            &testLocal.Id,
		CommunityServiceId: &testCommunityService1.Id, // Yoga
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Intentar crear otra sesión con el mismo local pero diferente actividad
	conflictStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC)
	conflictEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:              "Gym Session",
		Date:               sessionDate,
		StartTime:          conflictStartTime,
		EndTime:            conflictEndTime,
		Capacity:           15,
		ProfessionalId:     testProfessional2.Id,      // Diferente profesional
		LocalId:            &testLocal.Id,             // Mismo local
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con mismo local pero diferente actividad
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto (nueva lógica: solo una actividad por local)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionSameLocalSameActivity verifica que NO se pueda crear una sesión
// en el mismo local con la misma actividad
func TestCreateSessionSameLocalSameActivity(t *testing.T) {
	// GIVEN: Un local con una sesión existente de una actividad
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente (Yoga)
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional1.Id,
		LocalId:            &testLocal.Id,
		CommunityServiceId: &testCommunityService.Id, // Yoga
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Intentar crear otra sesión en el mismo local con la misma actividad
	createRequest := schemas.CreateSessionRequest{
		Title:              "Another Yoga Session",
		Date:               sessionDate,
		StartTime:          existingStartTime, // Mismo horario
		EndTime:            existingEndTime,   // Mismo horario
		Capacity:           10,
		ProfessionalId:     testProfessional2.Id,
		LocalId:            &testLocal.Id,            // Mismo local
		CommunityServiceId: &testCommunityService.Id, // Misma actividad
	}

	// WHEN: CreateSession es llamado con mismo local y misma actividad
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionTimeOverlapOnly verifica que se detecten conflictos
// solo por superposición de horarios (sin considerar profesional o local)
func TestCreateSessionTimeOverlapOnly(t *testing.T) {
	// GIVEN: Una sesión existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal1 := factories.NewLocalModel(db, factories.LocalModelF{})
	testLocal2 := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional1.Id,
		LocalId:            &testLocal1.Id,
		CommunityServiceId: &testCommunityService1.Id,
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Crear una sesión con superposición de horario pero diferente profesional, local y actividad
	overlapStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC)
	overlapEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:              "Different Session",
		Date:               sessionDate,
		StartTime:          overlapStartTime,
		EndTime:            overlapEndTime,
		Capacity:           10,
		ProfessionalId:     testProfessional2.Id,      // Diferente profesional
		LocalId:            &testLocal2.Id,            // Diferente local
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con superposición de horario pero todo diferente
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente (no hay conflicto real)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
}

// TestBulkCreateSessionsSameLocalDifferentActivities verifica que NO se pueda crear
// un lote con sesiones en el mismo local con diferentes actividades al mismo tiempo
func TestBulkCreateSessionsSameLocalDifferentActivities(t *testing.T) {
	// GIVEN: Un lote con sesiones en el mismo local pero diferentes actividades
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear lote con sesiones en el mismo local pero diferentes actividades
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	createSessionsData := []*schemas.CreateSessionRequest{
		{
			Title:              "Yoga Session",
			Date:               sessionDate,
			StartTime:          startTime,
			EndTime:            endTime,
			Capacity:           15,
			ProfessionalId:     testProfessional1.Id,
			LocalId:            &testLocal.Id,
			CommunityServiceId: &testCommunityService1.Id, // Yoga
		},
		{
			Title:              "Gym Session",
			Date:               sessionDate,
			StartTime:          startTime, // Mismo horario
			EndTime:            endTime,   // Mismo horario
			Capacity:           20,
			ProfessionalId:     testProfessional2.Id,
			LocalId:            &testLocal.Id,             // Mismo local
			CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
		},
	}

	// WHEN: BulkCreateSessions es llamado con conflicto de local
	result, err := controller.BulkCreateSessions(createSessionsData, "test_admin")

	// THEN: Se debe retornar un error de conflicto (nueva lógica: solo una actividad por local)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestBulkCreateSessionsSameLocalSameActivity verifica que NO se puedan crear
// sesiones masivas en el mismo local con la misma actividad
func TestBulkCreateSessionsSameLocalSameActivity(t *testing.T) {
	// GIVEN: Un local sin sesiones existentes
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear solicitudes de creación masiva en el mismo local con la misma actividad
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	bulkRequest := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{
			{
				Title:              "Yoga Session 1",
				Date:               sessionDate,
				StartTime:          startTime,
				EndTime:            endTime,
				Capacity:           10,
				ProfessionalId:     testProfessional1.Id,
				LocalId:            &testLocal.Id,
				CommunityServiceId: &testCommunityService.Id, // Yoga
			},
			{
				Title:              "Yoga Session 2",
				Date:               sessionDate,
				StartTime:          startTime, // Mismo horario
				EndTime:            endTime,   // Mismo horario
				Capacity:           15,
				ProfessionalId:     testProfessional2.Id,
				LocalId:            &testLocal.Id,            // Mismo local
				CommunityServiceId: &testCommunityService.Id, // Misma actividad
			},
		},
	}

	// WHEN: BulkCreateSessions es llamado con mismo local y misma actividad
	result, err := controller.BulkCreateSessions(bulkRequest.Sessions, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionSameProfessionalSameLocal verifica que NO se pueda crear una sesión
// con el mismo profesional en el mismo local al mismo tiempo
func TestCreateSessionSameProfessionalSameLocal(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente en un local
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional.Id,
		LocalId:            &testLocal.Id,
		CommunityServiceId: &testCommunityService1.Id, // Yoga
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Intentar crear otra sesión con el mismo profesional en el mismo local pero diferente actividad
	createRequest := schemas.CreateSessionRequest{
		Title:              "Gym Session",
		Date:               sessionDate,
		StartTime:          existingStartTime, // Mismo horario
		EndTime:            existingEndTime,   // Mismo horario
		Capacity:           15,
		ProfessionalId:     testProfessional.Id,       // Mismo profesional
		LocalId:            &testLocal.Id,             // Mismo local
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con mismo profesional y mismo local
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionSameProfessionalDifferentLocals verifica que NO se pueda crear una sesión
// con el mismo profesional en diferentes locales al mismo tiempo
func TestCreateSessionSameProfessionalDifferentLocals(t *testing.T) {
	// GIVEN: Un profesional con una sesión existente en un local
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal1 := factories.NewLocalModel(db, factories.LocalModelF{})
	testLocal2 := factories.NewLocalModel(db, factories.LocalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional.Id,
		LocalId:            &testLocal1.Id,
		CommunityServiceId: &testCommunityService1.Id,
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
	})

	// Intentar crear una sesión con el mismo profesional pero en diferente local
	createRequest := schemas.CreateSessionRequest{
		Title:              "Different Local Session",
		Date:               sessionDate,
		StartTime:          existingStartTime, // Mismo horario
		EndTime:            existingEndTime,   // Mismo horario
		Capacity:           15,
		ProfessionalId:     testProfessional.Id,       // Mismo profesional
		LocalId:            &testLocal2.Id,            // Diferente local
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con mismo profesional pero diferente local
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto (profesional no puede estar en dos lugares)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}

// TestCreateSessionVirtualSessionsSimultaneous verifica que se puedan crear
// múltiples sesiones virtuales simultáneamente (solo se valida conflicto de profesional)
func TestCreateSessionVirtualSessionsSimultaneous(t *testing.T) {
	// GIVEN: Un profesional con una sesión virtual existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión virtual existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	sessionLink1 := "https://meet.google.com/session1"

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional1.Id,
		LocalId:            nil, // Sesión virtual
		CommunityServiceId: &testCommunityService1.Id,
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
		SessionLink:        &sessionLink1,
	})

	// Crear otra sesión virtual simultánea con diferente profesional
	sessionLink2 := "https://meet.google.com/session2"

	createRequest := schemas.CreateSessionRequest{
		Title:              "Virtual Session 2",
		Date:               sessionDate,
		StartTime:          existingStartTime, // Mismo horario
		EndTime:            existingEndTime,   // Mismo horario
		Capacity:           20,
		SessionLink:        &sessionLink2,
		ProfessionalId:     testProfessional2.Id,      // Diferente profesional
		LocalId:            nil,                       // Sesión virtual
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con sesión virtual simultánea
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: La sesión se debe crear exitosamente (no hay conflicto de local para virtuales)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
	assert.Nil(t, result.LocalId) // Confirmar que es virtual
	assert.Equal(t, testProfessional2.Id, result.ProfessionalId)
}

// TestCreateSessionVirtualSessionsProfessionalConflict verifica que NO se pueda crear
// una sesión virtual si hay conflicto de profesional
func TestCreateSessionVirtualSessionsProfessionalConflict(t *testing.T) {
	// GIVEN: Un profesional con una sesión virtual existente
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Crear dependencias
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testCommunityService1 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})
	testCommunityService2 := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	// Crear una sesión virtual existente
	sessionDate := time.Now().Add(24 * time.Hour)
	existingStartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	existingEndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	sessionLink1 := "https://meet.google.com/session1"

	_ = factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId:     &testProfessional.Id,
		LocalId:            nil, // Sesión virtual
		CommunityServiceId: &testCommunityService1.Id,
		Date:               &sessionDate,
		StartTime:          &existingStartTime,
		EndTime:            &existingEndTime,
		SessionLink:        &sessionLink1,
	})

	// Intentar crear otra sesión virtual con el mismo profesional
	sessionLink2 := "https://meet.google.com/session2"

	createRequest := schemas.CreateSessionRequest{
		Title:              "Virtual Session 2",
		Date:               sessionDate,
		StartTime:          existingStartTime, // Mismo horario
		EndTime:            existingEndTime,   // Mismo horario
		Capacity:           20,
		SessionLink:        &sessionLink2,
		ProfessionalId:     testProfessional.Id,       // Mismo profesional
		LocalId:            nil,                       // Sesión virtual
		CommunityServiceId: &testCommunityService2.Id, // Diferente actividad
	}

	// WHEN: CreateSession es llamado con conflicto de profesional
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Se debe retornar un error de conflicto (profesional no puede estar en dos lugares)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_001", err.Code)
	assert.Contains(t, err.Message, "Session time conflict detected")
}
