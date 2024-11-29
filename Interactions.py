import json
import sys
import os
from Bio.PDB import PDBParser, NeighborSearch
from Bio.SeqUtils import seq1

pdb_file = sys.argv[1]

# Umbral de distancia para considerar una interacción (en angstroms)
distance_threshold = 4.0

# Archivo de salida
output_file = "Segmentos.json"

# Crear un parser para leer el archivo PDB
parser = PDBParser(QUIET=True)
structure = parser.get_structure("protein", pdb_file)

# Obtener el nombre de la proteína
protein_name = os.path.splitext(os.path.basename(pdb_file))[0]

# Inicializar variables
ligand_residues = []
protein_atoms = []
residues_list = []

# Procesar solo el primer modelo
model = structure[0]
for chain in model:
    for residue in chain:
        hetfield, resseq, icode = residue.id
        resname = residue.get_resname()
        if hetfield != " " and resname != "HOH":
            ligand_residues.append(residue)
        elif hetfield == " ":
            protein_atoms.extend(residue.get_atoms())
            residues_list.append(residue)

if not ligand_residues:
    print("No se encontraron ligandos en la estructura.")
    exit()

# Obtener códigos únicos de ligandos
ligand_codes = list(set(residue.get_resname() for residue in ligand_residues))
print("Ligandos encontrados en la estructura:")
for idx, lig_code in enumerate(ligand_codes):
    print(f"{idx + 1}. {lig_code}")

# Generar la secuencia de la proteína una sola vez
sequence = ""
for residue in residues_list:
    resname = residue.get_resname()
    one_letter = seq1(resname)
    sequence += one_letter.lower()

print("\nSecuencia de aminoácidos de la proteína:")
print(sequence)

# Cargar datos previos del archivo JSON
try:
    with open(output_file, "r") as file:
        data = json.load(file)
except FileNotFoundError:
    data = {}

# Asegurarse de que la proteína actual tenga una entrada
if protein_name not in data:
    data[protein_name] = {}

# Procesar cada residuo de ligando individualmente
for lig_residue in ligand_residues:
    lig_code = lig_residue.get_resname()
    print(
        f"\nProcesando ligando: {lig_code} en residuo número {lig_residue.id[1]} (cadena {lig_residue.get_parent().id})"
    )

    # Crear una clave única para cada ligando individual
    lig_id = f"{lig_code}_{lig_residue.get_parent().id}_{lig_residue.id[1]}"

    if lig_id not in data[protein_name]:
        data[protein_name][lig_id] = []

    # Obtener los átomos del ligando actual
    ligand_atoms = list(lig_residue.get_atoms())

    # Crear una búsqueda de vecinos
    ns = NeighborSearch(protein_atoms)

    # Identificar los aminoácidos que interactúan con el ligando
    interacting_residues = set()
    for atom in ligand_atoms:
        neighbors = ns.search(atom.get_coord(), distance_threshold, level="R")
        interacting_residues.update(neighbors)

    # Generar la secuencia resaltando los residuos que interactúan
    highlighted_sequence = ""
    for residue in residues_list:
        resname = residue.get_resname()
        one_letter = seq1(resname)
        if residue in interacting_residues:
            one_letter = one_letter.upper()
        else:
            one_letter = one_letter.lower()
        highlighted_sequence += one_letter

    # Obtener las posiciones de los residuos que interactúan
    interacting_positions = [
        i for i, c in enumerate(highlighted_sequence) if c.isupper()
    ]

    if interacting_positions:
        # Encontrar el segmento desde el primer hasta el último residuo interactuante
        i_min = min(interacting_positions)
        i_max = max(interacting_positions)
        interacting_segment = highlighted_sequence[i_min : i_max + 1]
        if interacting_segment not in data[protein_name][lig_id]:
            data[protein_name][lig_id].append(interacting_segment)
            print(f"Segmento interactuante para {lig_id}: {interacting_segment}")
    else:
        print(f"No se encontraron residuos interactuantes para {lig_id}")

# Guardar los datos actualizados en el archivo JSON
with open(output_file, "w") as file:
    json.dump(data, file, indent=4)

print(f"\nDatos guardados en {output_file}")
