# PatternFinder - Detección de Patrones en Secuencias de Aminoácidos

Sistema completo para la detección, análisis y comparación de patrones en secuencias de aminoácidos (proteínas). Incluye herramientas de búsqueda, procesamiento por lotes, análisis estadístico y visualización de rendimiento.

## 📋 Tabla de Contenidos

-   [Requisitos](#requisitos)
-   [Instalación](#instalación)
-   [Herramientas Disponibles](#herramientas-disponibles)
    -   [1. PatternFinder](#1-patternfinder)
    -   [2. BatchCompare](#2-batchcompare)
    -   [3. Generate Sequences](#3-generate-sequences)
    -   [4. Generate Plots](#4-generate-plots)
    -   [5. Test Batch Modes](#5-test-batch-modes)
-   [Flujo de Trabajo Típico](#flujo-de-trabajo-típico)
-   [Ejemplos Completos](#ejemplos-completos)
-   [Formato de Datos](#formato-de-datos)

---

## 🔧 Requisitos

### Para las herramientas Go:

-   Go 1.18 o superior
-   Sistema operativo: Linux, macOS o Windows

### Para las herramientas Python:

-   Python 3.7 o superior
-   Bibliotecas: `matplotlib`, `numpy`

```bash
pip install matplotlib numpy
```

---

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/lucckkas/Memoria.git
cd Memoria
```

### 2. Compilar las herramientas Go

```bash
# Compilar PatternFinder
go build -o build/patternfinder cmd/patternfinder/main.go

# Compilar BatchCompare
go build -o build/batchcompare cmd/batchcompare/main.go
```

---

## 🛠️ Herramientas Disponibles

### 1. PatternFinder

Encuentra patrones comunes entre dos secuencias de aminoácidos usando LCS (Longest Common Subsequence) y detección de gaps.

#### Uso básico:

```bash
./build/patternfinder "AxxBxxxC" "AyyyyBzzzC"
```

#### Opciones:

-   `-dp`: Muestra la matriz LCS (para debugging)
-   `-seq`: Usa versión secuencial del algoritmo LCS (por defecto usa paralelo)

#### Ejemplo:

```bash
# Comparar dos secuencias
./build/patternfinder "AxxBxxxCxxxxD" "AyyyyByyyyyyyyCzzzzzD"

# Ver matriz LCS
./build/patternfinder -dp "ABCD" "AXBXCXD"
```

#### Salida:

```
Patrón: ABCD
Patrones expandidos:
  A-x(2)-B-x(3)-C-x(4)-D
  A-x(2)-B-x(3)-C-x(5)-D
  A-x(3)-B-x(4)-C-x(4)-D
  ...
```

---

### 2. BatchCompare

Compara múltiples secuencias en lotes, ejecutando PatternFinder para cada par posible. Soporta ejecución **paralela** y **secuencial**.

#### Uso básico:

```bash
./build/batchcompare -f secuencias.txt -csv resultados.csv
```

#### Opciones principales:

| Opción           | Descripción                             | Default               |
| ---------------- | --------------------------------------- | --------------------- |
| `-f <archivo>`   | Archivo con secuencias (una por línea)  | **REQUERIDO**         |
| `-csv <archivo>` | Genera CSV con estadísticas de patrones | -                     |
| `-w <número>`    | Número de workers paralelos             | 6                     |
| `-seq`           | Modo secuencial (sin paralelización)    | false                 |
| `-o <archivo>`   | Archivo de salida para resultados       | stdout                |
| `-p <path>`      | Ruta al ejecutable patternfinder        | ./build/patternfinder |
| `-dp`            | Muestra matriz LCS (debug)              | false                 |

#### Ejemplos:

```bash
# Modo paralelo con 8 workers
./build/batchcompare -f sec.txt -w 8 -csv stats.csv

# Modo secuencial (útil para debugging)
./build/batchcompare -f sec.txt -seq -csv stats.csv

# Guardar resultados en archivo
./build/batchcompare -f sec.txt -w 4 -o resultados.txt -csv stats.csv
```

#### Formato del CSV generado:

```csv
Patron,Mayusculas,Secuencias,Porcentaje
ABCD,4,15,75.00
ABC,3,18,90.00
AB,2,20,100.00
```

-   **Patron**: Patrón detectado (letras mayúsculas del LCS)
-   **Mayusculas**: Número de caracteres en el patrón
-   **Secuencias**: Cuántas secuencias tienen este patrón
-   **Porcentaje**: % de secuencias con el patrón

---

### 3. Generate Sequences

Genera secuencias aleatorias de aminoácidos para pruebas y benchmarks.

#### Uso básico:

```bash
./generate_sequences.py -n 10 -l 100 -o secuencias.txt
```

#### Opciones:

| Opción            | Descripción                            | Default       |
| ----------------- | -------------------------------------- | ------------- |
| `-n <número>`     | Número de secuencias a generar         | 10            |
| `-l <longitud>`   | Longitud de cada secuencia             | 100           |
| `-u <porcentaje>` | % de letras en mayúsculas (0.0-1.0)    | 0.2           |
| `-o <archivo>`    | Archivo de salida                      | sequences.txt |
| `--min-len <n>`   | Longitud mínima (con --variable)       | 50            |
| `--max-len <n>`   | Longitud máxima (con --variable)       | 200           |
| `--variable`      | Genera secuencias de longitud variable | false         |
| `--seed <n>`      | Semilla para reproducibilidad          | -             |

#### Ejemplos:

```bash
# Generar 50 secuencias de 150 aminoácidos
./generate_sequences.py -n 50 -l 150 -o test.txt

# Secuencias con 30% de mayúsculas
./generate_sequences.py -n 20 -l 100 -u 0.3 -o high_upper.txt

# Longitud variable entre 80 y 200
./generate_sequences.py -n 30 --variable --min-len 80 --max-len 200 -o var.txt

# Con semilla para reproducibilidad
./generate_sequences.py -n 10 -l 100 --seed 42 -o reproducible.txt
```

#### Aminoácidos utilizados:

```
A C D E F G H I K L M N P Q R S T V W Y
```

(Los 20 aminoácidos estándar)

---

### 4. Generate Plots

Genera gráficos de rendimiento a partir de los resultados de benchmarks.

#### Uso básico:

```bash
./generate_plots2.py benchmark_output.txt
```

#### Opciones:

| Opción            | Descripción                              |
| ----------------- | ---------------------------------------- |
| `-o <directorio>` | Directorio de salida (default: `plots/`) |
| `--all`           | Genera todos los gráficos (default)      |
| `--time`          | Solo gráfico de tiempos de ejecución     |
| `--speedup`       | Solo gráfico de speedup                  |
| `--comparison`    | Gráfico comparativo                      |
| `--table`         | Tabla resumen                            |

#### Gráficos generados:

1. **execution_times.png** - Barras con tiempos de ejecución
2. **speedup.png** - Speedup real vs ideal
3. **comparison.png** - Comparativo tiempo + speedup
4. **summary_table.png** - Tabla resumen con todas las métricas

#### Ejemplos:

```bash
# Generar todos los gráficos
./generate_plots2.py benchmark_output.txt

# Solo speedup
./generate_plots2.py results.txt --speedup

# Guardar en directorio específico
./generate_plots2.py benchmark.txt -o graficos/
```

#### Métricas calculadas:

-   **Speedup**: $\text{Speedup} = \frac{\text{Tiempo Secuencial}}{\text{Tiempo Paralelo}}$

---

### 5. Test Batch Modes

Script de benchmark que compara el rendimiento de los modos secuencial y paralelo.

#### Uso:

```bash
./test_batch_modes.sh
```

#### Requisitos:

-   Archivo `sec.txt` con secuencias de prueba
-   Ejecutable `build/batchcompare` compilado

#### Qué hace:

1. Ejecuta BatchCompare en modo secuencial
2. Ejecuta BatchCompare en modo paralelo con 2, 4, 8 workers
3. Mide tiempos de ejecución
4. Calcula speedup
5. Genera CSVs de estadísticas
6. Guarda resultados en `benchmark_output.txt`

#### Salida ejemplo:

```
================================================
Benchmark BatchCompare: Secuencial vs Paralelo
================================================

Archivo: sec.txt
Secuencias: 20
Comparaciones: 190

=== MODO SECUENCIAL ===
Ejecutando... ✓ Completado en 4523ms

=== MODO PARALELO (2 workers) ===
Ejecutando... ✓ Completado en 2410ms

=== ANÁLISIS DE RESULTADOS ===
Secuencial          4523ms
Paralelo (2w)       2410ms  Speedup: 1.88x
Paralelo (4w)       1305ms  Speedup: 3.47x
Paralelo (8w)        892ms  Speedup: 5.07x
```

---

## 🔄 Flujo de Trabajo Típico

### 1. Generar datos de prueba

```bash
./generate_sequences.py -n 30 -l 120 -u 0.25 -o test_sequences.txt
```

### 2. Ejecutar análisis batch

```bash
./build/batchcompare -f test_sequences.txt -w 8 -csv resultados.csv
```

### 3. Ejecutar benchmark

```bash
# Asegúrate de tener sec.txt con tus secuencias
./test_batch_modes.sh > benchmark_output.txt
```

### 4. Generar visualizaciones

```bash
./generate_plots2.py benchmark_output.txt -o graficos/
```

### 5. Analizar resultados

```bash
# Ver CSV de patrones
cat resultados.csv

# Ver gráficos
xdg-open graficos/speedup.png
xdg-open graficos/efficiency.png
```

---

## 📊 Formato de Datos

### Archivo de secuencias (entrada)

Archivo de texto con una secuencia por línea:

```
AxxBxxxCxxxxD
AyyyyByyyyyyyyCzzzzzD
MxxxxxNxxxOxxP
...
```

-   **Mayúsculas**: Aminoácidos importantes (patrón)
-   **Minúsculas**: Gaps/espaciadores
-   Cada línea es una secuencia

### CSV de estadísticas (salida)

```csv
Patron,Mayusculas,Secuencias,Porcentaje
ABCD,4,25,83.33
ABC,3,28,93.33
```

---

## 📈 Interpretación de Resultados

### Speedup

-   **Valor ideal**: Igual al número de workers (lineal)
-   **Bueno**: 80-90% del ideal
-   **Aceptable**: 50-80% del ideal
-   **Bajo**: <50% del ideal (overhead elevado)

### Patrones encontrados

-   Analiza qué patrones son más comunes
-   Mayor **Porcentaje** = patrón más conservado
-   Mayor **Mayusculas** = patrón más largo/complejo

---

## 🐛 Debugging

### PatternFinder no encuentra patrones

```bash
# Verificar que hay mayúsculas en ambas secuencias
echo "Tu secuencia debe tener MAYUSCULAS"

# Ver matriz LCS para entender el proceso
./build/patternfinder -dp "ABCD" "AXBXCXD"
```

### BatchCompare muy lento

```bash
# Probar con menos workers
./build/batchcompare -f sec.txt -w 2

# O usar modo secuencial para debug
./build/batchcompare -f sec.txt -seq
```

### Generate Plots falla

```bash
# Verificar que matplotlib está instalado
pip install matplotlib numpy

# Verificar formato del archivo de entrada
cat benchmark_output.txt
```

---

## 📝 Notas Importantes

1. **Secuencias grandes**: Para secuencias >500 aminoácidos, considera usar `-seq` en patternfinder
2. **Memoria**: El algoritmo optimizado usa O(1) memoria por recursión
3. **Cores**: El speedup máximo está limitado por el número de cores físicos
4. **Formato**: Las mayúsculas son los aminoácidos del patrón, minúsculas son gaps

---

## 📚 Referencias

-   **LCS Algorithm**: Longest Common Subsequence (Dynamic Programming)
-   **Gap Detection**: DFS con poda temprana y búsqueda binaria
-   **Pattern Expansion**: Producto cartesiano de valores de gaps
-   **Paralelización**: Worker pools con goroutines de Go

---

## 👥 Autor

Luckas - Universidad de Chile

## 📄 Licencia

Este proyecto es parte de una tesis de memoria universitaria.
