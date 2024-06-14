package internal_test

import (
	"testing"

	"github.com/jonas27/kubectl-lsd/internal"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name  string
		stdin []byte
		want  string
	}{
		{"decode single yaml", []byte(mockYAML), mockYAMLDecoded},
		{"decode single json", []byte(mockJSON), mockJSONDecoded},
		{"decode list yaml", []byte(mockYAMLList), mockYAMLListDecoded},
		{"decode list json", []byte(mockJSONList), mockJSONListDecoded},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out, err := internal.Run(tt.stdin)
			require.NoError(t, err)
			require.Equal(t, tt.want, out)
		})
	}
}

const (
	mockYAML = `apiVersion: v1
data:
  password: c2VjcmV0
  app: bGlzdCBzZWNyZXQgZGVjb2Rl
kind: Secret
metadata:
  name: "list secret decode"
  namespace: lsd
type: Opaque`
	mockYAMLDecoded = `apiVersion: v1
kind: Secret
metadata:
  name: list secret decode
  namespace: lsd
stringData:
  app: list secret decode
  password: secret
type: Opaque
`
	mockJSON = `{
    "apiVersion": "v1",
    "data": {
        "password": "c2VjcmV0",
        "app": "bGlzdCBzZWNyZXQgZGVjb2Rl"
    },
    "kind": "Secret",
    "metadata": {
        "name": "list secret decode",
        "namespace": "lsd"
    },
    "type": "Opaque"
}`
	mockJSONDecoded = `{
    "apiVersion": "v1",
    "kind": "Secret",
    "metadata": {
        "name": "list secret decode",
        "namespace": "lsd"
    },
    "stringData": {
        "app": "list secret decode",
        "password": "secret"
    },
    "type": "Opaque"
}`
	mockYAMLList = `apiVersion: v1
items:
- data:
    password: c2VjcmV0
    app: bGlzdCBzZWNyZXQgZGVjb2Rl
- data:
    password: c2VjcmV0
    app: bGlzdCBzZWNyZXQgZGVjb2Rl
`
	mockYAMLListDecoded = `apiVersion: v1
items:
- stringData:
    app: list secret decode
    password: secret
- stringData:
    app: list secret decode
    password: secret
`
	mockJSONList = `{
  "apiVersion": "v1",
  "items": [
    {
      "data": {
        "password": "c2VjcmV0",
        "app": "bGlzdCBzZWNyZXQgZGVjb2Rl"
      }
    },
    {
      "data": {
        "password": "c2VjcmV0",
        "app": "bGlzdCBzZWNyZXQgZGVjb2Rl"
      }
    }
  ]
}`
	mockJSONListDecoded = `{
    "apiVersion": "v1",
    "items": [
        {
            "stringData": {
                "app": "list secret decode",
                "password": "secret"
            }
        },
        {
            "stringData": {
                "app": "list secret decode",
                "password": "secret"
            }
        }
    ]
}`
)
