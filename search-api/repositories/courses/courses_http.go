package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"search-api/domain/courses" // Importación correcta del paquete
)

type HTTPConfig struct {
	Host string
	Port string
}

type HTTP struct {
	baseURL func(courseID string) string
}

func NewHTTP(config HTTPConfig) HTTP {
	return HTTP{
		baseURL: func(courseID string) string {
			return fmt.Sprintf("http://%s:%s/courses/%s", config.Host, config.Port, courseID)
		},
	}
}

// GetCourseByID obtiene los detalles de un curso usando su ID
func (repository HTTP) GetCourseByID(ctx context.Context, id string) (courses.CourseUpdate, error) {
	resp, err := http.Get(repository.baseURL(id))
	if err != nil {
		return courses.CourseUpdate{}, fmt.Errorf("Error fetching course (%s): %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return courses.CourseUpdate{}, fmt.Errorf("Failed to fetch course (%s): received status code %d", id, resp.StatusCode)
	}

	// Lee el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return courses.CourseUpdate{}, fmt.Errorf("Error reading response body for course (%s): %w", id, err)
	}

	// Deserializa los datos del curso en la estructura CourseUpdate
	var course courses.CourseUpdate
	if err := json.Unmarshal(body, &course); err != nil {
		return courses.CourseUpdate{}, fmt.Errorf("Error unmarshaling course data (%s): %w", id, err)
	}

	return course, nil
}
