package routers

import (
	"server/pkg/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/ok", handlers.OKHandler).Methods("GET")
    router.HandleFunc("/test", handlers.TestHandler).Methods("GET")
    router.HandleFunc("/cityclimate", handlers.FetchCityClimate).Methods("GET")
    router.HandleFunc("/cityclimategps", handlers.ListCityClimateSensorsBasedOnDistance).Methods("POST")
    router.HandleFunc("/meteoblue", handlers.FetchMeteoBlue).Methods("GET")
    router.HandleFunc("/hotareas", handlers.FetchAndReportHotAreas).Methods("GET")
    router.HandleFunc("/hotareasgps", handlers.FetchAndReportHotAreasBasedOnLocation).Methods("POST")
    router.HandleFunc("/speak", handlers.TTSTest).Methods("POST")
    router.HandleFunc("/weather", handlers.FetchAndSpeakWeatherData).Methods("POST")
    router.HandleFunc("/weathergps", handlers.FetchAndSpeakWeatherBasedOnGPS).Methods("POST")
    return router
}
