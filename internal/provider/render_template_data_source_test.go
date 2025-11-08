package provider

import (
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		values      map[string]string
		expected    string
		expectError bool
	}{
		{
			name:     "basic replacement",
			template: "Hello @@NAME@@!",
			values: map[string]string{
				"NAME": "World",
			},
			expected:    "Hello World!",
			expectError: false,
		},
		{
			name:     "multiple placeholders",
			template: "@@GREETING@@ @@NAME@@, welcome to @@PLACE@@!",
			values: map[string]string{
				"GREETING": "Hello",
				"NAME":     "Alice",
				"PLACE":    "Wonderland",
			},
			expected:    "Hello Alice, welcome to Wonderland!",
			expectError: false,
		},
		{
			name:     "repeated placeholder",
			template: "@@VAR@@ and @@VAR@@ again",
			values: map[string]string{
				"VAR": "test",
			},
			expected:    "test and test again",
			expectError: false,
		},
		{
			name:        "missing placeholder",
			template:    "Hello @@NAME@@!",
			values:      map[string]string{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "no placeholders",
			template:    "Just a plain string",
			values:      map[string]string{},
			expected:    "Just a plain string",
			expectError: false,
		},
		{
			name:        "empty template",
			template:    "",
			values:      map[string]string{},
			expected:    "",
			expectError: false,
		},
		{
			name:     "extra values provided",
			template: "Hello @@NAME@@!",
			values: map[string]string{
				"NAME":  "World",
				"EXTRA": "ignored",
			},
			expected:    "Hello World!",
			expectError: false,
		},
		{
			name: "preserve Argo syntax",
			template: `apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: @@WORKFLOW_NAME@@
  namespace: @@NAMESPACE@@
spec:
  entrypoint: main
  arguments:
    parameters:
      - name: input
        value: "{{inputs.parameters.input}}"
  templates:
    - name: main
      inputs:
        parameters:
          - name: input
      container:
        image: alpine:latest
        command: [sh, -c]
        args: ["echo {{inputs.parameters.input}}"]`,
			values: map[string]string{
				"WORKFLOW_NAME": "test-workflow",
				"NAMESPACE":     "default",
			},
			expected: `apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: test-workflow
  namespace: default
spec:
  entrypoint: main
  arguments:
    parameters:
      - name: input
        value: "{{inputs.parameters.input}}"
  templates:
    - name: main
      inputs:
        parameters:
          - name: input
      container:
        image: alpine:latest
        command: [sh, -c]
        args: ["echo {{inputs.parameters.input}}"]`,
			expectError: false,
		},
		{
			name: "preserve shell syntax",
			template: `#!/bin/bash
NAMESPACE=@@NAMESPACE@@
IMAGE_TAG=@@IMAGE_TAG@@

# Use shell variables
echo "Namespace: ${NAMESPACE}"
echo "Tag: $(echo $IMAGE_TAG)"

# Bash conditionals
if [[ -n "$NAMESPACE" ]]; then
  echo "Valid"
fi

# Arithmetic
result=$((1 + 2))

# Logical operators
[[ -f file.txt ]] && echo "exists" || echo "not found"`,
			values: map[string]string{
				"NAMESPACE": "prod",
				"IMAGE_TAG": "v1.0.0",
			},
			expected: `#!/bin/bash
NAMESPACE=prod
IMAGE_TAG=v1.0.0

# Use shell variables
echo "Namespace: ${NAMESPACE}"
echo "Tag: $(echo $IMAGE_TAG)"

# Bash conditionals
if [[ -n "$NAMESPACE" ]]; then
  echo "Valid"
fi

# Arithmetic
result=$((1 + 2))

# Logical operators
[[ -f file.txt ]] && echo "exists" || echo "not found"`,
			expectError: false,
		},
		{
			name:     "underscore in placeholder name",
			template: "Value: @@MY_VAR@@",
			values: map[string]string{
				"MY_VAR": "test",
			},
			expected:    "Value: test",
			expectError: false,
		},
		{
			name:     "numbers in placeholder name",
			template: "Value: @@VAR123@@",
			values: map[string]string{
				"VAR123": "test",
			},
			expected:    "Value: test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderTemplate(tt.template, tt.values)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected:\n%s\n\ngot:\n%s", tt.expected, result)
				}
			}
		})
	}
}
