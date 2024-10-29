from Bio.PDB import PDBParser, NeighborSearch, PPBuilder
from Bio.SeqUtils import seq1

pdb_file = "5bxq.pdb"

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

# Recopilar información en una sola pasada
for model in structure:
    for chain in model:
        peptides = pp_builder.build_peptides(chain)
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

# Seleccionar el ligando (si hay más de uno)
if len(ligand_codes) > 1:
    lig_choice = int(input("Seleccione el número del ligando que desea usar: ")) - 1
    ligand_resname = ligand_codes[lig_choice]
else:
    ligand_resname = ligand_codes[0]
    print(f"Usando el ligando: {ligand_resname}")

# Obtener los átomos del ligando seleccionado
ligand_atoms = [
    atom
    for residue in ligand_residues
    if residue.get_resname() == ligand_resname
    for atom in residue.get_atoms()
]

# Crear una búsqueda de vecinos
ns = NeighborSearch(protein_atoms)

# Identificar los aminoácidos que interactúan con el ligando
interacting_residues = set()
for atom in ligand_atoms:
    neighbors = ns.search(atom.get_coord(), distance_threshold, level="R")
    interacting_residues.update(neighbors)

# Generar la secuencia y encontrar los índices
sequence = ""
first_interacting_index = None
last_interacting_index = None

for idx, residue in enumerate(residues_list):
    resname = residue.get_resname()
    one_letter = seq1(resname)
    if residue in interacting_residues:
        one_letter = one_letter.upper()
        if first_interacting_index is None:
            first_interacting_index = idx
        last_interacting_index = idx
    else:
        one_letter = one_letter.lower()
    sequence += one_letter

# Extraer el segmento interactuante
if first_interacting_index is not None:
    sequence_segment = sequence[first_interacting_index : last_interacting_index + 1]
    print("\nSecuencia de aminoácidos de la proteína:")
    print(sequence)
    print("\nSegmento que interactúa con el ligando:")
    print(sequence_segment)
else:
    print("No se encontraron residuos que interactúen con el ligando.")
