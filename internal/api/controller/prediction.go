package controller

import (
	"encoding/json"
	"errors"
	"github.com/code7unner/vk-scrapper/internal/api/repository"
	"github.com/code7unner/vk-scrapper/internal/app"
	"github.com/code7unner/vk-scrapper/vw"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PredictionController interface {
	GetWithFilter(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type PredictionCtrl struct {
	app *app.App
}

func NewPredictionController(app *app.App) PredictionController {
	return &PredictionCtrl{app}
}

func (p PredictionCtrl) Get(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.error(w, r, http.StatusBadRequest, err)
	}
	r.Body.Close()

	formatData, err := vw.FormatData(string(data))
	if err != nil {
		p.error(w, r, http.StatusInternalServerError, err)
	}

	predict, err := p.app.Vws.Predict(formatData)
	if err != nil {
		p.error(w, r, http.StatusInternalServerError, err)
	}

	p.respond(w, r, http.StatusOK, predict)
}

func (p PredictionCtrl) GetWithFilter(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["city"]
	if keys == nil {
		p.error(w, r, http.StatusBadRequest, errors.New("city query is required"))
		return
	}
	city, _ := strconv.Atoi(keys[0])
	//school := keys[1]
	vkUserImpl := repository.NewVkUserImpl(p.app.Conf)
	users, err := vkUserImpl.GetVkUsers(city)
	if err != nil {
		p.error(w, r, http.StatusInternalServerError, err)
	}
	p.respond(w, r, http.StatusOK, users)
}

// respond with error
func (p PredictionCtrl) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	p.respond(w, r, code, map[string]string{"error": err.Error()})
}

// abstract respond
func (p PredictionCtrl) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
