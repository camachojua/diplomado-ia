# Reporte de comparación

Hernández Vela Daniel

## Observaciones generales

- El dataset CIFAR10 es bastante pesado por si mismo, por lo que el proceso de reescalar las imagenes al tamaño requerido por la arquitectura de AlexNet (224x224) es bastante lento.
- El proceso mencionado anteriormente fue realizado de manera muy eficiente en pytorch por lo que no hubo problemas relevantes para reescalar las imágenes.
- El mismo proceso fue complejo en TensorFlow, pues la función proporcionada para reescalar imágenes es muy poco eficiente y consume demasiada memoria RAM. La solución utilizada fue escealar cada imagen conforme se utilizaba, en lugar de reescalar todo el dataset antes del entrenamiento.

## Resultados

- El modelo entrenado desde 0 en PyTorch fue el que da mejores resultados, se obtuvo una precisión de 65.97%.
- El modelo entrenado desde 0 en TensorFlow dio una precisión de 38%.
- El modelo preentrenado en PyTorch fue el peor, obteniendo una precisión mínima.

### Argumentos

La precisión de los modelos en general no es demasiado alta debido al poder de cómputo del que se dispone, pudiendo realizar un entrenamiento relativamente simple de 10 epochs y learning rate de 0.001. 

Es interesante notar que el modelo entrenado desde 0 en PyTorch alcanzo una precisión "aceptable" dado el modesto entrenamiento que tuvo.

Por otro lado, el modelo entrenado desde 0 en TensorFlow alcanzo una precisión más "pobre", esto se puede deber a factores como los pesos aleatorios iniciales y como se dividió el conjunto de datos para entrenamiento y prueba.

Finalmente, el modelo preentrenado en PyTorch obtuvo una precisión prácticamente nula, por el hecho de que este fue preentrenado con un dataset diferente, y también porque la resolución de CIFAR10 es muy diferente a aquella predefinida para la arquitectura de AlexNet.

## Notas finales

Las precisiones obtenidas son un indicador general del desempeño de los modelos, sin embargo, será útil analizar otras métricas de desempeño que permitan analizar porque cada modelo obtuvo su respectiva precisión, igualmente encontrando si existe un área de oportunidad en el ajuste de hiperparámetros, procesamiento de datos, etc.