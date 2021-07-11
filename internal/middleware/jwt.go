package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/RindangRamadhan/mnc-interview/internal/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/pt-mai/maihelper/mailog"
)

var unlessPath = []string{
	"/web/otp/send",
	"/web/setting/api-key/save",
	"/web/auth/register",
	"/web/auth/register/confirm/{zonid}/email",
	"/web/auth/register/{zonid}/confirm",
	"/web/auth/forgot-password",
	"/web/auth/change-password",
	"/web/auth/login",
	"/web/auth/logout",
}

var unlessPrefix = []string{
	"/v2",
	"/api-specs",
}

func (m *maiMiddleware) JWT(*mux.Router) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {

		return http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				var err error

				pathTpl, err := mux.CurrentRoute(r).GetPathTemplate()
				if err != nil {
					log.Println("Error get template router")
					helper.HttpHandler.ResponseJSON(rw, &helper.ResponseJSON{
						HTTPCode: http.StatusInternalServerError,
						Status:   "error",
						Message:  http.StatusText(http.StatusInternalServerError),
					})
					return
				}

				isUnlessPath := false
				isUnlessPrefix := false

				for _, val := range unlessPrefix {
					regex := regexp.MustCompile(`^` + val)
					if regex.MatchString(pathTpl) {
						isUnlessPrefix = true
						break
					}
				}

				// Unless Prefix
				if isUnlessPrefix {
					h.ServeHTTP(rw, r)
					return
				}

				for _, v := range unlessPath {
					if pathTpl == v {
						isUnlessPath = !isUnlessPath
						break
					}
				}

				// Unless Path
				if isUnlessPath {
					h.ServeHTTP(rw, r)
					return
				}

				cookie, err := r.Cookie(os.Getenv("AUTH_COOKIE_NAME"))
				if err != nil {
					log.Println("Error read auth cookie name :", mailog.Error(err))
					helper.HttpHandler.ResponseJSON(rw, &helper.ResponseJSON{
						HTTPCode: http.StatusUnauthorized,
						Status:   "error",
						Message:  http.StatusText(http.StatusUnauthorized),
					})
					return
				}

				cookieStr := cookie.Value

				token, err := jwt.Parse(cookieStr, func(token *jwt.Token) (interface{}, error) {
					// ! Don't forget to validate the alg what you expect
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(os.Getenv("AUTH_SECRET")), nil
				})

				if err != nil {
					log.Println("Error parse token :", mailog.Error(err))
					http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}

				if token == nil {
					helper.HttpHandler.ResponseJSON(rw, &helper.ResponseJSON{
						HTTPCode: http.StatusUnauthorized,
						Status:   "error",
						Message:  http.StatusText(http.StatusUnauthorized),
					})
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					helper.HttpHandler.ResponseJSON(rw, &helper.ResponseJSON{
						HTTPCode: http.StatusUnauthorized,
						Status:   "error",
						Message:  http.StatusText(http.StatusUnauthorized),
					})
					return
				}

				rw.Header().Set("X-MAI-User", claims["sub"].(string))

				h.ServeHTTP(rw, r)

			},
		)
	}
}
