import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate } from 'react-router-dom';
import '../assets/styles/MyCourses.css';

// Asegúrate de que las imágenes estén en la carpeta src/images
const images = [
    require('../../images/image1.png'),
    require('../../images/image2.jpeg'),
    require('../../images/image3.png'),
    require('../../images/image4.jpg'),
    require('../../images/image5.jpg'),
    require('../../images/image6.jpg'),
    require('../../images/image7.jpg'),
    require('../../images/image8.png'),
];

function MyCourses() {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchMyCourses = async () => {
            try {
                const userId = localStorage.getItem('userId');
                if (!userId) {
                    throw new Error('User ID not found');
                }
                
                // Agregar log para ver el userId
                console.log('userId:', userId);
    
                // Obtener inscripciones del usuario
                const inscriptionsResponse = await axios.get(`http://localhost:8081/users/${userId}/inscriptions`);
                
                // Agregar log para ver la respuesta de inscripciones
                console.log('inscriptionsResponse:', inscriptionsResponse.data);
    
                // Verificar que tenemos courseId válidos y filtrar cualquier undefined
                const courseIds = inscriptionsResponse.data
                    .filter(inscription => inscription.course_id)
                    .map(inscription => inscription.course_id);
    
                // Agregar log para ver los courseIds extraídos
                console.log('courseIds extraídos:', courseIds);
    
                if (courseIds.length === 0) {
                    console.log('No se encontraron courseIds');
                    setCourses([]);
                    return;
                }

            // Obtener detalles de los cursos usando los IDs
            const coursesData = await Promise.all(
                courseIds.map(async (id) => {
                    try {
                        const response = await axios.get(`http://localhost:8080/courses/${id}`);
                        return response.data;
                    } catch (error) {
                        console.error(`Error fetching course ${id}:`, error);
                        return null;
                    }
                })
            );

            // Filtrar cualquier curso null que haya fallado en la obtención
            setCourses(coursesData.filter(course => course !== null));
            
        } catch (err) {
            console.error('Error fetching my courses:', err);
            setError('Error fetching my courses');
        } finally {
            setLoading(false);
        }
    };

    fetchMyCourses();
}, []);


    const handleUploadClick = (courseId) => {
        navigate(`/upload/${courseId}`);
    };

    const handleCommentClick = (courseId) => {
        navigate(`/courses/${courseId}/comments`);
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div className="my-courses-container">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Mis Cursos</h1>
            {courses.length > 0 ? (
                <ul className="course-list">
                    {courses.map(course => {
                        // Selecciona una imagen aleatoria para cada curso
                        const randomIndex = Math.floor(Math.random() * images.length);
                        const randomImage = images[randomIndex];

                        return (
                            <li key={course.id} className="course-item">
                                {/* Mostrar la imagen aleatoria del curso */}
                                <img src={randomImage} alt={course.name} className="course-image" />
                                <div className="course-info">
                                    <h3>{course.name}</h3>
                                    <p>{course.description}</p>
                                </div>
                                <Link to={`/courses/${course.id}`} className="details-button">Click para conocer más detalles</Link>
                                <button onClick={() => handleUploadClick(course.id)} className="upload-button">Subir Archivo</button>
                                <button onClick={() => handleCommentClick(course.id)} className="comment-button">Realizar Comentario</button>
                            </li>
                        );
                    })}
                </ul>
            ) : (
                <p>No estás inscrito en ningún curso aún.</p>
            )}
        </div>
    );
}

export default MyCourses;
