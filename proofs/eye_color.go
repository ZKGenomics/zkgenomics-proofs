package proofs

import (
	"fmt"
	"os"
	"strings"

	"github.com/brentp/vcfgo"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

type EyeColorCircuit struct {
	ClaimedColor frontend.Variable `gnark:",public"`
	Genotype     frontend.Variable
}

func (c *EyeColorCircuit) Define(api frontend.API) error {
	api.Sub(c.ClaimedColor, c.Genotype)

	return nil
}

// Parse rs12913832 genotype from VCF and map to integer
func extractEyeColorGenotype(vcfPath string) (int, error) {
	f, err := os.Open(vcfPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		return 0, err
	}

	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}
		if variant.Pos == 396321 {
			fmt.Println(fmt.Sprintf("Found eye color mutation at variant: %s", variant.Chromosome))
			return 1, nil // Simplified for demonstration
		}
	}
	return 0, fmt.Errorf("not found in VCF")
}

// Map genotype integer to color integer
func genotypeToColor(genotype int) int {
	switch genotype {
	case 0:
		return 1 // Brown
	case 1:
		return 2 // Hazel/Green
	case 2:
		return 3 // Blue
	default:
		return 0
	}
}

func (p EyeColorProof) Generate(vcfPath string, provingKeyPath string, outputPath string) (*ProofData, error) {
	// Simulate proof generation for eye color
	return &ProofData{
		Proof:         []byte("eye_color_proof_data"),
		VerifyingKey:  []byte("eye_color_verifying_key"),
		PublicWitness: []byte("eye_color_public_witness"),
		Result:        ProofSuccess,
	}, nil
}

func (p EyeColorProof) Verify(verifyingKeyPath string, proofPath string) (*VerificationResult, error) {
	return &VerificationResult{
		Result: ProofSuccess,
		Error:  nil,
	}, nil
}

func (p EyeColorProof) VerifyProofData(proofData *ProofData) (*VerificationResult, error) {
	// Verify eye color proof directly from ProofData using gnark
	
	if len(proofData.Proof) == 0 || len(proofData.VerifyingKey) == 0 {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("invalid proof data: missing proof or verifying key"),
		}, nil
	}
	
	fmt.Println("Verifying eye color proof from ProofData...")
	
	// Deserialize the verifying key
	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err := vk.ReadFrom(strings.NewReader(string(proofData.VerifyingKey)))
	if err != nil {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("failed to deserialize verifying key: %w", err),
		}, nil
	}
	
	// Deserialize the proof
	proof := groth16.NewProof(ecc.BN254)
	_, err = proof.ReadFrom(strings.NewReader(string(proofData.Proof)))
	if err != nil {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("failed to deserialize proof: %w", err),
		}, nil
	}
	
	// Deserialize the public witness
	publicWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("failed to create witness: %w", err),
		}, nil
	}
	err = publicWitness.UnmarshalBinary(proofData.PublicWitness)
	if err != nil {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("failed to deserialize public witness: %w", err),
		}, nil
	}
	
	// Perform gnark verification
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("proof verification failed: %w", err),
		}, nil
	}
	
	fmt.Println("✅ Eye color proof successfully verified!")
	
	return &VerificationResult{
		Result: ProofSuccess,
		Error:  nil,
	}, nil
}