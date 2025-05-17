# Skaffold: Continuous Development for Kubernetes

**Streamlining the build-push-deploy workflow for Kubernetes applications**

---

## What is Skaffold?

Skaffold is an open-source CLI tool that automates the workflow for building, pushing, and deploying Kubernetes applications. It enables rapid, repeatable development cycles by handling the tedious parts of Kubernetes app development.

---

## Development Workflow

Skaffold manages the entire development loop:

1. **Code Changes** (edit your source code)
2. **Build Images** (containerize your app)
3. **Push to Registry** (upload images)
4. **Deploy to Kubernetes** (apply manifests)
5. **Repeat** (Skaffold watches for changes and restarts the loop)

> _Skaffold automatically detects changes and redeploys, so you can focus on coding._

---

## Key Features

- ⚡ **Fast Local Development:** Optimized source-to-deploy workflow
- 🔄 **File Synchronization:** Sync files directly to containers without full rebuilds
- 🧩 **Multi-Component Support:** Manage complex, multi-container applications
- 💻 **Client-Side Only:** No cluster-side components required
- 🔌 **Pluggable Architecture:** Supports various build and deploy tools (Docker, Jib, Helm, Kustomize, etc.)
- 📄 **Configuration as Code:** Declarative YAML-based configuration

---

## Architecture Overview

```
+-------------------+
|   CLI & Config    |
+-------------------+
         |
+-------------------------------+
| Builders | Renderers | Deployers |
+-------------------------------+
         |
+-------------------+
| Kubernetes Cluster |
+-------------------+
```

_Modular design with pluggable components for flexibility and extensibility._

---

## Developer Benefits

- ⏱️ **Reduced Feedback Loop:** See changes live in seconds
- 🔄 **Consistent Environments:** Same workflow for local and CI/CD
- 🛠️ **Tool Integration:** Works with popular build and deployment tools
- 📊 **Production-like Testing:** Test locally with Kubernetes before shipping

---

## Resources

- **GitHub:** [github.com/GoogleContainerTools/skaffold](https://github.com/GoogleContainerTools/skaffold)
- **Documentation:** [skaffold.dev](https://skaffold.dev)

---