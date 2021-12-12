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

//YeastService service for managing Yeasts
type YeastService struct {
	db     *gorm.DB
	config *viper.Viper
}

//NewYeastService creates an instance of the Yeast service
func NewYeastService() *YeastService {
	return &YeastService{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}
}

//GetAllYeasts returns a list of all Yeasts
func (s *YeastService) GetAllYeasts() []models.Yeast {
	yeasts := make([]models.Yeast, 0)
	s.db.Preload(clause.Associations).Find(&yeasts)
	return yeasts
}

//CreateYeast creates a new Yeast and returns it
func (s *YeastService) CreateYeast(name string, brand string, tolerance string, attenuation string, flocculation string, temperature string) models.Yeast {
	y := &models.Yeast{
		Name:         name,
		Brand:        brand,
		Tolerance:    tolerance,
		Attenuation:  attenuation,
		Flocculation: flocculation,
		Temperature:  temperature,
	}
	s.db.Create(y)
	return *y
}

//GetYeastByID returns a Yeast by its id and returns it and an error if not found
func (s *YeastService) GetYeastByID(id int) (*models.Yeast, error) {
	g := models.Yeast{}
	var c int64
	s.db.Preload(clause.Associations).Find(&g, id).Count(&c)
	if c > 0 {
		return &g, nil
	}
	return nil, errors.New("YEAST_NOT_FOUND")
}

//DeleteYeastByID deletes a Yeast by its id and returns the deleted item and an error is it cannot be found
func (s *YeastService) DeleteYeastByID(id int) (*models.Yeast, error) {
	g := models.Yeast{}
	var c int64
	s.db.Find(&g, id).Count(&c)
	if c < 1 {
		return nil, errors.New("YEAST_NOT_FOUND")
	}
	s.db.Delete(&g)
	return &g, nil
}
