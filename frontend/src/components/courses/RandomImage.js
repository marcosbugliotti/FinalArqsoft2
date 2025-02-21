import React, { useEffect, useState } from 'react';

// Asegúrate de que las imágenes estén en la carpeta src/images
const images = [
    require('../../images/image1.jpg'),
    require('../../images/image2.jpg'),
    require('../../images/image3.jpg'),
    require('../../images/image4.jpg'),
    require('../../images/image5.jpg'),
    require('../../images/image6.jpg'),
    require('../../images/image7.jpg'),
    require('../../images/image8.jpg'),
];

function RandomImage() {
    const [randomImage, setRandomImage] = useState('');

    useEffect(() => {
        // Selecciona una imagen aleatoria
        const randomIndex = Math.floor(Math.random() * images.length);
        setRandomImage(images[randomIndex]);
    }, []);

    return (
        <div>
            {randomImage && <img src={randomImage} alt="Imagen aleatoria" />}
        </div>
    );
}

export default RandomImage;
