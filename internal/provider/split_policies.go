package provider

import "fmt"

func splitPolicies(inputs []string, maxSize int) ([][]string, error) {
	var chunks = make([][]string, 0)

	for i, input := range inputs {
		if len(input) > maxSize {
			return nil, fmt.Errorf("input %v exceeeds max size", i)
		}
	}

	remainingInputs := inputs
	var chunk []string
	for len(remainingInputs) > 0 {
		chunk, remainingInputs = assembleChunk(remainingInputs, maxSize)
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}
	}

	return chunks, nil
}

// Assemble a chunk from the list of inputs, returns the chunk and everything that didn't fit
// Precondition: None of the inputs exceeds the maximum size
func assembleChunk(inputs []string, maxSize int) (out []string, overflow []string) {
	var l = 0
	for _, input := range inputs {
		if l+len(input) <= maxSize {
			out = append(out, input)
			l += len(input)
		} else {
			overflow = append(overflow, input)
		}
	}

	return out, overflow
}
