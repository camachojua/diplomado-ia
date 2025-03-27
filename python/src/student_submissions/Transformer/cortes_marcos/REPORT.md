# Ejercicio Transformer

## Objetivo principal

Correr el notebook Neural Machine Translation with Transformer

## Estudiante

* Marcos Cortes Valadez

## Actividades realizadas

* Se corrigió la ruta donde se descarga la información de Anki (Español - Ingles) http://storage.googleapis.com/download.tensorflow.org/data/spa-eng.zip
* Se entrenó el modelo usando varios parametros:
  * Para una época
  * Para 5 épocas
  * Para 30 épocas
  * Para 30 épocas configurando ngrams=2
  * Para 50 épocas
* Adicionalmente se corrió el notebook "neural_machine_translation_with_keras_hub"

## Conclusiones
* Es bastante interesante que este método Transformer se ocupe para todo el tema de los chat. Seguiré haciendo pruebas con la idea de generar un modelo con mejor accuracy
* Se me ocurre que en lugar usar datos de los idiomas ingles-español podriamos ingresar datos de español a nahuatl o algún lenguaje indigena
* Apesar de hacer varios entrenamientos con diferentes configuraciones el módelo nunca tuvo un accuracy mayor a .30
* Mientras que el modelo usando Keras Hub tuvo un acurracy de 0.81 incluso con una sola época
* Comparando las traducciones de ambos modelos a mi parecer traduce mejor el que tiene un accuracy de 0.30 que el de 0.81
* Dedicaré más esfuerzo a usar el Transformer para crear chatbots
