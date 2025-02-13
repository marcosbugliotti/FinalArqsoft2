package courses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"search-api/domain/courses"

	"github.com/stevenferrer/solr-go"
)

type SolrConfig struct {
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	// Construct the BaseURL using the provided host and port
	baseURL := fmt.Sprintf("http://%s:%s", config.Host, config.Port)
	client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
	}
}

// Index adds a new course document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, course courses.CourseUpdate) (string, error) {
	// Prepare the document for SolR
	doc := map[string]interface{}{
		"id":          course.CourseID,
		"name":        course.Name,
		"category":    course.Category,
		"description": course.Description,
	}

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in SolR
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling course document: %w", err)
	}

	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error indexing course: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("failed to index course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to SolR: %w", err)
	}

	return fmt.Sprintf("%d", course.CourseID), nil // Convert CourseID to string
}

// Update modifies an existing course document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, course courses.CourseUpdate) error {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"id":          course.CourseID,
		"name":        course.Name,
		"category":    course.Category,
		"description": course.Description,
	}

	// Prepare the update request
	updateRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Update the document in Solr
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	// Execute the update request
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating course: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

// Delete removes a course document from the Solr collection
func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	// Prepare the delete request
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	// Execute the delete request
	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling course document: %w", err)
	}

	// Execute the delete request
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting course: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to delete course: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

// Search searches for courses in the Solr collection based on a query
func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]courses.CourseUpdate, error) {
	// Asegúrate de que 'query' no esté vacío
	if query == "" {
		return nil, fmt.Errorf("la consulta no puede estar vacía")
	}

	// Prepare the Solr query with limit and offset
	solrQuery := fmt.Sprintf("q=(name:%s OR description:%s)&rows=%d&start=%d&wt=json", query, query, limit, offset)

	// Execute the search request
	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la consulta de búsqueda: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta de búsqueda: %v, consulta: %s", resp.Error, solrQuery)
	}

	// Parse the response and extract course documents
	var coursesList []courses.CourseUpdate
	for _, doc := range resp.Response.Documents {
		course := courses.CourseUpdate{
			CourseID:    getIntField(doc, "id"),
			Name:        getStringField(doc, "name"),
			Category:    getStringField(doc, "category"),
			Description: getStringField(doc, "description"),
		}
		coursesList = append(coursesList, course)
	}

	return coursesList, nil
}

// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get int64 fields from the document
func getIntField(doc map[string]interface{}, field string) int64 {
	if val, ok := doc[field].(float64); ok {
		return int64(val)
	}
	return 0
}
