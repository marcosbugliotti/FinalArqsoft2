package courses

// Course representa la estructura de un curso en la base de datos
type Course struct {
	ID          string `json:"id"`          // Identificador único del curso
	Name        string `json:"name"`        // Nombre del curso
	Category    string `json:"category"`    // Categoría del curso
	Description string `json:"description"` // Descripción del curso
}
