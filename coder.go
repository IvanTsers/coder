package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/ivantsers/fastautils"
	"log"
	"os"
)

var codons [256][3]byte
var nuclCodes [256]byte
var residues [64]byte

func init() {
	for i := range codons {
		codons[i] = [3]byte{'N', 'N', 'N'}
	}

	codons['A'] = [3]byte{'G', 'C', 'T'}
	codons['R'] = [3]byte{'C', 'G', 'T'}
	codons['N'] = [3]byte{'A', 'A', 'T'}
	codons['D'] = [3]byte{'G', 'A', 'T'}
	codons['C'] = [3]byte{'T', 'G', 'T'}
	codons['Q'] = [3]byte{'C', 'A', 'A'}
	codons['E'] = [3]byte{'G', 'A', 'A'}
	codons['G'] = [3]byte{'G', 'G', 'T'}
	codons['H'] = [3]byte{'C', 'A', 'T'}
	codons['I'] = [3]byte{'A', 'T', 'T'}
	codons['L'] = [3]byte{'C', 'T', 'G'}
	codons['K'] = [3]byte{'A', 'A', 'A'}
	codons['M'] = [3]byte{'A', 'T', 'G'}
	codons['F'] = [3]byte{'T', 'T', 'T'}
	codons['P'] = [3]byte{'C', 'C', 'T'}
	codons['S'] = [3]byte{'T', 'C', 'T'}
	codons['T'] = [3]byte{'A', 'C', 'T'}
	codons['W'] = [3]byte{'T', 'G', 'G'}
	codons['Y'] = [3]byte{'T', 'A', 'T'}
	codons['V'] = [3]byte{'G', 'T', 'T'}
	codons['*'] = [3]byte{'T', 'A', 'A'}
	for i := range nuclCodes {
		nuclCodes[i] = 0xFF
	}

	nuclCodes['A'] = 0
	nuclCodes['C'] = 1
	nuclCodes['G'] = 2
	nuclCodes['T'] = 3
	for i := range residues {
		residues[i] = 'X'
	}

	for _, r := range []byte("ARNDCQEGHILKMFPSTWYV*") {
		c := codons[r]
		c0 := nuclCodes[c[0]]
		c1 := nuclCodes[c[1]]
		c2 := nuclCodes[c[2]]

		idx := (c0 << 4) | (c1 << 2) | c2
		residues[idx] = r
	}
}

func main() {
	optE := flag.Bool("e", false, "encode peptides with DNA codons")
	optD := flag.Bool("d", false, "decode DNA codons as peptides")
	optX := flag.Bool("x", false, "print dot '.' instead of 'X'")
	u := "coder [option]..."
	p := "Encode or decode a biological sequence " +
		"using a non-degenerate table of DNA codons."
	e := "coder -e peptide.faa\n         coder -d codons.fna"
	clio.Usage(u, p, e)
	log.SetPrefix("coder: ")
	log.SetFlags(log.Lmsgprefix)
	flag.Parse()
	args := flag.Args()
	edBoth := *optE && *optD
	edNone := *optE == false && *optD == false
	if edNone || edBoth {
		log.Fatal("please use either encode or decode mode")
	}

	if len(args) != 1 {
		log.Fatal("please specify the input file")
	}
	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("couldn't open file %v: %v", path, err)
	}
	sequences := fastautils.ReadAll(file)
	var result []*fasta.Sequence

	if *optE {
		for _, seq := range sequences {
			h := seq.Header()
			d := seq.Data()
			encodedData := make([]byte, 3*len(d))
			j := 0
			for _, r := range d {
				c := codons[r]
				encodedData[j] = c[0]
				encodedData[j+1] = c[1]
				encodedData[j+2] = c[2]
				j += 3
			}
			encodedSeq := fasta.NewSequence(h, encodedData)
			result = append(result, encodedSeq)
		}
	}

	if *optD {
		for _, seq := range sequences {
			h := seq.Header()
			d := seq.Data()
			lenD := len(d)
			bestFrameScore := 0
			bestFrameData := make([]byte, lenD/3)
			bestFrameName := 0
			for f := 0; f < 3; f++ {
				frameLen := (lenD - f) / 3
				frameData := make([]byte, frameLen)
				j := 0
				frameScore := 0
				for i := f; i+2 < lenD; i += 3 {
					c0 := nuclCodes[d[i]]
					c1 := nuclCodes[d[i+1]]
					c2 := nuclCodes[d[i+2]]
					if c0 == 0xFF || c1 == 0xFF || c2 == 0xFF {
						frameData[j] = 'X'
						j++
						continue
					}
					idx := (c0 << 4) | (c1 << 2) | c2
					frameData[j] = residues[idx]
					j++
				}
				for _, r := range frameData {
					if r != 'X' && r != '*' {
						frameScore++
					}
				}
				if frameScore > bestFrameScore {
					bestFrameScore = frameScore
					bestFrameData = frameData
					bestFrameName = f
				}
			}
			newH := fmt.Sprintf("%s frame=%d", h, bestFrameName)
			decodedSeq := fasta.NewSequence(newH, bestFrameData)
			result = append(result, decodedSeq)
		}
	}
	if *optX {
		for _, res := range result {
			data := res.Data()
			for i, b := range data {
				if b == 'X' {
					data[i] = '.'
				}
			}
		}
	}

	for _, res := range result {
		fmt.Println(res)
	}
}
