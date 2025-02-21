import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/ManageCourses.css';

function ManageCourses() {
    const [cursos, setCursos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        fetchCursos();
    }, []);

    const fetchCursos = async () => {
        try {
            const response = await axios.get('http://localhost:8080/courses');
            setCursos(response.data);
            setLoading(false);
        } catch (error) {
            setError('Error fetching courses');
            setLoading(false);
        }
    };

    const handleDelete = async (courseId) => {
        if (window.confirm('¿Estás seguro de que deseas eliminar este curso?')) {
            try {
                await axios.delete(`http://localhost:8080/courses/${courseId}`);
                // Actualizar la lista de cursos después de eliminar
                const updatedCourses = cursos.filter(curso => curso.id !== courseId);
                setCursos(updatedCourses);
            } catch (error) {
                // Verificar si el error es por inscripciones activas
                if (error.response?.data?.error?.includes('inscripciones activas')) {
                    alert('No se puede eliminar el curso porque tiene estudiantes inscritos');
                } else {
                    alert('Error al eliminar el curso: ' + (error.response?.data?.error || 'Error desconocido'));
                }
            }
        }
    };

    return (
        <div className="manage-courses-container">
             <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Gestión de Cursos</h1>
            {loading ? <p>Loading...</p> : null}
            {error ? <p>{error}</p> : null}
            {cursos === null ? (
                <div>
                    <p>No hay cursos disponibles. Haz clic aquí para agregar un curso.</p>
                    <button onClick={() => navigate('/add-course')} className="add-course-button">Agregar un curso</button>
                </div>
            ) : (
                <ul className="course-list">
                <button onClick={() => navigate('/add-course')} className="add-course-button">Agregar un curso</button>
                    {cursos.map(curso => (
                        <li key={curso.id} className="course-item">
                            <div className="course-info">
                                <h3>{curso.name}</h3>
                                <p>Capacidad: {curso.capacity}</p>
                            </div>
                            <div className="course-actions">
                                <button onClick={() => handleDelete(curso.id)} className="delete-button">Eliminar</button>
                                <button onClick={() => navigate(`/edit-course/${curso.id}`)} className="edit-button">Editar</button>
                                <button onClick={() => navigate(`/courses/${curso.id}`)} className="details-button">Click para conocer más detalles</button>
                            </div>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default ManageCourses;
