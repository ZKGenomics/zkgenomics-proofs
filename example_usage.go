package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zkgenomics/zkgenomics-proofs/proofs"
)

func main() {
	// Example 1: Create a DynamicProof for eye color (rs12913832)
	fmt.Println("=== Eye Color Proof Example ===")
	eyeColorProof := proofs.NewDynamicProof(28356859, "G", "A")
	
	// Attempt to generate proof (will fail since we don't have a real VCF file)
	proofData, err := eyeColorProof.Generate("sample.vcf", "proving.key", "output.proof")
	if err != nil {
		fmt.Printf("Proof generation failed (expected): %v\n", err)
	}
	
	fmt.Printf("Proof Result: %s\n", proofData.Result.String())
	if proofData.Result == proofs.ProofFail {
		fmt.Println("✓ Correctly failed when variant not found")
	}
	
	// Example 2: Simulate successful proof verification
	fmt.Println("\n=== Proof Verification Example ===")
	verifyResult, err := eyeColorProof.Verify("verifying.key", "proof.data")
	if err != nil {
		log.Printf("Verification error: %v", err)
	} else {
		fmt.Printf("Verification Result: %s\n", verifyResult.Result.String())
		if verifyResult.Result == proofs.ProofSuccess {
			fmt.Println("✓ Verification succeeded")
		}
	}
	
	// Example 3: Show ProofData structure
	fmt.Println("\n=== ProofData Structure Example ===")
	// Create a mock successful proof for demonstration
	mockProofData := &proofs.ProofData{
		Proof:         []byte("mock_zk_proof_data_for_blue_eyes"),
		VerifyingKey:  []byte("mock_verifying_key_rs12913832"),
		PublicWitness: []byte("mock_witness_pos_28356859_G_A_genotype_2"),
		Result:        proofs.ProofSuccess,
	}
	
	// Convert to JSON for demonstration
	jsonData, err := json.MarshalIndent(mockProofData, "", "  ")
	if err != nil {
		log.Printf("JSON marshaling error: %v", err)
	} else {
		fmt.Printf("ProofData JSON:\n%s\n", string(jsonData))
	}
	
	// Example 4: Different proof results
	fmt.Println("\n=== Proof Result Types ===")
	results := []proofs.ProofResult{
		proofs.ProofSuccess,
		proofs.ProofFail,
		proofs.ProofUnknown,
	}
	
	for _, result := range results {
		fmt.Printf("%s -> %s\n", fmt.Sprintf("ProofResult(%d)", int(result)), result.String())
	}
}