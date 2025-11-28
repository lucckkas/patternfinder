#!/usr/bin/env python3
import argparse
import re
from pathlib import Path
from typing import List, Tuple

from Bio.PDB import PDBList  # pip install biopython

ID_REGEX = re.compile(r"([0-9][A-Za-z0-9]{3})$")


def parse_args():
    ap = argparse.ArgumentParser(
        description="Descarga mmCIFs de RCSB PDB a partir de una lista (con tokens tipo pdb_0000XXXX)."
    )
    ap.add_argument(
        "--ids-file",
        required=True,
        help="Archivo con IDs (separados por ; , espacios o saltos de línea).",
    )
    ap.add_argument(
        "--outdir", default="cifs", help="Directorio de salida (default: ./cifs)"
    )
    ap.add_argument(
        "--overwrite",
        action="store_true",
        help="Forzar redescarga si el .cif ya existe.",
    )
    return ap.parse_args()


def load_ids(path: Path) -> List[str]:
    raw = path.read_text(encoding="utf-8")
    toks = re.split(r"[;,\s]+", raw.strip())
    return [t for t in toks if t]


def to_pdbid(token: str) -> str:
    m = ID_REGEX.search(token.strip())
    if not m:
        return ""
    return m.group(1).lower()


def main():
    args = parse_args()
    outdir = Path(args.outdir)
    outdir.mkdir(parents=True, exist_ok=True)

    tokens = load_ids(Path(args.ids_file))
    ids = []
    for t in tokens:
        pid = to_pdbid(t)
        if pid:
            ids.append(pid)
    ids = sorted(set(ids))

    if not ids:
        print("No se encontraron PDB IDs válidos en el archivo.")
        return

    print(f"Descargando {len(ids)} entradas a: {outdir.resolve()}")
    pdb_list = PDBList()

    ok, skipped, fail = [], [], []
    for i, pid in enumerate(ids, 1):
        dest = outdir / f"{pid}.cif"
        if dest.exists() and not args.overwrite:
            skipped.append(pid)
            continue
        try:
            print(f"[{i}/{len(ids)}] {pid} …", flush=True)
            # Biopython guarda como <outdir>/<pid>.cif (minúsculas)
            fp = pdb_list.retrieve_pdb_file(pid, pdir=str(outdir), file_format="mmCif")
            # Según versión, ya se llama pid.cif; si no, normalizamos nombre
            got = Path(fp)
            if got.name != dest.name:
                got.rename(dest)
            ok.append(pid)
        except Exception as e:
            print(f"   ERROR {pid}: {e}")
            fail.append((pid, str(e)))

    print("\nResumen:")
    print(f"  OK:      {len(ok)}")
    print(f"  Skipped: {len(skipped)} (ya existían)")
    print(f"  Fails:   {len(fail)}")
    if fail:
        print("  Algunos fallos (primeros 10):")
        for pid, msg in fail[:10]:
            print(f"   - {pid}: {msg}")


if __name__ == "__main__":
    main()
