package courses

import (
	"context"
	coursesDAO "courses-api/DAO/courses"
	"courses-api/clients"
	"courses-api/domain/courses"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Repository interface para las operaciones de curso
type Repository interface {
	CreateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error)
	GetCourses(ctx context.Context) ([]coursesDAO.Course, error)
	GetCourseByID(ctx context.Context, id int64) (coursesDAO.Course, error)
	UpdateCourse(ctx context.Context, course coursesDAO.Course) (coursesDAO.Course, error)
	DeleteCourse(ctx context.Context, id int64) error
}

// CommentsRepository interface para las operaciones de comentarios
type CommentsRepository interface {
	DeleteCommentsByCourseID(ctx context.Context, courseID int64) error
}

type Queue interface {
	Publish(courseNew courses.CursosNew) error
}

// FilesRepository interface para las operaciones de archivos
type FilesRepository interface {
	DeleteFilesByCourseID(ctx context.Context, courseID int64) error
}

// Service estructura para el servicio de cursos
type Service struct {
	repository         Repository
	commentsRepository CommentsRepository
	filesRepository    FilesRepository
	eventsQueue        Queue
	httpClient         *clients.HTTPClient
}

// NewService constructor para el servicio de cursos
func NewService(repository Repository, commentsRepository CommentsRepository, filesRepository FilesRepository, eventsQueue Queue, httpClient *clients.HTTPClient) Service {
	return Service{
		repository:         repository,
		commentsRepository: commentsRepository,
		filesRepository:    filesRepository,
		eventsQueue:        eventsQueue,
		httpClient:         httpClient,
	}
}

// Inicializa el generador de números aleatorios
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Función para obtener una imagen aleatoria de la carpeta de imágenes
func GetRandomImage() (string, error) {
	// Cambia esta ruta a la ubicación correcta
	imageDir, err := filepath.Abs("./images")
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(imageDir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no hay imágenes disponibles")
	}

	// Selecciona un archivo aleatorio
	randomIndex := rand.Intn(len(files))
	return filepath.Join(imageDir, files[randomIndex].Name()), nil
}

func (s Service) CreateCourse(ctx context.Context, req courses.CreateCourseRequest) (courses.CourseResponse, error) {
	// Obtener una imagen aleatoria
	imagePath, err := GetRandomImage()
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("error al obtener imagen aleatoria: %v", err)
	}

	course := coursesDAO.Course{
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Duration:     req.Duration,
		InstructorID: req.InstructorID,
		ImageID:      imagePath, // Asignar la imagen aleatoria
		Capacity:     req.Capacity,
		Rating:       0, // Inicialmente, el rating es 0
	}

	createdCourse, err := s.repository.CreateCourse(ctx, course)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to create course: %v", err)
	}

	go func() {
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "POST",
			CourseID:  createdCourse.ID,
		}); err != nil {
			fmt.Println(fmt.Sprintf("Error al publicar nuevo curso: %v", err))
		}
	}()

	return courses.CourseResponse{
		ID:           createdCourse.ID,
		Name:         createdCourse.Name,
		Description:  createdCourse.Description,
		Category:     createdCourse.Category,
		Duration:     createdCourse.Duration,
		InstructorID: createdCourse.InstructorID,
		ImageID:      createdCourse.ImageID,
		Capacity:     createdCourse.Capacity,
		Rating:       createdCourse.Rating,
	}, nil
}

func (s Service) GetCourses(ctx context.Context) ([]courses.CourseResponse, error) {
	coursesDAO, err := s.repository.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	var coursesResponse []courses.CourseResponse
	for _, course := range coursesDAO {
		coursesResponse = append(coursesResponse, courses.CourseResponse{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Category:     course.Category,
			Duration:     course.Duration,
			InstructorID: course.InstructorID,
			ImageID:      course.ImageID,
			Capacity:     course.Capacity,
			Rating:       course.Rating,
		})
	}

	return coursesResponse, nil
}

func (s Service) GetCourseByID(ctx context.Context, id int64) (courses.CourseResponse, error) {
	course, err := s.repository.GetCourseByID(ctx, id)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to get course: %v", err)
	}

	return courses.CourseResponse{
		ID:           course.ID,
		Name:         course.Name,
		Description:  course.Description,
		Category:     course.Category,
		Duration:     course.Duration,
		InstructorID: course.InstructorID,
		ImageID:      course.ImageID,
		Capacity:     course.Capacity,
		Rating:       course.Rating,
	}, nil
}

func (s Service) UpdateCourse(ctx context.Context, id int64, req courses.UpdateCourseRequest) (courses.CourseResponse, error) {
	course, err := s.repository.GetCourseByID(ctx, id)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("course not found: %v", err)
	}

	if req.Name != "" {
		course.Name = req.Name
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Category != "" {
		course.Category = req.Category
	}
	if req.Duration != "" {
		course.Duration = req.Duration
	}
	if req.InstructorID != 0 {
		course.InstructorID = req.InstructorID
	}
	if req.Capacity != 0 {
		course.Capacity = req.Capacity
	}
	// No actualizamos el rating aquí, ya que se actualizará con los comentarios

	updatedCourse, err := s.repository.UpdateCourse(ctx, course)
	if err != nil {
		return courses.CourseResponse{}, fmt.Errorf("failed to update course: %v", err)
	}

	go func() {
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "PUT",
			CourseID:  updatedCourse.ID,
		}); err != nil {
			fmt.Println(fmt.Sprintf("Error al publicar actualización de curso: %v", err))
		}
	}()

	return courses.CourseResponse{
		ID:           updatedCourse.ID,
		Name:         updatedCourse.Name,
		Description:  updatedCourse.Description,
		Category:     updatedCourse.Category,
		Duration:     updatedCourse.Duration,
		InstructorID: updatedCourse.InstructorID,
		ImageID:      updatedCourse.ImageID,
		Capacity:     updatedCourse.Capacity,
		Rating:       updatedCourse.Rating,
	}, nil
}

func (s Service) DeleteCourse(ctx context.Context, id int64) error {
	// Verificar si hay inscripciones para este curso
	//inscriptions, err := s.httpClient.GetInscriptionsByCourse(uint(id))
	//if err != nil {
	//return fmt.Errorf("error al verificar inscripciones: %v", err)
	//}

	//if len(inscriptions) > 0 {
	//	return errors.New("no se puede eliminar el curso porque tiene inscripciones activas")
	//}

	// Eliminar los comentarios asociados al curso
	err := s.commentsRepository.DeleteCommentsByCourseID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar los comentarios del curso: %v", err)
	}

	// Eliminar los archivos asociados al curso
	err = s.filesRepository.DeleteFilesByCourseID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar los archivos del curso: %v", err)
	}

	// Eliminar el curso
	err = s.repository.DeleteCourse(ctx, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el curso: %v", err)
	}

	go func() {
		if err := s.eventsQueue.Publish(courses.CursosNew{
			Operation: "DELETE",
			CourseID:  id,
		}); err != nil {
			fmt.Println(fmt.Sprintf("Error al publicar eliminación de curso: %v", err))
		}
	}()

	return nil
}

// Agregar este método para actualizar el rating del curso
func (s Service) UpdateCourseRating(ctx context.Context, courseID int64, newRating float64) error {
	course, err := s.repository.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	course.Rating = newRating
	_, err = s.repository.UpdateCourse(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to update course rating: %v", err)
	}

	return nil
}
