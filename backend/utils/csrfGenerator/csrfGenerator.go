package csrfgenerator

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	csrfCookieName = "XSRF-TOKEN"
	csrfHeaderName = "X-XSRF-TOKEN"
)

func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HandleCsrf(c *gin.Context) {
	token, err := GenerateCSRFToken()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		csrfCookieName,
		token,
		3600,
		"/",
		"",
		false, // Secure
		false, // HttpOnly=false
	)

	c.Status(http.StatusNoContent)
}
