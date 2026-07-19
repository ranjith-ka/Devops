# .NET & C# for DevOps (beginner)

No prior .NET or C# required. This tutorial maps Microsoft’s toolchain to concepts you already know (Go, Python, npm, Linux/CLI) so you can open solutions in **Rider**, run `dotnet` commands, and read the small amount of C# used in our labs.

**Audience:** DevOps / platform engineers entering .NET  
**IDE:** JetBrains Rider (optional but recommended)  
**Next lab:** [Reqnroll Quickstart (Rider)](./reqnroll/REQNROLL_QUICKSTART.md)

**Time:** ~20–30 minutes

---

## What you will learn

1. What .NET, the SDK, the runtime, and C# are
2. Solution / project / `.csproj` / NuGet layout
3. Everyday `dotnet` CLI commands (and `global.json` traps)
4. Just enough C# to read and edit lab code
5. Essential Rider shortcuts

---

## Table of contents

1. [Prerequisites](#1-prerequisites)
2. [What is .NET?](#2-what-is-net)
3. [Solution, project, and files](#3-solution-project-and-files)
4. [Install the SDK](#4-install-the-sdk)
5. [CLI cheat sheet](#5-cli-cheat-sheet)
6. [`global.json` and SDK selection](#6-globaljson-and-sdk-selection)
7. [C# you need for our labs](#7-c-you-need-for-our-labs)
8. [Rider shortcuts](#8-rider-shortcuts)
9. [Tiny hands-on: create and run a console app](#9-tiny-hands-on-create-and-run-a-console-app)
10. [Where to go next](#10-where-to-go-next)
11. [Troubleshooting](#11-troubleshooting)

---

## 1. Prerequisites

| Tool | Why | Check |
|------|-----|--------|
| Terminal + Git | Clone / commands | `git --version` |
| [.NET 8 SDK](https://dotnet.microsoft.com/download) | Build & run | Install in [§4](#4-install-the-sdk) |
| [JetBrains Rider](https://www.jetbrains.com/rider/) (optional) | IDE for .NET/C# | — |

Assumed background: Linux shells, Git, and any mainstream language (Python, Go, Java, or JS).

---

## 2. What is .NET?

**.NET** is Microsoft’s platform for building and running applications. Most application code is written in **C#**.

| Piece | What it is | Closest to |
|-------|------------|------------|
| **.NET SDK** | Compiler, templates, **`dotnet` CLI** | `go` toolchain, or Python + pip + build tools |
| **.NET Runtime** | Runs compiled apps | JVM / Node |
| **C#** | Main language | Typed/compiled like Java or Go |
| **NuGet** | Package registry + client | npm / PyPI |
| **`dotnet` CLI** | restore / build / test / run / publish | `go` / `mvn` / `npm` |

Install the **SDK** on your laptop (it includes what you need to develop). In CI, use an SDK container image such as `mcr.microsoft.com/dotnet/sdk:8.0`.

**Version naming:** `net8.0` in a project file means “target .NET 8”. That is the *framework*, not the exact SDK patch (`8.0.422`).

---

## 3. Solution, project, and files

```text
MyApp.sln                         ← SOLUTION — open this in Rider
├── MyApp/                        ← PROJECT — production code
│   ├── MyApp.csproj              ← deps + target framework (≈ package.json / go.mod)
│   └── Program.cs                ← C# source
└── MyApp.Tests/                  ← PROJECT — tests (optional second project)
    ├── MyApp.Tests.csproj
    └── ...
```

Example from the Reqnroll lab:

```text
ReqnrollQuickstart.sln
├── ReqnrollQuickstart.App/       ← production
└── ReqnrollQuickstart.Specs/     ← BDD tests
```

| File / concept | Role | Analogy |
|----------------|------|---------|
| `.sln` | Lists related projects | Monorepo / multi-module workspace root |
| `.csproj` | Dependencies, `TargetFramework` | `package.json` / `go.mod` / `pom.xml` |
| `.cs` | Source | `.go` / `.py` / `.java` |
| `.dll` | Build output (assembly) | `.jar` / compiled package |
| NuGet package | External library | npm / PyPI package |

**Always open the `.sln` in Rider** when one exists — not a random folder — so build and Unit Tests discovery work correctly.

---

## 4. Install the SDK

### macOS (Homebrew)

```bash
brew install --cask dotnet-sdk
dotnet --list-sdks
dotnet --info
```

You want at least one **8.0.x** SDK listed for our labs (`TargetFramework` = `net8.0`).

### Verify

```bash
dotnet --version
# e.g. 8.0.422
```

---

## 5. CLI cheat sheet

| Command | Does | When |
|---------|------|------|
| `dotnet --list-sdks` | Show installed SDKs | Debug version issues |
| `dotnet restore` | Download NuGet packages | After clone / package change |
| `dotnet build` | Compile → `.dll` under `bin/` | After code changes |
| `dotnet test` | Build (if needed) + run tests | Local / CI |
| `dotnet run` | Build + run the startup project | Console / web apps |
| `dotnet new console -n Hello` | Scaffold a tiny app | Learning / smoke test |

Rider can Build and Run/Test from UI. Learn the CLI anyway — **CI uses these same commands**.

Typical pipeline sequence:

```bash
dotnet restore MyApp.sln
dotnet build MyApp.sln --configuration Release --no-restore
dotnet test MyApp.sln --configuration Release --no-build
```

---

## 6. `global.json` and SDK selection

`dotnet` walks up from the current directory looking for `global.json`. A pin there can **override** which SDK is used — even in unrelated repos under your home folder if you keep `~/global.json`.

Example problem:

```text
Requested SDK version: 10.0.107
A compatible .NET SDK was not found.
```

You might have `8.0.422` and `10.0.301` installed, but not exactly `10.0.107`.

**Fix** — point at an installed SDK and allow roll-forward:

```json
{
  "sdk": {
    "version": "8.0.422",
    "rollForward": "latestFeature"
  }
}
```

Use a version from *your* `dotnet --list-sdks`. For net8.0 labs, pinning **8.0.x** is enough.

| Setting | Meaning |
|---------|---------|
| `sdk.version` | Preferred SDK |
| `rollForward` | What to do if exact version missing (`latestFeature`, `latestMajor`, …) |

---

## 7. C# you need for our labs

You only need to **recognize and lightly edit** these patterns. You do not need to master C# before the Reqnroll lab.

### 7.1 File shape

```csharp
namespace MyCompany.MyApp;     // package-like grouping

public class PriceCalculator  // type: data + methods
{
    // fields and methods
}
```

| Keyword | Meaning |
|---------|---------|
| `namespace` | Logical package path |
| `class` | Type / blueprint |
| `public` | Visible to other projects |
| `private` | Only inside this class |
| `readonly` | Assigned once, then not reassigned |

### 7.2 Common types

| C# | Meaning | Analogy |
|----|---------|---------|
| `string` | Text | Python `str` |
| `int` | Whole number | `int` |
| `decimal` | Precise number (money) | Prefer over `float`/`double` for cash |
| `Dictionary<K,V>` | Map | Go `map`, Python `dict` |
| `void` | Returns nothing | |
| `var` | Infer type | Go `:=` / local type inference |

`180.0m` — the **`m`** suffix marks a **`decimal`** literal (not `double`).

### 7.3 Fields vs methods

```csharp
private readonly Dictionary<string, int> _basket = new();  // FIELD — object state

public void Reset()                                        // METHOD — does work
{
    _basket.Clear();                                       // call with .
}
```

Leading `_` on private fields is a common C# convention.

### 7.4 `new` — create an instance

```csharp
private readonly PriceCalculator _calc = new();
```

Same idea as constructing an object in Go/Python/Java.

### 7.5 Attributes — `[brackets]`

```csharp
[Binding]                                    // framework metadata
public class Steps
{
    [Given("the client started shopping")]   // links method to a Gherkin step
    public void GivenStarted() { }
}
```

Like Java annotations / Python decorators. Test and BDD frameworks discover methods from attributes.

### 7.6 Exceptions

```csharp
throw new NotImplementedException();   // “not written yet”
```

Uncaught exceptions fail the test (red in Rider).

### 7.7 Assertions

```csharp
Assert.AreEqual(expected, actual);
```

Same idea as `assert` / testify `require.Equal`.

### 7.8 Imports (`using`)

If Rider underlines a type in red: `⌥⏎` (macOS) / `Alt+Enter` → **import …**. You get:

```csharp
using Microsoft.VisualStudio.TestTools.UnitTesting;
```

Do not memorize namespaces — use the IDE.

### 7.9 Simple loop

```csharp
foreach (var item in basket)
{
    price += _priceTable[item.Key] * item.Value;
}
```

---

## 8. Rider shortcuts

| Action | macOS | Windows/Linux |
|--------|-------|----------------|
| Open solution | **File → Open** → `.sln` | same |
| Build | `⌘F9` | `Ctrl+F9` |
| Unit Tests | **View → Tool Windows → Unit Tests** | same |
| Quick-fix (import, Create Step, …) | `⌥⏎` | `Alt+Enter` |
| Go to definition | `⌘B` / `⌘Click` | `Ctrl+B` / `Ctrl+Click` |
| Find usages | `⌥F7` | `Alt+F7` |
| Terminal | **View → Tool Windows → Terminal** | same |

Rider feels like IntelliJ/WebStorm with a .NET toolchain.

---

## 9. Tiny hands-on: create and run a console app

Verify your SDK and mental model with a 2-minute app:

```bash
mkdir -p ~/code/tmp && cd ~/code/tmp
dotnet new console -n HelloDotnet
cd HelloDotnet
dotnet run
```

Expected output includes `Hello, World!`.

Optional in Rider: **File → Open** → `HelloDotnet.csproj` (or create a solution later). Edit `Program.cs`, Build, Run.

Clean up when done:

```bash
cd .. && rm -rf HelloDotnet
```

---

## 10. Where to go next

| Path | Description |
|------|-------------|
| **[Reqnroll Quickstart (Rider)](./reqnroll/REQNROLL_QUICKSTART.md)** | BDD lab: Gherkin → step definitions → price calculator |
| [Microsoft C# tour](https://learn.microsoft.com/dotnet/csharp/tour-of-csharp/) | Broader language intro (optional) |
| [dotnet CLI docs](https://learn.microsoft.com/dotnet/core/tools/) | Full command reference |

---

## 11. Troubleshooting

| Problem | Cause | Fix |
|---------|-------|-----|
| `Compatible .NET SDK was not found` | Pin in `global.json` | Match `dotnet --list-sdks`; see [§6](#6-globaljson-and-sdk-selection) |
| `restore`/`build` “application does not exist” | SDK resolution failed | Same as above — not a bad command spelling |
| Red squiggles in Rider | Missing import / wrong SDK | `⌥⏎` import; install/select net8 SDK |
| Opened folder but no tests/projects | Didn’t open `.sln` | **File → Open** the solution file |
| Confused by C# syntax | — | Re-read [§7](#7-c-you-need-for-our-labs); copy lab snippets as-is |

---

## One-page recap

```text
SDK  = toolchain (dotnet CLI + compiler)
C#   = language
.sln = workspace of projects
.csproj ≈ go.mod / package.json
NuGet ≈ npm

dotnet restore → build → test|run

Read C#: class, fields, methods, [attributes], Dictionary, decimal (m), Assert
Rider: open .sln, ⌘F9 build, ⌥⏎ quick-fix
```

Then continue with the [Reqnroll lab](./reqnroll/REQNROLL_QUICKSTART.md).
