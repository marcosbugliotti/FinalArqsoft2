import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import { useParams } from 'react-router-dom';
import '../assets/styles/CourseDetails.css';


function CourseDetails() {
    const { courseId } = useParams();
    const navigate = useNavigate();
    const [course, setCourse] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchCourse = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/courses/${courseId}`);
                setCourse(response.data);
                setLoading(false);
            } catch (err) {
                setError('Error fetching course details: ' + err.message);
                setLoading(false);
            }
        };

        fetchCourse();
    }, [courseId]);

    const handleEnroll = async () => {
        const userId = localStorage.getItem('userId');
        if (!userId) {
            alert('Por favor, inicia sesión para inscribirte en el curso.');
            return;
        }

        try {
            // Intentar inscribirse en el curso
            await axios.post(`http://localhost:8081/inscriptions`, {
                user_id: parseInt(userId),
                course_id: parseInt(courseId)
            });
            alert('Inscripción exitosa!');
            navigate('/my-courses'); // Redirigir a "Mis Cursos" después de inscribirse
        } catch (err) {
            alert('Error en la inscripción: ' + err.response?.data?.error || err.message);
        }
    };

    if (loading) return <div>Cargando...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div className="course-details-container">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <div className="course-details">
                <h1>{course.name}</h1>
                <img src={`frontend\src\images${course.image_id}`} alt={course.name} className="course-image" />
                <p><strong>Description:</strong> {course.description}</p>
                <p><strong>Category:</strong> {course.category}</p>
                <p><strong>Duration:</strong> {course.duration}</p>
                <p><strong>Instructor ID:</strong> {course.instructor_id}</p>
                <p><strong>Capacidad:</strong> {course.capacity}</p>
                <p><strong>Rating:</strong> {course.rating}</p>
                {course.capacity > 0 ? (
                    <button onClick={handleEnroll}>Inscribirse</button>
                ) : (
                    <p>El curso está lleno. No se puede inscribir.</p>
                )}
                <div className="course-files">
                    <h2>Archivos del Curso</h2>
                    <Link to={`/courses/${courseId}/files`} className="files-button">Ver Archivos</Link>
                </div>
                <div className="course-comments">
                    <h2>Comentarios del Curso</h2>
                    <Link to={`/courses/${courseId}/comments`} className="comments-button">Ver Comentarios</Link>
                </div>
            </div>
        </div>
    );
}

export default CourseDetails;
