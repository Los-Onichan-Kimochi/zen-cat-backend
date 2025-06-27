package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type TemplateModelF struct {
	Id             *uuid.UUID
	Link           *string
	ProfessionalId *uuid.UUID
}

// Create a new template on DB
func NewTemplateModel(db *gorm.DB, option ...TemplateModelF) *model.Template {
	// Create default professional if not provided
	professional := NewProfessionalModel(db)

	template := &model.Template{
		Id:             uuid.New(),
		Link:           "https://example.com/template",
		ProfessionalId: professional.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			template.Id = *parameters.Id
		}
		if parameters.Link != nil {
			template.Link = *parameters.Link
		}
		if parameters.ProfessionalId != nil {
			template.ProfessionalId = *parameters.ProfessionalId
		}
	}

	result := db.Create(template)
	if result.Error != nil {
		log.Fatalf("Error when trying to create template: %v", result.Error)
	}

	return template
}

// Create size number of new templates on DB
func NewTemplateModelBatch(
	db *gorm.DB,
	size int,
	option ...TemplateModelF,
) []*model.Template {
	templates := []*model.Template{}
	for i := 0; i < size; i++ {
		var template *model.Template
		if len(option) > 0 {
			template = NewTemplateModel(db, option[0])
		} else {
			template = NewTemplateModel(db)
		}
		templates = append(templates, template)
	}
	return templates
}
