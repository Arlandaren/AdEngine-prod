package advertisers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"service/internal/domains/advertisers/models"
)

type Controller struct {
	svc *Service
}

func NewAdvertisersController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (cont *Controller) Endpoints(r *gin.Engine) {
	r.GET("/advertisers/:advertiser_id", cont.GetAdvertiserByID)
	r.POST("/advertisers/bulk", cont.UpsertAdvertisers)
	r.POST("/ml-scores", cont.UpsertMLScore)
}

func (cont *Controller) GetAdvertiserByID(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid advertiser ID"})
		return
	}

	advertiser, err := cont.svc.GetAdvertiserByID(c.Request.Context(), advertiserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, advertiser)
}

func (cont *Controller) UpsertAdvertisers(c *gin.Context) {
	var req []models.AdvertiserUpsert

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and parse UUIDs
	for _, adv := range req {
		if adv.AdvertiserID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "advertiser_id is required"})
			return
		}
	}

	err := cont.svc.UpsertAdvertisers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (cont *Controller) UpsertMLScore(c *gin.Context) {
	var req models.MLScore

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate UUIDs
	if req.ClientID == uuid.Nil || req.AdvertiserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id and advertiser_id are required"})
		return
	}

	err := cont.svc.UpsertMLScore(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
