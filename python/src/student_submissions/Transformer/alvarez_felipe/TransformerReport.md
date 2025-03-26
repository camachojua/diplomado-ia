# Evaluación de un Modelo Transformer para Traducción Automática Inglés-Español

## Introducción

En este challenge, construí y entrené un modelo de traducción de inglés a español inspirado en el trabajo de François Chollet y el enfoque presentado en *Deep Learning with Python*. El modelo fue desarrollado y ejecutado en Google Colab, utilizando TensorFlow y Keras, y evaluado a través de métricas especializadas como **BLEU** y **ROUGE**.

La tarea se enfocó no solo en entrenar el modelo, sino también en analizar cómo ciertos hiperparámetros como el número de épocas, el optimizador ó el learning rate afectan el rendimiento del modelo.

---

## Entrenamiento Inicial (1 época)

Mi primer paso fue entrenar el modelo durante una sola época. Esta configuración fue útil para validar que todo el pipeline (desde la vectorización hasta la inferencia) funcionaba correctamente. Sin embargo, los resultados fueron limitados:

- **BLEU**: 0.18  
- **ROUGE-L**: 0.24

Estas puntuaciones reflejan un sistema apenas inicializado, con traducciones superficiales y frecuentes errores semánticos o sintácticos.

---

## Entrenamiento Completo (30 épocas)

Luego, ejecuté el modelo con una configuración más realista de **30 épocas**. Como era de esperarse, el rendimiento mejoró sustancialmente:

- **BLEU**: 0.47  
- **ROUGE-L**: 0.61  
- **ROUGE-1**: 0.66  
- **ROUGE-2**: 0.53

Las traducciones se volvieron significativamente más coherentes, gramaticalmente correctas y más fieles al significado original en inglés. Incluso oraciones complejas mostraban una buena estructura de sintaxis en español.

---

## Experimentación con Hiperparámetros

Para explorar la robustez del modelo, modifiqué distintos aspectos:


### **Optimización y tasa de aprendizaje**

Experimenté con distintos optimizadores:

- `RMSprop` (baseline)
- `Adam` (mejor desempeño)
- `SGD` (peor desempeño)

El uso de `Adam` con una tasa de aprendizaje de `1e-4` produjo mejores resultados en convergencia y estabilidad.


---

## Métricas de Evaluación: BLEU y ROUGE

### 🔷 ¿Qué es BLEU?

BLEU (Bilingual Evaluation Understudy) mide la superposición de *n-gramas* entre la traducción generada y una o más referencias humanas. Es una métrica ampliamente usada en traducción automática. Puntuaciones entre **0.4 y 0.6** ya indican traducciones razonablemente buenas.

### 🔶 ¿Qué es ROUGE?

ROUGE (Recall-Oriented Understudy for Gisting Evaluation) evalúa qué tan bien una traducción capta el contenido de la referencia. Mide *recall* y *F1-score* en n-gramas, siendo útil especialmente para tareas de resumen, aunque también se aplica a traducción. En este caso, usé `ROUGE-1`, `ROUGE-2` y `ROUGE-L` para complementar el análisis.

---

## Limitaciones Técnicas

Desafortunadamente, no pude experimentar tan a fondo como hubiera querido. El acceso a **Colab Pro** se agotó, y sin GPU disponible, el entrenamiento de modelos pesados como Transformers en CPU se vuelve extremadamente lento e impráctico. Esto limitó mi capacidad de ejecutar múltiples combinaciones de hiperparámetros y de extender el conjunto de pruebas.

---

## Conclusiones

Este experimento confirmó que el modelo Transformer es una arquitecturas eficaz para tareas de traducción. El entrenamiento durante una sola época fue suficiente para validar el pipeline, pero se requiere una mayor cantidad de iteraciones para obtener resultados competitivos.

El uso de métricas como BLEU y ROUGE me permitió evaluar de manera cuantitativa el progreso del modelo. Aunque no pude explorar todos los escenarios posibles, los resultados alcanzados fueron consistentes con la teoría y la literatura actual.
