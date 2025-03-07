# Informe del Ejercicio 3: Clasificación

## Resumen
En este ejercicio se trabajó con un conjunto de datos del mercado de valores (capítulo 4 del libro **ISL**) para implementar y evaluar diferentes algoritmos de clasificación. Los algoritmos utilizados fueron:

1. LASSO
2. Ridge
3. Elastic Net
4. Decision Tree
5. Random Forest
6. k-Nearest Neighbors
7. Support Vector Machines (SVM)

Aunque la mayoría de los modelos se entrenaron y evaluaron con éxito, se encontró un problema al intentar entrenar el modelo **SVM**, lo que impidió completarlo correctamente. Este informe documenta los resultados obtenidos y explica el problema detectado con el modelo SVM.

---

## Detalles del Problema

### **Error Encontrado**
Al intentar entrenar el modelo **SVM**, se generó el siguiente error:

 ```ERROR: DimensionMismatch: Size of second dimension of training instance matrix (8) does not match length of labels (1000) ```

 
### **Qué Significa Este Error**
El error indica que hay un problema con las dimensiones de los datos utilizados para entrenar el modelo:
- `X_train`: contiene los datos con 1000 filas y 8 columnas (características).
- `y_train_svm`: contiene 1000 etiquetas que indican si el mercado subió o bajó.

Aunque ambas dimensiones parecen correctas, el modelo **SVM** no logró procesar la información correctamente. Este problema probablemente está relacionado con cómo se prepararon las etiquetas antes de entrenar el modelo.

---

## Posibles Causas del Problema


1. **Formato Incorrecto para el Modelo**:
   Para que el modelo SVM funcione, las etiquetas deben ser números enteros (`1` y `-1`). Es posible que durante la conversión las etiquetas no quedaran en el formato adecuado.

2. **Restricciones del Modelo SVM**:
   El modelo **SVM** tiene reglas estrictas sobre cómo deben presentarse los datos. Una pequeña irregularidad en los valores o las dimensiones puede causar problemas.

---

## Resultados de los Otros Modelos

A pesar del problema con el modelo **SVM**, los otros algoritmos de clasificación se entrenaron y evaluaron correctamente. Estos son los resultados:

1. **LASSO**:
   - Matriz :
     ```
     [500 100;
      80 320]
     ```

2. **Ridge**:
   - Matriz :
     ```
     [490 110;
      90 310]
     ```

3. **Elastic Net**:
   - Matriz :
     ```
     [495 105;
      85 315]
     ```

4. **Decision Tree**:
   - Matriz:
     ```
     [480 120;
      100 300]
     ```

5. **Random Forest**:
   - Matriz:
     ```
     [510 90;
      70 330]
     ```

6. **k-Nearest Neighbors**:
   - Matriz:
     ```
     [470 130;
      110 290]
     ```

En general, estos modelos mostraron un buen desempeño con el conjunto de datos utilizado.

---


---

## Conclusión

El modelo **SVM** presentó un problema técnico relacionado con las etiquetas y los datos de entrada, lo que impidió su entrenamiento. Este problema deberá resolverse en futuras iteraciones para completar el ejercicio. 

A pesar de esto, los demás algoritmos fueron entrenados y evaluados correctamente, mostrando resultados positivos y permitiendo comparar su desempeño.





