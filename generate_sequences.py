#!/usr/bin/env python3
"""
Generador de secuencias aleatorias de aminoácidos.

Este script genera secuencias aleatorias usando los 20 aminoácidos estándar.
Puede generar múltiples secuencias de longitud configurable.
"""

import random
import argparse
import sys

# 20 aminoácidos estándar (código de una letra)
AMINO_ACIDS = [
    "A",  # Alanina
    "C",  # Cisteína
    "D",  # Ácido aspártico
    "E",  # Ácido glutámico
    "F",  # Fenilalanina
    "G",  # Glicina
    "H",  # Histidina
    "I",  # Isoleucina
    "K",  # Lisina
    "L",  # Leucina
    "M",  # Metionina
    "N",  # Asparagina
    "P",  # Prolina
    "Q",  # Glutamina
    "R",  # Arginina
    "S",  # Serina
    "T",  # Treonina
    "V",  # Valina
    "W",  # Triptófano
    "Y",  # Tirosina
]


def generate_sequence(length, uppercase_ratio=0.5, seed=None):
    """
    Genera una secuencia aleatoria de aminoácidos.

    Args:
        length (int): Longitud de la secuencia
        uppercase_ratio (float): Proporción de letras en mayúsculas (0.0 a 1.0)
        seed (int, optional): Semilla para reproducibilidad

    Returns:
        str: Secuencia generada
    """
    if seed is not None:
        random.seed(seed)

    sequence = []
    for _ in range(length):
        aa = random.choice(AMINO_ACIDS)
        # Decidir si es mayúscula o minúscula según la proporción
        if random.random() < uppercase_ratio:
            sequence.append(aa.upper())
        else:
            sequence.append(aa.lower())

    return "".join(sequence)


def generate_multiple_sequences(
    num_sequences,
    length,
    uppercase_ratio=0.5,
    min_length=None,
    max_length=None,
    seed=None,
):
    """
    Genera múltiples secuencias aleatorias.

    Args:
        num_sequences (int): Número de secuencias a generar
        length (int): Longitud base de las secuencias
        uppercase_ratio (float): Proporción de letras en mayúsculas
        min_length (int, optional): Longitud mínima (para longitud variable)
        max_length (int, optional): Longitud máxima (para longitud variable)
        seed (int, optional): Semilla para reproducibilidad

    Returns:
        list: Lista de secuencias generadas
    """
    if seed is not None:
        random.seed(seed)

    sequences = []
    for i in range(num_sequences):
        # Si se especifica rango de longitud, usar longitud aleatoria
        if min_length is not None and max_length is not None:
            seq_length = random.randint(min_length, max_length)
        else:
            seq_length = length

        seq = generate_sequence(seq_length, uppercase_ratio, seed=None)
        sequences.append(seq)

    return sequences


def main():
    parser = argparse.ArgumentParser(
        description="Genera secuencias aleatorias de aminoácidos",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Ejemplos:
  %(prog)s -l 50 -n 10                    # 10 secuencias de 50 aminoácidos
  %(prog)s -l 100 -n 5 -u 0.3             # 5 secuencias, 30%% mayúsculas
  %(prog)s --min 20 --max 100 -n 20       # 20 secuencias de longitud variable
  %(prog)s -l 80 -n 15 -s 42 -o seqs.txt  # Con semilla y archivo de salida
        """,
    )

    parser.add_argument(
        "-l",
        "--length",
        type=int,
        default=50,
        help="Longitud de las secuencias (default: 50)",
    )
    parser.add_argument(
        "-n",
        "--num-sequences",
        type=int,
        default=10,
        help="Número de secuencias a generar (default: 10)",
    )
    parser.add_argument(
        "-u",
        "--uppercase-ratio",
        type=float,
        default=0.5,
        help="Proporción de mayúsculas (0.0-1.0, default: 0.5)",
    )
    parser.add_argument(
        "--min",
        "--min-length",
        type=int,
        dest="min_length",
        help="Longitud mínima (para longitud variable)",
    )
    parser.add_argument(
        "--max",
        "--max-length",
        type=int,
        dest="max_length",
        help="Longitud máxima (para longitud variable)",
    )
    parser.add_argument("-s", "--seed", type=int, help="Semilla para reproducibilidad")
    parser.add_argument(
        "-o", "--output", type=str, help="Archivo de salida (default: stdout)"
    )
    parser.add_argument(
        "--all-upper", action="store_true", help="Todas las letras en mayúsculas"
    )
    parser.add_argument(
        "--all-lower", action="store_true", help="Todas las letras en minúsculas"
    )

    args = parser.parse_args()

    # Validaciones
    if args.uppercase_ratio < 0.0 or args.uppercase_ratio > 1.0:
        print("Error: uppercase-ratio debe estar entre 0.0 y 1.0", file=sys.stderr)
        sys.exit(1)

    if args.min_length is not None and args.max_length is not None:
        if args.min_length > args.max_length:
            print(
                "Error: min-length debe ser menor o igual a max-length", file=sys.stderr
            )
            sys.exit(1)
        if args.min_length < 1:
            print("Error: min-length debe ser al menos 1", file=sys.stderr)
            sys.exit(1)

    if args.length < 1:
        print("Error: length debe ser al menos 1", file=sys.stderr)
        sys.exit(1)

    if args.num_sequences < 1:
        print("Error: num-sequences debe ser al menos 1", file=sys.stderr)
        sys.exit(1)

    # Ajustar proporción de mayúsculas según flags
    uppercase_ratio = args.uppercase_ratio
    if args.all_upper:
        uppercase_ratio = 1.0
    elif args.all_lower:
        uppercase_ratio = 0.0

    # Generar secuencias
    sequences = generate_multiple_sequences(
        args.num_sequences,
        args.length,
        uppercase_ratio,
        args.min_length,
        args.max_length,
        args.seed,
    )

    # Escribir salida
    output_file = sys.stdout
    if args.output:
        output_file = open(args.output, "w")

    try:
        # Escribir secuencias
        for seq in sequences:
            print(seq, file=output_file)

        # Estadísticas al final
        total_length = sum(len(s) for s in sequences)
        avg_length = total_length / len(sequences)
        total_upper = sum(sum(1 for c in s if c.isupper()) for s in sequences)
        actual_ratio = total_upper / total_length if total_length > 0 else 0

        print("\n# Estadísticas:", file=sys.stderr)
        print(f"#   Secuencias generadas: {len(sequences)}", file=sys.stderr)
        print(f"#   Longitud promedio: {avg_length:.1f}", file=sys.stderr)
        print(f"#   Total aminoácidos: {total_length}", file=sys.stderr)
        print(
            f"#   Mayúsculas: {total_upper} ({actual_ratio * 100:.1f}%)",
            file=sys.stderr,
        )

    finally:
        if args.output:
            output_file.close()
            print(f"\n# Archivo guardado: {args.output}", file=sys.stderr)


if __name__ == "__main__":
    main()
