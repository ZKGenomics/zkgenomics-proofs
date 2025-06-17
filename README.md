# ZKGenomics Proofs

A Go package for generating and verifying zero-knowledge proofs of genomic traits from VCF files.

## Overview

This package provides cryptographic proof generation capabilities for genomic data, allowing users to prove they possess certain genetic traits without revealing the underlying genomic information.

## Supported Proof Types

- **Chromosome Proof**: Proves presence of specific chromosomes in genomic data
- **BRCA1 Proof**: Proves presence/absence of BRCA1 pathogenic variants  
- **HERC2 Proof**: Proves HERC2 gene variants related to eye color
- **Eye Color Proof**: Proves eye color traits based on genetic markers

## Installation

```bash
go get github.com/zkgenomics/zkgenomics-proofs
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	
	zkgenomics "github.com/zkgenomics/zkgenomics-proofs"
)

func main() {
	// Create a proof generator
	generator := zkgenomics.NewProofGenerator()
	
	// Generate a chromosome proof
	err := generator.GenerateProof(
		zkgenomics.ChromosomeProofType,
		"path/to/sample.vcf",
		"", // Empty string for new proving key
		"output/proof",
	)
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}
	
	// Verify the proof
	valid, err := generator.VerifyProof(
		zkgenomics.ChromosomeProofType,
		"output/proof.vk",
		"output/proof",
	)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	
	fmt.Printf("Proof is valid: %v\n", valid)
}
```

## Trait Data

The package includes trait definitions in `traits.json` with genomic positions for various genetic markers including:

- BRCA1/BRCA2 cancer susceptibility variants
- APOE Alzheimer's risk alleles  
- Drug metabolism markers (CYP2C19)
- Ancestry markers
- And more...

## API Reference

### ProofGenerator

The main interface for generating and verifying proofs.

#### Methods

- `GenerateProof(proofType ProofType, vcfPath, provingKeyPath, outputPath string) error`
- `VerifyProof(proofType ProofType, verifyingKeyPath, proofPath string) (bool, error)`
- `GetSupportedProofTypes() []ProofType`

### ProofType Constants

- `ChromosomeProofType`
- `EyeColorProofType` 
- `BRCA1ProofType`
- `HERC2ProofType`

## Dependencies

- [gnark](https://github.com/consensys/gnark) - Zero-knowledge proof framework
- [vcfgo](https://github.com/brentp/vcfgo) - VCF file parsing

## License

This project is part of the ZKGenomics ecosystem.