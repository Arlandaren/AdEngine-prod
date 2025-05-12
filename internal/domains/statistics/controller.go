package statistics

import (
	"errors"
	"net/http"
	"service/internal/domains/statistics/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	svc *Service
}

func NewStatisticsController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (cont *Controller) Endpoints(r *gin.Engine) {
	r.GET("/stats/campaigns/:campaign_id", cont.GetCampaignStats)
	r.GET("/stats/advertisers/:advertiser_id/campaigns", cont.GetAdvertiserCampaignsStats)
	r.GET("/stats/campaigns/:campaign_id/daily", cont.GetCampaignDailyStats)
	r.GET("/stats/advertisers/:advertiser_id/campaigns/daily", cont.GetAdvertiserDailyStats)
}

// GetCampaignStats GET /api/stats/campaigns/:campaign_id
func (cont *Controller) GetCampaignStats(c *gin.Context) {
	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	stats, err := cont.svc.GetCampaignStats(c.Request.Context(), campaignID)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := &models.StatsResponse{
		ImpressionsCount: stats.ImpressionsCount,
		ClicksCount:      stats.ClicksCount,
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}

	c.JSON(http.StatusOK, resp)
}

// GetAdvertiserCampaignsStats gET /api/stats/advertisers/:advertiser_id/campaigns
func (cont *Controller) GetAdvertiserCampaignsStats(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	stats, err := cont.svc.GetAdvertiserCampaignsStats(c.Request.Context(), advertiserID)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := &models.StatsResponse{
		ImpressionsCount: stats.ImpressionsCount,
		ClicksCount:      stats.ClicksCount,
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}

	c.JSON(http.StatusOK, resp)
}

// GetCampaignDailyStats GET /api/stats/campaigns/:campaign_id/daily
func (cont *Controller) GetCampaignDailyStats(c *gin.Context) {
	campaignIDStr := c.Param("campaign_id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	dailyStats, err := cont.svc.GetCampaignDailyStats(c.Request.Context(), campaignID)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]models.DailyStatsResponse, len(dailyStats))
	for i, obj := range dailyStats {
		resp[i] = models.DailyStatsResponse{
			Date:             obj.Date,
			ImpressionsCount: obj.Stats.ImpressionsCount,
			ClicksCount:      obj.Stats.ClicksCount,
			Conversion:       obj.Stats.Conversion,
			SpentImpressions: obj.Stats.SpentImpressions,
			SpentClicks:      obj.Stats.SpentClicks,
			SpentTotal:       obj.Stats.SpentTotal,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetAdvertiserDailyStats gtT /api/stats/advertisers/:advertiser_id/campaigns/daily
func (cont *Controller) GetAdvertiserDailyStats(c *gin.Context) {
	advertiserIDStr := c.Param("advertiser_id")
	advertiserID, err := uuid.Parse(advertiserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertiser ID"})
		return
	}

	dailyStats, err := cont.svc.GetAdvertiserDailyStats(c.Request.Context(), advertiserID)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := make([]models.DailyStatsResponse, len(dailyStats))
	for i, obj := range dailyStats {
		resp[i] = models.DailyStatsResponse{
			Date:             obj.Date,
			ImpressionsCount: obj.Stats.ImpressionsCount,
			ClicksCount:      obj.Stats.ClicksCount,
			Conversion:       obj.Stats.Conversion,
			SpentImpressions: obj.Stats.SpentImpressions,
			SpentClicks:      obj.Stats.SpentClicks,
			SpentTotal:       obj.Stats.SpentTotal,
		}
	}

	c.JSON(http.StatusOK, resp)
}
