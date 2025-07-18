package zkgenomics

import (
	"fmt"
	"github.com/zkgenomics/zkgenomics-proofs/proofs"
	"github.com/zkgenomics/zkgenomics-proofs/traits"
)

// Re-export important types for convenience
type ProofData = proofs.ProofData
type VerificationResult = proofs.VerificationResult
type ProofResult = proofs.ProofResult

// Re-export constants
const (
	ProofSuccess ProofResult = proofs.ProofSuccess
	ProofFail    ProofResult = proofs.ProofFail
	ProofUnknown ProofResult = proofs.ProofUnknown
)

// ProofType represents the type of genomic proof to generate
type ProofType string

const (
	ChromosomeProofType ProofType = "chromosome"
	EyeColorProofType   ProofType = "eye_color"
	BRCA1ProofType      ProofType = "brca1"
	HERC2ProofType      ProofType = "herc2"
	DynamicProofType    ProofType = "dynamic"
)

// ProofGenerator provides a unified interface for generating genomic proofs
type ProofGenerator struct{}

// NewProofGenerator creates a new proof generator instance
func NewProofGenerator() *ProofGenerator {
	return &ProofGenerator{}
}

// GenerateProof generates a proof of the specified type and returns the proof data
func (pg *ProofGenerator) GenerateProof(proofType ProofType, vcfPath, provingKeyPath, outputPath string) (*ProofData, error) {
	var proof proofs.Proof

	switch proofType {
	case ChromosomeProofType:
		proof = &proofs.ChromosomeProof{}
	case EyeColorProofType:
		proof = &proofs.EyeColorProof{}
	case BRCA1ProofType:
		proof = &proofs.BRCA1Proof{}
	case HERC2ProofType:
		proof = &proofs.HERC2Proof{}
	case DynamicProofType:
		proof = &proofs.DynamicProof{}
	default:
		return nil, &UnsupportedProofTypeError{Type: string(proofType)}
	}

	return proof.Generate(vcfPath, provingKeyPath, outputPath)
}

// VerifyProof verifies a proof of the specified type and returns the verification result
func (pg *ProofGenerator) VerifyProof(proofType ProofType, verifyingKeyPath, proofPath string) (*VerificationResult, error) {
	var proof proofs.Proof

	switch proofType {
	case ChromosomeProofType:
		proof = &proofs.ChromosomeProof{}
	case EyeColorProofType:
		proof = &proofs.EyeColorProof{}
	case BRCA1ProofType:
		proof = &proofs.BRCA1Proof{}
	case HERC2ProofType:
		proof = &proofs.HERC2Proof{}
	case DynamicProofType:
		proof = &proofs.DynamicProof{}
	default:
		return nil, &UnsupportedProofTypeError{Type: string(proofType)}
	}

	return proof.Verify(verifyingKeyPath, proofPath)
}

// VerifyProofData verifies a proof directly from ProofData without file operations
func (pg *ProofGenerator) VerifyProofData(proofType ProofType, proofData *ProofData) (*VerificationResult, error) {
	var proof proofs.Proof

	switch proofType {
	case ChromosomeProofType:
		proof = &proofs.ChromosomeProof{}
	case EyeColorProofType:
		proof = &proofs.EyeColorProof{}
	case BRCA1ProofType:
		proof = &proofs.BRCA1Proof{}
	case HERC2ProofType:
		proof = &proofs.HERC2Proof{}
	case DynamicProofType:
		proof = &proofs.DynamicProof{}
	default:
		return nil, &UnsupportedProofTypeError{Type: string(proofType)}
	}

	return proof.VerifyProofData(proofData)
}

// VerifyAnyProofData attempts to verify ProofData by trying all supported proof types
// This is useful when the proof type is unknown or not stored with the proof
func (pg *ProofGenerator) VerifyAnyProofData(proofData *ProofData) (ProofType, *VerificationResult, error) {
	supportedTypes := pg.GetSupportedProofTypes()
	
	for _, proofType := range supportedTypes {
		result, err := pg.VerifyProofData(proofType, proofData)
		if err != nil {
			continue // Try next proof type
		}
		
		if result.Result == ProofSuccess {
			return proofType, result, nil
		}
	}
	
	return "", &VerificationResult{
		Result: ProofFail,
		Error:  fmt.Errorf("proof verification failed for all supported types"),
	}, nil
}

// GetSupportedProofTypes returns a list of supported proof types
func (pg *ProofGenerator) GetSupportedProofTypes() []ProofType {
	return []ProofType{
		ChromosomeProofType,
		EyeColorProofType,
		BRCA1ProofType,
		HERC2ProofType,
		DynamicProofType,
	}
}

// TraitVariant re-exports the trait variant structure for convenience
type TraitVariant = traits.TraitVariant

// TraitRegion re-exports the trait region structure for convenience  
type TraitRegion = traits.TraitRegion

// TraitPanel re-exports the trait panel structure for convenience
type TraitPanel = traits.TraitPanel