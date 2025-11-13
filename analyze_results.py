#!/usr/bin/env python3
"""
Script simple para analizar resultados sin dependencias externas
"""

import csv
import sys


def load_csv(filename):
    """Carga el CSV y devuelve los datos"""
    data = []
    with open(filename, "r") as f:
        reader = csv.DictReader(f)
        for row in reader:
            data.append(
                {
                    "Length": int(row["Length"]),
                    "Seq_Total_ms": float(row["Seq_Total_ms"]),
                    "Par_Total_ms": float(row["Par_Total_ms"]),
                    "Speedup_Total": float(row["Speedup_Total"]),
                    "Speedup_DP": float(row["Speedup_DP"]),
                    "Speedup_BT": float(row["Speedup_BT"]),
                    "LCS_Count": int(row["LCS_Count"]),
                }
            )
    return data


def print_table(data):
    """Imprime tabla formateada"""
    print("\n" + "=" * 90)
    print("RESULTADOS DEL BENCHMARK: SECUENCIAL vs PARALELO")
    print("=" * 90)
    print()
    print(
        f"{'Len':<6} | {'Seq (ms)':<12} | {'Par (ms)':<12} | {'Speedup':<10} | {'LCS#':<8} | {'Observación':<20}"
    )
    print("-" * 90)

    for row in data:
        seq = row["Seq_Total_ms"]
        par = row["Par_Total_ms"]
        speedup = row["Speedup_Total"]
        lcs = row["LCS_Count"]
        length = row["Length"]

        # Determinar observación
        if speedup >= 1.0:
            obs = "✓ Paralelo mejor"
        elif speedup >= 0.5:
            obs = "≈ Similar"
        elif speedup >= 0.2:
            obs = "✗ Secuencial mejor"
        else:
            obs = "✗✗ Secuencial mucho mejor"

        print(
            f"{length:<6} | {seq:>12.2f} | {par:>12.2f} | {speedup:>10.2f}x | {lcs:<8} | {obs:<20}"
        )

    print("-" * 90)


def print_statistics(data):
    """Imprime estadísticas agregadas"""
    total_speedup = sum(row["Speedup_Total"] for row in data)
    avg_speedup = total_speedup / len(data)

    speedup_dp = sum(row["Speedup_DP"] for row in data) / len(data)
    speedup_bt = sum(row["Speedup_BT"] for row in data) / len(data)

    total_lcs = sum(row["LCS_Count"] for row in data)

    max_speedup_row = max(data, key=lambda x: x["Speedup_Total"])
    min_speedup_row = min(data, key=lambda x: x["Speedup_Total"])

    print()
    print("=" * 90)
    print("ESTADÍSTICAS")
    print("=" * 90)
    print()
    print(f"  Speedup promedio total:         {avg_speedup:>6.2f}x")
    print(f"  Speedup promedio DP Table:      {speedup_dp:>6.2f}x")
    print(f"  Speedup promedio Backtracking:  {speedup_bt:>6.2f}x")
    print()
    print(
        f"  Speedup máximo:  {max_speedup_row['Speedup_Total']:>6.2f}x  (longitud {max_speedup_row['Length']})"
    )
    print(
        f"  Speedup mínimo:  {min_speedup_row['Speedup_Total']:>6.2f}x  (longitud {min_speedup_row['Length']})"
    )
    print()
    print(f"  Total de LCS encontradas: {total_lcs:,}")
    print()
    print("=" * 90)


def print_analysis(data):
    """Imprime análisis de los resultados"""
    print()
    print("ANÁLISIS")
    print("=" * 90)
    print()
    print("El speedup promedio < 1.0 indica que la versión SECUENCIAL es más rápida")
    print("que la versión PARALELA para estas secuencias aleatorias.")
    print()
    print("Razones principales:")
    print("  1. Overhead de crear y sincronizar goroutines")
    print("  2. Sincronización con mutexes en el registro de caminos visitados")
    print("  3. El número de caminos crece exponencialmente con la longitud")
    print("  4. Las secuencias aleatorias generan muchas bifurcaciones")
    print()
    print("La paralelización podría beneficiar más en:")
    print("  • Secuencias con pocas bifurcaciones")
    print("  • Matrices muy grandes donde el DP Table domina el tiempo")
    print("  • Hardware con muchos cores disponibles")
    print()
    print("=" * 90)


def main():
    csv_file = "benchmark_results/results.csv"

    try:
        data = load_csv(csv_file)
    except FileNotFoundError:
        print(f"Error: No se encontró el archivo {csv_file}")
        print("Ejecuta primero: go run cmd/benchmark/main.go")
        sys.exit(1)

    print_table(data)
    print_statistics(data)
    print_analysis(data)

    print()
    print("Nota: Para generar gráficos, instala pandas y matplotlib:")
    print("      pip install pandas matplotlib")
    print("      python3 generate_plots.py")
    print()


if __name__ == "__main__":
    main()
