package campaigns

import (
	"errors"
	"net/http"
	"service/internal/infrastructure/storage/models/dto"
	"strconv"

	"service/internal/domains/campaigns/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	svc *Service
}

func NewCampaignsController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (cont *Controller) Endpoints(r *gin.Engine) {
	group := r.Group("/advertisers/:advertiser_id/campaigns")
	{
		group.POST("/", cont.CreateCampaign)
		group.GET("/", cont.ListCampaigns)
		group.GET("/:campaign_id", cont.GetCampaignByID)
		group.PUT("/:campaign_id", cont.UpdateCampaign)
		group.DELETE("/:campaign_id", cont.DeleteCampaign)
	}
}

func (cont *Controller) CreateCampaign(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	var campaignCreate models.CampaignCreate
	if err := c.ShouldBindJSON(&campaignCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = campaignCreate.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := cont.svc.CreateCampaign(c.Request.Context(), advertiserID, &campaignCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

func (cont *Controller) UpdateCampaign(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var campaignUpdate models.CampaignUpdate
	if err := c.ShouldBindJSON(&campaignUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := cont.svc.GetCampaignByID(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	if campaign.AdvertiserID != advertiserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Campaign does not belong to the advertiser"})
		return
	}

	updatedCampaign, err := cont.svc.UpdateCampaign(c.Request.Context(), campaignID, &campaignUpdate)
	if err != nil {
		if errors.Is(err, dto.ErrActionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCampaign)
}

func (cont *Controller) GetCampaignByID(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := cont.svc.GetCampaignByID(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	if campaign.AdvertiserID != advertiserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Campaign does not belong to the advertiser"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func (cont *Controller) ListCampaigns(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	sizeStr := c.Query("size")
	pageStr := c.Query("page")

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10 // Default size
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default page
	}

	campaigns, err := cont.svc.ListCampaigns(c.Request.Context(), advertiserID, int32(size), int32(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (cont *Controller) DeleteCampaign(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := cont.svc.GetCampaignByID(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	if campaign.AdvertiserID != advertiserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Campaign does not belong to the advertiser"})
		return
	}

	err = cont.svc.DeleteCampaign(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
