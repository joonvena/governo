package rules

import "context"

type ActionHandler interface {
	Type() string
	Handle(ctx context.Context, resource string, params map[string]any, attrs map[string]any) error
}

var actionHandlers = map[string]ActionHandler{}

func RegisterAction(h ActionHandler) {
	actionHandlers[h.Type()] = h
}

func ExecuteActions(ctx context.Context, actions []Action, resource string, attrs map[string]any) error {
	for _, a := range actions {
		if h, ok := actionHandlers[a.Type]; ok {
			_ = h.Handle(ctx, resource, a.Params, attrs)
		}
	}
	return nil
}

func UnregisterAllActions() {
	actionHandlers = make(map[string]ActionHandler)
}
