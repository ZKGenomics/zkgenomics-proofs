package proofs

import (
	"testing"
)

func TestNewDynamicProof(t *testing.T) {
	position := uint64(28356859)
	ref := "G"
	alt := "A"
	
	proof := NewDynamicProof(position, ref, alt)
	
	if proof.Position != position {
		t.Errorf("Expected position %d, got %d", position, proof.Position)
	}
	
	if proof.Reference != ref {
		t.Errorf("Expected reference %s, got %s", ref, proof.Reference)
	}
	
	if proof.Alternate != alt {
		t.Errorf("Expected alternate %s, got %s", alt, proof.Alternate)
	}
}

func TestParseGenotype(t *testing.T) {
	proof := &DynamicProof{}
	
	tests := []struct {
		genotype string
		ref      string
		alt      string
		expected int
		hasError bool
	}{
		{"0/0", "G", "A", 0, false}, // Homozygous reference
		{"0/1", "G", "A", 1, false}, // Heterozygous
		{"1/0", "G", "A", 1, false}, // Heterozygous (reversed)
		{"1/1", "G", "A", 2, false}, // Homozygous alternate
		{"0|0", "G", "A", 0, false}, // Phased homozygous reference
		{"0|1", "G", "A", 1, false}, // Phased heterozygous
		{"1|1", "G", "A", 2, false}, // Phased homozygous alternate
		{"invalid", "G", "A", 0, true}, // Invalid format
		{"2/2", "G", "A", 0, true},     // Unsupported genotype
	}
	
	for _, test := range tests {
		result, err := proof.parseGenotype(test.genotype, test.ref, test.alt)
		
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for genotype %s, but got none", test.genotype)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for genotype %s: %v", test.genotype, err)
			}
			if result != test.expected {
				t.Errorf("For genotype %s, expected %d, got %d", test.genotype, test.expected, result)
			}
		}
	}
}

func TestDynamicProofEyeColorExample(t *testing.T) {
	// Test the eye color example from the requirements
	// SNP: rs12913832, Position: 28356859, REF: G, ALT: A
	position := uint64(28356859)
	ref := "G"
	alt := "A"
	
	proof := NewDynamicProof(position, ref, alt)
	
	// Test genotype mapping for eye color
	// 0/0 = G/G → Brown (genotype 0)
	// 0/1 = G/A → Hazel/Green (genotype 1) 
	// 1/1 = A/A → Blue (genotype 2)
	
	tests := []struct {
		genotype     string
		expectedType int
		eyeColor     string
	}{
		{"0/0", 0, "Brown"},
		{"0/1", 1, "Hazel/Green"},
		{"1/1", 2, "Blue"},
	}
	
	for _, test := range tests {
		result, err := proof.parseGenotype(test.genotype, ref, alt)
		if err != nil {
			t.Errorf("Unexpected error for genotype %s: %v", test.genotype, err)
		}
		if result != test.expectedType {
			t.Errorf("For genotype %s (eye color: %s), expected type %d, got %d", 
				test.genotype, test.eyeColor, test.expectedType, result)
		}
	}
}

func TestDynamicProofInterfaces(t *testing.T) {
	proof := NewDynamicProof(12345, "C", "T")
	
	// Test that DynamicProof implements Proof interface
	var _ Proof = proof
	
	// Test that DynamicProof implements DynamicProofGenerator interface
	var _ DynamicProofGenerator = proof
}

func TestVerify(t *testing.T) {
	proof := NewDynamicProof(12345, "C", "T")
	
	result, err := proof.Verify("dummy_key_path", "dummy_proof_path")
	if err != nil {
		t.Errorf("Unexpected error during verification: %v", err)
	}
	if result.Result != ProofSuccess {
		t.Errorf("Expected verification to succeed, got %s", result.Result.String())
	}
	if result.Error != nil {
		t.Errorf("Expected no error in result, got: %v", result.Error)
	}
}

func TestProofDataStructure(t *testing.T) {
	proof := NewDynamicProof(28356859, "G", "A")
	
	// Test successful proof generation simulation
	proofData, err := proof.GenerateDynamic("dummy.vcf", "dummy.key", "dummy.out", 28356859, "G", "A")
	
	// Since we don't have a real VCF file, this should fail with ProofFail
	if err == nil {
		t.Error("Expected error when VCF file doesn't exist")
	}
	if proofData == nil {
		t.Error("Expected ProofData to be returned even on failure")
	}
	if proofData.Result != ProofFail {
		t.Errorf("Expected ProofFail result, got %s", proofData.Result.String())
	}
}

func TestProofResultString(t *testing.T) {
	tests := []struct {
		result   ProofResult
		expected string
	}{
		{ProofSuccess, "success"},
		{ProofFail, "fail"},
		{ProofUnknown, "unknown"},
		{ProofResult(999), "unknown"}, // Invalid value should return "unknown"
	}
	
	for _, test := range tests {
		result := test.result.String()
		if result != test.expected {
			t.Errorf("For ProofResult %d, expected %s, got %s", int(test.result), test.expected, result)
		}
	}
}