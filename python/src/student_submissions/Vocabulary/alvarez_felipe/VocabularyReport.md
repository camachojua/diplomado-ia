# Análisis del vocabulario construido en "Los Miserables"

## 1. Introducción

El objetivo de esta tarea es analizar el vocabulario del libro *Los Miserables* de Victor Hugo, realizando un procesamiento de texto para limpiar, estructurar y extraer estadísticas clave. Se sigue la metodología propuesta en el reto **Vocabulary Challenge**, en la cual se requiere:

- Convertir el libro, descargado de https://aprende.org/pruebasat?sectionId=6, a texto estructurado en **CSV**.
- Limpiar el texto eliminando mayúsculas, puntuación y caracteres especiales.
- Construir un vocabulario con las palabras del texto y almacenarlo en formato **Parquet**.
- Realizar análisis estadísticos sobre la cantidad de palabras, palabras únicas y su frecuencia.

## 2. Metodología

### 2.1. Adquisición del Texto

Se obtuvo el texto de *Los Miserables* en formato **EPUB** desde la Fundación Carlos Slim y se convirtió en un archivo de texto utilizando **Python y BeautifulSoup**.

### 2.2. Preprocesamiento de Texto

Para preparar los datos, se realizaron los siguientes pasos:

1. **Conversión a minúsculas**: Para evitar distinciones innecesarias entre palabras idénticas en distintas capitalizaciones.
2. **Eliminación de acentos**: Se reemplazaron caracteres acentuados (`á, é, í, ó, ú`) por sus versiones sin tilde (`a, e, i, o, u`).
3. **Eliminación de puntuación y caracteres especiales**: Se retiraron signos de puntuación y otros caracteres no alfanuméricos.
4. **Tokenización**: Se segmentó el texto en palabras individuales.
5. **Almacenamiento del vocabulario**: Se guardó en un archivo **CSV** y en formato **Parquet** para análisis posterior.

### 2.3. Análisis Estadístico

Se calcularon métricas clave sobre el vocabulario del libro, incluyendo:

- **Cantidad total de palabras** en el libro.
- **Cantidad de palabras únicas** en el vocabulario.
- **Las 100 palabras más frecuentes**.
- **Las 100 palabras menos frecuentes**.

Adicionalmente, se realizó una segunda fase de análisis eliminando **stopwords** en español para evaluar su impacto en el vocabulario resultante.

## 3. Resultados

### 3.1. Estadísticas Generales del Vocabulario

| Métrica | Valor |
|---------|-------|
| Total de palabras en el texto | 109,377 |
| Palabras únicas en el vocabulario | 13,175 |
| Palabras más frecuentes (top 100) | Se muestra en el Apéndice |
| Palabras menos frecuentes (bottom 100) | Se muestra en el Apéndice |

Las palabras más frecuentes incluyen artículos, preposiciones y pronombres, como **"el", "de", "la", "y"**, lo cual es esperado en textos en español.

### 3.2. Impacto de la Eliminación de Stopwords

Se aplicó una lista extendida de **stopwords en español** usando la librería **NLTK**. Tras su eliminación, se obtuvieron los siguientes cambios en las métricas:

| Métrica | Antes de eliminar stopwords | Después de eliminar stopwords |
|---------|-----------------------------|-------------------------------|
| Total de palabras en el texto | 109,377 | 57,301 |
| Palabras únicas en el vocabulario | 13,175 | 13,001 |

Se observa que la eliminación de stopwords **reduce significativamente el número de palabras en el texto**, ya que muchas palabras comunes han sido filtradas. También hay que notar que las palabras únicas se mantienen casi en el mismo número, lo que indica que las stopwords seleccionadas son las que mas aparecen en el texto.

## 4. Análisis y Discusión

Debido a mi experiencia en modelos e NLP, este desafío no presento particulares complicaciónes. Decidí complementar con el análisis de stopwords porque creo que puede ser útil para sesiones futuras.

- La distribución del vocabulario muestra que una gran proporción de palabras aparece solo unas pocas veces, mientras que unas pocas palabras (como artículos y preposiciones) dominan el conteo total.
- La eliminación de **stopwords** permite que el análisis se centre en palabras más representativas del contenido real del libro, eliminando palabras de uso común sin carga semántica significativa.
- El almacenamiento en **Parquet** ofrece una opción eficiente para manejar grandes volúmenes de datos textuales, facilitando futuras consultas y análisis.

## 5. Conclusiones

- Se logró limpiar y estructurar el texto de *Los Miserables*, generando un vocabulario útil para análisis de lenguaje natural.
- La eliminación de acentos y puntuación fue clave para una tokenización efectiva.
- La eliminación de **stopwords** impactó significativamente el número de palabras únicas, mostrando la importancia de este paso en análisis de texto.
- Este análisis puede extenderse para incluir técnicas más avanzadas, como lematización y extracción de temas.

## 6. Apéndice

### 6.1. Top 100 palabras más frecuentes (antes de eliminar stopwords)

Se presenta una lista de las 100 palabras mas frecuentes incluyendo su frecuencia de aparición.

de 5325
la 3918
que 3818
el 3394
y 3123
en 2836
a 2489
se 1681
un 1601
no 1498
los 1353
una 1319
su 1245
por 936
las 935
con 924
habia 858
del 813
al 756
es 749
lo 719
le 667
era 650
como 572
mas 513
para 504
senor 447
esta 414
pero 372
hombre 363
si 358
sus 344
todo 327
me 326
sin 311
obispo 286
dijo 281
cuando 274
estaba 273
sobre 269
dos 264
este 261
aquel 253
mi 244
ya 229
hacia 219
yo 218
esto 218
madeleine 214
tenia 212
jean 200
ha 199
fantine 194
valjean 192
aquella 190
hay 186
he 182
ser 181
muy 178
javert 175
nada 174
mismo 173
o 164
os 163
tan 158
poco 158
bien 157
ni 156
ella 155
quien 151
alcalde 149
vez 148
despues 146
fue 145
todos 141
puerta 137
anos 136
hubiera 133
cual 133
donde 131
dios 130
mujer 127
momento 125
tiempo 124
sido 124
casa 123
son 121
aqui 120
noche 119
hecho 118
tres 115
dia 114
luego 113
cabeza 113
decir 112
voz 111
ojos 107
alli 107
monsenor 105
aun 105