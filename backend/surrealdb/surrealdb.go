package surrealdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oldronald/surrealdb/backend/config"
)

// Función para obtener el token de SurrealDB
func GetSurrealToken(ns, db, user, pass string) (string, error) {
	// Configurar datos de autenticación
	authData := map[string]string{
		"ns":   ns,   // Namespace
		"db":   db,   // Base de datos
		"user": user, // Usuario
		"pass": pass, // Contraseña
	}

	// Verificar si las credenciales están vacías
	if authData["ns"] == "" || authData["db"] == "" || authData["user"] == "" || authData["pass"] == "" {
		return "", fmt.Errorf("credenciales incompletas: ns, db, user y pass son obligatorios")
	}

	// Convertir los datos de autenticación a JSON
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("error al convertir a JSON: %v", err)
	}

	// Construir la URL
	url := fmt.Sprintf("https://%s/signin", config.Credentials.Host) // Usar config.Credentials.Host directamente

	// Crear la solicitud POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al enviar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta inesperada: %s", resp.Status)
	}

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Decodificar la respuesta JSON
	type SurrealResponse struct {
		Token string `json:"token"`
	}

	var surrealResp SurrealResponse
	err = json.Unmarshal(body, &surrealResp)
	if err != nil {
		return "", fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
	}

	// Verificar si hay un token en la respuesta
	token := surrealResp.Token
	if token == "" {
		return "", fmt.Errorf("no se recibió token en la respuesta")
	}

	// Imprimir el token (opcional, para depuración)
	fmt.Println("Token:", token)

	// Almacenar el token en las credenciales
	config.Credentials.Token = token

	// Devolver el token
	return token, nil
}
