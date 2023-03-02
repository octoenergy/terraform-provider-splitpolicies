package provider

import "fmt"

func split_policies(inputs []string, max_size int) ([][]string, string) {
	var chunks = make([][]string, 0)

	for i, input := range inputs {
		if len(input) > max_size {
			return nil, fmt.Sprintf("Input %v exceeeds max size", i)
		}
	}

	for len(inputs) > 0 {
		var chunk, inps = assemble_chunk(inputs, max_size)
		inputs = inps
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}
	}
	
	return chunks, ""
}

// Assemble a chunk from the list of inputs, returns the chunk and everything that didn't fit
// Precondition: None of the inputs exceeds the maximum size
func assemble_chunk(inputs []string, max_size int) (out []string, overflow []string) {
	var l = 0
	for _, input := range inputs {
		if l + len(input) <= max_size {
			out = append(out, input)
			l += len(input)
		} else {
			overflow = append(overflow, input)
		}
	}

	return out, overflow
}
