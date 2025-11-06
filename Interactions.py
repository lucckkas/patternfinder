#!/usr/bin/env python3
import argparse
import json
import sys
from pathlib import Path

from Bio.PDB import NeighborSearch, MMCIFParser
from Bio.PDB.Polypeptide import is_aa
from Bio.SeqUtils import seq1


def parse_args():
    p = argparse.ArgumentParser(
        description="Detecta segmentos de AA que interactúan con ligandos en una estructura mmCIF."
    )
    p.add_argument("cif_file", help="Ruta al archivo .cif (mmCIF)")
    p.add_argument(
        "-d",
        "--distance",
        type=float,
        default=4.0,
        help="Umbral de distancia (Å) para considerar interacción (default: 4.0)",
    )
    p.add_argument(
        "-o",
        "--output",
        default="Segmentos.json",
        help='Archivo JSON de salida (default: "Segmentos.json")',
    )
    return p.parse_args()


def main():
    args = parse_args()
    cif_path = Path(args.cif_file)
    if not cif_path.exists():
        print(f"Archivo no encontrado: {cif_path}", file=sys.stderr)
        sys.exit(1)

    # Parser mmCIF (puedes cambiar a FastMMCIFParser si lo prefieres)
    parser = MMCIFParser(QUIET=True)
    structure = parser.get_structure("protein", str(cif_path))

    protein_name = cif_path.stem

    ligand_residues = []
    protein_atoms = []
    residues_list = []

    # Procesar solo el primer modelo
    model = structure[0]
    for chain in model:
        for residue in chain:
            hetfield, _, _ = residue.id
            resname = residue.get_resname()

            # AA estándar -> proteína
            if is_aa(residue, standard=True):
                protein_atoms.extend(residue.get_atoms())
                residues_list.append(residue)
                continue

            # Excluir agua
            if hetfield == "W" or resname in {"HOH", "WAT"}:
                continue

            # Todo lo demás (HETATM) -> ligando
            ligand_residues.append(residue)

    if not ligand_residues:
        print("No se encontraron ligandos en la estructura.")
        sys.exit(0)

    # Códigos únicos de ligandos
    ligand_codes = sorted({r.get_resname() for r in ligand_residues})
    print("Ligandos encontrados en la estructura:")
    for i, code in enumerate(ligand_codes, 1):
        print(f"{i}. {code}")

    # Generar secuencia una sola vez (AA estándar)
    seq_letters = []
    for residue in residues_list:
        resname = residue.get_resname()
        one_letter = seq1(resname)  # seguro porque filtramos standard=True
        seq_letters.append(one_letter.lower())
    sequence = "".join(seq_letters)

    print("\nSecuencia de aminoácidos de la proteína:")
    print(sequence)

    # Cargar/crear JSON
    output_file = Path(args.output)
    try:
        with output_file.open("r", encoding="utf-8") as fh:
            data = json.load(fh)
    except FileNotFoundError:
        data = {}

    if protein_name not in data:
        data[protein_name] = {}

    # Crear NeighborSearch solo una vez (optimización)
    ns = NeighborSearch(protein_atoms)

    # Mapa para comparación robusta (cadena, residue.id)
    def res_key(r):
        return (r.get_parent().id, r.id)

    for lig_residue in ligand_residues:
        lig_code = lig_residue.get_resname()
        chain_id = lig_residue.get_parent().id
        resseq = lig_residue.id[1]
        print(
            f"\nProcesando ligando: {lig_code} en residuo {resseq} (cadena {chain_id})"
        )

        lig_id = f"{lig_code}_{chain_id}_{resseq}"
        if lig_id not in data[protein_name]:
            data[protein_name][lig_id] = []

        ligand_atoms = list(lig_residue.get_atoms())

        # Buscar residuos de proteína cercanos al ligando
        interacting_keys = set()
        for atom in ligand_atoms:
            neighbors = ns.search(atom.get_coord(), args.distance, level="R")
            for r in neighbors:
                interacting_keys.add(res_key(r))

        # Secuencia resaltada (mayúsculas = interactúan)
        highlighted = []
        for r in residues_list:
            letter = seq1(r.get_resname())
            if res_key(r) in interacting_keys:
                highlighted.append(letter.upper())
            else:
                highlighted.append(letter.lower())
        highlighted_sequence = "".join(highlighted)

        # Posiciones interactuantes y segmento mínimo que las cubre
        interacting_positions = [
            i for i, c in enumerate(highlighted_sequence) if c.isupper()
        ]
        if interacting_positions:
            i_min, i_max = min(interacting_positions), max(interacting_positions)
            segment = highlighted_sequence[i_min : i_max + 1]
            if segment not in data[protein_name][lig_id]:
                data[protein_name][lig_id].append(segment)
                print(f"Segmento interactuante para {lig_id}: {segment}")
        else:
            print(f"No se encontraron residuos interactuantes para {lig_id}")

    with output_file.open("w", encoding="utf-8") as fh:
        json.dump(data, fh, indent=4, ensure_ascii=False)

    print(f"\nDatos guardados en {output_file}")


if __name__ == "__main__":
    main()
