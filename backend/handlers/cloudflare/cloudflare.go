package cloudflare

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

var turnstileSecret string

func InitSecret() {
	turnstileSecret = os.Getenv("TURNSTILESECRET")
}

func verifyTurnstile(response, remoteIP string) (bool, error) {
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify",
		url.Values{"secret": {turnstileSecret}, "response": {response}, "remoteip": {remoteIP}})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var tr TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return false, err
	}

	return tr.Success, nil
}

func CloudflareHandler(c *gin.Context) {
	turnstileResponse := c.PostForm("cf-turnstile-response")
	clientIP := c.ClientIP()

	success, err := verifyTurnstile(turnstileResponse, clientIP)
	if err != nil || !success {
		c.JSON(http.StatusForbidden, gin.H{"error": "Капча не пройдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Валидация капчи успешна"})
}
