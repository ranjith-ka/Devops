package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanDetectsProjectAndIgnoresBuildOutput(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/App.csproj", "<Project />")
	writeTestFile(t, root, "src/Program.cs", "")
	writeTestFile(t, root, "Dockerfile", "FROM scratch")
	writeTestFile(t, root, "k8s/deployment.yaml", "kind: Deployment")
	writeTestFile(t, root, "obj/project.assets.json", "ignored")

	facts, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if !reflect.DeepEqual(facts.Languages, []string{"dotnet"}) {
		t.Fatalf("languages = %v", facts.Languages)
	}
	if !reflect.DeepEqual(facts.ProjectFiles, []string{"src/App.csproj"}) {
		t.Fatalf("project files = %v", facts.ProjectFiles)
	}
	if !reflect.DeepEqual(facts.Dockerfiles, []string{"Dockerfile"}) {
		t.Fatalf("dockerfiles = %v", facts.Dockerfiles)
	}
	for _, file := range facts.Files {
		if file == "obj/project.assets.json" {
			t.Fatal("ignored obj directory was scanned")
		}
	}
}

func writeTestFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
