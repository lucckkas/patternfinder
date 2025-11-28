#!/usr/bin/env python3
"""
Script para generar gráficos de los resultados de benchmark
Secuencial vs Paralelo
"""

import pandas as pd
import matplotlib.pyplot as plt
import sys
import os


def load_data(csv_file):
    """Carga los datos del CSV"""
    if not os.path.exists(csv_file):
        print(f"Error: No se encontró el archivo {csv_file}")
        print("Ejecuta primero: go run cmd/benchmark/main.go")
        sys.exit(1)

    return pd.read_csv(csv_file)


def plot_execution_time(df):
    """Genera gráfico de tiempo de ejecución"""
    fig, ax = plt.subplots(figsize=(12, 6))

    ax.plot(
        df["Length"], df["Seq_Total_ms"], marker="o", label="Secuencial", linewidth=2
    )
    ax.plot(df["Length"], df["Par_Total_ms"], marker="s", label="Paralelo", linewidth=2)

    ax.set_xlabel("Longitud de Secuencia", fontsize=12)
    ax.set_ylabel("Tiempo de Ejecución (ms)", fontsize=12)
    ax.set_title(
        "Comparación de Tiempo de Ejecución: Secuencial vs Paralelo",
        fontsize=14,
        fontweight="bold",
    )
    ax.legend(fontsize=11)
    ax.grid(True, alpha=0.3)
    ax.set_yscale("log")

    plt.tight_layout()
    plt.savefig("benchmark_results/execution_time.png", dpi=300, bbox_inches="tight")
    print("Gráfico guardado: benchmark_results/execution_time.png")


def plot_speedup(df):
    """Genera gráfico de speedup"""
    fig, ax = plt.subplots(figsize=(12, 6))

    ax.plot(
        df["Length"],
        df["Speedup_Total"],
        marker="o",
        label="Speedup Total",
        linewidth=2,
        color="green",
    )
    ax.axhline(y=1.0, color="r", linestyle="--", label="Sin mejora (1x)", alpha=0.7)

    ax.set_xlabel("Longitud de Secuencia", fontsize=12)
    ax.set_ylabel("Speedup (veces)", fontsize=12)
    ax.set_title(
        "Speedup de Versión Paralela vs Secuencial", fontsize=14, fontweight="bold"
    )
    ax.legend(fontsize=11)
    ax.grid(True, alpha=0.3)

    # Anotar el promedio de speedup
    avg_speedup = df["Speedup_Total"].mean()
    ax.text(
        0.02,
        0.98,
        f"Speedup promedio: {avg_speedup:.2f}x",
        transform=ax.transAxes,
        fontsize=11,
        verticalalignment="top",
        bbox=dict(boxstyle="round", facecolor="wheat", alpha=0.5),
    )

    plt.tight_layout()
    plt.savefig("benchmark_results/speedup.png", dpi=300, bbox_inches="tight")
    print("Gráfico guardado: benchmark_results/speedup.png")


def plot_components(df):
    """Genera gráfico comparando componentes (DP vs Backtracking)"""
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, 6))

    # DP Table
    ax1.plot(df["Length"], df["Seq_DP_ms"], marker="o", label="Secuencial", linewidth=2)
    ax1.plot(df["Length"], df["Par_DP_ms"], marker="s", label="Paralelo", linewidth=2)
    ax1.set_xlabel("Longitud de Secuencia", fontsize=12)
    ax1.set_ylabel("Tiempo de Ejecución (ms)", fontsize=12)
    ax1.set_title("Construcción de Tabla DP", fontsize=13, fontweight="bold")
    ax1.legend(fontsize=11)
    ax1.grid(True, alpha=0.3)
    ax1.set_yscale("log")

    # Backtracking
    ax2.plot(df["Length"], df["Seq_BT_ms"], marker="o", label="Secuencial", linewidth=2)
    ax2.plot(df["Length"], df["Par_BT_ms"], marker="s", label="Paralelo", linewidth=2)
    ax2.set_xlabel("Longitud de Secuencia", fontsize=12)
    ax2.set_ylabel("Tiempo de Ejecución (ms)", fontsize=12)
    ax2.set_title("Backtracking", fontsize=13, fontweight="bold")
    ax2.legend(fontsize=11)
    ax2.grid(True, alpha=0.3)
    ax2.set_yscale("log")

    plt.tight_layout()
    plt.savefig("benchmark_results/components.png", dpi=300, bbox_inches="tight")
    print("Gráfico guardado: benchmark_results/components.png")


def plot_lcs_count(df):
    """Genera gráfico del número de LCS encontradas"""
    fig, ax = plt.subplots(figsize=(12, 6))

    ax.bar(df["Length"], df["LCS_Count"], color="steelblue", alpha=0.7)
    ax.set_xlabel("Longitud de Secuencia", fontsize=12)
    ax.set_ylabel("Número de LCS encontradas", fontsize=12)
    ax.set_title(
        "Número de Subsecuencias Comunes Más Largas (LCS)",
        fontsize=14,
        fontweight="bold",
    )
    ax.grid(True, alpha=0.3, axis="y")
    ax.set_yscale("log")

    plt.tight_layout()
    plt.savefig("benchmark_results/lcs_count.png", dpi=300, bbox_inches="tight")
    print("Gráfico guardado: benchmark_results/lcs_count.png")


def print_statistics(df):
    """Imprime estadísticas del benchmark"""
    print("\n" + "=" * 60)
    print("ESTADÍSTICAS DEL BENCHMARK")
    print("=" * 60)

    print(f"\nSpeedup promedio: {df['Speedup_Total'].mean():.2f}x")
    print(
        f"Speedup máximo: {df['Speedup_Total'].max():.2f}x (longitud {df.loc[df['Speedup_Total'].idxmax(), 'Length']:.0f})"
    )
    print(
        f"Speedup mínimo: {df['Speedup_Total'].min():.2f}x (longitud {df.loc[df['Speedup_Total'].idxmin(), 'Length']:.0f})"
    )

    print(f"\nTiempo secuencial total: {df['Seq_Total_ms'].sum():.2f} ms")
    print(f"Tiempo paralelo total: {df['Par_Total_ms'].sum():.2f} ms")

    print(f"\nLCS total encontradas: {df['LCS_Count'].sum()}")
    print(
        f"LCS máximas en una ejecución: {df['LCS_Count'].max()} (longitud {df.loc[df['LCS_Count'].idxmax(), 'Length']:.0f})"
    )

    print("\n" + "=" * 60)


def main():
    csv_file = "benchmark_results/results.csv"

    print("Generando gráficos de benchmark...")
    print("=" * 60)

    # Cargar datos
    df = load_data(csv_file)

    # Generar gráficos
    plot_execution_time(df)
    plot_speedup(df)
    plot_components(df)
    plot_lcs_count(df)

    # Imprimir estadísticas
    print_statistics(df)

    print("\n✓ Todos los gráficos han sido generados exitosamente")
    print("  Los archivos están en el directorio: benchmark_results/")


if __name__ == "__main__":
    main()
