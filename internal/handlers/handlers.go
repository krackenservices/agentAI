package handlers

import (
	"encoding/json"
	"net/http"
)

// HelloHandler godoc
// @Summary Returns a greeting message
// @Description Returns a simple hello world JSON message.
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello, world!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
