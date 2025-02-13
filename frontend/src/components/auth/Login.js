import React, { useState, useContext } from 'react';
import axios from 'axios';
import { UserContext } from '../context/UserContext';
import { useNavigate } from 'react-router-dom';
import './Login.css'; 

function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { setUser } = useContext(UserContext);
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8083/login', {
                username,
                password
            });
            const { user_id, user_type, token } = response.data;

            localStorage.setItem('userId', user_id);
            localStorage.setItem('usertype', user_type);
            localStorage.setItem('token', token);
            setUser(response.data);
            alert('Login successful: ');
            navigate('/home');
        } catch (error) {
            setError('Failed to login: ' + error.message);
        }
    };

    const handleRegisterRedirect = () => {
        navigate('/register');
    };

    return (
        <div className="login-container">
            <h1 className="welcome-title">Bienvenido al Sistema de Cursos</h1>
            <h2>Iniciar sesión</h2>
            <form onSubmit={handleLogin} className="login-form">
                <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="Usuario" required className="input-field" />
                <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="Contraseña" required className="input-field" />
                <button type="submit" className="login-button">Iniciar sesión</button>
                <button type="button" onClick={handleRegisterRedirect} className="register-button">Registrarse</button>
            </form>
            {error && <p className="error-message">{error}</p>}
        </div>
    );
}

export default Login;
