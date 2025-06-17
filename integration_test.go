package zkgenomics

import (
	"testing"
)

// TestProofGeneratorIntegration tests the main API of the library
func TestProofGeneratorIntegration(t *testing.T) {
	generator := NewProofGenerator()
	
	// Test getting supported proof types
	supportedTypes := generator.GetSupportedProofTypes()
	expectedTypes := []ProofType{
		ChromosomeProofType,
		EyeColorProofType,
		BRCA1ProofType,
		HERC2ProofType,
	}
	
	if len(supportedTypes) != len(expectedTypes) {
		t.Errorf("Expected %d supported types, got %d", len(expectedTypes), len(supportedTypes))
	}
	
	for i, expected := range expectedTypes {
		if supportedTypes[i] != expected {
			t.Errorf("Expected type %s, got %s", expected, supportedTypes[i])
		}
	}
}

// TestProofResultTypes tests the result type system
func TestProofResultTypes(t *testing.T) {
	tests := []struct {
		result   ProofResult
		expected string
	}{
		{ProofSuccess, "success"},
		{ProofFail, "fail"},
		{ProofUnknown, "unknown"},
	}
	
	for _, test := range tests {
		result := test.result.String()
		if result != test.expected {
			t.Errorf("For ProofResult %d, expected %s, got %s", int(test.result), test.expected, result)
		}
	}
}

// TestProofDataStructure tests the ProofData structure
func TestProofDataStructure(t *testing.T) {
	proofData := &ProofData{
		Proof:         []byte("test_proof"),
		VerifyingKey:  []byte("test_vk"),
		PublicWitness: []byte("test_witness"),
		Result:        ProofSuccess,
	}
	
	if string(proofData.Proof) != "test_proof" {
		t.Errorf("Expected proof data 'test_proof', got %s", string(proofData.Proof))
	}
	
	if proofData.Result != ProofSuccess {
		t.Errorf("Expected ProofSuccess, got %s", proofData.Result.String())
	}
}