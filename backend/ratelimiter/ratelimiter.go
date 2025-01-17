// rate_limiter.go

package ratelimiter

import (
	"sync"
	"time"
)

// Estructura para gestionar intentos fallidos y bloqueo
type LoginAttempt struct {
	Attempts          int       // Contador de intentos fallidos
	BlockedUntil      time.Time // Tiempo hasta que la IP está bloqueada
	MaxFailedAttempts int       // Número máximo de intentos fallidos permitidos
}

// Mapa para rastrear solicitudes por IP
var (
	rateLimiter          = make(map[string]*LoginAttempt)
	mu                   sync.Mutex
	maxRequestsPerMinute = 60 // Número máximo de solicitudes permitidas por minuto
	MaxFailedAttempts    = 5
)

// CheckRateLimit verifica si una IP ha excedido el límite de solicitudes
func CheckRateLimit(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	// Inicializa el registro si no existe

	if _, exists := rateLimiter[ip]; !exists {
		rateLimiter[ip] = &LoginAttempt{Attempts: 0, BlockedUntil: time.Time{}, MaxFailedAttempts: 3} // Ejemplo: máximo 3 intentos fallidos
	}

	// Verifica si se han excedido los intentos fallidos
	if rateLimiter[ip].Attempts >= rateLimiter[ip].MaxFailedAttempts {
		rateLimiter[ip].BlockedUntil = time.Now().Add(time.Minute) // Bloquear por 1 minuto
		return false                                               // Bloqueado por exceder el número máximo de intentos fallidos
	}

	// Si la IP está bloqueada, retorna false
	if rateLimiter[ip].Attempts >= maxRequestsPerMinute {
		if time.Since(rateLimiter[ip].BlockedUntil) < time.Minute {
			return false // Bloqueado por exceder límite de solicitudes
		}
		// Reiniciar contador si ya pasó 1 minuto
		rateLimiter[ip].Attempts = 0
		rateLimiter[ip].BlockedUntil = time.Time{} // Resetear el tiempo de bloqueo
	}

	rateLimiter[ip].Attempts++
	return true
}
