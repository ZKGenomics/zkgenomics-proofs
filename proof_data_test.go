package zkgenomics

import (
	"errors"
	"testing"
)

func TestProofGenerator_VerifyProofData(t *testing.T) {
	pg := NewProofGenerator()

	// Test with fake proof data (should fail with real gnark verification)
	fakeProofData := &ProofData{
		Proof:         []byte("test_proof_data"),
		VerifyingKey:  []byte("test_verifying_key"),
		PublicWitness: []byte("test_public_witness"),
		Result:        ProofSuccess,
	}

	// Test each proof type - expect failures since data is fake
	testCases := []struct {
		proofType ProofType
		name      string
	}{
		{EyeColorProofType, "EyeColor"},
		{BRCA1ProofType, "BRCA1"},
		{HERC2ProofType, "HERC2"},
		{ChromosomeProofType, "Chromosome"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := pg.VerifyProofData(tc.proofType, fakeProofData)
			// With real gnark verification, fake proof data should fail
			if err != nil {
				// This is expected - fake data should cause deserialization errors
				t.Logf("Expected failure for fake proof data in %s: %v", tc.proofType, err)
				return
			}
			if result.Result != ProofFail {
				t.Errorf("Expected ProofFail for fake data in %s, got %s", tc.proofType, result.Result.String())
			}
		})
	}
}

func TestProofGenerator_VerifyProofData_InvalidData(t *testing.T) {
	pg := NewProofGenerator()

	// Test with empty proof data
	emptyProofData := &ProofData{
		Proof:         []byte{},
		VerifyingKey:  []byte{},
		PublicWitness: []byte("test_witness"),
		Result:        ProofUnknown,
	}

	result, err := pg.VerifyProofData(EyeColorProofType, emptyProofData)
	if err != nil {
		t.Errorf("VerifyProofData should not return error for invalid data: %v", err)
	}
	if result.Result != ProofFail {
		t.Errorf("Expected ProofFail for empty proof data, got %s", result.Result.String())
	}
}

func TestProofGenerator_VerifyProofData_UnsupportedType(t *testing.T) {
	pg := NewProofGenerator()

	validProofData := &ProofData{
		Proof:         []byte("test_proof_data"),
		VerifyingKey:  []byte("test_verifying_key"),
		PublicWitness: []byte("test_public_witness"),
		Result:        ProofSuccess,
	}

	_, err := pg.VerifyProofData(ProofType("unsupported"), validProofData)
	if err == nil {
		t.Error("Expected error for unsupported proof type")
	}

	var unsupportedErr *UnsupportedProofTypeError
	if !errors.As(err, &unsupportedErr) {
		t.Error("Expected UnsupportedProofTypeError")
	}
}

func TestProofGenerator_VerifyAnyProofData(t *testing.T) {
	pg := NewProofGenerator()

	fakeProofData := &ProofData{
		Proof:         []byte("test_proof_data"),
		VerifyingKey:  []byte("test_verifying_key"),
		PublicWitness: []byte("test_public_witness"),
		Result:        ProofSuccess,
	}

	proofType, result, err := pg.VerifyAnyProofData(fakeProofData)
	// With real gnark verification, this should fail
	if err != nil {
		t.Logf("Expected failure with fake proof data: %v", err)
		return
	}
	if result.Result != ProofFail {
		t.Errorf("Expected ProofFail for fake data, got %s", result.Result.String())
	}
	if proofType != "" {
		t.Errorf("Expected empty proof type for failed verification, got %s", proofType)
	}
}

func TestProofGenerator_VerifyAnyProofData_InvalidData(t *testing.T) {
	pg := NewProofGenerator()

	invalidProofData := &ProofData{
		Proof:         []byte{},
		VerifyingKey:  []byte{},
		PublicWitness: []byte("test_witness"),
		Result:        ProofUnknown,
	}

	proofType, result, err := pg.VerifyAnyProofData(invalidProofData)
	if err != nil {
		t.Errorf("VerifyAnyProofData should not return error: %v", err)
	}
	if result.Result != ProofFail {
		t.Errorf("Expected ProofFail for invalid data, got %s", result.Result.String())
	}
	if proofType != "" {
		t.Errorf("Expected empty proof type for failed verification, got %s", proofType)
	}
}