# Go Authenticator

A simple CLI-based Google Authenticator–like TOTP app built in Go.  
Supports:
- Adding accounts manually with a Base32 secret
- Adding accounts by scanning QR code images (`otpauth://` format)
- Encrypted local storage (AES-GCM + scrypt key derivation)
- Generating current TOTP codes
- Validating TOTP codes

---

## Features
- **Secure**: Stores secrets encrypted with a master passphrase.
- **Cross-platform**: Works on Linux, macOS, and Windows.
- **Interoperable**: Works with any service supporting TOTP (e.g., Stripe, GitHub, Google, AWS).
- **QR Code Parsing**: Import accounts by scanning authenticator QR codes.
