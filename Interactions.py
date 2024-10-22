from Bio.PDB import PDBParser, NeighborSearch, PPBuilder
from Bio.SeqUtils import seq1

# Ruta al archivo PDB de la proteína
pdb_file = "1a2b.pdb"  # Reemplaza con la ruta a tu archivo PDB

# Umbral de distancia para considerar una interacción (en angstroms)
distance_threshold = 4.0

# Crear un parser para leer el archivo PDB
parser = PDBParser(QUIET=True)
structure = parser.get_structure("protein", pdb_file)

# Obtener todos los residuos no estándar (posibles ligandos)
ligand_residues = []

for model in structure:
    for chain in model:
        for residue in chain:
            hetfield, resseq, icode = residue.id
            if hetfield != " " and residue.get_resname() != "HOH":
                ligand_residues.append(residue)

if not ligand_residues:
    print("No se encontraron ligandos en la estructura.")
    exit()

# Si hay múltiples ligandos, opcionalmente puedes seleccionar uno
print("Ligandos encontrados en la estructura:")
ligand_codes = set()
for residue in ligand_residues:
    ligand_codes.add(residue.get_resname())

ligand_codes = list(ligand_codes)
for idx, lig_code in enumerate(ligand_codes):
    print(f"{idx + 1}. {lig_code}")

# Seleccionar el ligando (si hay más de uno)
if len(ligand_codes) > 1:
    lig_choice = int(input("Seleccione el número del ligando que desea usar: ")) - 1
    ligand_resname = ligand_codes[lig_choice]
else:
    ligand_resname = ligand_codes[0]
    print(f"Usando el ligando: {ligand_resname}")

# Obtener los átomos del ligando seleccionado
ligand_atoms = []

for residue in ligand_residues:
    if residue.get_resname() == ligand_resname:
        ligand_atoms.extend(residue.get_atoms())

# Obtener todos los átomos de aminoácidos
protein_atoms = []
residue_list = []

for model in structure:
    for chain in model:
        for residue in chain:
            hetfield, resseq, icode = residue.id
            if hetfield == " ":
                protein_atoms.extend(residue.get_atoms())
                residue_list.append(residue)

# Crear una búsqueda de vecinos
ns = NeighborSearch(protein_atoms)

# Identificar los aminoácidos que interactúan con el ligando
interacting_residues = set()

for atom in ligand_atoms:
    neighbors = ns.search(atom.get_coord(), distance_threshold, level="R")
    for neighbor in neighbors:
        interacting_residues.add(neighbor)

# Generar la secuencia de aminoácidos, marcando los que interactúan
ppb = PPBuilder()
sequence = ""

for pp in ppb.build_peptides(structure):
    for residue in pp:
        resname = residue.get_resname()
        one_letter = seq1(resname)
        # Verificar si el residuo interactúa con el ligando
        if residue in interacting_residues:
            sequence += one_letter.upper()
        else:
            sequence += one_letter.lower()

print("\nSecuencia anotada:")
print(sequence)
