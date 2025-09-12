from flask import Flask, request, jsonify
import numpy as np
from scipy.optimize import minimize, linprog
from datetime import datetime
import json
import math
from itertools import permutations

app = Flask(__name__)

@app.route("/health")
def health():
    return jsonify({"status": "ok"})


def haversine_distance(lat1, lon1, lat2, lon2):
    """Calculate haversine distance between two points"""
    R = 6371  # Earth's radius in kilometers

    lat1, lon1, lat2, lon2 = map(math.radians, [lat1, lon1, lat2, lon2])
    dlat = lat2 - lat1
    dlon = lon2 - lon1

    a = (
        math.sin(dlat / 2) ** 2
        + math.cos(lat1) * math.cos(lat2) * math.sin(dlon / 2) ** 2
    )
    c = 2 * math.asin(math.sqrt(a))

    return R * c


@app.route("/optimize_route", methods=["POST"])
def optimize_route():
    """Enhanced route optimization with multiple algorithms"""
    try:
        data = request.json

        locations = data.get("locations", [])
        start_location = data.get("start_location")
        vehicle_capacity = data.get("vehicle_capacity", 1000)
        algorithm = data.get(
            "algorithm", "nearest_neighbor"
        )  # nearest_neighbor, genetic, or brute_force

        if len(locations) < 2:
            return jsonify({"error": "At least 2 locations required"}), 400

        # Parse locations and calculate distance matrix
        coords = []
        location_names = []
        demands = []

        for loc in locations:
            if not all(key in loc for key in ["lat", "lng"]):
                return (
                    jsonify({"error": "Each location must have 'lat' and 'lng'"}),
                    400,
                )

            coords.append((loc["lat"], loc["lng"]))
            location_names.append(loc.get("name", f"Location {len(coords)}"))
            demands.append(loc.get("demand", 0))

        # Add start location if provided
        if start_location:
            if not all(key in start_location for key in ["lat", "lng"]):
                return (
                    jsonify({"error": "Start location must have 'lat' and 'lng'"}),
                    400,
                )
            coords.insert(0, (start_location["lat"], start_location["lng"]))
            location_names.insert(0, start_location.get("name", "Start"))
            demands.insert(0, 0)

        # Calculate distance matrix
        n = len(coords)
        distance_matrix = np.zeros((n, n))

        for i in range(n):
            for j in range(n):
                if i != j:
                    distance_matrix[i][j] = haversine_distance(
                        coords[i][0], coords[i][1], coords[j][0], coords[j][1]
                    )

        # Choose optimization algorithm
        if algorithm == "nearest_neighbor":
            route, total_distance = nearest_neighbor_tsp(
                distance_matrix, 0 if start_location else None
            )
        elif algorithm == "genetic":
            route, total_distance = genetic_algorithm_tsp(
                distance_matrix, 0 if start_location else None
            )
        elif (
            algorithm == "brute_force" and len(locations) <= 8
        ):  # Only for small problems
            route, total_distance = brute_force_tsp(
                distance_matrix, 0 if start_location else None
            )
        else:
            route, total_distance = nearest_neighbor_tsp(
                distance_matrix, 0 if start_location else None
            )

        # Calculate total demand and check capacity
        total_demand = sum(demands[i] for i in route)
        capacity_utilization = (
            (total_demand / vehicle_capacity) * 100 if vehicle_capacity > 0 else 0
        )

        # Create route details
        route_details = []
        cumulative_distance = 0
        cumulative_demand = 0

        for i, location_idx in enumerate(route):
            if i > 0:
                prev_idx = route[i - 1]
                segment_distance = distance_matrix[prev_idx][location_idx]
                cumulative_distance += segment_distance

            cumulative_demand += demands[location_idx]

            route_details.append(
                {
                    "order": i + 1,
                    "location_index": location_idx,
                    "name": location_names[location_idx],
                    "coordinates": {
                        "lat": coords[location_idx][0],
                        "lng": coords[location_idx][1],
                    },
                    "demand": demands[location_idx],
                    "cumulative_demand": cumulative_demand,
                    "cumulative_distance": round(cumulative_distance, 2),
                }
            )

        # Calculate savings compared to visiting each location from start
        if start_location:
            naive_distance = sum(
                distance_matrix[0][i] * 2 for i in range(1, n)
            )  # Round trips
            savings = naive_distance - total_distance
            savings_percent = (
                (savings / naive_distance) * 100 if naive_distance > 0 else 0
            )
        else:
            savings = 0
            savings_percent = 0

        return jsonify(
            {
                "optimized_route": route_details,
                "total_distance_km": round(total_distance, 2),
                "total_demand": total_demand,
                "capacity_utilization": round(capacity_utilization, 2),
                "algorithm_used": algorithm,
                "optimization_results": {
                    "savings_km": round(savings, 2),
                    "savings_percent": round(savings_percent, 2),
                    "feasible": capacity_utilization <= 100,
                },
                "processing_time": datetime.now().isoformat(),
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


def nearest_neighbor_tsp(distance_matrix, start_idx=None):
    """Nearest neighbor algorithm for TSP"""
    n = len(distance_matrix)

    if start_idx is None:
        start_idx = 0

    visited = [False] * n
    route = [start_idx]
    visited[start_idx] = True
    total_distance = 0

    current = start_idx
    for _ in range(n - 1):
        nearest_dist = float("inf")
        nearest_idx = -1

        for i in range(n):
            if not visited[i] and distance_matrix[current][i] < nearest_dist:
                nearest_dist = distance_matrix[current][i]
                nearest_idx = i

        if nearest_idx != -1:
            route.append(nearest_idx)
            visited[nearest_idx] = True
            total_distance += nearest_dist
            current = nearest_idx

    # Return to start if needed
    if start_idx is not None:
        total_distance += distance_matrix[current][start_idx]
        route.append(start_idx)

    return route, total_distance


def genetic_algorithm_tsp(
    distance_matrix, start_idx=None, population_size=50, generations=100
):
    """Genetic algorithm for TSP"""
    n = len(distance_matrix)

    if n <= 3:
        return nearest_neighbor_tsp(distance_matrix, start_idx)

    # Create initial population
    population = []
    for _ in range(population_size):
        if start_idx is not None:
            individual = [start_idx] + list(
                np.random.permutation([i for i in range(n) if i != start_idx])
            )
        else:
            individual = list(np.random.permutation(range(n)))
        population.append(individual)

    def fitness(route):
        total_dist = 0
        for i in range(len(route) - 1):
            total_dist += distance_matrix[route[i]][route[i + 1]]
        if start_idx is not None:
            total_dist += distance_matrix[route[-1]][route[0]]
        return 1 / (1 + total_dist)  # Higher fitness for shorter routes

    def crossover(parent1, parent2):
        if start_idx is not None:
            # Keep start fixed
            p1_cities = parent1[1:]
            p2_cities = parent2[1:]
            cut = len(p1_cities) // 2
            child_cities = p1_cities[:cut] + [
                city for city in p2_cities if city not in p1_cities[:cut]
            ]
            return [start_idx] + child_cities
        else:
            cut = len(parent1) // 2
            child = parent1[:cut] + [
                city for city in parent2 if city not in parent1[:cut]
            ]
            return child

    def mutate(route):
        if len(route) < 3:
            return route

        route_copy = route.copy()
        start_pos = 1 if start_idx is not None else 0

        i, j = np.random.choice(range(start_pos, len(route_copy)), 2, replace=False)
        route_copy[i], route_copy[j] = route_copy[j], route_copy[i]
        return route_copy

    # Evolution
    for generation in range(generations):
        # Calculate fitness
        fitness_scores = [fitness(individual) for individual in population]

        # Selection and reproduction
        new_population = []
        for _ in range(population_size):
            # Tournament selection
            tournament_size = 3
            tournament = np.random.choice(
                population_size, tournament_size, replace=False
            )
            winner = max(tournament, key=lambda x: fitness_scores[x])

            if np.random.random() < 0.8:  # Crossover probability
                parent2_idx = np.random.choice(population_size)
                child = crossover(population[winner], population[parent2_idx])
            else:
                child = population[winner].copy()

            if np.random.random() < 0.1:  # Mutation probability
                child = mutate(child)

            new_population.append(child)

        population = new_population

    # Return best solution
    fitness_scores = [fitness(individual) for individual in population]
    best_idx = np.argmax(fitness_scores)
    best_route = population[best_idx]

    total_distance = 0
    for i in range(len(best_route) - 1):
        total_distance += distance_matrix[best_route[i]][best_route[i + 1]]
    if start_idx is not None:
        total_distance += distance_matrix[best_route[-1]][best_route[0]]
        best_route.append(best_route[0])  # Complete the cycle

    return best_route, total_distance


def brute_force_tsp(distance_matrix, start_idx=None):
    """Brute force TSP for small problems"""
    n = len(distance_matrix)

    if start_idx is not None:
        cities = [i for i in range(n) if i != start_idx]
        min_distance = float("inf")
        best_route = None

        for perm in permutations(cities):
            route = [start_idx] + list(perm) + [start_idx]
            distance = sum(
                distance_matrix[route[i]][route[i + 1]] for i in range(len(route) - 1)
            )

            if distance < min_distance:
                min_distance = distance
                best_route = route

        return best_route, min_distance
    else:
        cities = list(range(n))
        min_distance = float("inf")
        best_route = None

        for perm in permutations(cities):
            route = list(perm)
            distance = sum(
                distance_matrix[route[i]][route[(i + 1) % n]] for i in range(n)
            )

            if distance < min_distance:
                min_distance = distance
                best_route = route

        return best_route, min_distance


@app.route("/optimize", methods=["POST"])
def optimize():
    """Legacy endpoint - redirects to optimize_route"""
    return optimize_route()


@app.route("/optimize_inventory", methods=["POST"])
def optimize_inventory():
    """Advanced inventory optimization"""
    try:
        data = request.json

        items = data.get("items", [])
        budget_constraint = data.get("budget", 10000)
        storage_constraint = data.get("storage_capacity", 1000)
        demand_uncertainty = data.get("demand_uncertainty", 0.1)

        if not items:
            return jsonify({"error": "No items provided"}), 400

        results = []
        total_cost = 0
        total_storage = 0
        total_value = 0

        for item in items:
            name = item.get("name", "Unknown")
            demand = item.get("demand", 0)
            unit_cost = item.get("unit_cost", 0)
            storage_per_unit = item.get("storage_per_unit", 1)
            lead_time = item.get("lead_time_days", 7)
            service_level = item.get("service_level", 0.95)

            # Calculate safety stock considering demand uncertainty
            if demand > 0:
                # Wilson's EOQ formula
                holding_cost_rate = 0.2  # 20% per year
                ordering_cost = unit_cost * 0.1  # 10% of unit cost

                eoq = math.sqrt(
                    (2 * demand * ordering_cost) / (unit_cost * holding_cost_rate)
                )

                # Safety stock calculation
                demand_std = demand * demand_uncertainty
                z_score = (
                    1.645 if service_level >= 0.95 else 1.28
                )  # 95% or 90% service level
                safety_stock = z_score * demand_std * math.sqrt(lead_time / 365)

                # Reorder point
                reorder_point = (demand * lead_time / 365) + safety_stock

                # Optimal order quantity considering constraints
                max_affordable = budget_constraint // unit_cost
                max_storage = storage_constraint // storage_per_unit

                optimal_quantity = min(eoq, max_affordable, max_storage)

                item_cost = optimal_quantity * unit_cost
                item_storage = optimal_quantity * storage_per_unit
                item_value = optimal_quantity * demand  # Simplified value calculation

                total_cost += item_cost
                total_storage += item_storage
                total_value += item_value

                results.append(
                    {
                        "item_name": name,
                        "optimal_quantity": round(optimal_quantity, 2),
                        "eoq": round(eoq, 2),
                        "safety_stock": round(safety_stock, 2),
                        "reorder_point": round(reorder_point, 2),
                        "total_cost": round(item_cost, 2),
                        "storage_required": round(item_storage, 2),
                        "expected_value": round(item_value, 2),
                        "turnover_rate": round(
                            demand / optimal_quantity if optimal_quantity > 0 else 0, 2
                        ),
                    }
                )

        # Calculate overall metrics
        budget_utilization = (total_cost / budget_constraint) * 100
        storage_utilization = (total_storage / storage_constraint) * 100

        return jsonify(
            {
                "optimization_results": results,
                "summary": {
                    "total_cost": round(total_cost, 2),
                    "total_storage_required": round(total_storage, 2),
                    "total_expected_value": round(total_value, 2),
                    "budget_utilization_percent": round(budget_utilization, 2),
                    "storage_utilization_percent": round(storage_utilization, 2),
                    "roi_estimate": (
                        round((total_value - total_cost) / total_cost * 100, 2)
                        if total_cost > 0
                        else 0
                    ),
                },
                "recommendations": [
                    (
                        "Monitor reorder points closely"
                        if any(
                            r["reorder_point"] > r["optimal_quantity"] for r in results
                        )
                        else "Inventory levels appear optimal"
                    ),
                    (
                        "Consider increasing budget"
                        if budget_utilization > 95
                        else "Budget utilization is healthy"
                    ),
                    (
                        "Optimize storage space"
                        if storage_utilization > 90
                        else "Storage capacity is adequate"
                    ),
                ],
            }
        )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/optimize_resource_allocation", methods=["POST"])
def optimize_resource_allocation():
    """Resource allocation optimization using linear programming"""
    try:
        data = request.json

        resources = data.get("resources", [])
        tasks = data.get("tasks", [])
        constraints = data.get("constraints", {})

        if not resources or not tasks:
            return jsonify({"error": "Both resources and tasks required"}), 400

        # Create allocation matrix
        n_resources = len(resources)
        n_tasks = len(tasks)

        # Objective: maximize total value or minimize total cost
        objective = []
        for i in range(n_resources):
            for j in range(n_tasks):
                # Value or negative cost for maximization
                value = tasks[j].get("value", 0) - tasks[j].get("cost", 0)
                objective.append(value)

        # Constraints
        A_ub = []
        b_ub = []

        # Resource capacity constraints
        for i, resource in enumerate(resources):
            constraint_row = [0] * (n_resources * n_tasks)
            capacity = resource.get("capacity", 1)

            for j in range(n_tasks):
                idx = i * n_tasks + j
                constraint_row[idx] = tasks[j].get("resource_requirement", 1)

            A_ub.append(constraint_row)
            b_ub.append(capacity)

        # Task requirements constraints
        for j, task in enumerate(tasks):
            constraint_row = [0] * (n_resources * n_tasks)
            requirement = task.get("min_allocation", 0)

            for i in range(n_resources):
                idx = i * n_tasks + j
                constraint_row[idx] = -1  # Negative for >= constraint

            if requirement > 0:
                A_ub.append(constraint_row)
                b_ub.append(-requirement)

        # Bounds (all variables >= 0)
        bounds = [(0, None) for _ in range(n_resources * n_tasks)]

        # Solve linear program
        try:
            result = linprog(
                c=[-x for x in objective],  # Negative for maximization
                A_ub=A_ub,
                b_ub=b_ub,
                bounds=bounds,
                method="highs",
            )

            if result.success:
                # Parse results
                allocation_matrix = np.array(result.x).reshape(n_resources, n_tasks)

                allocations = []
                total_value = 0

                for i, resource in enumerate(resources):
                    resource_allocations = []
                    resource_utilization = 0

                    for j, task in enumerate(tasks):
                        allocation = allocation_matrix[i][j]
                        if allocation > 0.01:  # Ignore very small allocations
                            resource_allocations.append(
                                {
                                    "task_name": task.get("name", f"Task {j+1}"),
                                    "allocation": round(allocation, 2),
                                    "value_generated": round(
                                        allocation * task.get("value", 0), 2
                                    ),
                                }
                            )
                            resource_utilization += allocation * task.get(
                                "resource_requirement", 1
                            )
                            total_value += allocation * task.get("value", 0)

                    allocations.append(
                        {
                            "resource_name": resource.get("name", f"Resource {i+1}"),
                            "capacity": resource.get("capacity", 1),
                            "utilization": round(resource_utilization, 2),
                            "utilization_percent": round(
                                (resource_utilization / resource.get("capacity", 1))
                                * 100,
                                2,
                            ),
                            "task_allocations": resource_allocations,
                        }
                    )

                return jsonify(
                    {
                        "optimization_successful": True,
                        "total_value": round(total_value, 2),
                        "resource_allocations": allocations,
                        "optimization_status": "optimal",
                        "processing_time": datetime.now().isoformat(),
                    }
                )
            else:
                return (
                    jsonify(
                        {
                            "optimization_successful": False,
                            "error": "No optimal solution found",
                            "status": result.message,
                        }
                    ),
                    400,
                )

        except Exception as solve_error:
            return (
                jsonify(
                    {
                        "optimization_successful": False,
                        "error": f"Solver error: {str(solve_error)}",
                    }
                ),
                500,
            )

    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8003, debug=True)
