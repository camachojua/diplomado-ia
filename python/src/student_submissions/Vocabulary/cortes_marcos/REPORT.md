# Creación de Volabulario para el libro "Los Miserables TOMO 1"

## Objetivo principal

Generar el vocabulario del libro Los Miserables de Victor Hugo

## Integrantes del equipo

* Marcos Cortes Valadez
* Cecilia Gomez Castañeda


## Actividades realizadas

* Se descargó el libro en formato PDF desde la página de fundación Slim (https://aprende.org/pruebasat?sectionId=6)
* Para convertirlo en formato CSV se uso una página en línea que obtuvo todo el texto del PDF
* Una vez que se tuvo el archivo en texto plano, se eliminaron las mayusculas, los acentos y los signos de puntuación, usando para esto expresiones regulares
* Se creo una lista con todas las palabras del texto, la cual se recorrió para hacer el recuento del número de apariciones por palabra. Dando los siguientes resultados:
  * Total de palabras: 109261
  * Total de palabras no repetidas: 13105
* Se ordenó el arreglo tomando en cuenta el número de apariciones por cada palabra, resultando las siguientes palabras como las 100 más frecuentes:

    [('de', 5325), ('la', 3918), ('que', 3818), ('el', 3394), ('y', 3123), ('en', 2836), ('a', 2489), ('se', 1681), ('un', 1601), ('no', 1498), ('los', 1353), ('una', 1319), ('su', 1245), ('por', 936), ('las', 935), ('con', 924), ('habia', 858), ('del', 813), ('al', 756), ('es', 749), ('lo', 719), ('le', 667), ('era', 650), ('como', 572), ('mas', 513), ('para', 504), ('senor', 447), ('esta', 414), ('pero', 372), ('hombre', 363), ('si', 358), ('sus', 344), ('todo', 327), ('me', 326), ('sin', 311), ('obispo', 286), ('dijo', 281), ('cuando', 274), ('estaba', 273), ('sobre', 269), ('dos', 264), ('este', 261), ('aquel', 253), ('mi', 244), ('ya', 229), ('hacia', 219), ('yo', 218), ('esto', 218), ('madeleine', 214), ('tenia', 212), ('jean', 200), ('ha', 199), ('fantine', 194), ('valjean', 192), ('aquella', 190), ('hay', 186), ('he', 182), ('ser', 181), ('muy', 178), ('javert', 175), ('nada', 174), ('mismo', 173), ('o', 164), ('os', 163), ('poco', 158), ('tan', 158), ('bien', 157), ('ni', 156), ('ella', 155), ('quien', 151), ('alcalde', 149), ('vez', 148), ('despues', 146), ('fue', 145), ('todos', 141), ('puerta', 137), ('anos', 136), ('hubiera', 133), ('cual', 133), ('donde', 131), ('dios', 130), ('mujer', 127), ('momento', 125), ('tiempo', 124), ('sido', 124), ('casa', 123), ('son', 121), ('aqui', 120), ('noche', 119), ('hecho', 118), ('tres', 115), ('dia', 114), ('luego', 113), ('cabeza', 113), ('decir', 112), ('voz', 111), ('alli', 107), ('ojos', 107), ('monsenor', 105), ('aun', 105)]

* Y como las menos frecuentes con empate de una sola aparición las siguientes

    [('reprobo', 1), ('emocionantes', 1), ('florecer', 1), ('palidos', 1), ('fulgor', 1), ('sepultura', 1), ('recluyo', 1), ('lugares', 1), ('conversaciones', 1), ('bejean', 1), ('bojean', 1), ('boujean', 1), ('almibarado', 1), ('rehusado', 1), ('perillanes', 1), ('sucia', 1), ('abundaron', 1), ('abonada', 1), ('drapeau', 1), ('blanc', 1), ('ensenara', 1), ('partidarios', 1), ('despavorida', 1), ('reflexionando', 1), ('puestos', 1), ('velaban', 1), ('colgo', 1), ('esperara', 1), ('inconscientemente', 1), ('ensimismamiento', 1), ('candela', 1), ('boquiabierta', 1), ('retenido', 1), ('embargada', 1), ('barrote', 1), ('guardaria', 1), ('maestra', 1), ('lateral', 1), ('registrado', 1), ('conducian', 1), ('peldanos', 1), ('deshecha', 1), ('huella', 1), ('penultima', 1), ('obtuvo', 1), ('envolvio', 1), ('embalaba', 1), ('mordiendo', 1), ('comprobado', 1), ('migas', 1), ('encontradas', 1), ('pesquisas', 1), ('enrojecidos', 1), ('violencias', 1), ('integros', 1), ('entranas', 1), ('obligan', 1), ('doblar', 1), ('leerlo', 1), ('servira', 1), ('sonidos', 1), ('inarticulados', 1), ('persiguiendome', 1), ('turbaria', 1), ('alboroto', 1), ('murmullos', 1), ('protestas', 1), ('ambiente', 1), ('respirable', 1), ('integramente', 1), ('objecion', 1), ('correcto', 1), ('amuralladas', 1), ('retirarse', 1), ('quedarse', 1), ('aventurar', 1), ('desfallecer', 1), ('insista', 1), ('evadido', 1), ('mintio', 1), ('seguidas', 1), ('holocausto', 1), ('valga', 1), ('reparo', 1), ('singularidad', 1), ('bujia', 1), ('marchando', 1), ('brumas', 1), ('alejaba', 1), ('devuelta', 1), ('posiblemente', 1), ('reservar', 1), ('simplifico', 1), ('estricto', 1), ('enterrada', 1), ('gratuito', 1), ('cementerio', 1), ('encontrados', 1), ('sufrio', 1), ('promiscuidad', 1)]

* Para finalizar se creo un parquet con el vocabulario y el número de apariciones de cada palabra

## Conclusiones

* El reto no fue especialmente díficil
* Fue una buena práctica para tener un código con el cual se pueda obtener un vocabulario apartir de algún texto
* Nos queda pendiente generar un vocabulario más amplio para poder ocuparlo en el entrenamiento de un modelo
