package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("Internal Server Error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Errorw("internal server error", " method: ", r.Method, " path: ", r.URL.Path, " error: ", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "The server encontered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("Bad Request Error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Warnf("Bad Request Error:", " method: ", r.Method, " path: ", r.URL.Path, " error: ", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) statusNotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("Status Not Found Error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Warnf("Status Not Found Error:", " method: ", r.Method, " path: ", r.URL.Path, " error: ", err.Error())
	writeJSONError(w, http.StatusNotFound, "Not found")
}
