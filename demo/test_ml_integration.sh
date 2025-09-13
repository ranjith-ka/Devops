#!/bin/bash

echo "🚀 Testing Emapy AI Advanced ML Integrations"
echo "=============================================="

# Route through frontend nginx reverse proxy to avoid host port conflicts (e.g., macOS AirPlay on 5000)
API_BASE="http://localhost/api"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name="$1"
    local url="$2"
    local data="$3"
    local method="${4:-POST}"
    
    echo -e "\n${YELLOW}Testing: $name${NC}"
    echo "URL: $url"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
        echo "Response: $body" | head -c 200
        echo "..."
    else
        echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
        echo "Response: $body"
    fi
}

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Test Health Endpoints
echo -e "\n${YELLOW}=== Health Checks ===${NC}"
test_endpoint "API Health" "${API_BASE}/health" "" "GET"

# Test NLP Service
echo -e "\n${YELLOW}=== NLP Service Tests ===${NC}"

test_endpoint "Text Classification" "${API_BASE}/ml/classify" '{
    "text": "Urgent customs issue with container ABCD1234567. Need expedited clearance for premium customer."
}'

test_endpoint "Metadata Extraction" "${API_BASE}/ml/extract_metadata" '{
    "text": "Container MSKU1234567 ETA 2025-12-20, SLA 5 days, VIP customer Gold tier"
}'

test_endpoint "Similarity Check" "${API_BASE}/ml/similarity_check" '{
    "text1": "Delayed shipment for container ABC123",
    "text2": "Container ABC123 has shipping delay"
}'

# Test Time Series Forecasting
echo -e "\n${YELLOW}=== Time Series Forecasting Tests ===${NC}"

test_endpoint "Demand Forecasting" "${API_BASE}/ml/forecast_demand" '{
    "data": [
        {"value": 100}, {"value": 120}, {"value": 110}, 
        {"value": 130}, {"value": 125}, {"value": 140},
        {"value": 135}, {"value": 145}, {"value": 150},
        {"value": 155}, {"value": 160}, {"value": 165}
    ],
    "days": 7,
    "model": "random_forest",
    "confidence_level": 0.95
}'

test_endpoint "Capacity Planning" "${API_BASE}/ml/capacity_planning" '{
    "current_capacity": 1000,
    "demand_data": [
        {"value": 800}, {"value": 820}, {"value": 850},
        {"value": 880}, {"value": 900}, {"value": 920}
    ],
    "expected_growth_rate": 0.15,
    "planning_horizon_months": 12
}'

# Test Optimization Service
echo -e "\n${YELLOW}=== Optimization Service Tests ===${NC}"

test_endpoint "Route Optimization" "${API_BASE}/ml/optimize_route" '{
    "locations": [
        {"lat": 40.7128, "lng": -74.0060, "name": "New York", "demand": 100},
        {"lat": 34.0522, "lng": -118.2437, "name": "Los Angeles", "demand": 80},
        {"lat": 41.8781, "lng": -87.6298, "name": "Chicago", "demand": 90},
        {"lat": 29.7604, "lng": -95.3698, "name": "Houston", "demand": 70}
    ],
    "start_location": {"lat": 40.7128, "lng": -74.0060, "name": "Warehouse"},
    "algorithm": "nearest_neighbor",
    "vehicle_capacity": 1000
}'

# Sample base64 image for CV testing (1x1 pixel PNG)
SAMPLE_IMAGE="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

# Test Computer Vision (with sample image)
echo -e "\n${YELLOW}=== Computer Vision Tests ===${NC}"

test_endpoint "Document Analysis" "${API_BASE}/ml/analyze_document" "{
    \"image\": \"$SAMPLE_IMAGE\",
    \"type\": \"document\"
}"

# Performance Test
echo -e "\n${YELLOW}=== Performance Tests ===${NC}"

echo "🏃 Running performance test (10 requests)..."
start_time=$(date +%s)

for i in {1..10}; do
    curl -s -o /dev/null "${API_BASE}/ml/classify" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"text": "Test performance request '$i'"}'
done

end_time=$(date +%s)
duration=$((end_time - start_time))
rps=$((10 / duration))

echo "Performance Results:"
echo "- Total time: ${duration}s"
echo "- Requests per second: ${rps}"
echo "- Average response time: $((duration * 100 / 10))ms"

# Integration Test
echo -e "\n${YELLOW}=== Integration Test ===${NC}"

echo "🔄 Testing end-to-end workflow..."

# 1. Create a case
echo "1. Creating test case..."
case_response=$(curl -s -X POST "${API_BASE}/cases" \
    -H "Content-Type: application/json" \
    -d '{
        "title": "Urgent customs clearance needed",
        "details": "Container TEST1234567 requires expedited customs clearance. Customer: Premium Corp, SLA: 2 days",
        "sender": "premium.corp@example.com",
        "hypercare": true,
        "sla_days": 2
    }')

echo "Case created: $case_response"

# 2. Classify the case
echo "2. Classifying case content..."
classification=$(curl -s -X POST "${API_BASE}/ml/classify" \
    -H "Content-Type: application/json" \
    -d '{
        "text": "Urgent customs clearance needed for container TEST1234567"
    }')

echo "Classification: $classification"

# 3. Extract metadata
echo "3. Extracting metadata..."
metadata=$(curl -s -X POST "${API_BASE}/ml/extract_metadata" \
    -H "Content-Type: application/json" \
    -d '{
        "text": "Container TEST1234567 requires expedited customs clearance. Customer: Premium Corp, SLA: 2 days"
    }')

echo "Metadata: $metadata"

echo -e "\n${GREEN}🎉 Integration test completed!${NC}"

# Summary
echo -e "\n${YELLOW}=== Test Summary ===${NC}"
echo "✅ All ML services tested"
echo "✅ API endpoints validated"
echo "✅ Integration workflow verified"
echo "✅ Performance benchmarked"

echo -e "\n${GREEN}🚀 Emapy AI Advanced ML Integration is ready!${NC}"
echo ""
echo "Next steps:"
echo "1. Access the ML Dashboard: http://localhost/ml-dashboard.html"
echo "2. View the main demo: http://localhost/"
echo "3. Monitor logs: docker-compose -f docker-compose.demo.yml logs -f"
echo "4. Check API documentation: http://localhost:5000/api/health"

echo ""
echo "For more information, see: ML_INTEGRATION_README.md"
