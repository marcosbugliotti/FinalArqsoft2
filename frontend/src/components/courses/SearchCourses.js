import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/SearchCourses.css';

const SearchCourses = () => {
    const [searchTerm, setSearchTerm] = useState('');
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const handleSearch = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const response = await axios.get(`http://localhost:8082/search?q=${searchTerm}`);
            console.log("resultado:", response.data);
            
            setCourses(response.data);

        } catch (err) {
            setError('Error fetching courses: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="search-courses-container">
            <h1>Buscar Cursos</h1>
            <form onSubmit={handleSearch} className="search-form">
                <input
                    type="text"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Ingrese el nombre del curso"
                    required
                    className="search-input"
                />
                <button type="submit" className="search-button">Buscar</button>
            </form>
            {loading && <p>Cargando...</p>}
            {error && <p className="error-message">{error}</p>}
            {courses.length > 0 ? (
                <ul className="course-list">
                    {courses.map(course => (
                        <li key={course.course_id} className="course-item">
                            <h3>{course.name}</h3>
                            <p>{course.description}</p>
                            <button className="button" onClick={() => navigate(`/courses/${course.course_id}`)}>Ver Detalles</button>
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No se encontraron cursos.</p>
            )}
        </div>
    );
};

export default SearchCourses;
