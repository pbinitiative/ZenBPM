// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sql

type ProcessInstance struct {
	Key                  int64  `json:"key"`
	ProcessDefinitionKey int64  `json:"process_definition_key"`
	CreatedAt            int64  `json:"created_at"`
	State                int64  `json:"state"`
	VariableHolder       string `json:"variable_holder"`
	CaughtEvents         string `json:"caught_events"`
	Activities           string `json:"activities"`
}
