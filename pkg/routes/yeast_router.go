package routes

import (
	"github.com/LucasCarioca/ibru-api/pkg/config"
	"github.com/LucasCarioca/ibru-api/pkg/datasource"
	"github.com/LucasCarioca/ibru-api/pkg/models"
	"github.com/LucasCarioca/ibru-api/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	AuditId      string `json:"audit_id" binding:"required"`
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
	approved, err := strconv.ParseBool(ctx.Query("approved"))
	if err != nil {
		approved = true
	}
	auditId := ctx.Query("audit_id")
	yeasts := make([]models.Yeast, 0)
	if auditId != "" {
		yeasts = r.gs.GetAllYeastsByAuditId(approved, auditId)
	} else {
		yeasts = r.gs.GetAllYeasts(approved)
	}
	ctx.JSON(http.StatusOK, yeasts)
}

func (r *YeastRouter) createYeast(ctx *gin.Context) {
	var data CreateYeastRequest
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body: some or all fields are missing or incorrect", "error": "INVALID_REQUEST_BODY"})
		return
	}
	g := r.gs.CreateYeast(data.Name, data.Brand, data.Temperature, data.Attenuation, data.Flocculation, data.Temperature, data.AuditId)
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
