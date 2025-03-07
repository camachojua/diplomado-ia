# AlexNet challenge

Elaboración del alexnet challenge que consiste en:

    Recrear el modelo AlexNet desde cero en tensorflow/keras y en pyTorch
    Usar el dataset CIFAR10 para entrenar el modelo y comparar los resultados con los modelos preentrenados de pyTorch/TF

#### Resultados

**PyTorch pre-entrenado**: Para el modelo importado de Torch y pre-entrenado se observó una exactitud general de __83.40%__, y una exactitud por categoria de

    Accuracy of plane : 91.07 %
    Accuracy of car   : 100.00 %
    Accuracy of bird  : 77.22 %
    Accuracy of cat   : 68.49 %
    Accuracy of deer  : 74.55 %
    Accuracy of dog   : 71.19 %
    Accuracy of frog  : 94.64 %
    Accuracy of horse : 85.94 %
    Accuracy of ship  : 84.48 %
    Accuracy of truck : 87.18 %

La ejecución de este modelo tanto en carga de datos como en entrenamiento fue bastante sencilla y los ajustes necesarios para obtener dicho resultado fueron mínimas, reduciéndose al ajuste de batch y entrenamiento.

**PyTroch desde cero**: Para el modelo generado a partir de cero en torch y con una sola ejecución a traves de todo el conjunto de entrenamiento, se obtuvo una exactitud general de __58.93 %__ con un valor de perdida de __loss: 1.218__, y una exactitud por categoria de

    Accuracy of plane : 66.07 %
    Accuracy of car   : 62.00 %
    Accuracy of bird  : 45.57 %
    Accuracy of cat   : 38.36 %
    Accuracy of deer  : 30.91 %
    Accuracy of dog   : 54.24 %
    Accuracy of frog  : 62.50 %
    Accuracy of horse : 59.38 %
    Accuracy of ship  : 68.97 %
    Accuracy of truck : 71.79 %

En el caso de este modelo, si bien el generar el modelo resultó más complicado que simplemente importarlo ya entrenado, la definición de la arquitectura fue relativamente sencilla, ya que no es una arquitectura demasiado extensa.

**Tensorflow/Keras desde cero**: El modelo de tensorflow/keras creado desde cero y con 3 entrenamientos de 10 épocas cada uno y recargando el mejor resultado del ciclo anterior, obtuvo una exactitud general de __67.3%__ y una perdida de __0.919__ durante el entrenamiento, pero registrando solamente __59.3%__ de exactitud y __1.207__ de perdida al evaluarlo sobre los datos de prueba, los cuales no había observado nunca antes.

#### Conclusiones

El modelo pre-entrenado de Alexnet fue sin lugar a dudas el que obtuvo los mejores resultados, esto tiene todo el sentido ya que es un modelo entrenado con una enorme cantidad de datos y durante un gran número de épocas, lo que le permite partir de un estado muy cómodo y adaptarse rápidamente y con buenos resultados a la clasificación de un grupo de imágenes.

Para el caso de los modelos generados desde cero, se obtuvieron resultados similares en ambos casos, rondando el 60% de exactitud general, esto es esperado ya que al ser un modelo entrenado a partir de cero y durante pocas épocas, se parte desde un punto de partida básico y aumentando poco a poco la exactitud del modelo, otro factor que afecta en su desempeño es la cantidad menor de datos de entrenamiento utilizados, ya que el estar utilizando el conjunto de CIFAR10 solamente tenemos 10 categorías de imágenes sobre las que se entrena y una cantidad relativamente pequeña, sabiendo esto podemos observar principalmente en el modelo de Tf/Keras un posible comportamiento de overfitting que se alcanza rápidamente y no permite generalizar de manera correcta para imágenes no vistas previamente, por esta razón podemos observar que la exactitud disminuye en el conjunto de test respecto al observado durante entrenamiento.

Como conclusión el modelo pre-entrando fue el ganador indiscutible, sin embargo los modelos generados desde cero si lograron tener una exactitud más allá de un valor aleatorio, que sin llegar al desempeño del modelo pre-entrenado si se tiene una clasificación aceptable. Con mayor cantidad de épocas al entrenar seguramente se logrará un mejor desempeño, también jugar con mayor variedad de hiperparametros como el batch size, learning rate, épocas, y utilidades de la red como dropout o padding también ayudará a mejorar el desempeño del modelo.
