import pdfplumber  # Para extraer texto de archivos PDF
import pandas as pd  # Para manejar datos y guardar en formato CSV
import re  # Para limpieza de texto (expresiones regulares)
import unicodedata  # Para normalizar caracteres (eliminar acentos)
from collections import Counter  # Para contar frecuencias de palabras
import pyarrow.parquet as pq  # Para guardar el vocabulario en formato Parquet
import pyarrow as pa  # Para convertir DataFrame a formato Parquet

# Función para limpiar el texto
def limpiar_texto(texto):
    texto = texto.lower()  # Convertir a minúsculas
    texto = unicodedata.normalize("NFKD", texto).encode("ascii", "ignore").decode("utf-8")  # Eliminar acentos
    texto = re.sub(r'\W+', ' ', texto)  # Eliminar puntuación
    texto = texto.strip()  # Eliminar espacios adicionales
    return texto.split()  # Dividir en palabras

# Función para construir el vocabulario
def construir_vocabulario(tokens):
    contador = Counter(tokens)
    vocabulario = {}
    palabras_comunes = contador.most_common()

    for indice in range(len(palabras_comunes)):
        palabra, _ = palabras_comunes[indice]
        vocabulario[palabra] = indice

    return vocabulario, contador

# Función para guardar el vocabulario en formato Parquet
def guardar_vocabulario_parquet(vocabulario, contador, nombre_archivo="vocabulario.parquet"):
    df = pd.DataFrame({'palabra': list(vocabulario.keys()), 
                       'indice': list(vocabulario.values()), 
                       'frecuencia': [contador[palabra] for palabra in vocabulario]})
    tabla = pa.Table.from_pandas(df)
    pq.write_table(tabla, nombre_archivo)
    print(f"vocabulario guardado en {nombre_archivo}.")

# Función para obtener estadísticas y organizar resultados en DataFrames
def obtener_estadisticas(tokens, vocabulario, contador):
    total_palabras = len(tokens)
    palabras_unicas = len(vocabulario)
    mas_frecuentes = contador.most_common(100)
    menos_frecuentes = contador.most_common()[:-101:-1]

    print("\n--- Resultados Generales ---\n")
    print(f"Total de palabras en el texto: {total_palabras}")
    print(f"Cantidad de palabras únicas: {palabras_unicas}")

    # Crear DataFrames
    df_mas_frecuentes = pd.DataFrame(mas_frecuentes, columns=["Palabra", "Frecuencia"])
    df_menos_frecuentes = pd.DataFrame(menos_frecuentes, columns=["Palabra", "Frecuencia"])
    return total_palabras, palabras_unicas, df_mas_frecuentes, df_menos_frecuentes

# Función principal (main)
def main():
    # Extraer texto del PDF
    path_pdf = "Los_miserables.pdf"
    contenido_texto = []

    print("Extrayendo texto del PDF...")
    try:
        with pdfplumber.open(path_pdf) as pdf:
            for pagina in pdf.pages:
                texto = pagina.extract_text()
                if texto:
                    # Dividir en líneas y filtrar líneas que no sean números
                    lineas_filtradas = [
                        linea for linea in texto.split("\n")
                        # Elimina líneas con solo números (números de página)
                        if not re.match(r'^\d+$', linea.strip())  
                    ]
                    contenido_texto.extend(lineas_filtradas)  # Añadir texto filtrado
    except FileNotFoundError:
        print(f"Error: No se encontró el archivo {path_pdf}. Verifica la ruta.")
        return
    except Exception as e:
        print(f"Ocurrió un error al leer el PDF: {e}")
        return
    

    # Guardar contenido como CSV
    print("Guardando contenido en archivo CSV...")
    df = pd.DataFrame({'texto': contenido_texto})
    df.to_csv("los_miserables.csv", index=False, encoding="utf-8")
    print(f"CSV guardado correctamente con {len(df)} filas.")

    # Unir el texto completo
    texto_completo = " ".join(contenido_texto)

    # Limpiar el texto
    print("Limpiando texto...")
    tokens = limpiar_texto(texto_completo)

    # Construir el vocabulario
    print("Construyendo vocabulario...")
    vocabulario, contador = construir_vocabulario(tokens)

    #  Guardar vocabulario en Parquet
    print("Guardando vocabulario en formato Parquet...")
    guardar_vocabulario_parquet(vocabulario, contador)

    # Obtener estadísticas
    print("Generando estadísticas...")
    total_palabras, palabras_unicas, df_mas_frecuentes, df_menos_frecuentes = obtener_estadisticas(tokens, vocabulario, contador)

    # Mostrar ejemplos de los DataFrames
    print("\n--- Palabras Más Frecuentes (Top 100) ---")
    print(df_mas_frecuentes)

    print("\n--- Palabras Menos Frecuentes (Últimas 100) ---")
    print(df_menos_frecuentes)

# Ejecutar el programa
if __name__ == "__main__":
    main()
