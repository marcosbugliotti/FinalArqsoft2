package router

import (
	"courses-api/controllers/comments"
	"courses-api/controllers/courses"
	"courses-api/controllers/files"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Función para configurar las rutas
func SetupRouter(courseController courses.Controller, commentController comments.Controller, fileController files.Controller) *gin.Engine {
	r := gin.Default()

	// Configuración de CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,                                                // Permitir todas las solicitudes de origen
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Encabezados permitidos
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas para cursos
	coursesGroup := r.Group("/courses")
	{
		coursesGroup.POST("", courseController.CreateCourse)       // Crear curso
		coursesGroup.GET("", courseController.GetCourses)          // Obtener todos los cursos
		coursesGroup.GET("/:id", courseController.GetCourseByID)   // Obtener curso por ID
		coursesGroup.PUT("/:id", courseController.UpdateCourse)    // Actualizar curso
		coursesGroup.DELETE("/:id", courseController.DeleteCourse) // Eliminar curso
		coursesGroup.POST("/:id/comments", commentController.AddCommentToCourse)
		coursesGroup.GET("/:id/comments", commentController.GetCommentsByCourseID)
		coursesGroup.POST("/:id/files", fileController.CreateFile)
		coursesGroup.GET("/:id/files", fileController.GetFilesByCourseID)
		coursesGroup.PUT("/:id/availability", courseController.UpdateCourseAvailability)
		coursesGroup.GET("/availability", courseController.CourseAvailability)
	}

	return r
}
