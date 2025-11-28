#!/bin/bash

# Script de benchmark de modos secuencial y paralelo del batchcompare

echo "================================================"
echo "Benchmark BatchCompare: Secuencial vs Paralelo"
echo "================================================"
echo ""

# Verificar que existe el archivo de secuencias
if [ ! -f "sec.txt" ]; then
    echo "Error: No se encuentra el archivo sec.txt"
    exit 1
fi

# Contar secuencias
NUM_SEQ=$(grep -v "^#" sec.txt | grep -v "^$" | wc -l)
NUM_COMP=$((NUM_SEQ * (NUM_SEQ - 1) / 2))

echo "Archivo: sec.txt"
echo "Secuencias: $NUM_SEQ"
echo "Comparaciones: $NUM_COMP"
echo ""

# Función para medir tiempo en milisegundos
measure_time() {
    local start=$(date +%s%N)
    "$@" > /dev/null 2>&1
    local end=$(date +%s%N)
    local elapsed=$(( (end - start) / 1000000 ))
    echo $elapsed
}

# Modo SECUENCIAL
echo "=== MODO SECUENCIAL ==="
echo -n "Ejecutando... "
TIME_SEQ=$(measure_time ./build/batchcompare -f sec.txt -seq -csv seq_stats.csv)
echo "✓ Completado en ${TIME_SEQ}ms"
echo "  - CSV: seq_stats.csv"
echo ""

# Modo PARALELO con 2 workers
echo "=== MODO PARALELO (2 workers) ==="
echo -n "Ejecutando... "
TIME_PAR2=$(measure_time ./build/batchcompare -f sec.txt -w 2 -csv par2_stats.csv)
echo "✓ Completado en ${TIME_PAR2}ms"
echo "  - CSV: par2_stats.csv"
echo ""

# Modo PARALELO con 4 workers
echo "=== MODO PARALELO (4 workers) ==="
echo -n "Ejecutando... "
TIME_PAR4=$(measure_time ./build/batchcompare -f sec.txt -w 4 -csv par4_stats.csv)
echo "✓ Completado en ${TIME_PAR4}ms"
echo "  - CSV: par4_stats.csv"
echo ""

# Modo PARALELO con 8 workers
echo "=== MODO PARALELO (8 workers) ==="
echo -n "Ejecutando... "
TIME_PAR8=$(measure_time ./build/batchcompare -f sec.txt -w 8 -csv par8_stats.csv)
echo "✓ Completado en ${TIME_PAR8}ms"
echo "  - CSV: par8_stats.csv"
echo ""

# Modo PARALELO con 12 workers
echo "=== MODO PARALELO (12 workers) ==="
echo -n "Ejecutando... "
TIME_PAR12=$(measure_time ./build/batchcompare -f sec.txt -w 12 -csv par12_stats.csv)
echo "✓ Completado en ${TIME_PAR12}ms"
echo "  - CSV: par12_stats.csv"
echo ""

# Comparar resultados
echo "=== RESUMEN DE RESULTADOS ==="
echo ""
echo "Patrones únicos encontrados:"
echo "  - Secuencial:     $(tail -n +2 seq_stats.csv 2>/dev/null | wc -l) patrones"
echo "  - Paralelo (2w):  $(tail -n +2 par2_stats.csv 2>/dev/null | wc -l) patrones"
echo "  - Paralelo (4w):  $(tail -n +2 par4_stats.csv 2>/dev/null | wc -l) patrones"
echo "  - Paralelo (8w):  $(tail -n +2 par8_stats.csv 2>/dev/null | wc -l) patrones"
echo "  - Paralelo (12w): $(tail -n +2 par12_stats.csv 2>/dev/null | wc -l) patrones"
echo ""

echo "Tiempos de ejecución:"
printf "  %-20s %8s\n" "Modo" "Tiempo"
printf "  %-20s %8s\n" "----" "------"
printf "  %-20s %7dms\n" "Secuencial" $TIME_SEQ
printf "  %-20s %7dms\n" "Paralelo (2 workers)" $TIME_PAR2
printf "  %-20s %7dms\n" "Paralelo (4 workers)" $TIME_PAR4
printf "  %-20s %7dms\n" "Paralelo (8 workers)" $TIME_PAR8
printf "  %-20s %7dms\n" "Paralelo (12 workers)" $TIME_PAR12
echo ""

# Calcular speedup
if [ $TIME_SEQ -gt 0 ]; then
    SPEEDUP2=$(awk "BEGIN {printf \"%.2f\", $TIME_SEQ / $TIME_PAR2}")
    SPEEDUP4=$(awk "BEGIN {printf \"%.2f\", $TIME_SEQ / $TIME_PAR4}")
    SPEEDUP8=$(awk "BEGIN {printf \"%.2f\", $TIME_SEQ / $TIME_PAR8}")
    SPEEDUP12=$(awk "BEGIN {printf \"%.2f\", $TIME_SEQ / $TIME_PAR12}")
    
    echo "Speedup (vs. secuencial):"
    printf "  %-20s %6sx\n" "2 workers" $SPEEDUP2
    printf "  %-20s %6sx\n" "4 workers" $SPEEDUP4
    printf "  %-20s %6sx\n" "8 workers" $SPEEDUP8
    printf "  %-20s %6sx\n" "12 workers" $SPEEDUP12
    echo ""
fi

echo "Archivos CSV generados:"
ls -lh *_stats.csv 2>/dev/null | awk '{printf "  %-25s %6s\n", $9, $5}'
