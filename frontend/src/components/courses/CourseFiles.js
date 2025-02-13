import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate, Link } from 'react-router-dom';
import '../assets/styles/CourseFiles.css';

const CourseFiles = () => {
  const { courseId } = useParams();
  const navigate = useNavigate();
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [course, setCourse] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Obtener información del curso
        const courseResponse = await axios.get(`http://localhost:8080/courses/${courseId}`);
        setCourse(courseResponse.data);

        // Obtener archivos del curso
        const filesResponse = await axios.get(`http://localhost:8080/courses/${courseId}/files`);
        setFiles(filesResponse.data || []);
        setLoading(false);
      } catch (err) {
        console.error("Error:", err);
        setError('Error al cargar los datos: ' + (err.response?.data?.message || err.message));
        setLoading(false);
      }
    };

    fetchData();
  }, [courseId]);

  const handleDownload = async (file) => {
    try {
      // Crear un blob desde el contenido base64
      const byteCharacters = atob(file.content);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);
      const blob = new Blob([byteArray]);

      // Crear URL y descargar
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = file.name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Error al descargar:', error);
      alert('Error al descargar el archivo');
    }
  };

  if (loading) {
    return (
      <div className="loading-container">
        <p>Cargando archivos del curso...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <p>{error}</p>
        <button onClick={() => navigate(-1)}>Volver</button>
      </div>
    );
  }

  return (
    <div className="course-files-container">
      <div className="header">
        <button className="back-button" onClick={() => navigate(-1)}>Volver</button>
        <h1>Archivos del Curso: {course?.name}</h1>
      </div>

      <div className="actions">
        <Link 
          to={`/upload/${courseId}`} 
          className="upload-button"
        >
          Subir Nuevo Archivo
        </Link>
      </div>

      <div className="files-list">
        {files.length === 0 ? (
          <div className="no-files">
            <p>No hay archivos disponibles en este curso.</p>
          </div>
        ) : (
          <ul>
            {files.map((file) => (
              <li key={file.id} className="file-item">
                <span className="file-name">{file.name}</span>
                <button 
                  className="download-button"
                  onClick={() => handleDownload(file)}
                >
                  Descargar
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default CourseFiles;
