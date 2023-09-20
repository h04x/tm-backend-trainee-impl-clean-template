package v1

import (
	"net/http"

	"tm-backend-trainee-impl-clean-template/internal/entity"
	"tm-backend-trainee-impl-clean-template/internal/usecase"
	"tm-backend-trainee-impl-clean-template/pkg/logger"

	"github.com/gin-gonic/gin"
)

type statisticsRoutes struct {
	t usecase.Statistics
	l logger.Interface
}

func newStatisticsRoutes(handler *gin.RouterGroup, t usecase.Statistics, l logger.Interface) {
	r := &statisticsRoutes{t, l}

	h := handler.Group("/statistics")
	{
		h.POST("/save", r.save)
		h.POST("/get", r.get)
		h.DELETE("/clear", r.clear)
	}
}

// @Summary     Save statistics
// @Description Save statistics
// @ID          save
// @Tags  	    statistics
// @Accept      json
// @Produce     json
// @Param       request body entity.Metrics true "Save statistics"
// @Success     200
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /statistics/save [post]
func (r *statisticsRoutes) save(c *gin.Context) {
	request := entity.Metrics{Cost: "0"}
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - save")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.t.Save(c, request); err != nil {
		r.l.Error(err, "http - v1 - save")
		errorResponse(c, http.StatusInternalServerError, "error occur while saving")
		return
	}

	c.Status(http.StatusAccepted)
}

// @Summary     Get statistics
// @Description Get statistics
// @ID          get
// @Tags  	    statistics
// @Accept      json
// @Produce     json
// @Param       request body entity.DoGetRequest true "Get statistics"
// @Success     200 {object} []entity.Statistics
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /statistics/get [post]
func (r *statisticsRoutes) get(c *gin.Context) {
	request := entity.DoGetRequest{Order: "Date"}
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	s, err := r.t.Get(c, request)
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, http.StatusInternalServerError, "error occure while get statitstics")
		return
	}

	c.JSON(http.StatusOK, s)
}

// @Summary     Clear statistics
// @Description Clear statistics
// @ID          clear
// @Tags  	    statistics
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response
// @Router      /statistics/clear [delete]
func (r *statisticsRoutes) clear(c *gin.Context) {
	if err := r.t.Clear(c); err != nil {
		r.l.Error(err, "http - v1 - clean")
		errorResponse(c, http.StatusInternalServerError, "error occure while clean")
		return
	}
}
