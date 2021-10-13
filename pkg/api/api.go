package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

const logfile = "./logfile.txt"

// Программный интерфейс приложения
type API struct {
	r *mux.Router
}

// Конструктор объекта API
func New() *API {
	api := API{}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	//метод проверки комментария
	api.r.HandleFunc("/cens", api.cens).Methods(http.MethodPost)
	//мидлварь для сквозной идентификации и логгирования
	api.r.Use(api.idLogger)
}

//мидлварь для сквозной идентификации и логгирования
//?request_id=327183798123
func (api *API) idLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logfile, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			http.Error(w, fmt.Sprintf("os.OpenFile error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer logfile.Close()
		id := r.URL.Query().Get("request_id")
		if id == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				http.Error(w, fmt.Sprintf("uuid.NewV4 error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			id = uid.String()
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "request_id", id)
		r = r.WithContext(ctx)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)

		fmt.Fprintf(logfile, "Request ID:%s\n", id)
		fmt.Fprintf(logfile, "Time:%s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(logfile, "Remote IP address:%s\n", r.RemoteAddr)
		fmt.Fprintf(logfile, "HTTP Status:%d\n", rec.Result().StatusCode)
		fmt.Fprintln(logfile)
	})
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

//список запрещенных слов
var wordList = []string{"qwerty", "йцукен", "zxvbnm"}

//проверка наличия подстрок subs в строке str
func censured(str string, subs ...string) bool {
	censured := false
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			censured = true
		}
	}
	return censured
}

//метод проверки комментария
//localhost:8083/cens
func (api *API) cens(w http.ResponseWriter, r *http.Request) {
	bComment, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("censrv ReadAll error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if censured(string(bComment), wordList...) {
		http.Error(w, "censrv censured error", http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
