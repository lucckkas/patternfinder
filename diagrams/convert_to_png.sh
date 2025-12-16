#!/bin/bash
# filepath: /home/luckas/Repositorios/Memoria/diagrams/convert_to_png.sh

# Script para convertir todos los diagramas LaTeX a PNG
# Uso: ./convert_to_png.sh [dpi]
# Ejemplo: ./convert_to_png.sh 300

DPI=${1:-300}  # DPI por defecto: 300
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Convirtiendo diagramas a PNG (${DPI} DPI) ==="

# Verificar que pdftoppm esté instalado
if ! command -v pdftoppm &> /dev/null; then
    echo "Error: pdftoppm no está instalado."
    echo "Instálalo con: sudo apt install poppler-utils"
    exit 1
fi

# Crear carpeta para las imágenes
mkdir -p "$SCRIPT_DIR/png"

# Primero compilar todos los .tex a PDF
for texfile in "$SCRIPT_DIR"/*.tex; do
    if [ -f "$texfile" ]; then
        filename=$(basename "$texfile" .tex)
        echo "Compilando $filename.tex..."
        pdflatex -interaction=nonstopmode -output-directory="$SCRIPT_DIR" "$texfile" > /dev/null 2>&1
        
        # Limpiar archivos auxiliares
        rm -f "$SCRIPT_DIR/$filename.aux" "$SCRIPT_DIR/$filename.log" "$SCRIPT_DIR/$filename.fls" "$SCRIPT_DIR/$filename.fdb_latexmk"
    fi
done

# Convertir todos los PDF a PNG
for pdffile in "$SCRIPT_DIR"/*.pdf; do
    if [ -f "$pdffile" ]; then
        filename=$(basename "$pdffile" .pdf)
        echo "Convirtiendo $filename.pdf a PNG..."
        pdftoppm -png -r "$DPI" -singlefile "$pdffile" "$SCRIPT_DIR/png/$filename"
        
        if [ -f "$SCRIPT_DIR/png/$filename.png" ]; then
            echo "  ✓ Creado: png/$filename.png"
        fi
    fi
done

echo ""
echo "=== Conversión completada ==="
echo "Las imágenes PNG están en: $SCRIPT_DIR/png/"
ls -lh "$SCRIPT_DIR/png/"*.png 2>/dev/null