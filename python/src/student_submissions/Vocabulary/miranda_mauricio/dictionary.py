import pandas as pd
import string
import re

#ruta del archivo
file_path = 'data/Los-miserables.csv'

#diccionario de signos de puntuacion
tabla = str.maketrans("", "", string.punctuation)

#diccionario para almecenar las palabras y su frecuencia
diccionario = {}

contador = 0
with open(file_path, "r") as archivo:

    for linea in archivo:
        #Eliminar caracteres de puntuacion
        cadena_limpia = linea.translate(tabla).strip()

        #Eliminar caracteres especiales que contiene el archivo csv
        cadena_limpia = re.sub(r"[—\d«»¿¡]+", "", cadena_limpia)

        #Eliminar espacios extra
        cadena_limpia = " ".join(cadena_limpia.split())

        #convertir en minusculas
        cadena_limpia = cadena_limpia.lower().split(" ")

        for palabra in cadena_limpia:
            if(palabra != ''):
                if(palabra in diccionario):
                    diccionario[palabra] = diccionario[palabra] + 1
                else:
                    diccionario[palabra] = 1
                contador = contador + 1

#Imprimir el numero total de palabras en el texto
print(f"Numero total de palabras: {contador}")

tamanio_diccionario = len(diccionario)

#Imprimir el numero total de palabras diferentes que se guardaron en el diccionario
print(f"Tamaño del diccionario: {tamanio_diccionario} palabras distintas")

diccionario_ordenado_asc = dict(sorted(diccionario.items(), key=lambda item: item[1]))
diccionario_ordenado_desc = dict(sorted(diccionario.items(), key=lambda item: item[1], reverse=True))

#Imprimir las 100 palabras menos repetidas en el diccionario
texto = "100 palabras menos repetidas"
resultado = f"{'-' * 20} {texto} {'-' * 20}"
print(resultado)

i = 0
for clave, valor in diccionario_ordenado_asc.items():
    if(i<100):
        print(f"{clave} => {valor}")
    else:
        break
    i = i + 1

#Imprimir las 100 palabras mas repetidas en el diccionario
texto = "100 palabras mas repetidas"
resultado = f"{'-' * 20} {texto} {'-' * 20}"
print(resultado)

i = 0
for clave, valor in diccionario_ordenado_desc.items():
    if(i<100):
        print(f"{clave} => {valor}")
    else:
        break
    i = i + 1

df = pd.DataFrame(list(diccionario.items()), columns=["clave", "valor"])

df.to_parquet("dictionary.parquet", engine="pyarrow", index=False)