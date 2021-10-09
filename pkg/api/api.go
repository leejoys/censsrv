package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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
