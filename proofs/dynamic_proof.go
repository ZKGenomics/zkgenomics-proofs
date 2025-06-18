package proofs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/brentp/vcfgo"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// stringToInt converts nucleotide strings to integers for circuit use
func stringToInt(s string) int {
	switch strings.ToUpper(s) {
	case "A":
		return 0
	case "T":
		return 1
	case "G":
		return 2
	case "C":
		return 3
	default:
		return 0 // Default to A
	}
}

type DynamicCircuit struct {
	ClaimedRef       frontend.Variable `gnark:",public"`
	ClaimedAlt       frontend.Variable `gnark:",public"`
	ClaimedGenotype  frontend.Variable `gnark:",public"`
	ActualRef        frontend.Variable
	ActualAlt        frontend.Variable
	ActualGenotype   frontend.Variable
}

func (c *DynamicCircuit) Define(api frontend.API) error {
	// Verify that the claimed reference matches actual reference
	api.AssertIsEqual(c.ClaimedRef, c.ActualRef)
	
	// Verify that the claimed alternate matches actual alternate
	api.AssertIsEqual(c.ClaimedAlt, c.ActualAlt)
	
	// Verify that the claimed genotype matches actual genotype
	api.AssertIsEqual(c.ClaimedGenotype, c.ActualGenotype)

	return nil
}

// NewDynamicProof creates a new DynamicProof with specified genomic parameters
func NewDynamicProof(position uint64, reference string, alternate string) *DynamicProof {
	return &DynamicProof{
		Position:  position,
		Reference: reference,
		Alternate: alternate,
	}
}

// Generate implements the Proof interface for DynamicProof
func (p *DynamicProof) Generate(vcfPath string, provingKeyPath string, outputPath string) (*ProofData, error) {
	return p.GenerateDynamic(vcfPath, provingKeyPath, outputPath, p.Position, p.Reference, p.Alternate)
}

// GenerateDynamic implements the DynamicProofGenerator interface
func (p *DynamicProof) GenerateDynamic(vcfPath string, provingKeyPath string, outputPath string, position uint64, ref string, alt string) (*ProofData, error) {
	genotype, actualRef, actualAlt, err := p.extractGenotypeAtPosition(vcfPath, position, ref, alt)
	if err != nil {
		// Return ProofData with Fail result
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("failed to extract genotype: %w", err)
	}

	fmt.Printf("Found variant at position %d:\n", position)
	fmt.Printf("  Reference: %s (expected: %s)\n", actualRef, ref)
	fmt.Printf("  Alternate: %s (expected: %s)\n", actualAlt, alt)
	fmt.Printf("  Genotype: %d\n", genotype)

	// Verify that the found variant matches expected reference and alternate
	if actualRef != ref {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("reference mismatch: expected %s, found %s", ref, actualRef)
	}
	if actualAlt != alt {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("alternate mismatch: expected %s, found %s", alt, actualAlt)
	}

	// Generate actual zk-SNARK proof using gnark
	fmt.Printf("Generating proof for position %d with genotype %d\n", position, genotype)
	
	// Convert string values to integers for circuit
	refInt := stringToInt(actualRef)
	altInt := stringToInt(actualAlt)
	
	// Compile the circuit
	fmt.Println("Compiling dynamic circuit...")
	var circuit DynamicCircuit
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("circuit compilation error: %w", err)
	}

	// Setup proving system
	fmt.Println("Setting up proving system...")
	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("setup error: %w", err)
	}

	// Create witness
	fmt.Println("Creating witness...")
	witness := DynamicCircuit{
		ClaimedRef:      refInt,
		ClaimedAlt:      altInt,
		ClaimedGenotype: genotype,
		ActualRef:       refInt,
		ActualAlt:       altInt,
		ActualGenotype:  genotype,
	}

	w, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField())
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("witness creation error: %w", err)
	}

	publicWitness, err := w.Public()
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("public witness error: %w", err)
	}

	// Generate proof
	fmt.Println("Generating cryptographic proof...")
	proof, err := groth16.Prove(cs, pk, w)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("proving error: %w", err)
	}

	// Serialize proof to bytes
	proofBytes := make([]byte, 0)
	proofWriter := &bytesWriter{data: &proofBytes}
	_, err = proof.WriteTo(proofWriter)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("serializing proof: %w", err)
	}

	// Serialize verifying key to bytes
	vkBytes := make([]byte, 0)
	vkWriter := &bytesWriter{data: &vkBytes}
	_, err = vk.WriteTo(vkWriter)
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("serializing verifying key: %w", err)
	}

	// Serialize public witness to bytes
	publicWitnessData, err := publicWitness.MarshalBinary()
	if err != nil {
		return &ProofData{
			Proof:         nil,
			VerifyingKey:  nil,
			PublicWitness: nil,
			Result:        ProofFail,
		}, fmt.Errorf("serializing public witness: %w", err)
	}

	fmt.Printf("✅ Dynamic proof successfully generated for position %d!\n", position)

	return &ProofData{
		Proof:         proofBytes,
		VerifyingKey:  vkBytes,
		PublicWitness: publicWitnessData,
		Result:        ProofSuccess,
	}, nil
}

// Verify implements the Proof interface for DynamicProof
func (p *DynamicProof) Verify(verifyingKeyPath string, proofPath string) (*VerificationResult, error) {
	// Here you would implement the actual zk-SNARK proof verification
	// For now, we'll simulate the verification process
	fmt.Printf("Verifying proof for position %d\n", p.Position)
	
	// Simulate different verification outcomes based on simple heuristics
	// In a real implementation, this would involve cryptographic verification
	
	// For demonstration, we'll simulate successful verification
	result := &VerificationResult{
		Result: ProofSuccess,
		Error:  nil,
	}
	
	return result, nil
}

func (p *DynamicProof) VerifyProofData(proofData *ProofData) (*VerificationResult, error) {
	// Verify dynamic proof directly from ProofData using gnark
	
	if len(proofData.Proof) == 0 || len(proofData.VerifyingKey) == 0 {
		return &VerificationResult{
			Result: ProofFail,
			Error:  fmt.Errorf("invalid proof data: missing proof or verifying key"),
		}, nil
	}
	
	fmt.Printf("Verifying dynamic proof for position %d from ProofData...\n", p.Position)
	
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
	
	fmt.Printf("✅ Dynamic proof for position %d successfully verified!\n", p.Position)
	
	return &VerificationResult{
		Result: ProofSuccess,
		Error:  nil,
	}, nil
}

// extractGenotypeAtPosition searches for a specific genomic position in the VCF file
// and returns the genotype, reference, and alternate alleles
func (p *DynamicProof) extractGenotypeAtPosition(vcfPath string, position uint64, expectedRef string, expectedAlt string) (int, string, string, error) {
	f, err := os.Open(vcfPath)
	if err != nil {
		return 0, "", "", err
	}
	defer f.Close()

	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		return 0, "", "", err
	}

	fmt.Printf("Searching for position %d in VCF file...\n", position)
	
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}

		// Debug: print progress every 10000 positions
		if variant.Pos%10000 == 0 {
			fmt.Printf("Searching position: %d\n", variant.Pos)
		}

		if uint64(variant.Pos) == position {
			fmt.Printf("Found variant at position %d\n", position)
			
			// Extract genotype from the first sample
			if len(variant.Samples) == 0 {
				return 0, "", "", fmt.Errorf("no samples found in VCF")
			}
			
			sample := variant.Samples[0]
			genotypeInts := sample.GT
			
			// Handle Reference and Alternate which can be strings or slices
			ref := variant.Reference
			alt := ""
			if len(variant.Alternate) > 0 {
				alt = variant.Alternate[0] // Use first alternate allele
			}
			
			genotype, err := p.parseGenotypeFromInts(genotypeInts)
			if err != nil {
				return 0, "", "", fmt.Errorf("failed to parse genotype: %w", err)
			}
			
			return genotype, ref, alt, nil
		}
	}
	
	return 0, "", "", fmt.Errorf("position %d not found in VCF file", position)
}

// parseGenotypeFromInts converts VCF genotype from integer slice to genotype integer
func (p *DynamicProof) parseGenotypeFromInts(genotypeInts []int) (int, error) {
	if len(genotypeInts) != 2 {
		return 0, fmt.Errorf("expected diploid genotype, got %d alleles", len(genotypeInts))
	}
	
	allele1 := genotypeInts[0]
	allele2 := genotypeInts[1]
	
	// Handle missing data
	if allele1 < 0 || allele2 < 0 {
		return 0, fmt.Errorf("missing genotype data")
	}
	
	// Convert to genotype integer:
	// 0/0 (homozygous reference) = 0
	// 0/1 or 1/0 (heterozygous) = 1  
	// 1/1 (homozygous alternate) = 2
	if allele1 == 0 && allele2 == 0 {
		return 0, nil // Homozygous reference
	} else if (allele1 == 0 && allele2 == 1) || (allele1 == 1 && allele2 == 0) {
		return 1, nil // Heterozygous
	} else if allele1 == 1 && allele2 == 1 {
		return 2, nil // Homozygous alternate
	}
	
	return 0, fmt.Errorf("unsupported genotype: %v", genotypeInts)
}

// parseGenotype converts VCF genotype format (e.g., "0/0", "0/1", "1/1") to integer
// This method is kept for testing purposes
func (p *DynamicProof) parseGenotype(genotypeStr string, ref string, alt string) (int, error) {
	// Handle different genotype separators
	var alleles []string
	if strings.Contains(genotypeStr, "/") {
		alleles = strings.Split(genotypeStr, "/")
	} else if strings.Contains(genotypeStr, "|") {
		alleles = strings.Split(genotypeStr, "|")
	} else {
		return 0, fmt.Errorf("invalid genotype format: %s", genotypeStr)
	}
	
	if len(alleles) != 2 {
		return 0, fmt.Errorf("expected diploid genotype, got: %s", genotypeStr)
	}
	
	allele1, err := strconv.Atoi(alleles[0])
	if err != nil {
		return 0, fmt.Errorf("invalid allele: %s", alleles[0])
	}
	
	allele2, err := strconv.Atoi(alleles[1])
	if err != nil {
		return 0, fmt.Errorf("invalid allele: %s", alleles[1])
	}
	
	// Convert to genotype integer:
	// 0/0 (homozygous reference) = 0
	// 0/1 or 1/0 (heterozygous) = 1  
	// 1/1 (homozygous alternate) = 2
	if allele1 == 0 && allele2 == 0 {
		return 0, nil // Homozygous reference
	} else if (allele1 == 0 && allele2 == 1) || (allele1 == 1 && allele2 == 0) {
		return 1, nil // Heterozygous
	} else if allele1 == 1 && allele2 == 1 {
		return 2, nil // Homozygous alternate
	}
	
	return 0, fmt.Errorf("unsupported genotype: %s", genotypeStr)
}