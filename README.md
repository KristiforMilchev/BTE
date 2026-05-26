# BTE

BTE is the Blockcert Trading Environment: a terminal-based crypto trading terminal for safe, seamless Ledger interactions. BTE is designed to run as a kiosk application under the Blockcert Trading distribution, where it is part of the broader application suite. It is also released as a standalone application for users who want to run it on Linux, macOS, or Windows.

With the BTE distribution, you can turn a computer into a secure trading terminal using a simple USB device. For now, you can build and burn the image yourself. Once the distribution is fully tested and signed for Secure Boot, this repository will be updated with links to the latest image.

It is built for operators who want a focused keyboard-driven workflow for managing networks, assets, contacts, transfers, and contract interaction checks without leaving the terminal.

## Features

- Terminal UI for wallet operations, built with Bubble Tea and Lip Gloss.
- Ledger wallet integration for signing native and ERC-20 transfers.
- Network management with saved RPC endpoints and active network switching.
- Native asset and imported token asset list, with native token pinned at the top.
- Automatic balance refresh after sends and token imports.
- Contact book with centered popup creation flow.
- Transaction history for native and token transfers, scoped by account and network.
- ERC-20 token import by contract address.
- Imported token persistence per network.
- ERC-20 metadata loading for name, symbol, decimals, and balances.
- Smart contract import view with callable/readable method panels and 24h interaction panel.
- Fork-backed simulation flow through the Blockcert API.
- Ledger-signed simulation transactions against cloned RPCs, not the live network.
- Function-level simulation reports showing calldata, pass/revert status, and observed consequences.
- Behavioral approval reporting that distinguishes expected direct approvals from unexpected approval-like behavior.
- Local SQLite persistence for networks, contacts, tokens, and transactions.

## What BTE Is For

BTE is designed around one practical goal: operational security. No browser extensions, no wallet popups, and no extra web surface area; just a direct terminal interface for Ledger-backed trading workflows.

## Current Workflows

### Assets

The assets panel shows the native network token first, followed by imported tokens saved for the active network. Imported ERC-20 tokens can be added from the contract import screen and are saved locally.

### Contacts

Contacts can be created from the dashboard with a simple popup containing:

- Name
- Address

Contacts are used as transfer recipients.

### Sending

BTE supports:

- Native network-token sends.
- ERC-20 `transfer(address,uint256)` sends.
- NFT transfers for ERC-721 and ERC-1155 are planned.

Both flows use the same confirmation screen and Ledger signing path. After a successful send, balances and transaction history refresh automatically.

### Contract Import And Simulation

The import contract screen lets the user enter a contract address, inspect placeholder callable/readable method groups, load recent interactions, and run simulation with `S`.

Simulation is API-backed:

- The API spins up a forked RPC.
- The caller can be seeded with fake native token on the fork for gas, hidden from the final user report.
- BTE signs fork transactions with Ledger.
- The API executes function calls on the fork and returns a report.
- The report returns to the import page so a clean result can be saved with `Y`.

Simulation features are intended for members who support the project.

## Persistence
Saved data includes:

- Accounts
- Networks
- Tokens
- Contacts
- Native transfer history
- Token transfer history

By default the app uses:

```text
data/bte.db
```

## Development

Run the terminal app:

```sh
go run .
```

Run the test suite:

```sh
go test ./...
```

The Go module is named:

```text
bte
```

## Status

BTE is an active terminal-first trading environment. Core wallet, asset, contact, token import, transaction persistence, and simulation workflows are implemented, with deeper bytecode-only behavioral analysis continuing to evolve.
