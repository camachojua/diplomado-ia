import ebooklib
from ebooklib import epub
import pandas as pd
import re
import nltk
from nltk.corpus import stopwords
from collections import Counter
import requests
import os

nltk.download('stopwords')  # Descargar la lista de stopwords si no la tienes

def convertir_epub_a_texto(archivo_epub):
    """Convierte un archivo .epub a texto plano.

    Args:
        archivo_epub: La ruta al archivo .epub.

    Returns:
        Una cadena de texto con el contenido del libro.
    """
    libro = epub.read_epub(archivo_epub)
    texto = ''
    for item in libro.get_items():
        if item.media_type == 'application/xhtml+xml':
            contenido = item.get_content().decode('utf-8')
            # Eliminar etiquetas HTML usando expresiones regulares
            contenido_limpio = re.sub('<[^<]+?>', '', contenido)
            texto += contenido_limpio + '\n'  # Agregar una nueva línea entre capítulos
    return texto

def limpiar_texto(texto):
    """Limpia el texto: minúsculas, sin espacios en blanco adicionales, sin puntuación, sin acentos.

    Args:
        texto: La cadena de texto a limpiar.

    Returns:
        Una lista de palabras limpias.
    """
    # a. Minúsculas
    texto = texto.lower()
    # b. Sin espacios en blanco adicionales
    texto = ' '.join(texto.split())
    # c. Sin puntuación
    texto = re.sub(r'[^\w\s]', '', texto)
    # d. Sin acentos 
    texto = texto.replace('á', 'a').replace('é', 'e').replace('í', 'i').replace('ó', 'o').replace('ú', 'u')
    texto = texto.replace('ü', 'u')  
    
    # Eliminar números
    texto = re.sub(r'\d+', '', texto)
    
    # Toquenización
    palabras = texto.split()
    stop_words = set(stopwords.words('spanish')) 
    palabras_limpias = [palabra for palabra in palabras if palabra not in stop_words]

    return palabras_limpias

def crear_vocabulario(lista_palabras):
    """Crea un vocabulario (diccionario) con la frecuencia de cada palabra.

    Args:
        lista_palabras: Una lista de palabras.

    Returns:
        Un diccionario donde las claves son las palabras y los valores son sus frecuencias.
    """
    contador = Counter(lista_palabras)
    return contador

def guardar_vocabulario_parquet(vocabulario, archivo_parquet):
    """Guarda el vocabulario en formato Parquet.

    Args:
        vocabulario: El diccionario del vocabulario.
        archivo_parquet: La ruta al archivo Parquet de salida.
    """
    df = pd.DataFrame(list(vocabulario.items()), columns=['palabra', 'frecuencia'])
    df.to_parquet(archivo_parquet)

def main():
    """Función principal para ejecutar el proceso."""

    archivo_epub = '/content/Los-miserables.epub'  # Nombre del archivo local
    archivo_csv = 'los_miserables.csv'
    archivo_parquet = 'vocabulario_los_miserables.parquet'

    # 1. Convertir a texto
    texto = convertir_epub_a_texto(archivo_epub)

    # 2. Limpiar el texto
    palabras_limpias = limpiar_texto(texto)

     # Guarda el texto en CSV
    df = pd.DataFrame({'palabra': palabras_limpias})
    df.to_csv(archivo_csv, index=False, encoding='utf-8')

    # 3. Crear el vocabulario
    vocabulario = crear_vocabulario(palabras_limpias)

    # 4. Guardar el vocabulario en Parquet
    guardar_vocabulario_parquet(vocabulario, archivo_parquet)

    # 5. Estadísticas
    num_palabras_original = len(texto.split()) # Contar palabras en el texto original antes de la limpieza
    num_palabras_vocabulario = len(vocabulario)

    print(f"Número de palabras en el texto original: {num_palabras_original}")
    print(f"Número de palabras diferentes en el vocabulario: {num_palabras_vocabulario}")

    print("\n100 palabras más frecuentes:")
    for palabra, frecuencia in vocabulario.most_common(100):
        print(f"{palabra}: {frecuencia}")

    print("\n100 palabras menos frecuentes:")
    for palabra, frecuencia in vocabulario.most_common()[:-101:-1]: # Obtener las 100 menos comunes
        print(f"{palabra}: {frecuencia}")

if __name__ == "__main__":
    main()
