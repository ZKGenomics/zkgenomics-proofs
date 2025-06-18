package proofs

// ProofResult represents the possible outcomes of proof operations
type ProofResult int

const (
	ProofSuccess ProofResult = iota
	ProofFail
	ProofUnknown
)

// String returns string representation of ProofResult
func (r ProofResult) String() string {
	switch r {
	case ProofSuccess:
		return "success"
	case ProofFail:
		return "fail"
	case ProofUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// ProofData contains all necessary data for verification
type ProofData struct {
	Proof         []byte      `json:"proof"`
	VerifyingKey  []byte      `json:"verifying_key"`
	PublicWitness []byte      `json:"public_witness"`
	Result        ProofResult `json:"result"`
}

// VerificationResult contains the result of proof verification
type VerificationResult struct {
	Result ProofResult `json:"result"`
	Error  error       `json:"error,omitempty"`
}

type Proof interface {
	Generate(vcfPath string, provingKeyPath string, outputPath string) (*ProofData, error)
	Verify(verifyingKeyPath string, proofPath string) (*VerificationResult, error)
	VerifyProofData(proofData *ProofData) (*VerificationResult, error)
}

// DynamicProofGenerator interface for proofs that can be configured with specific genomic parameters
type DynamicProofGenerator interface {
	Proof
	GenerateDynamic(vcfPath string, provingKeyPath string, outputPath string, position uint64, ref string, alt string) (*ProofData, error)
}

type ChromosomeProof struct {
	Proof
}

type EyeColorProof struct {
	Proof
}

type BRCA1Proof struct {
	Proof
}

type HERC2Proof struct {
	Proof
}

type DynamicProof struct {
	Position uint64
	Reference string
	Alternate string
}

const HERC2Pos uint64 = 28365618
