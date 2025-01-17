package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"surrealdb/backend/config"
	"surrealdb/backend/encriptacion"
	"surrealdb/backend/ratelimiter"
	"surrealdb/backend/surrealdb"
)

type Response struct {
	Message     string                      `json:"message"`
	Credentials config.SurrealDBCredentials `json:"credentials"`
}

// Handler principal de la API
func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	key := r.Header.Get("Authorization")
	if key == "" {
		resp := Response{Message: "No se proporcionó ninguna clave"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Verificar que la clave comience con "Bearer "
	if !strings.HasPrefix(key, "Bearer ") {
		resp := Response{Message: "Formato de clave inválido"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Extraer la clave secreta
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

	// Validar las credenciales (esto es un ejemplo simple, ajusta según tus necesidades)
	if requestBody.Ns != config.Credentials.Ns || requestBody.Db != config.Credentials.Db || requestBody.User != config.Credentials.User || requestBody.Pass != config.Credentials.Pass {
		resp := Response{Message: "Credenciales inválidas"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tokens, err := surrealdb.GetSurrealToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener el token: %v", err), http.StatusInternalServerError)
		return
	}

	encryptedToken, err := encriptacion.Encrypter(tokens, config.EncryptionKey) // Asegúrate de que `tokens` sea un string
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al encriptar el token: %v", err), http.StatusInternalServerError)
		return
	}

	config.Credentials.Token = encryptedToken

	errcredential := encriptacion.EncryptAllCredentials(&config.Credentials, config.EncryptionKey)
	if errcredential != nil {
		fmt.Printf("Error al encriptar las credenciales: %v\n", err)
		return
	}

	fmt.Println("Tokenencryted:", encryptedToken)

	dencryptedToken, err := encriptacion.Decrypter(encryptedToken, config.EncryptionKey) // Asegúrate de que `tokens` sea un string
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al encriptar el token: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("TokenDencryted:", dencryptedToken)

	resp := Response{
		Message:     "Clave valida",
		Credentials: config.Credentials,
	}
	json.NewEncoder(w).Encode(resp)
}
