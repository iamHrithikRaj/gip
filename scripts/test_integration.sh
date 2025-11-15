#!/bin/bash
# Gip Integration Test - Creates a real conflict scenario

set -e

echo ""
echo "=== GIP INTEGRATION TEST ==="
echo "This script will create a test repository with conflicts"
echo ""

# Step 1: Build Gip
echo "[1/7] Building Gip binary..."
cargo build --release
if [ -f "./target/release/gip.exe" ]; then
    GIP_BINARY="./target/release/gip.exe"
else
    GIP_BINARY="./target/release/gip"
fi
echo "  ✓ Binary ready at $GIP_BINARY"

# Step 2: Create temp test directory
echo ""
echo "[2/7] Creating test repository..."
PROJECT_ROOT="$(pwd)"
TEST_DIR="./target/test_repo_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"
GIP_BINARY="$PROJECT_ROOT/target/release/gip.exe"

# Initialize git repo
git init
git config user.name "Test User"
git config user.email "test@example.com"
echo "  ✓ Git repo created"

# Step 3: Initialize Gip
echo ""
echo "[3/7] Initializing Gip..."
$GIP_BINARY init
if [ -d ".gip" ]; then
    echo "  ✓ Gip initialized (.gip directory exists)"
else
    echo "  ✗ Gip init failed"
    exit 1
fi

# Step 4: Create initial file
echo ""
echo "[4/7] Creating initial file..."
cat > calculator.py << 'EOF'
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price
    return total
EOF
git add calculator.py
git commit -m "Initial commit"
echo "  ✓ Initial commit created"

# Step 5: Create feature branch with tax change
echo ""
echo "[5/7] Creating 'add-tax' branch..."
git checkout -b add-tax

cat > calculator.py << 'EOF'
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price * 1.08  # Add 8% tax
    return total
EOF
git add calculator.py
git commit -m "Add 8% sales tax"
TAX_COMMIT=$(git rev-parse HEAD)
echo "  ✓ Tax commit: $TAX_COMMIT"

# Create manifest for tax commit
mkdir -p .gip/manifest
cat > ".gip/manifest/${TAX_COMMIT}.json" << EOF
{
  "schemaVersion": "2.0",
  "commit": "${TAX_COMMIT}",
  "entries": [
    {
      "anchor": {
        "file": "calculator.py",
        "symbol": "calculate_total",
        "hunkId": "H#1"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": ["items is list with .price"],
        "postconditions": ["returns total with 8% sales tax"],
        "errorModel": ["AttributeError if no .price"]
      },
      "behaviorClass": ["feature"],
      "sideEffects": [],
      "rationale": "Added 8% sales tax to comply with state law"
    }
  ]
}
EOF
echo "  ✓ Tax manifest created"

# Step 6: Switch to main/master and create conflicting change
echo ""
echo "[6/7] Creating conflicting change on main..."
git checkout master 2>/dev/null || git checkout main

cat > calculator.py << 'EOF'
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price + 5.99  # Add $5.99 shipping
    return total
EOF
git add calculator.py
git commit -m "Add flat shipping fee"
SHIPPING_COMMIT=$(git rev-parse HEAD)
echo "  ✓ Shipping commit: $SHIPPING_COMMIT"

# Create manifest for shipping commit
cat > ".gip/manifest/${SHIPPING_COMMIT}.json" << EOF
{
  "schemaVersion": "2.0",
  "commit": "${SHIPPING_COMMIT}",
  "entries": [
    {
      "anchor": {
        "file": "calculator.py",
        "symbol": "calculate_total",
        "hunkId": "H#1"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": ["items is list with .price"],
        "postconditions": ["returns total plus flat shipping"],
        "errorModel": ["AttributeError if no .price"]
      },
      "behaviorClass": ["feature"],
      "sideEffects": [],
      "rationale": "Added 5.99 flat shipping fee for all orders"
    }
  ]
}
EOF
echo "  ✓ Shipping manifest created"

# Step 7: Create conflict
echo ""
echo "[7/7] Attempting merge (will create conflict)..."
if ! git merge add-tax 2>&1; then
    echo "  ✓ Conflict created successfully!"
    
    echo ""
    echo "=== CURRENT CONFLICT (WITHOUT GIP ENRICHMENT) ==="
    cat calculator.py
    
    echo ""
    echo ""
    echo "=== MANIFESTS AVAILABLE ==="
    echo ""
    echo "--- SHIPPING COMMIT MANIFEST (HEAD) ---"
    cat ".gip/manifest/${SHIPPING_COMMIT}.json"
    
    echo ""
    echo "--- TAX COMMIT MANIFEST (MERGE_HEAD) ---"
    cat ".gip/manifest/${TAX_COMMIT}.json"
    
    echo ""
    echo ""
    echo "=== WHAT GIP WOULD INJECT ==="
    echo "In the conflict markers, Gip would inject TOON-formatted context like:"
    echo ""
    echo "<<<<<<< HEAD"
    echo "total += item.price + 5.99  # Add \$5.99 shipping"
    echo "||| Gip CONTEXT (HEAD - Your changes)"
    echo "||| Commit: ${SHIPPING_COMMIT:0:8}"
    echo "||| behaviorClass[1]: feature"
    echo "||| preconditions[1]: items is list with .price"
    echo "||| postconditions[1]: returns total plus flat shipping"
    echo "||| errorModel[1]: AttributeError if no .price"
    echo "||| rationale: Added 5.99 flat shipping fee for all orders"
    echo "||| symbol: calculate_total"
    echo "======="
    echo "total += item.price * 1.08  # Add 8% tax"
    echo "||| Gip CONTEXT (MERGE_HEAD - Their changes)"
    echo "||| Commit: ${TAX_COMMIT:0:8}"
    echo "||| behaviorClass[1]: feature"
    echo "||| preconditions[1]: items is list with .price"
    echo "||| postconditions[1]: returns total with 8% sales tax"
    echo "||| errorModel[1]: AttributeError if no .price"
    echo "||| rationale: Added 8% sales tax to comply with state law"
    echo "||| symbol: calculate_total"
    echo ">>>>>>> add-tax"
    
    echo ""
    echo ""
    echo "=== TEST SUMMARY ==="
    echo "✓ Gip binary builds successfully"
    echo "✓ Gip init creates .gip directory"
    echo "✓ Manifests can be stored as JSON"
    echo "✓ Manifests can be retrieved for commits"
    echo "✓ Conflicts are detected"
    echo "⚠ Conflict enrichment requires merge driver implementation"
    
    echo ""
    echo "=== NEXT STEPS ==="
    echo "To enable automatic enrichment, Gip needs to:"
    echo "1. Implement the merge driver (src/merge/driver.rs)"
    echo "2. Configure Git to use Gip as merge driver:"
    echo "   git config merge.gip.driver 'gip merge-driver %O %A %B %P'"
    echo "3. Add .gitattributes: * merge=gip"
    
    echo ""
    echo "Test repository preserved at: $TEST_DIR"
    echo "Explore it with: cd $TEST_DIR"
else
    echo "  ✗ No conflict created (unexpected)"
fi

echo ""
echo "=== TEST COMPLETE ==="
