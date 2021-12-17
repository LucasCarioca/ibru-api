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

//SugarRouter router for Sugar CRUD operations
type SugarRouter struct {
	db     *gorm.DB
	config *viper.Viper
	gs     *services.SugarService
}

//CreateSugarRequest structure of the create request for Sugars
type CreateSugarRequest struct {
	Name            string  `json:"name" binding:"required"`
	GravityPerPound float32 `json:"gravity_per_pound" binding:"required"`
	AuditId         string  `json:"audit_id" binding:"required"`
}

//NewSugarRouter creates a new instance of the Sugar router
func NewSugarRouter(app *gin.Engine) {
	r := SugarRouter{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
		gs:     services.NewSugarService(),
	}

	app.GET("/api/v1/sugars", r.getAllSugars)
	app.GET("/api/v1/sugars/:id", r.getSugar)
	app.POST("/api/v1/sugars", r.createSugar)
	app.DELETE("/api/v1/sugars/:id", r.deleteSugar)
}

func (r *SugarRouter) getAllSugars(ctx *gin.Context) {
	approved, err := strconv.ParseBool(ctx.Query("approved"))
	if err != nil {
		approved = true
	}
	auditId := ctx.Query("audit_id")
	sugars := make([]models.Sugar, 0)
	if auditId != "" {
		sugars = r.gs.GetAllSugarsByAuditId(approved, auditId)
	} else {
		sugars = r.gs.GetAllSugars(approved)
	}
	ctx.JSON(http.StatusOK, sugars)
}

func (r *SugarRouter) createSugar(ctx *gin.Context) {
	var data CreateSugarRequest
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body: some or all fields are missing or incorrect", "error": "INVALID_REQUEST_BODY"})
		return
	}
	g := r.gs.CreateSugar(data.Name, data.GravityPerPound, data.AuditId)
	ctx.JSON(http.StatusOK, g)
}

func (r *SugarRouter) getSugar(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "sugar not found", "error": "SUGAR_NOT_FOUND"})
		return
	}

	g, err := r.gs.GetSugarByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "sugar not found", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, g)
}

func (r *SugarRouter) deleteSugar(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "sugar not found", "error": "SUGAR_NOT_FOUND"})
		return
	}

	g, err := r.gs.DeleteSugarByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "sugar not found", "error": "SUGAR_NOT_FOUND"})
		return
	}
	ctx.JSON(http.StatusOK, g)
}
