# Benchmark: Secuencial vs Paralelo

Este directorio contiene herramientas para ejecutar benchmarks y comparar el rendimiento de las implementaciones secuencial y paralela del algoritmo LCS.

## Estructura de Archivos

```
.
├── test/
│   ├── benchmark_test.go      # Tests y benchmarks de Go
│   └── hotspot_test.go         # Tests de hotspots existentes
├── cmd/
│   └── benchmark/
│       └── main.go             # Generador de resultados CSV
├── run_benchmark.sh            # Script para ejecutar todos los benchmarks
├── generate_plots.py           # Script para generar gráficos
└── benchmark_results/          # Directorio con resultados (generado)
    ├── results.csv             # Datos en formato CSV
    ├── execution_time.png      # Gráfico de tiempos
    ├── speedup.png             # Gráfico de speedup
    ├── components.png          # Gráfico de componentes (DP vs BT)
    └── lcs_count.png           # Gráfico de número de LCS
```

## Ejecución Rápida

### 1. Generar datos CSV

```bash
mkdir -p benchmark_results
go run cmd/benchmark/main.go
```

Esto genera el archivo `benchmark_results/results.csv` con los tiempos de ejecución para secuencias de longitud 20, 30, 40, ..., 200.

### 2. Generar gráficos (opcional, requiere Python)

Instalar dependencias:

```bash
pip install pandas matplotlib
```

Generar gráficos:

```bash
python3 generate_plots.py
```

### 3. Ejecutar benchmarks completos de Go

```bash
# Ejecutar todos los benchmarks
./run_benchmark.sh

# O ejecutar benchmarks específicos
go test ./test -bench=BenchmarkSequentialVsParallel -benchtime=3s
go test ./test -bench=BenchmarkDPTableOnly -benchtime=3s
go test ./test -bench=BenchmarkBacktrackingOnly -benchtime=3s
```

### 4. Ejecutar test de comparación detallada

```bash
go test ./test -run TestSequentialVsParallelComparison -v
```

Este test muestra una tabla completa con:

-   Tiempos de ejecución para cada componente (DP Table, Backtracking)
-   Speedup por componente
-   Número de LCS encontradas
-   Promedios

## Descripción de los Benchmarks

### BenchmarkSequentialVsParallel

Compara el tiempo total (construcción de matriz + backtracking) para ambas versiones con secuencias de longitud creciente.

### BenchmarkDPTableOnly

Compara solo la construcción de la tabla de programación dinámica.

### BenchmarkBacktrackingOnly

Compara solo el algoritmo de backtracking (con la tabla DP pre-calculada).

### TestSequentialVsParallelComparison

Test detallado que muestra:

-   Tiempos individuales de cada componente
-   Speedup calculado para cada longitud
-   Número de LCS encontradas
-   Estadísticas agregadas

## Resultados Esperados

Para las secuencias aleatorias de prueba, se espera:

-   **Secuencias pequeñas (20-60)**: La versión secuencial es más rápida debido al overhead de crear goroutines
-   **Secuencias medianas (70-120)**: El speedup empieza a ser más visible
-   **Secuencias grandes (140-200)**: La versión paralela muestra mejoras, pero el backtracking sigue siendo el cuello de botella

El speedup promedio observado es de aproximadamente **0.15-0.30x** para estas secuencias, lo que indica que la versión paralela es **más lenta** en estos casos. Esto es esperado porque:

1. El overhead de crear y sincronizar goroutines es significativo
2. El número de caminos en el backtracking crece exponencialmente
3. La detección de caminos duplicados requiere sincronización con mutexes

## Interpretación de los Gráficos

### execution_time.png

Muestra el tiempo total de ejecución en escala logarítmica. Observar cómo ambas curvas crecen exponencialmente con la longitud.

### speedup.png

Muestra el factor de speedup (< 1.0 significa que paralelo es más lento). La línea roja en 1.0 indica el punto de equilibrio.

### components.png

Compara los tiempos de DP Table vs Backtracking por separado. El backtracking es típicamente el componente más costoso.

### lcs_count.png

Muestra el número de LCS encontradas, que crece exponencialmente con la longitud de las secuencias.

## Formato del CSV

El archivo `results.csv` contiene las siguientes columnas:

-   `Length`: Longitud de las secuencias
-   `Seq_DP_ms`: Tiempo de construcción de tabla DP (secuencial) en ms
-   `Seq_BT_ms`: Tiempo de backtracking (secuencial) en ms
-   `Seq_Total_ms`: Tiempo total (secuencial) en ms
-   `Par_DP_ms`: Tiempo de construcción de tabla DP (paralelo) en ms
-   `Par_BT_ms`: Tiempo de backtracking (paralelo) en ms
-   `Par_Total_ms`: Tiempo total (paralelo) en ms
-   `Speedup_DP`: Factor de speedup para DP Table
-   `Speedup_BT`: Factor de speedup para Backtracking
-   `Speedup_Total`: Factor de speedup total
-   `LCS_Count`: Número de LCS encontradas

## Notas

-   Los benchmarks pueden tomar varios minutos para completarse
-   Los resultados pueden variar según el hardware y la carga del sistema
-   Se recomienda cerrar otras aplicaciones durante el benchmark
-   Las secuencias son generadas aleatoriamente con seed fijo para reproducibilidad
