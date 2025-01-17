package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/oldronald/surrealdb/backend/config"
	"github.com/oldronald/surrealdb/backend/encriptacion"
	"github.com/oldronald/surrealdb/backend/ratelimiter"
	"github.com/oldronald/surrealdb/backend/surrealdb"
)

type Response struct {
	Message     string                      `json:"message"`
	Credentials config.SurrealDBCredentials `json:"credentials"`
}

// Handler principal de la API
func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar el límite de tasa
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Error obteniendo la dirección IP", http.StatusInternalServerError)
		return
	}

	if !ratelimiter.CheckRateLimit(ip) {
		resp := Response{Message: "Excediste el número de solicitudes permitidas. Intenta más tarde."}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Verificar el encabezado Authorization
	key := r.Header.Get("Authorization")
	if key == "" {
		resp := Response{Message: "No se proporcionó ninguna clave"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !strings.HasPrefix(key, "Bearer ") {
		resp := Response{Message: "Formato de clave inválido"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	secretKey := strings.TrimPrefix(key, "Bearer ")
	if secretKey != config.JWTSecret {
		resp := Response{Message: "Clave inválida"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Leer el cuerpo de la solicitud
	var requestBody struct {
		Ns   string `json:"ns"`
		Db   string `json:"db"`
		User string `json:"user"`
		Pass string `json:"pass"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Asignar las credenciales del cuerpo de la solicitud a config.Credentials
	config.Credentials.Ns = requestBody.Ns
	config.Credentials.Db = requestBody.Db
	config.Credentials.User = requestBody.User
	config.Credentials.Pass = requestBody.Pass

	// Obtener el token de SurrealDB
	tokens, err := surrealdb.GetSurrealToken(requestBody.Ns, requestBody.Db, requestBody.User, requestBody.Pass)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener el token: %v", err), http.StatusInternalServerError)
		return
	}

	// Encriptar el token
	encryptedToken, err := encriptacion.Encrypter(tokens, config.EncryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al encriptar el token: %v", err), http.StatusInternalServerError)
		return
	}

	// Asignar el token encriptado a config.Credentials
	config.Credentials.Token = encryptedToken

	// Responder con un mensaje de éxito y las credenciales
	resp := Response{
		Message:     "Autenticación exitosa",
		Credentials: config.Credentials,
	}
	json.NewEncoder(w).Encode(resp)
}
