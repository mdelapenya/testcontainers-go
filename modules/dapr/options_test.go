package dapr

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestComponent_Render(t *testing.T) {
	yamlFilePath := filepath.Join("testdata", "component-statestore-state.in-memory.yaml")
	yamlFile, err := os.Open(yamlFilePath)
	if err != nil {
		t.Fatal(err)
	}

	content, err := io.ReadAll(yamlFile)
	if err != nil {
		t.Fatal(err)
	}

	component := NewComponent("statestore", "state.in-memory", map[string]string{
		"foo1": "bar1",
		"foo2": "bar2",
	})

	rendered, err := component.Render()
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != string(rendered) {
		t.Fatalf("expected %s, got %s", string(content), string(rendered))
	}
}
