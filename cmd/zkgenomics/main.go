package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/zkgenomics/zkgenomics-proofs"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	
	switch command {
	case "generate":
		handleGenerate()
	case "verify":
		handleVerify()
	case "list":
		handleList()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("zkgenomics - Zero-Knowledge Genomics Proof Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  zkgenomics generate <proof-type> <vcf-path> [proving-key] [output]")
	fmt.Println("  zkgenomics verify <proof-type> <verifying-key> <proof-path>")
	fmt.Println("  zkgenomics list")
	fmt.Println()
	fmt.Println("Proof Types:")
	fmt.Println("  chromosome  - Prove chromosome presence")
	fmt.Println("  eye_color   - Prove eye color trait")
	fmt.Println("  brca1       - Prove BRCA1 variant")
	fmt.Println("  herc2       - Prove HERC2 variant")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  zkgenomics generate eye_color sample.vcf")
	fmt.Println("  zkgenomics verify eye_color verifying.key proof.data")
	fmt.Println("  zkgenomics list")
}

func handleGenerate() {
	if len(os.Args) < 4 {
		fmt.Println("Error: generate requires at least proof-type and vcf-path")
		printUsage()
		os.Exit(1)
	}

	proofType := zkgenomics.ProofType(os.Args[2])
	vcfPath := os.Args[3]
	
	var provingKeyPath, outputPath string
	if len(os.Args) > 4 {
		provingKeyPath = os.Args[4]
	}
	if len(os.Args) > 5 {
		outputPath = os.Args[5]
	} else {
		outputPath = fmt.Sprintf("%s_proof.json", proofType)
	}

	generator := zkgenomics.NewProofGenerator()
	
	fmt.Printf("Generating %s proof from %s...\n", proofType, vcfPath)
	
	proofData, err := generator.GenerateProof(proofType, vcfPath, provingKeyPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate proof: %v", err)
	}

	fmt.Printf("Proof generation result: %s\n", proofData.Result.String())
	
	if proofData.Result == zkgenomics.ProofSuccess {
		// Save proof data to JSON file
		jsonData, err := json.MarshalIndent(proofData, "", "  ")
		if err != nil {
			log.Fatalf("Failed to serialize proof data: %v", err)
		}
		
		err = os.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write proof data to file: %v", err)
		}
		
		fmt.Printf("✅ Proof successfully generated and saved to: %s\n", outputPath)
		fmt.Printf("Proof size: %d bytes\n", len(proofData.Proof))
		fmt.Printf("Verifying key size: %d bytes\n", len(proofData.VerifyingKey))
		fmt.Printf("Public witness size: %d bytes\n", len(proofData.PublicWitness))
	} else {
		fmt.Printf("❌ Proof generation failed\n")
		os.Exit(1)
	}
}

func handleVerify() {
	if len(os.Args) < 5 {
		fmt.Println("Error: verify requires proof-type, verifying-key, and proof-path")
		printUsage()
		os.Exit(1)
	}

	proofType := zkgenomics.ProofType(os.Args[2])
	verifyingKeyPath := os.Args[3]
	proofPath := os.Args[4]

	generator := zkgenomics.NewProofGenerator()
	
	fmt.Printf("Verifying %s proof...\n", proofType)
	
	result, err := generator.VerifyProof(proofType, verifyingKeyPath, proofPath)
	if err != nil {
		log.Fatalf("Failed to verify proof: %v", err)
	}

	fmt.Printf("Verification result: %s\n", result.Result.String())
	
	if result.Result == zkgenomics.ProofSuccess {
		fmt.Println("✅ Proof verification succeeded!")
	} else {
		fmt.Println("❌ Proof verification failed!")
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
		}
		os.Exit(1)
	}
}

func handleList() {
	generator := zkgenomics.NewProofGenerator()
	supportedTypes := generator.GetSupportedProofTypes()
	
	fmt.Println("Supported proof types:")
	for _, proofType := range supportedTypes {
		fmt.Printf("  - %s\n", proofType)
	}
}