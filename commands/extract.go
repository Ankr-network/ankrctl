package commands

import (
	"fmt"
	"strconv"
)

func extractTaskIDs(s []string) ([]int, error) {
	taskIDs := []int{}

	for _, e := range s {
		i, err := strconv.Atoi(e)
		if err != nil {
			return nil, fmt.Errorf("Provided value [%v] for task id is not of type int", e)
		}
		taskIDs = append(taskIDs, i)
	}

	return taskIDs, nil
}
