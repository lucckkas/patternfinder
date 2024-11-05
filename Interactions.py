import re
from Bio.PDB import PDBParser, NeighborSearch, PPBuilder
from Bio.SeqUtils import seq1

pdb_file = "1znf.pdb"

# Umbral de distancia para considerar una interacción (en angstroms)
distance_threshold = 4.0

# Crear un parser para leer el archivo PDB
parser = PDBParser(QUIET=True)
structure = parser.get_structure("protein", pdb_file)

# Inicializar variables
ligand_residues = []
protein_atoms = []
residues_list = []
pp_builder = PPBuilder()

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

# Procesar cada ligando individualmente
for lig_code in ligand_codes:
    print(f"\nProcesando ligando: {lig_code}")

    # Obtener los átomos del ligando actual
    ligand_atoms = [
        atom
        for residue in ligand_residues
        if residue.get_resname() == lig_code
        for atom in residue.get_atoms()
    ]

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

    # Extraer el segmento interactuante
    match = re.search(r"[A-Z].*[A-Z]", highlighted_sequence)
    if match:
        interacting_segment = match.group()
    else:
        interacting_segment = ""

    print("Segmento que interactúa con el ligando:")
    print(interacting_segment)
