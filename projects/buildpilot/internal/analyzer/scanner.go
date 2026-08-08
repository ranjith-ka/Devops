package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ranjith-ka/buildpilot/internal/domain"
)

const maxFiles = 2000

var ignoredDirectories = map[string]bool{
	".git": true, ".idea": true, ".vscode": true,
	"bin": true, "obj": true, "node_modules": true, "vendor": true,
}

func Scan(root string) (domain.RepositoryFacts, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return domain.RepositoryFacts{}, fmt.Errorf("resolve repository path: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return domain.RepositoryFacts{}, fmt.Errorf("inspect repository path: %w", err)
	}
	if !info.IsDir() {
		return domain.RepositoryFacts{}, fmt.Errorf("repository path is not a directory")
	}

	facts := domain.RepositoryFacts{Root: abs}
	languages := map[string]bool{}
	err = filepath.WalkDir(abs, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && path != abs && ignoredDirectories[entry.Name()] {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		if len(facts.Files) >= maxFiles {
			return nil
		}
		relative, err := filepath.Rel(abs, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		facts.Files = append(facts.Files, relative)
		name := strings.ToLower(entry.Name())
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".go":
			languages["go"] = true
		case ".cs", ".fs", ".vb":
			languages["dotnet"] = true
		case ".js", ".jsx", ".ts", ".tsx":
			languages["javascript"] = true
		case ".py":
			languages["python"] = true
		}
		if ext == ".csproj" || ext == ".fsproj" || ext == ".vbproj" || name == "go.mod" || name == "package.json" {
			facts.ProjectFiles = append(facts.ProjectFiles, relative)
		}
		if name == "dockerfile" || strings.HasPrefix(name, "dockerfile.") {
			facts.Dockerfiles = append(facts.Dockerfiles, relative)
		}
		if (ext == ".yaml" || ext == ".yml") && (strings.Contains(relative, "k8s/") || strings.Contains(relative, "kubernetes/") || strings.Contains(relative, "charts/")) {
			facts.KubernetesFiles = append(facts.KubernetesFiles, relative)
		}
		if name == "skaffold.yaml" || name == "skaffold.yml" {
			facts.HasSkaffold = true
		}
		return nil
	})
	if err != nil {
		return domain.RepositoryFacts{}, fmt.Errorf("scan repository: %w", err)
	}
	for language := range languages {
		facts.Languages = append(facts.Languages, language)
	}
	sort.Strings(facts.Files)
	sort.Strings(facts.Languages)
	sort.Strings(facts.ProjectFiles)
	sort.Strings(facts.Dockerfiles)
	sort.Strings(facts.KubernetesFiles)
	return facts, nil
}
