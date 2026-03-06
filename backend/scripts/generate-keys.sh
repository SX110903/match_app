#!/bin/bash
# Generate RSA 4096-bit keypair for JWT RS256 signing
set -euo pipefail

KEYS_DIR="$(dirname "$0")/../keys"
mkdir -p "$KEYS_DIR"

echo "Generating RSA 4096-bit private key..."
openssl genrsa -out "$KEYS_DIR/private.pem" 4096

echo "Extracting public key..."
openssl rsa -in "$KEYS_DIR/private.pem" -pubout -out "$KEYS_DIR/public.pem"

echo "Generating 32-byte AES encryption key for TOTP secrets..."
ENCRYPTION_KEY=$(openssl rand -hex 32)
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo ""
echo "Add this to your .env file:"
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo ""
echo "Keys generated in $KEYS_DIR"
echo "IMPORTANT: Add /keys/ to .gitignore and never commit these files!"
