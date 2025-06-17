package zkgenomics

import "fmt"

// UnsupportedProofTypeError represents an error when an unsupported proof type is requested
type UnsupportedProofTypeError struct {
	Type string
}

func (e *UnsupportedProofTypeError) Error() string {
	return fmt.Sprintf("unsupported proof type: %s", e.Type)
}

// ProofGenerationError represents an error during proof generation
type ProofGenerationError struct {
	ProofType string
	Err       error
}

func (e *ProofGenerationError) Error() string {
	return fmt.Sprintf("failed to generate %s proof: %v", e.ProofType, e.Err)
}

func (e *ProofGenerationError) Unwrap() error {
	return e.Err
}

// ProofVerificationError represents an error during proof verification
type ProofVerificationError struct {
	ProofType string
	Err       error
}

func (e *ProofVerificationError) Error() string {
	return fmt.Sprintf("failed to verify %s proof: %v", e.ProofType, e.Err)
}

func (e *ProofVerificationError) Unwrap() error {
	return e.Err
}