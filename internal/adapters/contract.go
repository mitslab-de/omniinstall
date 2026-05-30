package adapter

// ContractRequirements documents the behavioral guarantees that all Adapter
// implementations must satisfy. These are validated by RunAdapterContract in
// the test package.
//
// An adapter MUST:
//   - return a non-empty string from Name()
//   - return false from CanHandle for unknown source types
//   - return a non-nil Result (not nil pointer) from Install and Remove
//   - return a non-nil VerificationResult from Verify
//   - populate ErrorCategory on failure results
