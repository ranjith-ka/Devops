# Spec-Driven Development and Behavior-Driven Development

This tutorial introduces two complementary practices that help teams turn requirements into reliable software:

- Spec-Driven Development (SDD): start from a clear contract or specification.
- Behavior-Driven Development (BDD): describe expected behavior in business-friendly language before implementation.

Together, they improve alignment between product, engineering, and QA.

---

## What you will learn

- The difference between SDD and BDD
- How to write a useful specification
- How to turn requirements into executable examples
- How to keep implementation aligned with business intent

---

## 1. Understand the two approaches

### Spec-Driven Development (SDD)

SDD starts with a contract that defines what the system must provide. This contract can be an OpenAPI document, an API schema, a design note, or a technical specification.

Benefits:
- Makes requirements explicit
- Reduces ambiguity
- Gives developers a shared source of truth
- Makes validation and testing easier

### Behavior-Driven Development (BDD)

BDD focuses on describing user-visible behavior in a plain-language format such as Gherkin.

Example:

```gherkin
Feature: Order placement
  Scenario: A customer places a valid order
    Given the customer is signed in
    When they submit a valid order
    Then the order should be created successfully
```

Benefits:
- Improves collaboration with non-technical stakeholders
- Makes acceptance criteria easy to review
- Helps connect requirements to automated tests

---

## 2. When to use each practice

| Practice | Best for | Typical artifact |
| --- | --- | --- |
| SDD | APIs, contracts, interfaces, architecture decisions | OpenAPI, JSON schema, RFCs, design docs |
| BDD | Business rules, user journeys, acceptance criteria | Gherkin scenarios, feature files |

In practice, teams often use both:
- SDD defines the contract.
- BDD captures the expected behavior around that contract.

---

## 3. A simple workflow

### Step 1: Start with the requirement

Write the requirement in a form that is easy to discuss and review.

Example:

- A customer can create an order
- The order must include at least one item
- The API must return HTTP 201 on success

### Step 2: Create a specification

For an API, create an OpenAPI document that describes:
- endpoints
- request body
- response codes
- validation rules

Example idea:

```yaml
openapi: 3.1.0
info:
  title: Order API
  version: 1.0.0
paths:
  /orders:
    post:
      summary: Create an order
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [items]
              properties:
                items:
                  type: array
                  minItems: 1
```

### Step 3: Write BDD scenarios

Turn the requirement into examples the team can discuss.

```gherkin
Feature: Creating an order
  Scenario: Create an order with one item
    Given a customer is signed in
    When they submit an order with one item
    Then the order is created
    And the response status is 201
```

### Step 4: Implement against the spec

Build the service or API so it satisfies both:
- the contract in the specification
- the behavior described in the scenarios

### Step 5: Validate continuously

Use automated checks such as:
- contract tests
- API validation tests
- acceptance tests based on BDD scenarios

---

## 4. Recommended tools

- OpenAPI / Swagger for API contracts
- Gherkin with Cucumber, Behave, or Reqnroll for BDD
- Postman or Bruno for API examples
- Contract testing tools such as Pact for service-to-service validation
- CI pipelines to run validation automatically

---

## 5. Practical exercise

Try this mini exercise in your team:

1. Pick one feature such as user signup or order creation.
2. Write a short specification for the feature.
3. Convert the main requirement into 2 to 3 Gherkin scenarios.
4. Implement the feature and test it against those scenarios.
5. Review whether the spec and the scenarios still reflect the business need.

---

## 6. Best practices

- Keep the specification understandable for both technical and non-technical readers.
- Prefer examples over abstract descriptions.
- Make acceptance criteria testable.
- Update the spec and scenarios whenever requirements change.
- Use automation so the team gets feedback quickly.

---

## Summary

Spec-Driven Development and Behavior-Driven Development help teams build the right thing in a disciplined and collaborative way.

- SDD gives you a shared contract.
- BDD turns that contract into observable business behavior.
- Together, they reduce ambiguity and improve software quality.
