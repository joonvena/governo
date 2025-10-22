package log

import (
	"context"
	"fmt"

	"governo/internal/plugins/common"
)

type LogAction struct{}

func (l *LogAction) Type() string { return "log" }

func (l *LogAction) Handle(ctx context.Context, resource string, params map[string]any, attrs map[string]any) error {
	msg, _ := params["message"].(string)

	processedMsg := common.ProcessTemplate(msg, attrs)

	fmt.Printf("[LOG] resource=%s msg=%s\n", resource, processedMsg)
	return nil
}
