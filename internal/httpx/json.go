package httpx

import (
	"encoding/json"
	"net/http"
)
type errorEnvelope struct { Error struct { Code string `json:"code"`; Message string `json:"message"` } `json:"error"` }
func JSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json; charset=utf-8");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func Error(w http.ResponseWriter,status int,code,message string){var e errorEnvelope;e.Error.Code=code;e.Error.Message=message;JSON(w,status,e)}
