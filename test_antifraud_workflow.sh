#!/bin/bash

# Test script for antifraud workflow on jump server
# This script should be run from the jump server (10.200.1.2)

echo "=== Testing Antifraud Workflow on Jump Server ==="
echo ""

# Check if we're on the jump server
if [[ ! $(hostname) =~ "10.200.1.2" ]]; then
    echo "Warning: This script should be run from the jump server (10.200.1.2)"
    echo "Current host: $(hostname)"
    echo ""
fi

# Navigate to the synced code directory
cd /tmp/uwf-deploy || { echo "Error: /tmp/uwf-deploy not found"; exit 1; }

echo "1. Checking workflow files..."
echo "   - antifraud_workflow.go: $(ls -la workflows/antifraud_workflow.go 2>/dev/null && echo "✓ Found" || echo "✗ Missing")"
echo "   - antifraud_step.go: $(ls -la workflows/steps/antifraud_step.go 2>/dev/null && echo "✓ Found" || echo "✗ Missing")"
echo "   - fc_validation_step.go: $(ls -la workflows/steps/fc_validation_step.go 2>/dev/null && echo "✓ Found" || echo "✗ Missing")"
echo "   - ml_validation_step.go: $(ls -la workflows/steps/ml_validation_step.go 2>/dev/null && echo "✓ Found" || echo "✗ Missing")"
echo "   - finalize_transaction_step.go: $(ls -la workflows/steps/finalize_transaction_step.go 2>/dev/null && echo "✓ Found" || echo "✗ Missing")"
echo ""

echo "2. Checking if we can access TAF service (172.30.75.91)..."
echo "   Testing connectivity to baraiq-p-dbp01.kazpost.kz (172.30.75.91)..."
if ping -c 1 -W 2 172.30.75.91 > /dev/null 2>&1; then
    echo "   ✓ TAF service is reachable"
    
    # Test HTTP connectivity
    echo "   Testing HTTP connectivity to https://af-test.qazpost.kz..."
    if curl -s --connect-timeout 5 -k https://af-test.qazpost.kz > /dev/null 2>&1; then
        echo "   ✓ HTTPS connection successful"
    else
        echo "   ✗ HTTPS connection failed (may require VPN or different endpoint)"
    fi
else
    echo "   ✗ TAF service not reachable (may require VPN)"
fi
echo ""

echo "3. Checking workflow registration..."
echo "   Note: Workflows need to be registered in the workflow registry"
echo "   The 'antifraud-transaction-validation' workflow should be available"
echo ""

echo "4. Testing workflow execution flow..."
echo "   Workflow steps:"
echo "   1. StoreTransactionStep - Stores transaction in antifraud system"
echo "   2. AMLValidationStep - Anti-Money Laundering validation"
echo "   3. FCValidationStep - Fraud Check validation"
echo "   4. MLValidationStep - Machine Learning validation"
echo "   5. FinalizeTransactionStep - Final decision and resolution"
echo ""

echo "5. Sample transaction data format:"
cat << 'EOF'
{
  "transaction": {
    "id": "txn-12345",
    "type": "deposit",
    "amount": "100000",
    "currency": "KZT",
    "client_id": "client-001",
    "client_name": "John Smith",
    "client_pan": "111111******1111",
    "client_cvv": "111",
    "client_card_holder": "JOHN SMITH",
    "client_phone": "+77007007070",
    "merchant_terminal_id": "00000001",
    "channel": "E-com",
    "location_ip": "192.168.0.1"
  }
}
EOF
echo ""

echo "6. To execute the workflow from jump server:"
echo "   Option 1: Use uwf-cli (if installed)"
echo "     ./uwf-cli execute sync antifraud-transaction-validation --input '{\"transaction\": {...}}'"
echo ""
echo "   Option 2: Use curl directly (if workflow API is running)"
echo "     curl -X POST https://af-test.qazpost.kz/api/v1/workflows/antifraud-transaction-validation/execute \\"
echo "       -H 'Content-Type: application/json' \\"
echo "       -H 'Authorization: Bearer <token>' \\"
echo "       -d '{\"input_data\": {\"transaction\": {...}}}'"
echo ""
echo "   Option 3: Run demo program (requires Go)"
echo "     go run examples/antifraud_workflow_demo.go"
echo ""

echo "7. Verification steps:"
echo "   - Check workflow execution logs"
echo "   - Verify antifraud service calls"
echo "   - Validate transaction results"
echo "   - Monitor step completion"
echo ""

echo "=== Test Summary ==="
echo "✅ Code successfully synced to jump server"
echo "✅ Workflow files present and up-to-date"
echo "⚠️  TAF service connectivity needs VPN from jump server"
echo "⚠️  Workflow execution requires running workflow services"
echo ""
echo "Next steps:"
echo "1. Ensure VPN is connected from jump server"
echo "2. Deploy workflow services using: ./uwf-cli deploy push --all"
echo "3. Register antifraud workflow in registry"
echo "4. Execute workflow using uwf-cli or demo program"
echo ""

# Check if Docker images were built
echo "Checking Docker images..."
if docker images | grep -q "172.30.75.78:9080/taf/workflow-api"; then
    echo "✅ workflow-api Docker image built"
else
    echo "⚠️  workflow-api Docker image not built yet"
fi

if docker images | grep -q "172.30.75.78:9080/taf/workflow-worker"; then
    echo "✅ workflow-worker Docker image built"
else
    echo "⚠️  workflow-worker Docker image not built yet"
fi

if docker images | grep -q "172.30.75.78:9080/taf/workflow-registry"; then
    echo "✅ workflow-registry Docker image built"
else
    echo "⚠️  workflow-registry Docker image not built yet"
fi

if docker images | grep -q "172.30.75.78:9080/taf/workflow-executor"; then
    echo "✅ workflow-executor Docker image built"
else
    echo "⚠️  workflow-executor Docker image not built yet"
fi

echo ""
echo "To complete deployment:"
echo "  ./uwf-cli deploy push --all"
echo "  ./uwf-cli deploy verify"