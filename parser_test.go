package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCalls(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expCalls []Call
		expErr   error
	}{
		{
			name:     "empty",
			input:    "",
			expCalls: nil,
			expErr:   nil,
		},
		{
			name:     "no function",
			input:    "invalid",
			expCalls: nil,
			expErr:   nil,
		},
		{
			name:  "no args",
			input: "<function>function1()",
			expCalls: []Call{
				{
					Tool:  "function1",
					Input: map[string]any{},
				},
			},
			expErr: nil,
		},
		{
			name:  "one arg",
			input: "<function>function1(arg='value')",
			expCalls: []Call{
				{
					Tool: "function1",
					Input: map[string]any{
						"arg": "value",
					},
				},
			},
			expErr: nil,
		},
		{
			name:  "two arg",
			input: "<function>function1(one='1', two='2')",
			expCalls: []Call{
				{
					Tool: "function1",
					Input: map[string]any{
						"one": "1",
						"two": "2",
					},
				},
			},
			expErr: nil,
		},
		{
			name:  "two functions",
			input: "<function>add(one='1', two='2')<<function>>sub(one='3', two='4')",
			expCalls: []Call{
				{
					Tool: "add",
					Input: map[string]any{
						"one": "1",
						"two": "2",
					},
				},
				{
					Tool: "sub",
					Input: map[string]any{
						"one": "3",
						"two": "4",
					},
				},
			},
			expErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			act, err := ParseCalls(tc.input)
			rq.ErrorIs(tc.expErr, err)
			rq.Equal(tc.expCalls, act)
		})
	}
}
