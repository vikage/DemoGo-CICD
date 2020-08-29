package restful

import (
	"context"
	"fmt"
	"go-cicd/app/authenticate"
	"go-cicd/app/domain/model"
	"go-cicd/app/logger"
	"go-cicd/app/restful/base"
	"go-cicd/app/restful/handlers/account"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
}

// StartWebAPI start http server
func StartWebAPI() {
	mainRouter := mux.NewRouter()
	mainRouter.NotFoundHandler = NotFoundHandler()

	v1Router := mainRouter.PathPrefix("/api/v1").Subrouter()

	preAuthenticateRouter := v1Router.NewRoute().Subrouter()
	preAuthenticateRouter.Use(JwtVerify)

	account.RegisterHandlerForAccountAPI(v1Router, preAuthenticateRouter)

	port := os.Getenv("REST_API_PORT")
	addr := fmt.Sprintf(":%s", port)
	go func() {
		log.Fatal(http.ListenAndServe(addr, RequestHandler(mainRouter)))
	}()

	logger.Debug("Started API at port %s", port)
}

// RequestHandler handle http handler
func RequestHandler(handler *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, PATCH, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.WriteHeader(200)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")

		recorder := httptest.NewRecorder()
		responseWriter := statusWriter{
			ResponseWriter: recorder,
			status:         200,
		}

		startTime := time.Now()
		handler.ServeHTTP(&responseWriter, r)
		for key, headerValues := range recorder.Header() {
			for _, value := range headerValues {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(responseWriter.status)
		w.Write(recorder.Body.Bytes())

		duration := int((time.Now().UnixNano() - startTime.UnixNano()) / 1000000)

		logger.LogReq(r, duration, responseWriter.status)
	})
}

// JwtVerify verify token
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr = r.Header.Get("token")

		if len(tokenStr) == 0 {
			logger.ErrorReq(r, "Missing token")
			w.WriteHeader(http.StatusBadRequest)
			base.ResponseError(model.APIErrorBadRequest, "Missing token", w)
			return
		}

		decoder := authenticate.ResolveTokenDecoder()
		user, err := decoder.UserFromToken(tokenStr)
		if err != nil {
			logger.ErrorReq(r, "Jwt Verify fail withe error %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), base.UserAuthenticatedKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NotFoundHandler handle case 404
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Request at path %s not found", r.URL.Path)
		response := model.APIResponse{ErrorCode: 404, Message: message}
		w.WriteHeader(http.StatusNotFound)

		base.ResponseToClient(&response, w)
	})
}
