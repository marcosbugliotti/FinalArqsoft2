import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/AddCourse.css';

function AddCourse() {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [category, setCategory] = useState('');
    const [duration, setDuration] = useState('');
    const [instructorId, setInstructorId] = useState('');
    const [capacity, setCapacity] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!name || !description || !category || !duration || !instructorId || !capacity) {
            setError('Todos los campos son obligatorios');
            return;
        }

        const parsedInstructorId = parseInt(instructorId);
        const parsedCapacity = parseInt(capacity);
        if (isNaN(parsedInstructorId) || isNaN(parsedCapacity)) {
            setError('ID del instructor y capacidad deben ser números');
            return;
        }

        const courseData = {
            name,
            description,
            category,
            duration,
            instructor_id: parsedInstructorId,
            capacity: parsedCapacity,
        };

        console.log('Datos enviados:', courseData);

        try {
            const response = await axios.post('http://localhost:8080/courses', courseData, {
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            setError('');
            alert('Curso agregado con éxito');
            navigate('/manage-courses');
        } catch (error) {
            setError('Error al agregar curso: ' + (error.response?.data?.message || error.message));
        }
    };

    return (
        <div className="add-course-container">
            <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
            <h1>Agregar nuevo curso</h1>
            {error && <p className="error-message">{error}</p>}
            <form onSubmit={handleSubmit} className="add-course-form">
                <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Nombre del curso" />
                <input type="text" value={description} onChange={e => setDescription(e.target.value)} placeholder="Descripción" />
                <input type="text" value={category} onChange={e => setCategory(e.target.value)} placeholder="Categoría" />
                <input type="text" value={duration} onChange={e => setDuration(e.target.value)} placeholder="Duración" />
                <input 
                    type="number" 
                    value={instructorId} 
                    onChange={e => setInstructorId(e.target.value)} 
                    placeholder="ID del instructor" 
                    min="0"
                />
                <input 
                    type="number" 
                    value={capacity} 
                    onChange={e => setCapacity(e.target.value)} 
                    placeholder="Capacidad del curso" 
                    min="1"
                />
                <button type="submit" className="submit-button">Agregar curso</button>
            </form>
        </div>
    );
}

export default AddCourse;
