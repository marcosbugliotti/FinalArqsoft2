import React, { useState } from 'react';
import axios from 'axios';
import '../assets/styles/FileUploadComponent.css';

function UploadFile({ courseId }) {
    const [file, setFile] = useState(null);
    const [error, setError] = useState(null);

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        setFile(selectedFile);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!file) {
            setError('Por favor, selecciona un archivo para subir.');
            return;
        }

        const reader = new FileReader();
        reader.onloadend = async () => {
            const base64String = reader.result.split(',')[1]; // Obtener solo la parte base64

            try {
                const response = await axios.post(`http://localhost:8080/files/${courseId}`, {
                    name: file.name,
                    content: base64String,
                    userId: Number(localStorage.getItem('userId')), // Asegúrate de que el userId esté en el localStorage
                }, {
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${localStorage.getItem('token')}`
                    }
                });

                alert('Archivo subido correctamente');
            } catch (error) {
                setError('Error al subir el archivo: ' + error.message);
            }
        };

        reader.readAsDataURL(file); // Leer el archivo como Data URL
    };

    return (
        <div>
            <h2>Subir Archivo</h2>
            <form onSubmit={handleSubmit}>
                <input type="file" onChange={handleFileChange} />
                <button type="submit">Subir Archivo</button>
            </form>
            {error && <p>{error}</p>}
        </div>
    );
}

export default UploadFile;
