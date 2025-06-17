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

### Basic Example

```go
package main

import (
	"fmt"
	"log"
	
	zkgenomics "github.com/zkgenomics/zkgenomics-proofs"
	"github.com/zkgenomics/zkgenomics-proofs/proofs"
)

func main() {
	// Create a proof generator
	generator := zkgenomics.NewProofGenerator()
	
	// Generate a chromosome proof
	proofData, err := generator.GenerateProof(
		zkgenomics.ChromosomeProofType,
		"path/to/sample.vcf",
		"", // Empty string for new proving key
		"output/proof",
	)
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}
	
	fmt.Printf("Proof generation result: %s\n", proofData.Result.String())
	if proofData.Result == zkgenomics.ProofSuccess {
		fmt.Printf("Proof size: %d bytes\n", len(proofData.Proof))
		fmt.Printf("Verifying key size: %d bytes\n", len(proofData.VerifyingKey))
	}
	
	// Verify the proof
	result, err := generator.VerifyProof(
		zkgenomics.ChromosomeProofType,
		"",  // Using embedded verifying key
		"",  // Using embedded proof data
	)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}
	
	fmt.Printf("Verification result: %s\n", result.Result.String())
}
```

### Dynamic Proofs for Custom Variants

```go
// Create a dynamic proof for a specific genomic position
eyeColorProof := proofs.NewDynamicProof(28356859, "G", "A") // rs12913832

// Generate proof for eye color variant
proofData, err := eyeColorProof.Generate("sample.vcf", "", "")
if err != nil {
	log.Fatalf("Failed to generate proof: %v", err)
}

// The proof will succeed only if the person has the G->A variant at position 28356859
if proofData.Result == zkgenomics.ProofSuccess {
	fmt.Println("Person has the blue eyes variant!")
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

- `GenerateProof(proofType ProofType, vcfPath, provingKeyPath, outputPath string) (*ProofData, error)`
- `VerifyProof(proofType ProofType, verifyingKeyPath, proofPath string) (*VerificationResult, error)`
- `GetSupportedProofTypes() []ProofType`

### ProofData Structure

Contains all necessary data for proof verification:

```go
type ProofData struct {
    Proof         []byte      `json:"proof"`         // ZK-SNARK proof bytes
    VerifyingKey  []byte      `json:"verifying_key"` // Verification key bytes
    PublicWitness []byte      `json:"public_witness"`// Public inputs bytes
    Result        ProofResult `json:"result"`        // success/fail/unknown
}
```

### VerificationResult Structure

Contains the result of proof verification:

```go
type VerificationResult struct {
    Result ProofResult `json:"result"` // success/fail/unknown
    Error  error       `json:"error"`  // Optional error details
}
```

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