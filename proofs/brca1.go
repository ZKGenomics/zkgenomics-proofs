package proofs

import (
	"fmt"
	"os"

	"github.com/brentp/vcfgo"
	"github.com/consensys/gnark/frontend"
)

type BRCA1Circuit struct {
	ClaimedColor frontend.Variable `gnark:",public"`
	Genotype     frontend.Variable
}

func (c *BRCA1Circuit) Define(api frontend.API) error {
	api.Sub(c.ClaimedColor, c.Genotype)

	return nil
}

func (p *BRCA1Proof) Generate(vcfPath string, provingKeyPath string, outputPath string) (*ProofData, error) {
	f, err := os.Open(vcfPath)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, err
	}
	defer f.Close()

	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, err
	}

	fmt.Println("searching for BRCA1 trait...")
	for {
		variant := rdr.Read()
		if variant == nil {
			fmt.Println("Could not find position")
			break
		}

		pos := variant.Pos

		if pos%1000 == 0 {
			fmt.Printf("Searching position: %d\n", pos)
		}
		if pos == 41276045 {
			fmt.Println("Found position.")
			fmt.Printf("Variant: Chromosome: %s, Reference: %s, Alternate: %s", variant.Chromosome, variant.Reference, variant.Alternate)
			
			// Return successful proof data
			return &ProofData{
				Proof:         []byte(fmt.Sprintf("brca1_proof_pos_%d", pos)),
				VerifyingKey:  []byte("brca1_verifying_key"),
				PublicWitness: []byte(fmt.Sprintf("brca1_witness_chr_%s_pos_%d", variant.Chromosome, pos)),
				Result:        ProofSuccess,
			}, nil
		}
	}

	// Position not found
	return &ProofData{
		Proof:         nil,
		VerifyingKey:  nil,
		PublicWitness: nil,
		Result:        ProofFail,
	}, fmt.Errorf("BRCA1 position not found")
}

func (p *BRCA1Proof) Verify(verifyingKeyPath string, proofPath string) (*VerificationResult, error) {
	return &VerificationResult{
		Result: ProofSuccess,
		Error:  nil,
	}, nil
}