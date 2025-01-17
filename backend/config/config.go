package config

import (
	"os"
)

// Secreto JWT (cargado desde las variables de entorno)
var JWTSecret string
var EncryptionKey []byte

func init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "default_jwt_secret" // Valor predeterminado si la variable de entorno no está definida
	}

	EncryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(EncryptionKey) == 0 {
		EncryptionKey = []byte("default_encryption_key") // Valor predeterminado si la variable de entorno no está definida
	}
}

// SurrealDBCredentials almacena las credenciales para conectarse a la base de datos
type SurrealDBCredentials struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	Ns        string `json:"ns"`
	Db        string `json:"db"`
	Protocolo string `json:"protocolo"`
	Token     string `json:"token"`
}

// Credentials contiene las credenciales cargadas desde las variables de entorno
var Credentials = SurrealDBCredentials{
	Host:      getEnv("DB_HOST", "localhost"),
	Port:      getEnv("DB_PORT", "8000"),
	User:      getEnv("DB_USER", "default_user"),
	Pass:      getEnv("DB_PASS", "default_pass"),
	Protocolo: getEnv("DB_PROTOCOLO", "http"),
	Ns:        getEnv("DB_NS", "default_ns"),
	Db:        getEnv("DB_NAME", "default_db"),
	Token:     "", // Inicialmente vacío, se llenará tras la autenticación
}

// getEnv devuelve el valor de una variable de entorno o un valor predeterminado si no está definida
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
