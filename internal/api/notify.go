package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eztwokey/l3-serv/internal/logic"
	"github.com/eztwokey/l3-serv/internal/models"
	"github.com/eztwokey/l3-serv/internal/storage"
)

func (a *Api) createNotify(c *gin.Context) {
	var req models.CreateNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("notify.create bind error", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	n, err := a.logic.CreateNotify(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusCreated, n)
}

func (a *Api) getNotify(c *gin.Context) {
	id := c.Param("id")

	n, err := a.logic.GetNotify(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, n)
}

func (a *Api) cancelNotify(c *gin.Context) {
	id := c.Param("id")

	n, err := a.logic.CancelNotify(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, n)
}
