package handlers

import (
	"cozy-doc-api/models"
	"cozy-doc-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BulkInsertDocs(s services.DocsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(models.DocumentRequest)
		err := ctx.Bind(req)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()})
			return
		}

		if err := req.Validate(); err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()})
			return
		}
		err = s.BulkInsertDocs(req)

		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": true})
	}
}
