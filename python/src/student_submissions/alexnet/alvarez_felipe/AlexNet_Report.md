# Comparación de Modelos "From Scratch" vs. Preentrenados en CIFAR-10

## 1. Introducción

Este trabajo tuvo como objetivo comparar el rendimiento de modelos de redes neuronales AlexNet aplicados a CIFAR-10, implementados de dos formas:
- **Desde cero ("from scratch")**: Se implementó la arquitectura clásica de AlexNet tanto en PyTorch como en Keras.
- **Preentrenados (fine-tuning)**: Se realizó fine-tuning de un modelo preentrenado de AlexNet en PyTorch y, en Keras, se utilizó VGG16 preentrenado (dado que no existe una versión oficial de AlexNet preentrenada en Keras) como proxy, modificándolo para clasificación en 10 clases.

Todos los experimentos se ejecutaron en Google Colab utilizando una GPU A100.

## 2. Metodología

### 2.1. Conjunto de Datos y Preprocesamiento
- **Dataset**: CIFAR-10, compuesto por 60,000 imágenes de 32×32 píxeles en 10 clases.
- **Preprocesamiento**:
  - Las imágenes se redimensionaron a 224×224 o 227×227 (según el modelo) utilizando transformaciones homogéneas.
  - Se aplicó la normalización con los parámetros de ImageNet (por ejemplo, restando la media y dividiendo por la desviación estándar) para garantizar la compatibilidad con los pesos preentrenados.

### 2.2. Modelos Implementados
- **Modelos "from scratch"**:
  - *AlexNet en PyTorch*: Arquitectura clásica con 5 capas convolucionales, 3 MaxPooling y 3 capas densas, entrenada desde cero.
  - *AlexNet en Keras*: Implementación similar, entrenada desde cero.
- **Modelos preentrenados (fine-tuning)**:
  - *AlexNet preentrenada en PyTorch*: Se cargó el modelo preentrenado en ImageNet y se sustituyó la última capa para adaptarlo a 10 clases.
  - *VGG16 en Keras*: Se usó VGG16 preentrenado en ImageNet (debido a la ausencia de AlexNet preentrenada en Keras) y se añadió una cabeza de clasificación para CIFAR-10.

### 2.3. Configuración del Entrenamiento
- **Optimización**: Se empleó el optimizador SGD; para el fine-tuning se ajustaron las tasas de aprendizaje (más bajas) para mantener la estabilidad.
- **Hiperparámetros**:  
  - Número de épocas: 10  
  - Tamaño del batch: 64  
- **Hardware**: Google Colab con GPU A100.

## 3. Resultados

### 3.1. Modelos "From Scratch"
- **AlexNet en PyTorch**  
  - Época 1: Accuracy = 27.47%  
  - Época 2: Accuracy = 53.16%  
  - Época 3: Accuracy = 65.47%  
  - Época 4: Accuracy = 72.23%  
  - Época 5: Accuracy = 76.50%  
  - Época 6: Accuracy = 79.99%  
  - Época 7: Accuracy = 82.81%  
  - Época 8: Accuracy = 85.29%  
  - Época 9: Accuracy = 87.12%  
  - Época 10: Accuracy = 88.60%  
  - **Tiempo total de entrenamiento**: 438.89 s  
  - **Precisión en test**: 82.78%

- **AlexNet en Keras**  
  - Época 1: Accuracy = 30.96%, Val_accuracy = 56.33%  
  - Época 2: Accuracy = 58.80%, Val_accuracy = 69.73%  
  - Época 3: Accuracy = 70.21%, Val_accuracy = 72.53%  
  - Época 4: Accuracy = 75.69%, Val_accuracy = 76.30%  
  - Época 5: Accuracy = 79.59%, Val_accuracy = 74.53%  
  - Época 6: Accuracy = 82.72%, Val_accuracy = 79.46%  
  - Época 7: Accuracy = 85.51%, Val_accuracy = 78.30%  
  - Época 8: Accuracy = 87.95%, Val_accuracy = 79.49%  
  - Época 9: Accuracy = 90.15%, Val_accuracy = 79.07%  
  - Época 10: Accuracy = 91.28%, Val_accuracy = 79.13%  
  - **Tiempo total de entrenamiento**: 200.03 s  
  - **Precisión en test**: 79.13%

### 3.2. Modelos Preentrenados (Fine-tuning)
- **AlexNet preentrenada en PyTorch**  
  - Época 1: Accuracy = 69.11%  
  - Época 2: Accuracy = 80.45%  
  - Época 3: Accuracy = 84.88%  
  - Época 4: Accuracy = 87.04%  
  - Época 5: Accuracy = 88.89%  
  - Época 6: Accuracy = 90.48%  
  - Época 7: Accuracy = 91.83%  
  - Época 8: Accuracy = 92.89%  
  - Época 9: Accuracy = 93.37%  
  - Época 10: Accuracy = 93.96%  
  - **Tiempo total de entrenamiento**: 473.11 s  
  - **Precisión en test**: 86.11%

- **VGG16 en Keras (Fine-tuning)**  
  - Época 1: Accuracy = 53.42%, Val_accuracy = 83.25%  
  - Época 2: Accuracy = 76.72%, Val_accuracy = 84.71%  
  - Época 3: Accuracy = 81.47%, Val_accuracy = 85.49%  
  - Época 4: Accuracy = 84.64%, Val_accuracy = 86.16%  
  - Época 5: Accuracy = 86.84%, Val_accuracy = 86.56%  
  - Época 6: Accuracy = 88.60%, Val_accuracy = 86.92%  
  - Época 7: Accuracy = 90.18%, Val_accuracy = 87.33%  
  - Época 8: Accuracy = 91.33%, Val_accuracy = 87.42%  
  - Época 9: Accuracy = 92.26%, Val_accuracy = 87.69%  
  - Época 10: Accuracy = 93.29%, Val_accuracy = 87.96%  
  - **Tiempo total de entrenamiento**: 377.25 s  
  - **Precisión en test**: 87.96%

## 4. Análisis de Resultados

- **Comportamiento Inicial**:  
  Los modelos "from scratch" empezaron con una precisión muy baja en la primera época, lo que indica una necesidad de entrenamiento prolongado y ajuste cuidadoso de hiperparámetros. En contraste, los modelos preentrenados mostraron una ventaja inicial significativa debido a la transferencia de conocimientos adquiridos en ImageNet.

- **Evolución y Convergencia**:  
  - En el entrenamiento "from scratch", tanto PyTorch como Keras lograron mejoras sustanciales en precisión a lo largo de las 10 épocas, alcanzando cerca del 88.60% (PyTorch) y 91.28% (Keras) en las últimas épocas.  
  - El fine-tuning de los modelos preentrenados presentó una convergencia más rápida, alcanzando precisiones de test de 86.11% (AlexNet en PyTorch) y 87.96% (VGG16 en Keras).

- **Tiempo de Entrenamiento**:  
  Aunque los tiempos de entrenamiento varían entre frameworks (438.89 s y 200.03 s para los modelos "from scratch" en PyTorch y Keras, respectivamente), el uso de la GPU A100 facilitó el procesamiento de modelos complejos y aceleró el fine-tuning.

- **Comparabilidad de Resultados**:  
  Es importante destacar que, aunque se aplicaron preprocesamientos y configuraciones similares, los resultados no son directamente comparables debido a:
  - **Diferencias en la arquitectura**: VGG16 posee una estructura más profunda y compleja que AlexNet, lo que puede contribuir a diferencias en la capacidad de extracción de características.
  - **Preentrenamiento**: El modelo preentrenado de AlexNet en PyTorch se entrenó originalmente en ImageNet, mientras que en Keras se optó por VGG16 como proxy debido a la falta de un modelo preentrenado oficial de AlexNet.
  - **Condiciones de Entrenamiento y Framework**: Variaciones en optimizadores, tasas de aprendizaje y sobrecarga de cada framework pueden influir en el rendimiento final.

## 5. Conclusiones

- **Efectividad del Fine-Tuning**:  
  Los modelos preentrenados presentan una ventaja clara, mostrando una precisión elevada desde las primeras épocas y superando a los modelos entrenados desde cero en el conjunto de test.

- **Desafíos del Entrenamiento "From Scratch"**:  
  Comenzar el entrenamiento desde cero requiere mayor tiempo y ajuste de hiperparámetros para lograr convergencia y evitar problemas de inestabilidad.

- **Consideraciones en la Comparación**:  
  Aunque se intentó homogeneizar el preprocesamiento y las configuraciones de entrenamiento, las diferencias inherentes en la arquitectura y la disponibilidad de pesos preentrenados (especialmente en Keras) hacen que la comparación directa de los resultados de ambos enfoques deba interpretarse con cautela.

- **Impacto del Hardware**:  
  La ejecución en Google Colab con GPU A100 permitió realizar experimentos complejos en tiempos razonables, lo que resalta la importancia del hardware adecuado para el entrenamiento de modelos de deep learning.
