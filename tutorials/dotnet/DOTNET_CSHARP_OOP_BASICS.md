# C# Basics and OOP for DevOps — From Syntax to Encapsulation

This is the next chapter after [.NET & C# for DevOps (beginner)](https://github.com/ranjith-ka/Devops/blob/main/tutorials/dotnet/DOTNET_FOR_DEVOPS.md). It moves from reading small C# snippets to writing a simple object-oriented program.

**Audience:** DevOps / platform engineers learning C#  
**Prerequisite:** Complete `DOTNET_FOR_DEVOPS.md` first  
**Target framework:** .NET 8 (`net8.0`)  
**Time:** ~60–90 minutes

---

## What you will learn

1. How a C# program is structured
2. Variables, types, operators, conditions, and loops
3. How to declare and call methods
4. The difference between a class and an object
5. Fields, properties, and constructors
6. Reference types and object state
7. Access modifiers
8. Encapsulation and why it matters

By the end, you will build a small `Deployment` class that protects its own state.

---

## Table of contents

1. [Create the learning project](#1-create-the-learning-project)
2. [Understand basic C# syntax](#2-understand-basic-c-syntax)
3. [Variables and common types](#3-variables-and-common-types)
4. [Operators and string interpolation](#4-operators-and-string-interpolation)
5. [Conditions](#5-conditions)
6. [Loops and collections](#6-loops-and-collections)
7. [Methods](#7-methods)
8. [Classes and objects](#8-classes-and-objects)
9. [Fields, properties, and constructors](#9-fields-properties-and-constructors)
10. [Object state and reference types](#10-object-state-and-reference-types)
11. [Access modifiers](#11-access-modifiers)
12. [Encapsulation](#12-encapsulation)
13. [Hands-on lab: encapsulated deployment](#13-hands-on-lab-encapsulated-deployment)
14. [Exercises](#14-exercises)
15. [Recap](#15-recap)

---

## 1. Create the learning project

```bash
mkdir CSharpOopBasics
cd CSharpOopBasics
dotnet new console
dotnet run
```

Expected output:

```text
Hello, World!
```

The important files are:

```text
CSharpOopBasics/
├── CSharpOopBasics.csproj   ← project configuration
├── Program.cs               ← application entry point
├── bin/                     ← compiled output
└── obj/                     ← temporary build files
```

You can run the project after each example with:

```bash
dotnet run
```

---

## 2. Understand basic C# syntax

Replace `Program.cs` with:

```csharp
string serviceName = "payments-api";
int replicaCount = 3;

Console.WriteLine($"Deploying {serviceName} with {replicaCount} replicas");
```

Notice these syntax rules:

| Syntax | Meaning |
|--------|---------|
| `string serviceName` | Declare a variable and its type |
| `=` | Assign a value |
| `;` | End a statement |
| `Console.WriteLine(...)` | Call a method |
| `()` | Supply arguments to a method |
| `{}` | Group statements into a block |
| `// text` | Single-line comment |
| `$"...{value}..."` | Insert values into a string |

C# is case-sensitive. `serviceName`, `ServiceName`, and `servicename` are different names.

### Naming conventions

| Item | Convention | Example |
|------|------------|---------|
| Local variable / parameter | `camelCase` | `replicaCount` |
| Method / class / property | `PascalCase` | `StartDeployment` |
| Private field | `_camelCase` | `_status` |
| Constant | `PascalCase` | `MaximumRetries` |

---

## 3. Variables and common types

```csharp
string environment = "production";
int replicaCount = 3;
bool isHealthy = true;
double cpuPercent = 72.5;
decimal monthlyCost = 125.50m;
char deploymentGrade = 'A';
DateTime deployedAt = DateTime.UtcNow;
```

| Type | Stores | Example |
|------|--------|---------|
| `string` | Text | `"production"` |
| `int` | Whole numbers | `3` |
| `bool` | `true` or `false` | `true` |
| `double` | General decimal values | `72.5` |
| `decimal` | Precise decimal values, such as money | `125.50m` |
| `char` | One character | `'A'` |
| `DateTime` | Date and time | `DateTime.UtcNow` |

Use `var` when the value makes the type obvious:

```csharp
var region = "centralindia"; // compiler infers string
var retries = 3;             // compiler infers int
```

`var` does not make C# dynamically typed. The compiler still fixes the type at compile time.

### Constants

Use `const` for a value that must not change:

```csharp
const int MaximumRetries = 3;
```

### Nullable values

In a project with nullable reference types enabled, `string` should contain text, while `string?` may contain `null`:

```csharp
string serviceName = "orders-api";
string? commitSha = null;
```

`null` means “no object or value is present.” Check it before use:

```csharp
if (commitSha is not null)
{
    Console.WriteLine(commitSha.ToUpper());
}
```

---

## 4. Operators and string interpolation

### Arithmetic

```csharp
int runningPods = 3;
int pendingPods = 2;
int totalPods = runningPods + pendingPods;
int podsAfterScaleDown = totalPods - 1;
```

### Comparison

```csharp
bool hasEnoughPods = runningPods >= 3;
bool isExactlyThree = runningPods == 3;
bool isNotZero = runningPods != 0;
```

`=` assigns a value; `==` compares two values.

### Logical operators

```csharp
bool deploymentSucceeded = isHealthy && runningPods >= 3; // AND
bool needsAttention = !isHealthy || pendingPods > 0;      // NOT, OR
```

### String interpolation

```csharp
Console.WriteLine($"Running: {runningPods}, pending: {pendingPods}");
```

This is usually clearer than joining strings with `+`.

---

## 5. Conditions

Use `if`, `else if`, and `else` to choose which code runs:

```csharp
int healthyReplicas = 2;
int desiredReplicas = 3;

if (healthyReplicas == desiredReplicas)
{
    Console.WriteLine("Deployment is healthy");
}
else if (healthyReplicas > 0)
{
    Console.WriteLine("Deployment is degraded");
}
else
{
    Console.WriteLine("Deployment is unavailable");
}
```

Use `switch` when one value has several known cases:

```csharp
string environment = "staging";

switch (environment)
{
    case "development":
        Console.WriteLine("Use one replica");
        break;
    case "staging":
        Console.WriteLine("Use two replicas");
        break;
    case "production":
        Console.WriteLine("Use at least three replicas");
        break;
    default:
        Console.WriteLine("Unknown environment");
        break;
}
```

---

## 6. Loops and collections

### A list

```csharp
var environments = new List<string>
{
    "development",
    "staging",
    "production"
};
```

`List<string>` means “a list whose items are strings.” The type inside `< >` is a generic type argument.

### `foreach`

Use `foreach` when you want to process every item:

```csharp
foreach (var environment in environments)
{
    Console.WriteLine($"Checking {environment}");
}
```

### `for`

Use `for` when you need an index or exact number of iterations:

```csharp
for (int attempt = 1; attempt <= 3; attempt++)
{
    Console.WriteLine($"Attempt {attempt}");
}
```

### `while`

Use `while` while a condition remains true:

```csharp
int pendingPods = 2;

while (pendingPods > 0)
{
    Console.WriteLine($"Waiting for {pendingPods} pod(s)");
    pendingPods--;
}
```

Make sure a `while` loop can eventually become false, or it will run forever.

---

## 7. Methods

A method gives a name to reusable behavior.

```csharp
static void PrintDeployment(string serviceName, int replicas)
{
    Console.WriteLine($"{serviceName}: {replicas} replica(s)");
}

PrintDeployment("payments-api", 3);
```

Read the declaration from left to right:

| Part | Meaning |
|------|---------|
| `static` | Belongs to the class rather than an object; explained later |
| `void` | Returns no value |
| `PrintDeployment` | Method name |
| `string serviceName` | First parameter |
| `int replicas` | Second parameter |

### Returning a value

```csharp
static bool IsHealthy(int healthyReplicas, int desiredReplicas)
{
    return healthyReplicas == desiredReplicas;
}

bool result = IsHealthy(3, 3);
Console.WriteLine(result);
```

`bool` before the method name is its return type. Every path through this method must return a `bool`.

### Parameters versus arguments

```csharp
static void Scale(int replicas) // replicas is a parameter
{
    Console.WriteLine(replicas);
}

Scale(5); // 5 is an argument
```

A **parameter** is the variable in the method declaration. An **argument** is the value passed by the caller.

---

## 8. Classes and objects

Object-oriented programming groups related data and behavior into objects.

- A **class** is a blueprint that defines data and behavior.
- An **object** is one instance created from that class.
- A **method** defines behavior.
- A **field** or **property** holds state.

Create a file named `Service.cs`:

```csharp
public class Service
{
    public string Name { get; set; } = "";
    public int Replicas { get; set; }

    public void PrintStatus()
    {
        Console.WriteLine($"{Name} has {Replicas} replica(s)");
    }
}
```

Use it in `Program.cs`:

```csharp
var payments = new Service();
payments.Name = "payments-api";
payments.Replicas = 3;
payments.PrintStatus();

var orders = new Service();
orders.Name = "orders-api";
orders.Replicas = 2;
orders.PrintStatus();
```

`payments` and `orders` are separate objects created from the same `Service` class. Changing one does not change the other.

### Instance versus static members

`payments.PrintStatus()` is an **instance method**: it operates on one object's state.

`Console.WriteLine()` is called through the `Console` class because `WriteLine` is **static**: no `Console` object is required.

---

## 9. Fields, properties, and constructors

### Fields

A field is a variable declared inside a class:

```csharp
public class Service
{
    private int _restartCount;
}
```

Fields are usually private so callers cannot change internal state directly.

### Properties

A property provides controlled access to data:

```csharp
public string Name { get; set; } = "";
public int Replicas { get; private set; }
```

| Property syntax | Caller can read? | Caller can change? |
|-----------------|------------------|--------------------|
| `{ get; set; }` | Yes | Yes |
| `{ get; private set; }` | Yes | No; only the class can change it |
| `{ get; }` | Yes | Only during initialization or construction |

Although a property is used with field-like syntax, its accessors can enforce rules.

### Constructors

A constructor initializes a new object. It has the same name as the class and no return type:

```csharp
public class Service
{
    public string Name { get; }
    public int Replicas { get; private set; }

    public Service(string name, int replicas)
    {
        Name = name;
        Replicas = replicas;
    }
}
```

Create the object by passing constructor arguments:

```csharp
var service = new Service("payments-api", 3);
```

The constructor ensures that required values are supplied when the object is created.

---

## 10. Object state and reference types

Most class objects are reference types. Two variables can refer to the same object:

```csharp
var first = new Service("payments-api", 3);
var second = first;

// If Scale changes the object through second,
// first observes the same changed object.
```

The variables do not contain two independent `Service` objects. They both point to one object.

This differs from simple value types such as `int`:

```csharp
int firstCount = 3;
int secondCount = firstCount;
secondCount = 5;

Console.WriteLine(firstCount);  // 3
Console.WriteLine(secondCount); // 5
```

Understanding shared object references helps explain why an object's state can change after a method call.

---

## 11. Access modifiers

Access modifiers control where a type or member can be used.

| Modifier | Accessible from |
|----------|-----------------|
| `public` | Any code that can reference the type |
| `private` | Only the containing class |
| `protected` | The class and derived classes |
| `internal` | The same project/assembly |
| `protected internal` | Same assembly **or** a derived class |
| `private protected` | Derived classes in the same assembly |

For now, focus on:

- `public` for the safe operations a caller should use.
- `private` for implementation details and state that the class must protect.
- `internal` when a type should be usable only inside its project.

Example:

```csharp
public class Deployment
{
    private int _successfulChecks;

    public string Name { get; }

    public void RecordSuccessfulCheck()
    {
        _successfulChecks++;
    }
}
```

Callers can read `Name` and call `RecordSuccessfulCheck`, but cannot directly set `_successfulChecks`.

---

## 12. Encapsulation

**Encapsulation** means keeping an object's internal state private and exposing a small set of safe operations that preserve its rules.

### A class without encapsulation

```csharp
public class Deployment
{
    public string Name = "";
    public int Replicas;
}
```

Any caller can create an invalid deployment:

```csharp
var deployment = new Deployment();
deployment.Name = "";
deployment.Replicas = -50;
```

The class has no control over its own state.

### An encapsulated class

```csharp
public class Deployment
{
    public string Name { get; }
    public int Replicas { get; private set; }

    public Deployment(string name, int replicas)
    {
        if (string.IsNullOrWhiteSpace(name))
        {
            throw new ArgumentException("Name is required", nameof(name));
        }

        if (replicas < 1)
        {
            throw new ArgumentOutOfRangeException(
                nameof(replicas),
                "A deployment needs at least one replica");
        }

        Name = name;
        Replicas = replicas;
    }

    public void ScaleTo(int replicas)
    {
        if (replicas < 1)
        {
            throw new ArgumentOutOfRangeException(
                nameof(replicas),
                "A deployment needs at least one replica");
        }

        Replicas = replicas;
    }
}
```

Use it like this:

```csharp
var deployment = new Deployment("payments-api", 3);
deployment.ScaleTo(5);

Console.WriteLine($"{deployment.Name}: {deployment.Replicas}");
```

The caller cannot write `deployment.Replicas = -1` because the setter is private. It must call `ScaleTo`, where the class validates the new value.

### Why encapsulation matters

Encapsulation:

- prevents invalid state;
- keeps validation in one place;
- makes callers use meaningful operations such as `ScaleTo`;
- allows internal implementation to change without breaking callers;
- makes classes easier to test and reason about.

Encapsulation is not merely making every field private. The important idea is that the object controls changes to its state.

---

## 13. Hands-on lab: encapsulated deployment

The final program models a basic deployment lifecycle.

### Step 1: Create `Deployment.cs`

```csharp
public class Deployment
{
    private readonly List<string> _events = new();

    public string ServiceName { get; }
    public string Environment { get; }
    public int Replicas { get; private set; }
    public bool IsRunning { get; private set; }
    public IReadOnlyList<string> Events => _events.AsReadOnly();

    public Deployment(string serviceName, string environment, int replicas)
    {
        if (string.IsNullOrWhiteSpace(serviceName))
        {
            throw new ArgumentException(
                "Service name is required",
                nameof(serviceName));
        }

        if (string.IsNullOrWhiteSpace(environment))
        {
            throw new ArgumentException(
                "Environment is required",
                nameof(environment));
        }

        ValidateReplicaCount(replicas);

        ServiceName = serviceName;
        Environment = environment;
        Replicas = replicas;
        _events.Add($"Created with {replicas} replica(s)");
    }

    public void Start()
    {
        if (IsRunning)
        {
            throw new InvalidOperationException("Deployment is already running");
        }

        IsRunning = true;
        _events.Add("Started");
    }

    public void ScaleTo(int replicas)
    {
        if (!IsRunning)
        {
            throw new InvalidOperationException(
                "Start the deployment before scaling it");
        }

        ValidateReplicaCount(replicas);
        Replicas = replicas;
        _events.Add($"Scaled to {replicas} replica(s)");
    }

    public void Stop()
    {
        if (!IsRunning)
        {
            throw new InvalidOperationException("Deployment is already stopped");
        }

        IsRunning = false;
        _events.Add("Stopped");
    }

    public void PrintSummary()
    {
        string status = IsRunning ? "Running" : "Stopped";
        Console.WriteLine(
            $"{ServiceName} [{Environment}] - {status}, {Replicas} replica(s)");
    }

    private static void ValidateReplicaCount(int replicas)
    {
        if (replicas < 1)
        {
            throw new ArgumentOutOfRangeException(
                nameof(replicas),
                "Replica count must be at least one");
        }
    }
}
```

### Step 2: Replace `Program.cs`

```csharp
var deployment = new Deployment("payments-api", "staging", 2);

deployment.PrintSummary();
deployment.Start();
deployment.ScaleTo(4);
deployment.PrintSummary();
deployment.Stop();

Console.WriteLine("Event history:");

foreach (var deploymentEvent in deployment.Events)
{
    Console.WriteLine($"- {deploymentEvent}");
}
```

### Step 3: Run it

```bash
dotnet run
```

Expected output:

```text
payments-api [staging] - Stopped, 2 replica(s)
payments-api [staging] - Running, 4 replica(s)
Event history:
- Created with 2 replica(s)
- Started
- Scaled to 4 replica(s)
- Stopped
```

### Step 4: Identify the encapsulation

| Code | Rule it protects |
|------|------------------|
| `private readonly List<string> _events` | Callers cannot replace or directly modify the internal list |
| `Replicas { get; private set; }` | Only the class can change replica count |
| `IsRunning { get; private set; }` | Only lifecycle methods change status |
| Constructor validation | An object cannot begin with invalid required data |
| `ScaleTo` validation | Replica count cannot become invalid later |
| `Start` / `Stop` checks | Invalid lifecycle transitions are rejected |
| `ValidateReplicaCount` | Reuses one private validation rule |
| `IReadOnlyList<string> Events` | Callers can inspect history without changing internal state |

---

## 14. Exercises

Work through these in order.

### Exercise 1: Add a version

Add a read-only `Version` property. Require it in the constructor and reject an empty version.

Example:

```csharp
var deployment = new Deployment(
    "payments-api",
    "staging",
    "1.4.0",
    2);
```

### Exercise 2: Prevent large staging deployments

If `Environment` is `"staging"`, reject a replica count greater than 5. Apply the rule both during construction and scaling.

### Exercise 3: Add `Restart`

Add a `Restart` method that:

1. works only when the deployment is running;
2. records `"Restarted"` in the event history;
3. does not expose `IsRunning` for public modification.

### Exercise 4: Test invalid operations manually

Try each statement separately and read the exception:

```csharp
new Deployment("", "staging", 2);
new Deployment("payments-api", "staging", 0);

var deployment = new Deployment("payments-api", "staging", 2);
deployment.ScaleTo(3); // should fail because it has not started
```

### Exercise 5: Explain the design

Answer these questions in your own words:

1. Why is `_events` private?
2. Why does `Replicas` use a private setter?
3. Why is validation inside the class instead of only in `Program.cs`?
4. What invalid states can the class prevent?
5. What is the difference between the `Deployment` class and the `deployment` object?

---

## 15. Recap

```text
statement       = an instruction ending in ;
variable        = a named value with a fixed type
condition       = chooses a path with if/switch
loop            = repeats work with foreach/for/while
method          = named, reusable behavior
class           = blueprint containing state and behavior
object          = one instance of a class
field           = data stored directly inside a class
property        = controlled access to data
constructor     = initializes a new object
access modifier = controls where a member can be used
encapsulation   = the object protects its state through safe operations
```

The central OOP idea from this chapter is:

> Keep state private, expose meaningful operations, and validate every change that could break the object's rules.

After this chapter, the next OOP topics are **abstraction, inheritance, and polymorphism**, followed by interfaces, dependency injection, and unit testing.
