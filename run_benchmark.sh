#!/bin/bash

# Script para ejecutar benchmarks y generar reporte de comparación
# Secuencial vs Paralelo

echo "=========================================="
echo "Benchmark: Secuencial vs Paralelo"
echo "=========================================="
echo ""

# Crear directorio para resultados si no existe
mkdir -p benchmark_results

# Ejecutar test de comparación detallada
echo "Ejecutando comparación detallada..."
go test ./test -run TestSequentialVsParallelComparison -v > benchmark_results/comparison_detailed.txt 2>&1

echo ""
echo "Resultados guardados en: benchmark_results/comparison_detailed.txt"
echo ""

# Ejecutar benchmarks oficiales de Go
echo "Ejecutando benchmarks de Go (esto puede tomar varios minutos)..."
go test ./test -bench=BenchmarkSequentialVsParallel -benchtime=2s -timeout=30m > benchmark_results/benchmark_full.txt 2>&1

echo "Resultados guardados en: benchmark_results/benchmark_full.txt"
echo ""

# Ejecutar benchmarks solo de DP Table
echo "Ejecutando benchmarks de construcción de tabla DP..."
go test ./test -bench=BenchmarkDPTableOnly -benchtime=2s > benchmark_results/benchmark_dp.txt 2>&1

echo "Resultados guardados en: benchmark_results/benchmark_dp.txt"
echo ""

# Ejecutar benchmarks solo de Backtracking
echo "Ejecutando benchmarks de backtracking..."
go test ./test -bench=BenchmarkBacktrackingOnly -benchtime=2s > benchmark_results/benchmark_bt.txt 2>&1

echo "Resultados guardados en: benchmark_results/benchmark_bt.txt"
echo ""

# Extraer resumen de la comparación detallada
echo "=========================================="
echo "RESUMEN DE RESULTADOS"
echo "=========================================="
echo ""

if [ -f benchmark_results/comparison_detailed.txt ]; then
    echo "Speedup promedio (de la comparación detallada):"
    grep "Speedup promedio" benchmark_results/comparison_detailed.txt
    echo ""
    
    echo "Tabla de resultados por longitud:"
    grep -A 20 "Len    | Seq Total" benchmark_results/comparison_detailed.txt | head -n 18
fi

echo ""
echo "=========================================="
echo "Todos los resultados han sido guardados en:"
echo "  - benchmark_results/comparison_detailed.txt"
echo "  - benchmark_results/benchmark_full.txt"
echo "  - benchmark_results/benchmark_dp.txt"
echo "  - benchmark_results/benchmark_bt.txt"
echo "=========================================="
