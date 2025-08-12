package payments

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

var merchantLogin string
var firstPass string
var secondPass string

func InitRobokassa() {
	merchantLogin = os.Getenv("ROBOKASSA_LOGIN")
	firstPass = os.Getenv("ROBOKASSA_FIRST_PASSWORD")
	secondPass = os.Getenv("ROBOKASSA_SECOND_PASSWORD")
}

func HandlePayment(c *gin.Context) {
	println(merchantLogin, firstPass, secondPass)
	grade, gradeErr := strconv.Atoi(c.Param("grade"))
	println(grade)
	if gradeErr != nil {
		println("GRADE ERROR")
		c.HTML(http.StatusBadRequest, "404.html", gin.H{})
		return
	}
	var amount float64
	if grade <= 8 {
		amount = 1590.0
	} else {
		amount = 2590.0
	}
	access_cookie, cookieErr := c.Cookie("access_token")
	if cookieErr != nil {
		println("COOKIE ERR")
		c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
		return
	}
	println(access_cookie)
	user_info, userErr := jwtconfigurator.ValidateAccessToken(access_cookie)
	if userErr != nil {
		log.Println(userErr)
		c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
		return
	}
	invID := generateRoboID()
	println(invID)
	signature := generateSignature(invID, amount, user_info.UserID, grade)
	println(signature)
	description := fmt.Sprintf("Покупка %s класс", c.Param("grade"))
	paymentURL := fmt.Sprintf("https://auth.robokassa.ru/Merchant/Index.aspx?"+
		"MerchantLogin=%s&"+
		"InvId=%d&"+
		"OutSum=%.2f&"+
		"Description=%s&"+
		"SignatureValue=%s&"+
		"IsTest=1&"+ // 1 для теста, 0 для боевого режима
		"Shp_grade=%d&"+
		"Shp_userID=%s",
		merchantLogin,
		invID,
		amount,
		url.QueryEscape(description),
		signature,
		grade,
		url.QueryEscape(user_info.UserID),
	)
	println(paymentURL)
	c.Redirect(http.StatusFound, paymentURL)
}

func generateSignature(invID int, amount float64, userID string, grade int) string {
	shpParams := map[string]string{
		"grade":  strconv.Itoa(grade),
		"userID": userID,
	}

	// Получаем и сортируем ключи параметров
	keys := make([]string, 0, len(shpParams))
	for k := range shpParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	data := fmt.Sprintf("%s:%.2f:%d:%s",
		merchantLogin, amount, invID, firstPass)

	for _, key := range keys {
		data += fmt.Sprintf(":Shp_%s=%s", key, shpParams[key])
	}
	println(data)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateRoboID() int {
	return rand.Intn(900000) + 100000
}
