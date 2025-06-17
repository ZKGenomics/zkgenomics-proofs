package zkgenomics

import (
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
	default:
		return nil, &UnsupportedProofTypeError{Type: string(proofType)}
	}

	return proof.Verify(verifyingKeyPath, proofPath)
}

// GetSupportedProofTypes returns a list of supported proof types
func (pg *ProofGenerator) GetSupportedProofTypes() []ProofType {
	return []ProofType{
		ChromosomeProofType,
		EyeColorProofType,
		BRCA1ProofType,
		HERC2ProofType,
	}
}

// TraitVariant re-exports the trait variant structure for convenience
type TraitVariant = traits.TraitVariant

// TraitRegion re-exports the trait region structure for convenience  
type TraitRegion = traits.TraitRegion

// TraitPanel re-exports the trait panel structure for convenience
type TraitPanel = traits.TraitPanel