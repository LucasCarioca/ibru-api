package routes

import (
	"fmt"
	"github.com/LucasCarioca/ibru-api/pkg/config"
	"github.com/LucasCarioca/ibru-api/pkg/datasource"
	"github.com/LucasCarioca/ibru-api/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
)

//YeastRouter router for Yeast CRUD operations
type YeastRouter struct {
	db     *gorm.DB
	config *viper.Viper
	gs     *services.YeastService
}

//CreateYeastRequest structure of the create request for Yeasts
type CreateYeastRequest struct {
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand" binding:"required"`
	Tolerance    string `json:"tolerance" binding:"required"`
	Attenuation  string `json:"attenuation" binding:"required"`
	Flocculation string `json:"flocculation" binding:"required"`
	Temperature  string `json:"temperature" binding:"required"`
}

//NewYeastRouter creates a new instance of the Yeast router
func NewYeastRouter(app *gin.Engine) {
	r := YeastRouter{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
		gs:     services.NewYeastService(),
	}

	app.GET("/api/v1/yeasts", r.getAllYeasts)
	app.GET("/api/v1/yeasts/:id", r.getYeast)
	app.POST("/api/v1/yeasts", r.createYeast)
	app.DELETE("/api/v1/yeasts/:id", r.deleteYeast)
}

func (r *YeastRouter) getAllYeasts(ctx *gin.Context) {
	yeasts := r.gs.GetAllYeasts()
	ctx.JSON(http.StatusOK, yeasts)
}

func (r *YeastRouter) createYeast(ctx *gin.Context) {
	var data CreateYeastRequest
	ctx.BindJSON(&data)
	fmt.Println(data)
	g := r.gs.CreateYeast(data.Name, data.Brand, data.Temperature, data.Attenuation, data.Flocculation, data.Temperature)
	ctx.JSON(http.StatusOK, g)
}

func (r *YeastRouter) getYeast(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "yeast not found", "error": "YEAST_NOT_FOUND"})
		return
	}

	g, err := r.gs.GetYeastByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "yeast not found", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, g)
}

func (r *YeastRouter) deleteYeast(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "yeast not found", "error": "YEAST_NOT_FOUND"})
		return
	}

	g, err := r.gs.DeleteYeastByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "yeast not found", "error": "YEAST_NOT_FOUND"})
		return
	}
	ctx.JSON(http.StatusOK, g)
}
