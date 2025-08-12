package tokenchecker

import (
	"context"
	godb "formular/backend/database/SQL_postgre"
	user "formular/backend/models/userConfig"
	"formular/backend/utils/jwtconfigurator"
	"time"

	"github.com/gin-gonic/gin"
)

// Проверяет access токен и обновляет его, возвращая пользователя
func ValidateAccessTokenWithRefresh(c *gin.Context) (user.User, error) {
	access_cookie, accessCookieErr := c.Cookie("access_token")
	if accessCookieErr != nil && access_cookie == "" {
		refresh_cookie, refreshCookieErr := c.Cookie("refresh_token")
		if refreshCookieErr != nil {
			return user.User{}, refreshCookieErr
		}
		new_access, genErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
		if genErr != nil {
			return user.User{}, genErr
		}
		c.SetCookie("access_token", new_access, 8*60*60, "/", "127.0.0.1", false, true)
		user_info_db, usErr := getUserInfo(new_access)
		if usErr != nil {
			return user.User{}, usErr
		}
		return user_info_db, nil

	} else {
		user_info_db, usErr := getUserInfo(access_cookie)
		if usErr != nil {
			return user.User{}, usErr
		}
		return user_info_db, nil
	}
}

func getUserInfo(access_cookie string) (user.User, error) {
	user_info, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookie)
	if jwtErr != nil {
		return user.User{}, jwtErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user_info_by_id, dbErr := godb.GetUserInfoByID(ctx, user_info.UserID)
	if dbErr != nil {
		return user.User{}, dbErr
	}
	return *user_info_by_id, nil
}
