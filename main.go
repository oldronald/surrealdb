package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oldronald/surrealdb/backend/handler"
)

func main() {

	// Configurar el puerto desde la variable de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto si no se especifica
	}

	// Leer variables de entorno para la configuración de la base de datos
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	apiKey := os.Getenv("API_KEY")

	// Imprimir las variables de entorno para verificar
	fmt.Printf("DB Host: %s\n", dbHost)
	fmt.Printf("DB User: %s\n", dbUser)
	fmt.Printf("API Key: %s\n", apiKey)

	// Calcular la longitud de la contraseña (sin imprimirla)
	passLength := len(dbPass)
	fmt.Printf("Longitud de la contraseña: %d\n", passLength)

	http.HandleFunc("/signin", handler.MainHandler) // Usa el MainHandler del paquete handlers

	// Configurar los manejadores para las rutas
	fmt.Printf("Servidor HTTPS corriendo en https://localhost:%s\n", port)

	// Iniciar el servidor HTTPS
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor HTTPS: %v", err)
	}
}
