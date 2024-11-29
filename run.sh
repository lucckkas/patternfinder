#!/bin/bash

# Read the Segmentos.txt file
input_file="Segmentos.txt"
lines=($(cat $input_file))

# Iterate over the lines in pairs
for ((i=0; i<${#lines[@]}; i+=2)); do
    seq1=${lines[$i]}
    seq2=${lines[$i+1]}
    
    # Execute Discovery.java with the sequences
    java Discovery "$seq1" "$seq2"
done
