package omnis

import (
	"github.com/gin-gonic/gin"
)

type ICorrelationService interface {
	WithContext(ctx gin.Context) ICorrelationService
	SetCorrelationID() (string, error)
	GetCorrelationID() string
}
