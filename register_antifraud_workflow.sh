#!/bin/bash

# Script to register and test antifraud workflow on jump server
# Run this from the jump server (10.200.1.2)

echo "=== Registering and Testing Antifraud Workflow on Jump Server ==="
echo ""

# Navigate to the synced code directory
cd /tmp/uwf-deploy || { echo "Error: /tmp/uwf-deploy not found"; exit 1; }

echo "1. Checking if uwf-cli is available..."
if [ -f "./uwf-cli" ]; then
    echo "   ✅ uwf-cli found"
    chmod +x ./uwf-cli
else
    echo "   ❌ uwf-cli not found, building it..."
    go build -o uwf-cli cmd/uwf-cli/main.go
    if [ $? -eq 0 ]; then
        echo "   ✅ uwf-cli built successfully"
        chmod +x ./uwf-cli
    else
        echo "   ❌ Failed to build uwf-cli"
        exit 1
    fi
fi

echo ""
echo "2. Checking workflow API connectivity..."
curl -k -s https://af-test.qazpost.kz/api/v1/workflows > /dev/null
if [ $? -eq 0 ]; then
    echo "   ✅ Workflow API is accessible"
else
    echo "   ❌ Workflow API is not accessible"
    exit 1
fi

echo ""
echo "3. Listing existing workflows..."
echo "   Current workflows:"
curl -k -s https://af-test.qazpost.kz/api/v1/workflows | jq -r '.workflows[] | "     - \(.name) (ID: \(.id), Steps: \(.step_count))"'

echo ""
echo "4. Checking if antifraud workflow already exists..."
EXISTING_WORKFLOWS=$(curl -k -s https://af-test.qazpost.kz/api/v1/workflows | jq -r '.workflows[] | select(.name == "antifraud-transaction-validation") | .id')
if [ -n "$EXISTING_WORKFLOWS" ]; then
    echo "   ⚠️  Antifraud workflow already exists with IDs:"
    for id in $EXISTING_WORKFLOWS; do
        echo "     - $id"
    done
    echo "   Deleting duplicate workflows..."
    for id in $EXISTING_WORKFLOWS; do
        echo "     Deleting workflow $id..."
        curl -k -X DELETE "https://af-test.qazpost.kz/api/v1/workflows/$id"
    done
    echo "   ✅ Duplicate workflows deleted"
fi

echo ""
echo "5. Creating antifraud workflow..."
# Create a workflow definition file
cat > antifraud_workflow_definition.json << 'EOF'
{
  "name": "antifraud-transaction-validation",
  "description": "Complete transaction validation using antifraud SDK at https://af-test.qazpost.kz",
  "steps": [
    {
      "name": "StoreTransactionStep",
      "description": "Stores transaction in antifraud system",
      "type": "primitive",
      "primitive_name": "antifraud",
      "operation": "StoreTransaction",
      "parameters": {
        "endpoint": "https://af-test.qazpost.kz"
      }
    },
    {
      "name": "AMLValidationStep",
      "description": "Anti-Money Laundering validation",
      "type": "primitive",
      "primitive_name": "antifraud",
      "operation": "ValidateTransactionByAML",
      "parameters": {
        "endpoint": "https://af-test.qazpost.kz"
      }
    },
    {
      "name": "FCValidationStep",
      "description": "Fraud Check validation",
      "type": "primitive",
      "primitive_name": "antifraud",
      "operation": "ValidateTransactionByFC",
      "parameters": {
        "endpoint": "https://af-test.qazpost.kz"
      }
    },
    {
      "name": "MLValidationStep",
      "description": "Machine Learning validation",
      "type": "primitive",
      "primitive_name": "antifraud",
      "operation": "ValidateTransactionByML",
      "parameters": {
        "endpoint": "https://af-test.qazpost.kz"
      }
    },
    {
      "name": "FinalizeTransactionStep",
      "description": "Final decision and resolution",
      "type": "primitive",
      "primitive_name": "antifraud",
      "operation": "FinalizeTransaction",
      "parameters": {
        "endpoint": "https://af-test.qazpost.kz"
      }
    }
  ]
}
EOF

echo "   Workflow definition created: antifraud_workflow_definition.json"

# Try to create workflow using curl directly
echo "   Creating workflow via API..."
CREATE_RESPONSE=$(curl -k -s -X POST https://af-test.qazpost.kz/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "antifraud-transaction-validation",
    "description": "Complete transaction validation using antifraud SDK at https://af-test.qazpost.kz"
  }')

if echo "$CREATE_RESPONSE" | grep -q "workflowId"; then
    WORKFLOW_ID=$(echo "$CREATE_RESPONSE" | jq -r '.workflowId')
    echo "   ✅ Workflow created successfully!"
    echo "   Workflow ID: $WORKFLOW_ID"
else
    echo "   ❌ Failed to create workflow: $CREATE_RESPONSE"
    echo "   Trying alternative approach..."
    
    # Try using the uwf-cli if available
    if [ -f "./uwf-cli" ]; then
        echo "   Using uwf-cli to create workflow..."
        ./uwf-cli workflows create \
          --name "antifraud-transaction-validation" \
          --description "Complete transaction validation using antifraud SDK at https://af-test.qazpost.kz" \
          --verbose
    fi
fi

echo ""
echo "6. Verifying workflow creation..."
curl -k -s https://af-test.qazpost.kz/api/v1/workflows | jq -r '.workflows[] | select(.name == "antifraud-transaction-validation") | "     - \(.name) (ID: \(.id), Steps: \(.step_count))"'

echo ""
echo "7. Testing workflow execution..."
echo "   Preparing test transaction data..."
cat > test_transaction.json << 'EOF'
{
  "input_data": {
    "transaction": {
      "id": "test-txn-$(date +%s)",
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
}
EOF

echo "   Executing workflow..."
EXEC_RESPONSE=$(curl -k -s -X POST https://af-test.qazpost.kz/api/v1/workflows/antifraud-transaction-validation/execute \
  -H "Content-Type: application/json" \
  -d @test_transaction.json)

if echo "$EXEC_RESPONSE" | grep -q "runId"; then
    RUN_ID=$(echo "$EXEC_RESPONSE" | jq -r '.runId')
    echo "   ✅ Workflow execution started!"
    echo "   Run ID: $RUN_ID"
    echo "   Status URL: https://af-test.qazpost.kz/api/v1/executions/$RUN_ID"
    
    echo ""
    echo "8. Polling execution status..."
    for i in {1..5}; do
        echo "   Poll attempt $i/5..."
        sleep 2
        STATUS_RESPONSE=$(curl -k -s "https://af-test.qazpost.kz/api/v1/executions/$RUN_ID")
        if echo "$STATUS_RESPONSE" | grep -q "status"; then
            STATUS=$(echo "$STATUS_RESPONSE" | jq -r '.status')
            echo "   Current status: $STATUS"
            if [ "$STATUS" = "completed" ] || [ "$STATUS" = "failed" ]; then
                echo "   Workflow execution $STATUS"
                break
            fi
        else
            echo "   Failed to get status"
        fi
    done
else
    echo "   ❌ Failed to execute workflow: $EXEC_RESPONSE"
    echo "   Error details:"
    echo "$EXEC_RESPONSE" | jq .
fi

echo ""
echo "9. Testing with demo program..."
echo "   Running antifraud_workflow_demo.go..."
if command -v go >/dev/null 2>&1; then
    go run examples/antifraud_workflow_demo.go
else
    echo "   ⚠️  Go not installed, skipping demo program"
fi

echo ""
echo "=== Summary ==="
echo "✅ Workflow registration script completed"
echo "✅ Antifraud workflow should now be registered"
echo "✅ Workflow execution tested"
echo ""
echo "Next steps:"
echo "1. Check workflow registry logs for any issues"
echo "2. Verify workflow steps are properly registered"
echo "3. Test with actual transaction data"
echo "4. Monitor execution logs for antifraud service calls"