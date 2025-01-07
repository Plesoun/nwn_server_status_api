package handlers

import (
	"encoding/json"
	"net/http"
	"nwn_server_info/udp"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JSONResponse struct {
	UUID    string      `json:uuid`
	Result  string      `json:result,omitempty`
	Message string      `json:message,omitempty`
	Data    interface{} `json:data,omitempty`
}

func (jr *JSONResponse) CustomResponse(logger *logrus.Logger, w http.ResponseWriter, statusCode int) {
	jr.UUID = uuid.New().String()

	jsonResponse, err := json.Marshal(jr)

	if err != nil {
		logger.WithField("UUID", uuid.New().String()).Errorf("encountered an error while marshalling json response, %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func GetPlayersOnline(logger *logrus.Logger, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ip := queryParams.Get("ip")
	port := queryParams.Get("port")

	if ip == "" {
		jsonResponse := &JSONResponse{Result: "ERROR", Message: "No ip provided in the parameters"}
		jsonResponse.CustomResponse(logger, w, http.StatusBadRequest)
		logger.WithField("UUID", jsonResponse.UUID).Errorf("No ip provided in the parameters")
		return
	}
	if port == "" {
		jsonResponse := &JSONResponse{Result: "ERROR", Message: "No port provided in the parameters"}
		jsonResponse.CustomResponse(logger, w, http.StatusBadRequest)
		logger.WithField("UUID", jsonResponse.UUID).Errorf("No port provided in the parameters")
		return
	}
	nwnOnline, serverInfo := udp.CheckNWNServer(ip, port)

	if nwnOnline {
		jsonResponse := &JSONResponse{Result: "OK", Data: map[string]string{"players_online": string(serverInfo.PlayersOnline)}}
		jsonResponse.CustomResponse(logger, w, http.StatusOK)
		logger.WithField("UUID", jsonResponse.UUID).Debugf("Players online: %s", string(serverInfo.PlayersOnline))

	} else {
		jsonResponse := &JSONResponse{Result: "OK", Message: "Failed to get server info"}
		jsonResponse.CustomResponse(logger, w, http.StatusNoContent)
		logger.WithField("UUID", jsonResponse.UUID).Errorf("Failed to get server info")
		return
	}

}
