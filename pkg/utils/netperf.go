package utils

import (
	"strings"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

func ParseNetPerfOutput(str string, output *v1alpha1.NetworkPerformanceTestOutput) {

	lines := strings.Split(str, "\n")
	headers := strings.Split(lines[0], ",")

	for index, value := range headers {
		headers[index] = strings.TrimSpace(value)
	}

	output.Bandwidth = make(map[string]map[string]string)
	for _, line := range lines[1:] {
		values := strings.Split(line, ",")
		name := strings.TrimSpace(values[0])

		output.Bandwidth[name] = make(map[string]string)
		for i := 2; i < len(values)-1; i++ {
			output.Bandwidth[name][headers[i]] = values[i]
		}
	}

}
