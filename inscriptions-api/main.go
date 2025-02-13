package main

import (
	"fmt"
	dao "inscriptions-api/DAOs/inscriptions"
	"inscriptions-api/clients"
	controller "inscriptions-api/controllers/inscriptions"
	repositories "inscriptions-api/repositories/inscriptions"
	router "inscriptions-api/router/inscriptions"
	service "inscriptions-api/services/inscriptions"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	// Intentar conectar a la base de datos con reintentos
	var db, err = connectWithRetry(30, 5*time.Second)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	// Crear el cliente HTTP
	httpClient := clients.NewHTTPClient(
		getEnv("USERS_API_URL", "http://localhost:8081/mock"),
		getEnv("COURSES_API_URL", "http://localhost:8080"), // Asegúrate de que esta URL sea correcta
	)

	// Inicialización de DAO, repositorio, servicio y controlador.
	inscriptionDAO := dao.NewInscriptionDAO(db)
	inscriptionRepository := repositories.NewInscriptionRepository(inscriptionDAO)
	inscriptionService := service.NewService(inscriptionRepository, httpClient)
	inscriptionController := controller.NewController(inscriptionService)

	// Configuración del router.
	r := gin.Default()
	router.MapRoutes(r, inscriptionController)

	// Agregar un endpoint mock para verificar usuarios
	r.GET("/mock/users/:id", func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		// Simular que los usuarios con ID par existen
		if userID%2 == 0 {
			c.JSON(http.StatusOK, gin.H{"id": userID, "name": "Mock User"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	})

	// Agregar un endpoint de prueba para verificar la conexión con la API de cursos
	r.GET("/test-course-api", func(c *gin.Context) {
		url := fmt.Sprintf("%s/courses/1", getEnv("COURSES_API_URL", "http://courses-api:8080"))
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al conectar con la API de cursos: %v", err)})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al leer la respuesta: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": resp.StatusCode,
			"body":   string(body),
		})
	})

	// Asegúrate de que la aplicación use el puerto correcto
	port := getEnv("PORT", "8081")
	log.Printf("Servidor ejecutándose en el puerto %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func connectWithRetry(attempts int, sleep time.Duration) (*gorm.DB, error) {
	for i := 0; i < attempts; i++ {
		db, err := repositories.Connect()
		if err == nil {
			return db, nil
		}
		log.Printf("Intento %d: Error al conectar a la base de datos: %v. Reintentando en %v...", i+1, err, sleep)
		time.Sleep(sleep)
	}
	return nil, fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos", attempts)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
