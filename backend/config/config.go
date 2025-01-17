package config

import (
	"log"
	"os"
)

// Secreto JWT (cargado desde las variables de entorno)
var JWTSecret string
var EncryptionKey []byte

func init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET no está definido en las variables de entorno")
	}

	EncryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(EncryptionKey) == 0 {
		log.Fatal("ENCRYPTION_KEY no está definido en las variables de entorno")
	}
}

// SurrealDBCredentials almacena las credenciales para conectarse a la base de datos
type SurrealDBCredentials struct {
	Host  string `json:"host"`
	Ns    string `json:"ns"`
	Db    string `json:"db"`
	User  string `json:"user"`
	Pass  string `json:"pass"`
	Token string `json:"token"`
}

// Credentials contiene las credenciales cargadas desde las variables de entorno
var Credentials = SurrealDBCredentials{
	Host: os.Getenv("DB_HOST"), // La dirección del host se carga desde la variable de entorno
}
