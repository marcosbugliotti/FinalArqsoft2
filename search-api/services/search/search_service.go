package search

import (
	"context"
	"fmt"
	"log"
	domain "search-api/domain/courses"     // Alias para los tipos de dominio
	repo "search-api/repositories/courses" // Alias para los repositorios
	"strconv"
)

// Repository define las operaciones necesarias en el índice de SolR
type Repository interface {
	Index(ctx context.Context, course domain.CourseUpdate) (string, error)
	Update(ctx context.Context, course domain.CourseUpdate) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int, offset int) ([]domain.CourseUpdate, error)
}

// Service representa el servicio de búsqueda
type Service struct {
	repository Repository
	httpClient repo.HTTP // Cliente HTTP para interactuar con la API de Cursos
}

// NewService crea una nueva instancia del servicio de búsqueda
func NewService(repository Repository, httpClient repo.HTTP) Service {
	return Service{
		repository: repository,
		httpClient: httpClient,
	}
}

// HandleCourseUpdate procesa las actualizaciones de cursos recibidas desde RabbitMQ
func (service Service) HandleCourseUpdate(courseUpdate domain.CourseUpdate) {
	ctx := context.Background()

	// Convertir CourseID a string
	courseIDStr := strconv.FormatInt(courseUpdate.CourseID, 10)

	// Llamar a GetCourseByID y almacenar el resultado en 'curso'
	curso, err := service.httpClient.GetCourseByID(ctx, courseIDStr)
	if err != nil {
		log.Printf("Error al obtener el curso (%s): %v", courseIDStr, err)
		return // Salir de la función si hay un error
	}

	log.Printf("Curso obtenido, listo para procesar la operación: %d", curso.CourseID) // Nueva línea de log

	switch courseUpdate.Operation {
	case "POST":
		fmt.Println("Course update: ", courseUpdate)
		// Indexar el nuevo curso en SolR
		if _, err := service.repository.Index(ctx, curso); err != nil { // Usar 'curso' en lugar de 'courseUpdate'
			log.Printf("Error al indexar el curso (%d): %v", curso.CourseID, err)
		} else {
			log.Printf("Curso indexado exitosamente: %d", curso.CourseID)
		}

	case "UPDATE":
		// Actualizar el curso existente en SolR
		if err := service.repository.Update(ctx, curso); err != nil { //
			log.Printf("Error al actualizar el curso (%d): %v", courseUpdate.CourseID, err)
		} else {
			log.Printf("Curso actualizado exitosamente: %d", courseUpdate.CourseID)
		}

	case "DELETE":
		// Eliminar el curso del índice de SolR
		if err := service.repository.Delete(ctx, fmt.Sprintf("%d", courseUpdate.CourseID)); err != nil { // Convierte courseUpdate.CourseID a string
			log.Printf("Error al eliminar el curso (%d): %v", courseUpdate.CourseID, err)
		} else {
			log.Printf("Curso eliminado exitosamente: %d", courseUpdate.CourseID)
		}

	default:
		log.Printf("Operación desconocida: %s", courseUpdate.Operation)
	}
}

// Search busca cursos en SolR según el término de búsqueda, límite y desplazamiento
func (service Service) Search(ctx context.Context, query string, limit int, offset int) ([]domain.CourseUpdate, error) {
	results, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error en la búsqueda de cursos: %w", err)
	}
	return results, nil
}
