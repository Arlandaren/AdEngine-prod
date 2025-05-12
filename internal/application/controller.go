package application

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"service/internal/domains/ads"
	"service/internal/domains/advertisers"
	"service/internal/domains/campaigns"
	"service/internal/domains/clients"
	"service/internal/domains/statistics"

	"time"

	"service/internal/domains/api"
)

type Controller struct {
	api         *api.Controller
	clients     *clients.Controller
	advertisers *advertisers.Controller
	campaigns   *campaigns.Controller
	ads         *ads.Controller
	stats       *statistics.Controller
	Router      *gin.Engine
}

func NewController(svc *Service, r *gin.Engine) *Controller {
	return &Controller{
		Router:      r,
		api:         api.NewApiController(svc.api),
		clients:     clients.NewClientsController(svc.clients),
		advertisers: advertisers.NewAdvertisersController(svc.advertisers),
		campaigns:   campaigns.NewCampaignsController(svc.campaigns),
		stats:       statistics.NewStatisticsController(svc.stats),
		ads:         ads.NewAdsController(svc.ads),
	}
}

func (c *Controller) InitRouter() {
	c.api.Endpoints(c.Router)
	c.clients.Endpoints(c.Router)
	c.advertisers.Endpoints(c.Router)
	c.campaigns.Endpoints(c.Router)
	c.ads.Endpoints(c.Router)
	c.stats.Endpoints(c.Router)
}

func (c *Controller) Run(addr string, ctx context.Context) {
	c.InitRouter()

	server := &http.Server{
		Addr:    addr,
		Handler: c.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server exited with error: %v", err)
		}
	}()
	log.Printf("HTTP server listening at %v", addr)

	<-ctx.Done()
	log.Println("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorf("HTTP server Shutdown Failed:%+v", err)
	} else {
		log.Println("HTTP server gracefully stopped")
	}
}
