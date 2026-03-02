# [`coder`](https://owncloud.gwdg.de/index.php/s/82p093Aoi3AhW5x): encode peptides or decode DNA codons

The program coder converts biological sequences between amino acid
residues and DNA using a non-degenerate genetic code, where each
residue maps to a single unique codon. The one-to-one correspondence
enables unambiguous round-trip translation that can be used in
bioinformatics for reversible encoding that ignores silent mutations.

This program uses the following table of codons:

| Residue | Codon |
| ------- | ----- |
| `A`     | `GCT` |
| `R`     | `CGT` |
| `N`     | `AAT` |
| `D`     | `GAT` |
| `C`     | `TGT` |
| `Q`     | `CAA` |
| `E`     | `GAA` |
| `G`     | `GGT` |
| `H`     | `CAT` |
| `I`     | `ATT` |
| `L`     | `CTG` |
| `K`     | `AAA` |
| `M`     | `ATG` |
| `F`     | `TTT` |
| `P`     | `CCT` |
| `S`     | `TCT` |
| `T`     | `ACT` |
| `W`     | `TGG` |
| `Y`     | `TAT` |
| `V`     | `GTT` |
| `*`     | `TAA` |

which is based on https://www.ncbi.nlm.nih.gov/Taxonomy/Utils/wprintgc.cgi#SG11.


Installation:

	git clone https://github.com/IvanTsers/coder
	cd coder
	make

To encode peptides as non-degenerate DNA codons:

   coder -e peptide.faa

To decode DNA codons as peptides using the same genetic code:

   coder -d codons.fna