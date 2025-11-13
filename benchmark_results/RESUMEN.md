# Resumen de Benchmark: Secuencial vs Paralelo

## 📊 Resultados Principales

Se ejecutaron benchmarks comparando las versiones **secuencial** y **paralela** del algoritmo LCS con secuencias aleatorias de longitud creciente (20, 30, 40, ..., 200).

### Hallazgos Clave

-   ✗ **La versión secuencial es más rápida** en todos los casos probados
-   📉 Speedup promedio: **0.17x** (paralelo es ~6x más lento)
-   🔍 El backtracking es el componente dominante en tiempo de ejecución
-   📈 El número de LCS encontradas crece exponencialmente con la longitud

## 📈 Estadísticas Detalladas

### Speedup por Componente

-   **DP Table**: 0.03x (paralelo es ~33x más lento)
-   **Backtracking**: 0.21x (paralelo es ~5x más lento)
-   **Total**: 0.17x (paralelo es ~6x más lento)

### Ejemplos de Tiempos

| Longitud | Secuencial | Paralelo | Speedup | LCS# |
| -------- | ---------- | -------- | ------- | ---- |
| 20       | 0.05 ms    | 0.48 ms  | 0.11x   | 4    |
| 50       | 0.21 ms    | 1.75 ms  | 0.12x   | 11   |
| 100      | 4.58 ms    | 25.90 ms | 0.18x   | 179  |
| 200      | 15.3 s     | 50.0 s   | 0.31x   | 186k |

## 🔍 Análisis

### ¿Por qué la versión paralela es más lenta?

1. **Overhead de goroutines**: Crear y gestionar goroutines tiene un costo significativo
2. **Sincronización**: Los mutexes para el registro de caminos visitados añaden latencia
3. **Naturaleza del problema**: El backtracking con detección de caminos duplicados requiere acceso compartido al estado
4. **Secuencias aleatorias**: Generan muchas bifurcaciones que requieren sincronización constante

### Componente más costoso

El **backtracking** domina el tiempo de ejecución, representando >95% del tiempo total para secuencias largas. Esto se debe a que:

-   El número de caminos crece exponencialmente
-   Cada camino debe verificarse contra el registro de visitados
-   La complejidad en el peor caso es O(2^n)

## 🎯 Cuándo usar cada versión

### Usar Versión Secuencial:

-   ✅ Secuencias pequeñas a medianas (< 200)
-   ✅ Secuencias aleatorias con muchas bifurcaciones
-   ✅ Cuando la latencia es crítica
-   ✅ Hardware con pocos cores

### Usar Versión Paralela:

-   ⚠️ Secuencias muy grandes (> 500) donde el DP Table es significativo
-   ⚠️ Secuencias con pocas bifurcaciones
-   ⚠️ Hardware con muchos cores (16+)
-   ⚠️ Cuando se procesan múltiples pares de secuencias en paralelo

## 🛠️ Archivos Generados

```
benchmark_results/
├── results.csv                    # Datos en formato CSV
├── comparison_detailed.txt        # Comparación detallada
├── benchmark_full.txt             # Benchmarks completos de Go
├── benchmark_dp.txt               # Benchmarks solo DP Table
├── benchmark_bt.txt               # Benchmarks solo Backtracking
├── execution_time.png             # Gráfico de tiempos (si se generó)
├── speedup.png                    # Gráfico de speedup (si se generó)
├── components.png                 # Gráfico de componentes (si se generó)
└── lcs_count.png                  # Gráfico de número de LCS (si se generó)
```

## 🚀 Cómo Ejecutar

### Ejecución rápida (solo datos CSV):

```bash
go run cmd/benchmark/main.go
python3 analyze_results.py
```

### Ejecución completa con benchmarks de Go:

```bash
./run_benchmark.sh
```

### Generar gráficos (requiere matplotlib):

```bash
pip install pandas matplotlib
python3 generate_plots.py
```

## 📝 Conclusiones

1. **La paralelización no siempre mejora el rendimiento**: El overhead puede superar los beneficios
2. **La detección de caminos duplicados requiere sincronización**: Esto añade latencia en la versión paralela
3. **El backtracking es inherentemente secuencial**: Los caminos dependen unos de otros
4. **La implementación secuencial es más simple y eficiente**: Para este caso de uso específico

### Recomendación

Para el caso de uso actual (secuencias de proteínas y ligandos), se recomienda:

-   Usar la **versión secuencial** para el procesamiento individual de pares
-   Considerar **paralelizar a nivel superior**: procesar múltiples pares de secuencias en paralelo
-   Optimizar el **algoritmo de backtracking** antes de paralelizar

## 📚 Referencias

-   Código fuente: `internal/lcs/lcs.go`
-   Tests: `test/benchmark_test.go`
-   Documentación: `BENCHMARK_README.md`
