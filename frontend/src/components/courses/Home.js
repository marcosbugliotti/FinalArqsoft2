import React, { useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';
import '../assets/styles/Home.css';

function Home() {
    const [cursos, setCursos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useContext(UserContext);
    const navigate = useNavigate();
    const userType = localStorage.getItem('usertype');
    
    const isAdmin = userType && userType.toLowerCase() === 'administrador';

    useEffect(() => {
        const fetchCursos = async () => {
            try {
                const response = await axios.get('http://localhost:8080/courses');
                setCursos(response.data);
            } catch (error) {
                setError('Error fetching data: ' + error.message);
            } finally {
                setLoading(false);
            }
        };

        fetchCursos();
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('userId');
        localStorage.removeItem('usertype');
        localStorage.removeItem('token');
        navigate('/login');
    };

    if (loading) {
        return <div className="home-container">Cargando...</div>;
    }

    if (error) {
        return <div className="home-container">{error}</div>;
    }

    return (
        <div className="home-container">
            <div className="logout-container">
                <button className="logout-button" onClick={handleLogout}>Logout</button>
            </div>
            <div className="welcome-message">Bienvenido al Portal de Cursos</div>
            <p className="description">Explora y administra tus cursos con facilidad. Aquí puedes encontrar información detallada sobre todos los cursos disponibles y gestionar tus cursos activos.</p>
            <div className="header">
                <button className="button" onClick={() => navigate('/search')}>Buscar un curso</button>
                <button className="button" onClick={() => navigate('/my-courses')}>Mis Cursos</button>
                {isAdmin && <button className="button" onClick={() => navigate('/manage-courses')}>Gestión de Cursos</button>}
            </div>
            <h2>Cursos Disponibles</h2>
            {cursos.length === 0 ? (
                <p>No hay cursos disponibles.</p>
            ) : (
                <ul className="course-list">
                    {cursos.map(curso => (
                        <li key={curso.id} className="course-item">
                            {curso.name}
                            <button className="button" onClick={() => navigate(`/courses/${curso.id}`)}>Click para conocer más detalles</button>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Home;
