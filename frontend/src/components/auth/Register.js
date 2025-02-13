import React, { useState, useContext } from 'react';
import axios from 'axios';
import { UserContext } from '../context/UserContext';
import { useNavigate } from 'react-router-dom';
import './Register.css'; 

function Register() {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [user_type, setUserType] = useState('alumno');
    const [error, setError] = useState('');
    const { setUser } = useContext(UserContext);
    const navigate = useNavigate();

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8083/users', {
                username,
                email,
                password,
                user_type
            });
            const { id } = response.data;

            localStorage.setItem('userId', id);
            localStorage.setItem('usertype', user_type);
            console.log(user_type);
            setUser(response.data);
            alert('Registration successful');
            navigate('/login');
        } catch (error) {
            setError('Failed to register: ' + error.message);
        }
    };

    return (
        <div className="register-container">
            <h2>Register</h2>
            <form onSubmit={handleRegister} className="register-form">
                <input type="text" value={username} onChange={e => setUsername(e.target.value)} placeholder="Username" required className="input-field" />
                <input type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="Email" required className="input-field" />
                <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="Password" required className="input-field" />
                <select value={user_type} onChange={e => setUserType(e.target.value)} required className="input-field">
                    <option value="alumno">Alumno</option>
                    <option value="administrador">Administrador</option>
                </select>
                <button type="submit" className="register-button">Register</button>
            </form>
            {error && <p className="error-message">{error}</p>}
        </div>
    );
}

export default Register;
