// surrealdb/surrealdb.go

package surrealdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os" // Usamos solo el paquete os para obtener las variables de entorno

	"github.com/oldronald/surrealdb/backend/config"
)

// Función para obtener el token de SurrealDB
func GetSurrealToken() (string, error) {
	// Configurar datos de autenticación
	authData := map[string]string{
		"ns":   os.Getenv("DB_NS"),   // Cargar valor de la variable de entorno DB_NS
		"db":   os.Getenv("DB_NAME"), // Cargar valor de la variable de entorno DB_NAME
		"user": os.Getenv("DB_USER"), // Cargar valor de la variable de entorno DB_USER
		"pass": os.Getenv("DB_PASS"), // Cargar valor de la variable de entorno DB_PASS
	}

	// Verificar si las variables de entorno están vacías
	if authData["ns"] == "" || authData["db"] == "" || authData["user"] == "" || authData["pass"] == "" {
		return "", fmt.Errorf("una o más variables de entorno no están configuradas correctamente")
	}

	// Convertir los datos de autenticación a JSON
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("error al convertir a JSON: %v", err)
	}

	// Crear la solicitud POST
	req, err := http.NewRequest("POST", "https://"+config.Credentials.Host+"/signin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err != nil {
		return "", fmt.Errorf("Error al crear la solicitud: %v", err)
	}

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al enviar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Decodificar la respuesta JSON
	var surrealResp map[string]interface{}
	err = json.Unmarshal(body, &surrealResp)
	if err != nil {
		return "", fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
	}

	// Verificar si hay un token en la respuesta
	token, ok := surrealResp["token"].(string)
	if !ok {
		return "", fmt.Errorf("no se recibió token en la respuesta")
	}

	// Imprimir el token
	fmt.Println("Token:", token)

	// Almacenar el token en las credenciales
	config.Credentials.Token = token

	// Devuelve el token
	return token, nil
}
