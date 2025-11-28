#!/usr/bin/env python3
"""
Script para graficar los resultados de benchmarks de test_batch_modes.sh

Genera gráficos de rendimiento comparando los diferentes modos de ejecución:
- Modo Secuencial
- Modo Paralelo con diferentes números de workers
"""

import argparse
import re
import sys
import matplotlib.pyplot as plt
import numpy as np
from pathlib import Path


def parse_benchmark_output(filename):
    """
    Parse el archivo de salida de test_batch_modes.sh

    Args:
        filename: Ruta al archivo de texto con los resultados

    Returns:
        dict con la información parseada
    """
    with open(filename, "r") as f:
        content = f.read()

    data = {
        "sequences": 0,
        "comparisons": 0,
        "modes": [],
        "times": [],
        "patterns": [],
        "speedups": [],
    }

    # Extraer número de secuencias y comparaciones
    seq_match = re.search(r"Secuencias:\s*(\d+)", content)
    if seq_match:
        data["sequences"] = int(seq_match.group(1))

    comp_match = re.search(r"Comparaciones:\s*(\d+)", content)
    if comp_match:
        data["comparisons"] = int(comp_match.group(1))

    # Extraer tiempos de ejecución
    time_pattern = r"(Secuencial|Paralelo \((\d+) workers\))\s+(\d+)ms"
    for match in re.finditer(time_pattern, content):
        mode = match.group(1)
        time_ms = int(match.group(3))

        data["modes"].append(mode)
        data["times"].append(time_ms)

    # Extraer número de patrones
    pattern_pattern = r"(Secuencial|Paralelo \(\d+w\)):\s+(\d+)\s+patrones"
    for match in re.finditer(pattern_pattern, content):
        patterns = int(match.group(2))
        data["patterns"].append(patterns)

    # Extraer speedups
    speedup_pattern = r"(\d+)\s+workers\s+([\d.]+)x"
    for match in re.finditer(speedup_pattern, content):
        workers = int(match.group(1))
        speedup = float(match.group(2))
        data["speedups"].append((workers, speedup))

    return data


def create_execution_time_plot(data, output_dir):
    """Crea gráfico de barras con tiempos de ejecución"""
    fig, ax = plt.subplots(figsize=(10, 6))

    colors = ["#d62728", "#2ca02c", "#1f77b4", "#ff7f0e", "#9467bd"]

    x = np.arange(len(data["modes"]))
    bars = ax.bar(
        x,
        data["times"],
        color=colors[: len(data["modes"])],
        alpha=0.8,
        edgecolor="black",
    )

    # Añadir valores sobre las barras
    for i, (bar, time) in enumerate(zip(bars, data["times"])):
        height = bar.get_height()
        ax.text(
            bar.get_x() + bar.get_width() / 2.0,
            height,
            f"{time}ms",
            ha="center",
            va="bottom",
            fontweight="bold",
            fontsize=10,
        )

    ax.set_xlabel("Modo de Ejecución", fontsize=12, fontweight="bold")
    ax.set_ylabel("Tiempo (ms)", fontsize=12, fontweight="bold")
    ax.set_title(
        f"Tiempo de Ejecución por Modo\n({data['sequences']} secuencias, {data['comparisons']} comparaciones)",
        fontsize=14,
        fontweight="bold",
    )
    ax.set_xticks(x)
    ax.set_xticklabels(data["modes"], rotation=15, ha="right")
    ax.grid(axis="y", alpha=0.3, linestyle="--")

    plt.tight_layout()
    output_file = output_dir / "execution_times.png"
    plt.savefig(output_file, dpi=300, bbox_inches="tight")
    print(f"✓ Gráfico guardado: {output_file}")
    plt.close()


def create_speedup_plot(data, output_dir):
    """Crea gráfico de línea con speedup"""
    if not data["speedups"]:
        print("⚠ No hay datos de speedup disponibles")
        return

    fig, ax = plt.subplots(figsize=(10, 6))

    workers = [w for w, _ in data["speedups"]]
    speedups = [s for _, s in data["speedups"]]

    # Línea de speedup real
    ax.plot(
        workers,
        speedups,
        "o-",
        linewidth=2,
        markersize=10,
        color="#1f77b4",
        label="Speedup Real",
    )

    # Añadir valores
    for w, s in zip(workers, speedups):
        ax.text(
            w, s, f"{s:.2f}x", ha="center", va="bottom", fontweight="bold", fontsize=10
        )

    ax.set_xlabel("Número de Workers", fontsize=12, fontweight="bold")
    ax.set_ylabel("Speedup (vs. Secuencial)", fontsize=12, fontweight="bold")
    ax.set_title(
        f"Speedup vs Número de Workers\n({data['sequences']} secuencias, {data['comparisons']} comparaciones)",
        fontsize=14,
        fontweight="bold",
    )
    ax.legend(fontsize=11)
    ax.grid(alpha=0.3, linestyle="--")
    ax.set_xticks(workers)

    plt.tight_layout()
    output_file = output_dir / "speedup.png"
    plt.savefig(output_file, dpi=300, bbox_inches="tight")
    print(f"✓ Gráfico guardado: {output_file}")
    plt.close()


def create_efficiency_plot(data, output_dir):
    """Crea gráfico de eficiencia paralela"""
    if not data["speedups"]:
        print("⚠ No hay datos de eficiencia disponibles")
        return

    fig, ax = plt.subplots(figsize=(10, 6))

    workers = [w for w, _ in data["speedups"]]
    speedups = [s for _, s in data["speedups"]]
    efficiencies = [(s / w) * 100 for w, s in zip(workers, speedups)]

    bars = ax.bar(
        workers, efficiencies, color="#2ca02c", alpha=0.8, edgecolor="black", width=0.6
    )

    # Añadir valores sobre las barras
    for bar, eff in zip(bars, efficiencies):
        height = bar.get_height()
        ax.text(
            bar.get_x() + bar.get_width() / 2.0,
            height,
            f"{eff:.1f}%",
            ha="center",
            va="bottom",
            fontweight="bold",
            fontsize=10,
        )

    ax.set_xlabel("Número de Workers", fontsize=12, fontweight="bold")
    ax.set_ylabel("Eficiencia (%)", fontsize=12, fontweight="bold")
    ax.set_title(
        f"Eficiencia Paralela\n({data['sequences']} secuencias, {data['comparisons']} comparaciones)",
        fontsize=14,
        fontweight="bold",
    )
    ax.set_xticks(workers)
    ax.set_ylim([0, max(110, max(efficiencies) + 10)])
    ax.legend(fontsize=11)
    ax.grid(axis="y", alpha=0.3, linestyle="--")

    plt.tight_layout()
    output_file = output_dir / "efficiency.png"
    plt.savefig(output_file, dpi=300, bbox_inches="tight")
    print(f"✓ Gráfico guardado: {output_file}")
    plt.close()


def create_comparison_plot(data, output_dir):
    """Crea gráfico comparativo con múltiples métricas"""
    if len(data["times"]) < 2:
        print("⚠ No hay suficientes datos para comparación")
        return

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 6))

    # Subplot 1: Tiempos
    colors = ["#d62728", "#2ca02c", "#1f77b4", "#ff7f0e", "#9467bd"]
    x = np.arange(len(data["modes"]))
    bars1 = ax1.bar(
        x,
        data["times"],
        color=colors[: len(data["modes"])],
        alpha=0.8,
        edgecolor="black",
    )

    for bar, time in zip(bars1, data["times"]):
        height = bar.get_height()
        ax1.text(
            bar.get_x() + bar.get_width() / 2.0,
            height,
            f"{time}ms",
            ha="center",
            va="bottom",
            fontweight="bold",
            fontsize=9,
        )

    ax1.set_xlabel("Modo de Ejecución", fontsize=11, fontweight="bold")
    ax1.set_ylabel("Tiempo (ms)", fontsize=11, fontweight="bold")
    ax1.set_title("Tiempo de Ejecución", fontsize=12, fontweight="bold")
    ax1.set_xticks(x)
    ax1.set_xticklabels(data["modes"], rotation=15, ha="right", fontsize=9)
    ax1.grid(axis="y", alpha=0.3, linestyle="--")

    # Subplot 2: Speedup
    if data["speedups"]:
        workers = [w for w, _ in data["speedups"]]
        speedups = [s for _, s in data["speedups"]]

        ax2.plot(
            workers,
            speedups,
            "o-",
            linewidth=2,
            markersize=10,
            color="#1f77b4",
            label="Speedup Real",
        )

        for w, s in zip(workers, speedups):
            ax2.text(
                w,
                s,
                f"{s:.2f}x",
                ha="center",
                va="bottom",
                fontweight="bold",
                fontsize=9,
            )

        ax2.set_xlabel("Número de Workers", fontsize=11, fontweight="bold")
        ax2.set_ylabel("Speedup", fontsize=11, fontweight="bold")
        ax2.set_title("Speedup vs Workers", fontsize=12, fontweight="bold")
        ax2.legend(fontsize=10)
        ax2.grid(alpha=0.3, linestyle="--")
        ax2.set_xticks(workers)

    fig.suptitle(
        f"Análisis de Rendimiento BatchCompare\n{data['sequences']} secuencias, {data['comparisons']} comparaciones",
        fontsize=14,
        fontweight="bold",
        y=1.02,
    )

    plt.tight_layout()
    output_file = output_dir / "comparison.png"
    plt.savefig(output_file, dpi=300, bbox_inches="tight")
    print(f"✓ Gráfico guardado: {output_file}")
    plt.close()


def create_summary_table(data, output_dir):
    """Crea una tabla resumen en formato de imagen"""
    fig, ax = plt.subplots(figsize=(12, 6))
    ax.axis("tight")
    ax.axis("off")

    # Preparar datos de la tabla
    table_data = []
    table_data.append(["Modo", "Tiempo (ms)", "Patrones", "Speedup", "Eficiencia"])

    for i, mode in enumerate(data["modes"]):
        time = data["times"][i]
        patterns = data["patterns"][i] if i < len(data["patterns"]) else "-"

        # Encontrar speedup y eficiencia si es modo paralelo
        speedup = "-"
        efficiency = "-"

        if i > 0 and data["speedups"]:  # No para secuencial
            workers_match = re.search(r"(\d+)\s+workers", mode)
            if workers_match:
                workers = int(workers_match.group(1))
                for w, s in data["speedups"]:
                    if w == workers:
                        speedup = f"{s:.2f}x"
                        efficiency = f"{(s / w) * 100:.1f}%"
                        break

        table_data.append([mode, str(time), str(patterns), speedup, efficiency])

    # Crear tabla
    table = ax.table(
        cellText=table_data,
        cellLoc="center",
        loc="center",
        colWidths=[0.3, 0.15, 0.15, 0.15, 0.15],
    )

    table.auto_set_font_size(False)
    table.set_fontsize(10)
    table.scale(1, 2)

    # Estilo de la tabla
    for i in range(len(table_data)):
        for j in range(len(table_data[0])):
            cell = table[(i, j)]
            if i == 0:  # Header
                cell.set_facecolor("#4CAF50")
                cell.set_text_props(weight="bold", color="white")
            else:
                if i % 2 == 0:
                    cell.set_facecolor("#f0f0f0")
                else:
                    cell.set_facecolor("white")
                cell.set_edgecolor("gray")

    plt.title(
        f"Resumen de Benchmark\n{data['sequences']} secuencias, {data['comparisons']} comparaciones",
        fontsize=14,
        fontweight="bold",
        pad=20,
    )

    output_file = output_dir / "summary_table.png"
    plt.savefig(output_file, dpi=300, bbox_inches="tight")
    print(f"✓ Tabla guardada: {output_file}")
    plt.close()


def main():
    parser = argparse.ArgumentParser(
        description="Genera gráficos a partir de resultados de test_batch_modes.sh",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Ejemplos:
  %(prog)s benchmark_results.txt
  %(prog)s results.txt -o graficos/
  %(prog)s benchmark.txt --all
        """,
    )

    parser.add_argument(
        "input", type=str, help="Archivo de texto con los resultados del benchmark"
    )
    parser.add_argument(
        "-o",
        "--output-dir",
        type=str,
        default="plots",
        help="Directorio de salida para los gráficos (default: plots)",
    )
    parser.add_argument(
        "--all", action="store_true", help="Generar todos los gráficos posibles"
    )
    parser.add_argument(
        "--time", action="store_true", help="Generar solo gráfico de tiempos"
    )
    parser.add_argument(
        "--speedup", action="store_true", help="Generar solo gráfico de speedup"
    )
    parser.add_argument(
        "--efficiency", action="store_true", help="Generar solo gráfico de eficiencia"
    )
    parser.add_argument(
        "--comparison", action="store_true", help="Generar gráfico comparativo"
    )
    parser.add_argument("--table", action="store_true", help="Generar tabla resumen")

    args = parser.parse_args()

    # Verificar que el archivo existe
    input_path = Path(args.input)
    if not input_path.exists():
        print(f"❌ Error: El archivo {args.input} no existe", file=sys.stderr)
        sys.exit(1)

    # Crear directorio de salida
    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    print(f"📊 Generando gráficos desde: {args.input}")
    print(f"📁 Directorio de salida: {output_dir}")
    print()

    # Parsear datos
    try:
        data = parse_benchmark_output(args.input)
        print(
            f"✓ Datos parseados: {data['sequences']} secuencias, {data['comparisons']} comparaciones"
        )
        print(f"✓ Modos encontrados: {len(data['modes'])}")
        print()
    except Exception as e:
        print(f"❌ Error al parsear el archivo: {e}", file=sys.stderr)
        sys.exit(1)

    # Determinar qué gráficos generar
    if not (
        args.time or args.speedup or args.efficiency or args.comparison or args.table
    ):
        args.all = True

    # Generar gráficos
    try:
        if args.all or args.time:
            create_execution_time_plot(data, output_dir)

        if args.all or args.speedup:
            create_speedup_plot(data, output_dir)

        if args.all or args.efficiency:
            create_efficiency_plot(data, output_dir)

        if args.all or args.comparison:
            create_comparison_plot(data, output_dir)

        if args.all or args.table:
            create_summary_table(data, output_dir)

        print()
        print("✅ Todos los gráficos generados exitosamente!")
        print(f"📂 Ver gráficos en: {output_dir.absolute()}")

    except Exception as e:
        print(f"❌ Error al generar gráficos: {e}", file=sys.stderr)
        import traceback

        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
