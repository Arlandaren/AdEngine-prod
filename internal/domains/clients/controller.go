package clients

import (
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	"service/internal/domains/clients/models"
)

type Controller struct {
	svc *Service
}

func NewClientsController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (cont *Controller) Endpoints(r *gin.Engine) {
	r.GET("/clients/:client_id", cont.GetClientById)
	r.POST("/clients/bulk", cont.UpsertClients)
}

func (cont *Controller) GetClientById(c *gin.Context) {
	clientIDStr := c.Param("client_id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format of id doesnt match an uuid format"})
		return
	}

	client, err := cont.svc.GetClientByID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (cont *Controller) UpsertClients(c *gin.Context) {
	var req []models.ClientUpsert

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cont.svc.UpsertClients(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}
