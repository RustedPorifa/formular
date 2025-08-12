package kims

import (
	"context"
	"errors"
	godb "formular/backend/database/SQL_postgre"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleGrade(c *gin.Context) {
	grade := c.Param(":grade")
	type_of := c.Param(":type")
	if strings.Contains(type_of, "solved") {
		is_allowed, allowErr := haveAccess(grade, type_of, c)
		if allowErr != nil && !is_allowed {
			log.Println(allowErr)
			c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
			return
		} else {
			path_to_kims := filepath.Join("frontend", "templates", "math", grade, type_of)
			_, readErr := os.ReadDir(path_to_kims) //folders добавить
			if readErr != nil {
				c.JSON(http.StatusInternalServerError, readErr.Error())
				return
			}

		}
	}

}

func haveAccess(grade string, type_of string, c *gin.Context) (bool, error) {
	access_cookie, cookieErr := c.Cookie("access_token")
	if cookieErr != nil {
		return false, cookieErr
	}
	user_info, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookie)
	if jwtErr != nil {
		return false, jwtErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	access_to_grades, dbErr := godb.GetPurchasedGradesByUserID(ctx, user_info.UserID)
	if dbErr != nil {
		return false, dbErr
	}
	if slices.Contains(access_to_grades, grade) {
		return true, nil
	}
	return false, errors.New("нет доступа")
}
