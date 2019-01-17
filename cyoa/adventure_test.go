package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseAdventure(t *testing.T) {
	tests := []struct {
		desc   string
		in     string
		expect Adventure
	}{
		{
			desc: "Works successfully",
			in: `{
				"intro": {
					"title": "Yet Another Story",
					"story": [
						"This is a story.",
						"It's a really great story."
					],
					"options": [
						{
							"text": "Option 1",
							"arc": "one"
						},
						{
							"text": "Option 2",
							"arc": "two"
						}
					]
				},
				"one": {
					"title": "One",
					"story": ["The end 1!"]
				},
				"two": {
					"title": "Two",
					"story": ["The end 2!"]
				}
			}`,
			expect: Adventure{
				arcs: map[string]Arc{
					"intro": Arc{
						Title: "Yet Another Story",
						Story: []string{
							"This is a story.",
							"It's a really great story.",
						},
						Options: []Option{
							Option{"Option 1", "one"},
							Option{"Option 2", "two"},
						},
					},
					"one": Arc{
						Title: "One",
						Story: []string{"The end 1!"},
					},
					"two": Arc{
						Title: "Two",
						Story: []string{"The end 2!"},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			adv, err := parseAdventure(strings.NewReader(tc.in))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(adv, tc.expect) {
				t.Fatalf("Want: %v\nExpect: %v", tc.expect, adv)
			}
		})
	}
}
