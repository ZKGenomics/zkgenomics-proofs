package proofs

import (
	"os"
	"testing"
)

func TestBRCA1Proof_Generate(t *testing.T) {
	// Create a temporary VCF file for testing
	vcfContent := `##fileformat=VCFv4.2
##INFO=<ID=DP,Number=1,Type=Integer,Description="Approximate read depth">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO
17	41276045	.	A	G	60	PASS	DP=30
`

	tmpFile, err := os.CreateTemp("", "test*.vcf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(vcfContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	proof := &BRCA1Proof{}
	proofData, err := proof.Generate(tmpFile.Name(), "", "")
	if err != nil {
		t.Errorf("Generate should not return error: %v", err)
	}
	if proofData.Result != ProofSuccess {
		t.Errorf("Expected ProofSuccess, got %s", proofData.Result.String())
	}
}

func TestBRCA1Proof_GenerateWithMissingPosition(t *testing.T) {
	// Create a temporary VCF file without the target position
	vcfContent := `##fileformat=VCFv4.2
##INFO=<ID=DP,Number=1,Type=Integer,Description="Approximate read depth">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO
17	12345678	.	A	G	60	PASS	DP=30
`

	tmpFile, err := os.CreateTemp("", "test*.vcf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(vcfContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	proof := &BRCA1Proof{}
	proofData, err := proof.Generate(tmpFile.Name(), "", "")
	if err == nil {
		t.Errorf("Generate should return error when position not found")
	}
	if proofData.Result != ProofFail {
		t.Errorf("Expected ProofFail, got %s", proofData.Result.String())
	}
}

func TestBRCA1Proof_Verify(t *testing.T) {
	proof := &BRCA1Proof{}
	result, err := proof.Verify("", "")
	if err != nil {
		t.Errorf("Verify should not return error: %v", err)
	}
	if result.Result != ProofSuccess {
		t.Errorf("Expected ProofSuccess, got %s", result.Result.String())
	}
}