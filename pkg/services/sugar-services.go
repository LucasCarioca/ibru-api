package services

import (
	"errors"

	"github.com/LucasCarioca/ibru-api/pkg/config"
	"github.com/LucasCarioca/ibru-api/pkg/datasource"
	"github.com/LucasCarioca/ibru-api/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//SugarService service for managing Sugars
type SugarService struct {
	db     *gorm.DB
	config *viper.Viper
}

//NewSugarService creates an instance of the Sugar service
func NewSugarService() *SugarService {
	return &SugarService{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}
}

//GetAllSugars returns a list of all Sugars
func (s *SugarService) GetAllSugars(approvedOnly bool) []models.Sugar {
	sugars := make([]models.Sugar, 0)
	if approvedOnly {
		s.db.Preload(clause.Associations).Find(&sugars, "approved = ?", approvedOnly)
	} else {
		s.db.Preload(clause.Associations).Find(&sugars)
	}
	return sugars
}

//GetAllSugarsByAuditId returns a list of all Sugars for a specific audit id
func (s *SugarService) GetAllSugarsByAuditId(approvedOnly bool, auditId string) []models.Sugar {
	sugars := make([]models.Sugar, 0)
	if approvedOnly {
		s.db.Preload(clause.Associations).Find(&sugars, "approved = ? AND audit_id = ?", approvedOnly, auditId)
	} else {
		s.db.Preload(clause.Associations).Find(&sugars, "audit_id = ?", auditId)
	}
	return sugars
}

//CreateSugar creates a new Sugar and returns it
func (s *SugarService) CreateSugar(name string, gravityPerPound float32, auditId string) models.Sugar {
	y := &models.Sugar{
		Name:            name,
		GravityPerPound: gravityPerPound,
		Base:            models.Base{AuditId: auditId, Approved: false},
	}
	s.db.Create(y)
	return *y
}

//GetSugarByID returns a Sugar by its id and returns it and an error if not found
func (s *SugarService) GetSugarByID(id int) (*models.Sugar, error) {
	g := models.Sugar{}
	var c int64
	s.db.Preload(clause.Associations).Find(&g, id).Count(&c)
	if c > 0 {
		return &g, nil
	}
	return nil, errors.New("SUGAR_NOT_FOUND")
}

//DeleteSugarByID deletes a Sugar by its id and returns the deleted item and an error is it cannot be found
func (s *SugarService) DeleteSugarByID(id int) (*models.Sugar, error) {
	g := models.Sugar{}
	var c int64
	s.db.Find(&g, id).Count(&c)
	if c < 1 {
		return nil, errors.New("SUGAR_NOT_FOUND")
	}
	s.db.Delete(&g)
	return &g, nil
}
