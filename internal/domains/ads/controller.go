package ads

import (
	"errors"
	"net/http"
	"service/internal/infrastructure/storage/models/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	svc *Service
}

func NewAdsController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (cont *Controller) Endpoints(r *gin.Engine) {
	r.GET("/ads", cont.GetAdForClient)
	r.POST("/ads/:ad_id/click", cont.RecordAdClick)
}

func (cont *Controller) GetAdForClient(c *gin.Context) {
	clientIDStr := c.Query("client_id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	ad, err := cont.svc.GetAdForClient(c.Request.Context(), clientID)
	if err != nil {
		if errors.Is(err, dto.ErrNoAds) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ad)
}

func (cont *Controller) RecordAdClick(c *gin.Context) {
	adIDStr := c.Param("ad_id")
	adID, err := uuid.Parse(adIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ad ID"})
		return
	}

	var req struct {
		ClientID string `json:"client_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	err = cont.svc.RecordAdClick(c.Request.Context(), adID, clientID)
	if err != nil {
		if err.Error() == "ad not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
