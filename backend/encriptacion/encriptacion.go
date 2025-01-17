package encriptacion

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"surrealdb/backend/config"
)

// Constantes y errores
const (
	KeySize = 32 // Para AES-256
)

var (
	ErrInvalidKeySize    = errors.New("la clave de encriptación debe ser de 32 bytes para AES-256")
	ErrInvalidCiphertext = errors.New("texto cifrado inválido")
	ErrInvalidBlockSize  = errors.New("el texto cifrado es más corto que el tamaño del bloque")
)

// Encrypt encripta el texto plano usando AES-256-CBC con padding PKCS7
func Encrypter(plaintext string, key []byte) (string, error) {
	// Validar tamaño de la clave
	if len(key) != KeySize {
		return "", ErrInvalidKeySize
	}

	// Crear cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generar IV aleatorio
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Aplicar padding PKCS7
	plainBytes := []byte(plaintext)
	paddingSize := aes.BlockSize - (len(plainBytes) % aes.BlockSize)
	paddingBytes := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	paddedText := append(plainBytes, paddingBytes...)

	// Encriptar
	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// Combinar IV y ciphertext y codificar en base64
	combined := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt desencripta el texto cifrado usando AES-256-CBC
func Decrypter(ciphertext string, key []byte) (string, error) {
	// Validar tamaño de la clave
	if len(key) != KeySize {
		return "", ErrInvalidKeySize
	}

	// Decodificar base64
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Validar longitud mínima
	if len(decoded) < aes.BlockSize {
		return "", ErrInvalidBlockSize
	}

	// Crear cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Extraer IV y ciphertext
	iv := decoded[:aes.BlockSize]
	ciphertextBytes := decoded[aes.BlockSize:]

	// Validar que el ciphertext sea múltiplo del tamaño de bloque
	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return "", ErrInvalidCiphertext
	}

	// Desencriptar
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	// Remover padding PKCS7
	paddingSize := int(plaintext[len(plaintext)-1])
	if paddingSize > aes.BlockSize || paddingSize == 0 {
		return "", ErrInvalidCiphertext
	}

	// Validar padding
	for i := len(plaintext) - paddingSize; i < len(plaintext); i++ {
		if plaintext[i] != byte(paddingSize) {
			return "", ErrInvalidCiphertext
		}
	}

	return string(plaintext[:len(plaintext)-paddingSize]), nil
}

func EncryptAllCredentials(creds *config.SurrealDBCredentials, key []byte) error {
	var err error

	// Encriptar el campo Host
	creds.Host, err = Encrypter(creds.Host, key)
	if err != nil {
		return fmt.Errorf("error al encriptar Host: %v", err)
	}

	// Encriptar el campo Port
	creds.Port, err = Encrypter(creds.Port, key)
	if err != nil {
		return fmt.Errorf("error al encriptar Port: %v", err)
	}

	// Encriptar el campo User
	creds.User, err = Encrypter(creds.User, key)
	if err != nil {
		return fmt.Errorf("error al encriptar User: %v", err)
	}

	// Encriptar el campo Pass
	creds.Pass, err = Encrypter(creds.Pass, key)
	if err != nil {
		return fmt.Errorf("error al encriptar Pass: %v", err)
	}

	// Encriptar el campo Protocolo
	creds.Protocolo, err = Encrypter(creds.Protocolo, key)
	if err != nil {
		return fmt.Errorf("error al encriptar Protocolo: %v", err)
	}

	// Encriptar el campo Protocolo
	creds.Ns, err = Encrypter(creds.Ns, key)
	if err != nil {
		return fmt.Errorf("error al encriptar NS: %v", err)
	}

	creds.Db, err = Encrypter(creds.Db, key)
	if err != nil {
		return fmt.Errorf("error al encriptar DB: %v", err)
	}

	return nil
}
