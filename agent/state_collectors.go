package agent

import (
	"encoding/json"
	"fmt"

	"github.com/bducha/assistagent/system"
)

// Return the memory state of the machine
func GetMemoryState() (string, error) {
	memory, err := system.GetMemory()

	if err != nil {
		return "", err
	}

	payload := struct{
		TotalMemory uint64 `json:"total_memory"`
		FreeMemory uint64 `json:"free_memory"`
		UsedMemory uint64 `json:"used_memory"`
	}{
		TotalMemory: memory.TotalMemory,
		FreeMemory: memory.FreeMemory,
		UsedMemory: memory.UsedMemory,
	}


	payloadJson, err := json.Marshal(payload)

	fmt.Println(payloadJson)

	if err != nil {
		return "", err
	}

	return string(payloadJson), nil
}