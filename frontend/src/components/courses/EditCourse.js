import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import '../assets/styles/EditCourse.css';

function EditCourse() {
    const { courseId } = useParams();
    const navigate = useNavigate();
    const [courseData, setCourseData] = useState({
        id: parseInt(courseId),
        name: '',
        description: '',
        category: '',
        duration: '',
        instructor_id: '',
        capacity: ''
    });
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchCourseData = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/courses/${courseId}`);
                setCourseData({
                    id: response.data.id,
                    name: response.data.name,
                    description: response.data.description,
                    category: response.data.category,
                    duration: response.data.duration,
                    instructor_id: response.data.instructor_id,
                    capacity: response.data.capacity
                });
            } catch (err) {
                setError('Error fetching course details: ' + err.message);
            }
        };
        fetchCourseData();
    }, [courseId]);

    const handleChange = (e) => {
        setCourseData({ ...courseData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await axios.put(`http://localhost:8080/courses/${courseId}`, { ...courseData});
            alert('Curso actualizado con éxito');
            navigate('/manage-courses');
        } catch (error) {
            setError(`Error updating course: ${error.response?.data?.message || error.message}`);
        }
    };

    return (
        <div className="edit-course-container">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Editar Curso</h1>
            {error && <p className="error-message">{error}</p>}
            <form onSubmit={handleSubmit} className="edit-course-form">
                <label>Nombre del curso</label>
                <input name="name" value={courseData.name} onChange={handleChange} />

                <label>Descripción</label>
                <input name="description" value={courseData.description} onChange={handleChange} />

                <label>Categoría</label>
                <input name="category" value={courseData.category} onChange={handleChange} />

                <label>Duración</label>
                <input name="duration" value={courseData.duration} onChange={handleChange} />

                <label>ID del Instructor</label>
                <input 
                    name="instructor_id" 
                    type="number" 
                    value={courseData.instructor_id} 
                    onChange={handleChange} 
                    min="0" 
                />

                <label>Capacidad</label>
                <input 
                    name="capacity" 
                    type="number" 
                    value={courseData.capacity} 
                    onChange={handleChange} 
                    min="1"
                />

                <label>Imagen Actual</label>
                <img src={`http://localhost:8080/images/${courseData.imageID}`} alt={courseData.name} />

                <button type="submit" className="submit-button">Actualizar Curso</button>
            </form>
        </div>
    );
}

export default EditCourse;
