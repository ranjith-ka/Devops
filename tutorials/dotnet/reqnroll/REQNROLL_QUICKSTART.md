# Reqnroll Quickstart — DevOps + Rider

Hands-on BDD with **Reqnroll** (SpecFlow’s successor) in **JetBrains Rider**: Gherkin feature → C# step definitions → green price-calculator tests.

**If you don’t know .NET/C# yet:** complete [**.NET & C# for DevOps**](../DOTNET_FOR_DEVOPS.md) first (~20–30 min), then return here.

**Links:** [Official Quickstart](https://docs.reqnroll.net/latest/quickstart/index.html) · [Setup Rider](https://docs.reqnroll.net/latest/installation/setup-ide.html#setup-rider) · [Reqnroll for Rider](https://plugins.jetbrains.com/plugin/24012-reqnroll-for-rider) · [Starter repo](https://github.com/reqnroll/Quickstart) · [`completed` branch](https://github.com/reqnroll/Quickstart/tree/completed)

**Time:** ~45–60 minutes

---

## What you will learn

1. Install **Reqnroll for Rider** and open the Quickstart solution
2. Run scenarios in the Unit Tests window
3. Create step definitions (Create Step / hand-written stubs)
4. Implement `PriceCalculator` test-first (simple basket → multi-item DataTable)
5. Run the same checks with `dotnet test` for CI

---

## Table of contents

1. [Prerequisites](#1-prerequisites)
2. [What is Reqnroll?](#2-what-is-reqnroll)
3. [Step 0 — SDK check](#3-step-0--sdk-check)
4. [Step 1 — Install Reqnroll for Rider](#4-step-1--install-reqnroll-for-rider)
5. [Step 2 — Clone and open in Rider](#5-step-2--clone-and-open-in-rider)
6. [Step 3 — Run tests (expect “undefined”)](#6-step-3--run-tests-expect-undefined)
7. [Step 4 — Create step definitions](#7-step-4--create-step-definitions)
8. [Step 5 — Add state fields](#8-step-5--add-state-fields)
9. [Step 6 — Automate the steps](#9-step-6--automate-the-steps)
10. [Step 7 — Implement the app (minimal)](#10-step-7--implement-the-app-minimal)
11. [Step 8 — Multi-item scenario (DataTable)](#11-step-8--multi-item-scenario-datatable)
12. [Step 9 — Finish `CalculatePrice`](#12-step-9--finish-calculateprice)
13. [CI checklist](#13-ci-checklist)
14. [Next steps](#14-next-steps)
15. [Troubleshooting](#15-troubleshooting)

---

## 1. Prerequisites

| Tool | Why | Check |
|------|-----|--------|
| [.NET 8 SDK](https://dotnet.microsoft.com/download) | Quickstart is `net8.0` | `dotnet --list-sdks` shows `8.0.x` — see [dotnet tutorial](../DOTNET_FOR_DEVOPS.md#4-install-the-sdk) |
| [JetBrains Rider](https://www.jetbrains.com/rider/) | IDE + Unit Tests | Compatible with [Reqnroll plugin versions](https://plugins.jetbrains.com/plugin/24012-reqnroll-for-rider/versions) |
| Git | Clone | `git --version` |
| [.NET & C# for DevOps](../DOTNET_FOR_DEVOPS.md) | Vocab + C# patterns | Skim §§2–8 if new to .NET |

Disable **SpecFlow for Rider** if installed — it conflicts with Reqnroll on `.feature` files.

Quick reminder of layout (details in the [dotnet tutorial](../DOTNET_FOR_DEVOPS.md#3-solution-project-and-files)):

```text
ReqnrollQuickstart.sln
├── ReqnrollQuickstart.App/       ← PriceCalculator (production)
└── ReqnrollQuickstart.Specs/     ← Features + StepDefinitions (tests)
```

---

## 2. What is Reqnroll?

**Reqnroll** = BDD for .NET. Flow:

```text
PriceCalculation.feature  →  StepDefinitions/*.cs  →  PriceCalculator.cs
(Gherkin / English)            (glue)                   (business logic)
```

| Term | Meaning |
|------|---------|
| **Feature** | Area under test |
| **Scenario** | One example (one test) |
| **Given / When / Then** | Arrange / Act / Assert |
| **Step definition** | C# method bound to a step line |
| **Undefined** | No C# binding yet |
| **Pending** | Stub throws `PendingStepException` |
| **DataTable** | Table of data under a step |

---

## 3. Step 0 — SDK check

```bash
dotnet --list-sdks    # need 8.0.x
```

If you see *Requested SDK version … was not found*, fix `~/global.json` — [dotnet tutorial §6](../DOTNET_FOR_DEVOPS.md#6-globaljson-and-sdk-selection).

```bash
# macOS install if needed
brew install --cask dotnet-sdk
```

---

## 4. Step 1 — Install Reqnroll for Rider

1. Rider → **Plugins** → **Marketplace** → **Reqnroll for Rider** → Install
2. Disable **SpecFlow for Rider**
3. Restart Rider

You will use: Gherkin highlighting, **Create Step** (`⌥⏎` / `Alt+Enter`), Go to Declaration, **Unit Tests** window.

---

## 5. Step 2 — Clone and open in Rider

```bash
cd ~/code/github/ranjith   # your path
git clone https://github.com/reqnroll/Quickstart.git
cd Quickstart
```

1. **File → Open…** → **`ReqnrollQuickstart.sln`**
2. Wait for NuGet restore, or:

```bash
dotnet restore ReqnrollQuickstart.sln
dotnet build ReqnrollQuickstart.sln
```

### Rider plugin: include features as Content

In `ReqnrollQuickstart.Specs/ReqnrollQuickstart.Specs.csproj`, ensure:

```xml
<ItemGroup>
  <Content Include="**/*.feature" />
</ItemGroup>
```

Restart Rider if `.feature` files look like plain text. Optional: right-click → **Associate with File Type…** → **Reqnroll file**.

### Starting files

**App** (`ReqnrollQuickstart.App/PriceCalculator.cs`) — unfinished:

```csharp
namespace ReqnrollQuickstart.App;

public class PriceCalculator
{
    private readonly Dictionary<string, decimal> _priceTable = new()
    {
        { "Electric guitar", 180.0 },
        { "Guitar pick", 1.5 }
    };

    public decimal CalculatePrice(Dictionary<string, int> basket)
    {
        throw new NotImplementedException();
    }
}
```

**Feature** (`…/Features/PriceCalculation.feature`):

```gherkin
Feature: Price calculation

Rule: The price for a basket with items can be calculated based on the item prices

Scenario: Client has a simple basket
    Given the client started shopping
    And the client added 1 pcs of "Electric guitar" to the basket
    When the basket is prepared
    Then the basket price should be $180
```

---

## 6. Step 3 — Run tests (expect “undefined”)

1. **View → Tool Windows → Unit Tests**
2. Build (`⌘F9` / `Ctrl+F9`) if needed
3. Run all Specs tests

Or:

```bash
dotnet test ReqnrollQuickstart.Specs/ReqnrollQuickstart.Specs.csproj --logger "console;verbosity=detailed"
```

**Expected:** scenario found, steps **undefined**.

| State | Meaning | Action |
|-------|---------|--------|
| Undefined | No binding | Create Step |
| Pending | Stub only | Fill method body |
| Failed | Wrong result / exception | Fix glue or app |
| Passed | OK | Continue |

---

## 7. Step 4 — Create step definitions

A step definition is a C# method whose `[Given]` / `[When]` / `[Then]` string matches the feature line. `{int}` / `{string}` capture parameters.

### Rider (preferred)

1. Open `PriceCalculation.feature`
2. Caret on undefined step → `⌥⏎` / `Alt+Enter` → **Create Step**
3. Class under `StepDefinitions/`, e.g. `PriceCalculationStepDefinitions.cs`
4. Rename params to `quantity`, `product`, `expectedPrice` if needed

### Hand-written fallback

`ReqnrollQuickstart.Specs/StepDefinitions/PriceCalculationStepDefinitions.cs`:

```csharp
namespace ReqnrollQuickstart.Specs.StepDefinitions;

[Binding]
public class PriceCalculationStepDefinitions
{
    [Given("the client started shopping")]
    public void GivenTheClientStartedShopping()
    {
        throw new PendingStepException();
    }

    [Given("the client added {int} pcs of {string} to the basket")]
    public void GivenTheClientAddedPcsOfToTheBasket(int quantity, string product)
    {
        throw new PendingStepException();
    }

    [When("the basket is prepared")]
    public void WhenTheBasketIsPrepared()
    {
        throw new PendingStepException();
    }

    [Then("the basket price should be ${float}")]
    public void ThenTheBasketPriceShouldBe(decimal expectedPrice)
    {
        throw new PendingStepException();
    }
}
```

Build (`⌘F9`), re-run → **pending**. (`[Binding]` / attributes explained in [dotnet tutorial §7.5](../DOTNET_FOR_DEVOPS.md#75-attributes--brackets).)

---

## 8. Step 5 — Add state fields

At the top of the class (above the methods):

```csharp
[Binding]
public class PriceCalculationStepDefinitions
{
    private readonly PriceCalculator _priceCalculator = new();
    private readonly Dictionary<string, int> _basket = new();
    private decimal _calculatedPrice;

    // ... methods below ...
}
```

`PriceCalculator` red? `⌥⏎` → import `ReqnrollQuickstart.App`.

| Field | Job |
|-------|-----|
| `_priceCalculator` | Object under test |
| `_basket` | Product → quantity |
| `_calculatedPrice` | Set in When, asserted in Then |

---

## 9. Step 6 — Automate the steps

Replace each `PendingStepException` with:

```csharp
namespace ReqnrollQuickstart.Specs.StepDefinitions;

[Binding]
public class PriceCalculationStepDefinitions
{
    private readonly PriceCalculator _priceCalculator = new();
    private readonly Dictionary<string, int> _basket = new();
    private decimal _calculatedPrice;

    [Given("the client started shopping")]
    public void GivenTheClientStartedShopping()
    {
        _basket.Clear();
        _calculatedPrice = 0.0m;
    }

    [Given("the client added {int} pcs of {string} to the basket")]
    public void GivenTheClientAddedPcsOfToTheBasket(int quantity, string product)
    {
        _basket.Add(product, quantity);
    }

    [When("the basket is prepared")]
    public void WhenTheBasketIsPrepared()
    {
        _calculatedPrice = _priceCalculator.CalculatePrice(_basket);
    }

    [Then("the basket price should be ${float}")]
    public void ThenTheBasketPriceShouldBe(decimal expectedPrice)
    {
        Assert.AreEqual(expectedPrice, _calculatedPrice);
    }
}
```

| Step | C# does |
|------|---------|
| Given started | Reset state |
| Given added N of “X” | Fill `_basket` |
| When prepared | Call `CalculatePrice` |
| Then price $… | `Assert.AreEqual` |

`Assert` red? import MSTest via quick-fix.

---

## 10. Step 7 — Implement the app (minimal)

Re-run tests → `NotImplementedException` (glue works; app empty).

Minimal implementation for scenario 1 only:

```csharp
namespace ReqnrollQuickstart.App;

public class PriceCalculator
{
    private readonly Dictionary<string, decimal> _priceTable = new()
    {
        { "Electric guitar", 180.0m },
        { "Guitar pick", 1.5m }
    };

    public decimal CalculatePrice(Dictionary<string, int> basket)
    {
        var item = basket.First();
        return _priceTable[item.Key];
    }
}
```

(`m` = decimal literal — [dotnet tutorial §7.2](../DOTNET_FOR_DEVOPS.md#72-common-types).)

Run tests → **Client has a simple basket** green.

---

## 11. Step 8 — Multi-item scenario (DataTable)

Add under the same `Rule` in the feature file:

```gherkin
Scenario: Client has multiple items in their basket
    Given the client started shopping
    And the client added
        | product         | quantity |
        | Electric guitar | 1        |
        | Guitar pick     | 10       |
    When the basket is prepared
    Then the basket price should be $195.0
```

\(180 + 10×1.5 = 195\). Create Step for `And the client added`, or paste:

```csharp
[Given("the client added")]
public void GivenTheClientAdded(DataTable itemsTable)
{
    var items = itemsTable.CreateSet<(string Product, int Quantity)>();
    foreach (var item in items)
    {
        _basket.Add(item.Product, item.Quantity);
    }
}
```

Run: scenario 1 green; scenario 2 fails (~195 vs 180) — expected.

---

## 12. Step 9 — Finish `CalculatePrice`

```csharp
public decimal CalculatePrice(Dictionary<string, int> basket)
{
    decimal price = 0;
    foreach (var item in basket)
    {
        price += _priceTable[item.Key] * item.Value;
    }
    return price;
}
```

Unit Tests / `dotnet test ReqnrollQuickstart.sln` → **both green**.

---

## 13. CI checklist

```bash
dotnet restore ReqnrollQuickstart.sln
dotnet build ReqnrollQuickstart.sln --no-restore --configuration Release
dotnet test ReqnrollQuickstart.sln --no-build --configuration Release \
  --logger "trx;LogFileName=reqnroll.trx"
```

```yaml
- uses: actions/setup-dotnet@v4
  with:
    dotnet-version: '8.0.x'
- run: dotnet test ReqnrollQuickstart.sln --configuration Release
```

---

## 14. Next steps

From the [official Next Steps](https://docs.reqnroll.net/latest/quickstart/index.html): currency transforms, Background catalog, discount `Rule`, REST API + hooks.

Stuck? Compare with [`completed`](https://github.com/reqnroll/Quickstart/tree/completed).

---

## 15. Troubleshooting

| Problem | Fix |
|---------|-----|
| SDK / `global.json` errors | [dotnet tutorial §6](../DOTNET_FOR_DEVOPS.md#6-globaljson-and-sdk-selection) |
| No Create Step / no coloring | Install Reqnroll for Rider; disable SpecFlow; restart |
| `.feature` plain text | Associate → Reqnroll file; add `Content Include="**/*.feature"` |
| Tests missing | Open `.sln`, Build, refresh Unit Tests |
| Still Undefined | Match step text exactly; rebuild |
| Red types (`Assert`, …) | `⌥⏎` import |
| C# confusing | Re-read [dotnet tutorial §7](../DOTNET_FOR_DEVOPS.md#7-c-you-need-for-our-labs) |

---

## One-page flow

```text
[.NET for DevOps primer] → SDK 8 + Reqnroll plugin
  → Open ReqnrollQuickstart.sln
  → Unit Tests (Undefined) → Create Steps (Pending)
  → Fill glue → NotImplemented → minimal CalculatePrice → green #1
  → DataTable scenario → full CalculatePrice → both green
  → CI: dotnet test
```
