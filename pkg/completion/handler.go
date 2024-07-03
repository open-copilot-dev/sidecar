package completion

import "fmt"

func ProcessRequest(request *CompletionRequest) (*CompletionResult, error) {
	fmt.Println(request)
	return &CompletionResult{}, nil
}
