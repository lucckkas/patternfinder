#!/usr/bin/env python3
import argparse
import csv
import re
from pathlib import Path
from typing import List, Tuple

from Bio.PDB import MMCIFParser, PPBuilder
from Bio.Align import PairwiseAligner

C2H2_REGEX = re.compile(r"C.{2,4}C.{10,14}H.{2,6}H")  # PROSITE-like, un poco laxo


def parse_args():
    ap = argparse.ArgumentParser(
        description="Selecciona un subconjunto diverso de proteínas C2H2 zinc finger desde mmCIF."
    )
    ap.add_argument("--cif-dir", required=True, help="Directorio con archivos .cif")
    ap.add_argument(
        "--ids-file",
        required=True,
        help="Archivo de texto con IDs (separados por ; , espacio o salto de línea). Ej: 'pdb_00001a1f; pdb_00001a1g; ...'",
    )
    ap.add_argument(
        "--mode",
        choices=["motif", "chain"],
        default="motif",
        help="Comparar diversidad por regiones de motivo (motif) o por cadena completa (chain).",
    )
    ap.add_argument(
        "--identity-threshold",
        type=float,
        default=0.40,
        help="Máxima identidad permitida entre seleccionadas (0.0–1.0). Default=0.40",
    )
    ap.add_argument(
        "--max",
        type=int,
        default=0,
        help="Máximo de resultados a seleccionar (0 = sin límite).",
    )
    ap.add_argument("--out", default="diverse_selection.csv", help="CSV de salida.")
    return ap.parse_args()


def load_ids(ids_file: Path) -> List[str]:
    raw = ids_file.read_text(encoding="utf-8")
    tokens = re.split(r"[;,\s]+", raw.strip())
    return [t for t in tokens if t]


def guess_pdb_id(token: str) -> str:
    # Obtiene un PDB ID de 4 chars al final del token (e.g., 'pdb_00001a1f' -> '1a1f')
    m = re.search(r"([0-9][A-Za-z0-9]{3})$", token)
    return m.group(1).lower() if m else token.lower()


def find_cif_path(cif_dir: Path, token: str) -> Path:
    # Busca mmCIF por dos nombres comunes:
    #   1) <pdbid>.cif  (1a1f.cif)
    #   2) <token>.cif  (pdb_00001a1f.cif)
    pdbid = guess_pdb_id(token)
    cand1 = cif_dir / f"{pdbid}.cif"
    if cand1.exists():
        return cand1
    cand2 = cif_dir / f"{token}.cif"
    if cand2.exists():
        return cand2
    # No encontrado: devolvemos la primera ruta por diagnóstico
    return cand1


def extract_chain_sequences(cif_path: Path):
    parser = MMCIFParser(QUIET=True)
    structure = parser.get_structure(cif_path.stem, str(cif_path))
    ppb = PPBuilder()
    out = []
    for model in structure:
        for chain in model:
            polypeps = ppb.build_peptides(chain)
            if not polypeps:
                continue
            seq = "".join(str(pp.get_sequence()) for pp in polypeps).upper()
            if not seq:
                continue
            out.append((chain.id, seq))
        break  # solo primer modelo
    return out  # List[(chain_id, seq)]


def find_c2h2_windows(seq: str) -> List[Tuple[int, int, str]]:
    wins = []
    for m in C2H2_REGEX.finditer(seq):
        start, end = m.span()
        wins.append((start, end, seq[start:end]))
    return wins


def signature_for_chain(seq: str, mode: str, wins: List[Tuple[int, int, str]]) -> str:
    if mode == "motif":
        if not wins:
            return ""
        return "".join(w[2] for w in wins)
    return seq  # mode == "chain"


def identity(a: str, b: str) -> float:
    # Global alignment; identidad ~ matches / max(len(a), len(b))
    if not a or not b:
        return 0.0
    aligner = PairwiseAligner()
    aligner.mode = "global"
    # Puntuaciones simples para alinear; no necesitamos matriz de sustitución para identidad
    aligner.match_score = 1.0
    aligner.mismatch_score = 0.0
    aligner.open_gap_score = -1.0
    aligner.extend_gap_score = -0.5
    aln = aligner.align(a, b)[0]

    # Contar matches comparando bloques alineados (sin gaps)
    matches = 0
    Ablocks, Bblocks = aln.aligned
    for (ai, aj), (bi, bj) in zip(Ablocks, Bblocks):
        # mismos tamaños por definición
        for k in range(aj - ai):
            if a[ai + k] == b[bi + k]:
                matches += 1
    denom = max(len(a), len(b))
    return matches / denom if denom else 0.0


def greedy_diverse(items, thr: float, mode: str, max_n: int):
    # items: list of dicts con keys: pdb_id, chain_id, seq, wins, sig
    selected = []
    for cand in sorted(items, key=lambda x: (-(len(x["wins"])), -len(x["sig"]))):
        if not cand["sig"]:
            continue
        ok = True
        for rep in selected:
            if identity(cand["sig"], rep["sig"]) > thr:
                ok = False
                break
        if ok:
            selected.append(cand)
            if max_n and len(selected) >= max_n:
                break
    return selected


def main():
    args = parse_args()
    ids = load_ids(Path(args.ids_file))
    cif_dir = Path(args.cif_dir)

    items = []
    missing = []

    for tok in ids:
        pdb_id = guess_pdb_id(tok)
        cif_path = find_cif_path(cif_dir, tok)
        if not cif_path.exists():
            missing.append((tok, str(cif_path)))
            continue

        try:
            chains = extract_chain_sequences(cif_path)
        except Exception as e:
            print(f"[WARN] No pude parsear {cif_path.name}: {e}")
            continue

        for chain_id, seq in chains:
            wins = find_c2h2_windows(seq)
            if not wins:
                continue
            sig = signature_for_chain(seq, args.mode, wins)
            items.append(
                {
                    "pdb_id": pdb_id,
                    "raw_id": tok,
                    "cif": cif_path.name,
                    "chain_id": chain_id,
                    "seq": seq,
                    "wins": wins,
                    "sig": sig,
                }
            )

    if not items:
        print("No se encontraron cadenas con motivo C2H2 en los CIF provistos.")
        if missing:
            print("Archivos faltantes (ejemplos):", missing[:5])
        return

    selected = greedy_diverse(items, args.identity_threshold, args.mode, args.max)

    # CSV
    out = Path(args.out)
    with out.open("w", newline="", encoding="utf-8") as fh:
        w = csv.writer(fh)
        w.writerow(
            [
                "PDB ID",
                "ID original",
                "Archivo CIF",
                "Cadena",
                "Longitud cadena",
                "Cantidad motivos C2H2",
                "Spans motivos (start-end)",
                "Firma usada para diversidad (acortada)",
                "Modo firma",
            ]
        )
        for it in selected:
            spans = ";".join([f"{s}-{e}" for (s, e, _) in it["wins"]])
            sig_show = (
                it["sig"]
                if len(it["sig"]) <= 120
                else it["sig"][:60] + "…" + it["sig"][-60:]
            )
            w.writerow(
                [
                    it["pdb_id"],
                    it["raw_id"],
                    it["cif"],
                    it["chain_id"],
                    len(it["seq"]),
                    len(it["wins"]),
                    spans,
                    sig_show,
                    args.mode,
                ]
            )

    print(f"Seleccionadas {len(selected)} / {len(items)} cadenas con C2H2. CSV: {out}")
    if missing:
        print(
            f"[INFO] {len(missing)} IDs no encontraron CIF en {cif_dir} (revisa nombres/descargas)."
        )


if __name__ == "__main__":
    main()
