package proofs

type Proof interface {
	Generate(vcfPath string, provingKeyPath string, outputPath string) error
	Verify(verifyingKeyPath string, proofPath string) (bool, error)
}

// DynamicProofGenerator interface for proofs that can be configured with specific genomic parameters
type DynamicProofGenerator interface {
	Proof
	GenerateDynamic(vcfPath string, provingKeyPath string, outputPath string, position uint64, ref string, alt string) error
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
