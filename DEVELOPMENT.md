# Development Guide

## Architecture Overview

money-backword follows a layered architecture pattern:

```
┌─────────────────────────────┐
│   CLI / HTTP API Layer      │
│   (cmd/, internal/api/)     │
├─────────────────────────────┤
│   Business Logic Layer      │
│   (internal/ledger/)        │
├─────────────────────────────┤
│   Storage Abstraction       │
│   (internal/storage/)       │
├─────────────────────────────┤
│   Persistence Layer         │
│   (JSON, SQLite planned)    │
└─────────────────────────────┘
```

## Package Structure

### `internal/ledger`
Core domain models representing the business domain:
- `Account`: Represents a financial account
- `Transaction`: A single monetary transaction
- `Category`: Expense/income classification
- `Budget`: Monthly spending limits

These are pure Go structs with validation methods.

### `internal/storage`
Storage abstraction layer. Defines a `Store` interface with implementations:
- `JSONStore`: File-based JSON persistence

Future implementations:
- `SQLiteStore`: Database-backed persistence

### `internal/api`
HTTP API handlers and routing. Implements RESTful endpoints for:
- `/api/transactions` — transaction CRUD
- `/api/accounts` — account management
- `/api/budgets` — budget management
- `/api/report` — financial reporting

### `internal/report`
Report generation and financial analysis:
- Monthly summaries
- Category analysis
- Trend calculations
- Future: CSV/PDF export

### `cmd/moneyback`
CLI application. Parses command-line flags and delegates to business logic.

## Development Workflow

### Adding a New Feature

1. **Define the domain model** in `internal/ledger/`
2. **Add storage methods** to the `Store` interface and implementations
3. **Add CLI commands** or API handlers
4. **Write tests** for the feature
5. **Update documentation** in README and this file

### Testing

```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -cover ./...

# Coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Style

Follow [Effective Go](https://golang.org/doc/effective_go):
- Use `gofmt` for formatting
- Use `go vet` to check code quality
- Prefer explicit error handling
- Use table-driven tests
- Keep functions small and focused

### Database Schema (for future SQLite implementation)

```sql
-- Accounts
CREATE TABLE accounts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    balance DECIMAL(15, 2),
    currency TEXT DEFAULT 'USD',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Transactions
CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL REFERENCES accounts(id),
    category TEXT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    description TEXT,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- Categories
CREATE TABLE categories (
    name TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    color TEXT,
    icon TEXT,
    active BOOLEAN DEFAULT TRUE
);

-- Budgets
CREATE TABLE budgets (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    monthly_limit DECIMAL(15, 2),
    alert_threshold DECIMAL(3, 2),
    start_date DATE,
    end_date DATE,
    active BOOLEAN DEFAULT TRUE
);
```

## Known TODOs

- [ ] Implement CSV export for reports
- [ ] Add SQLite support as alternative storage
- [ ] Implement budget spending tracking
- [ ] Add forecasting based on historical trends
- [ ] Support recurring transactions
- [ ] Multi-currency support
- [ ] Import from CSV/OFX format
- [ ] Transaction reconciliation
- [ ] API authentication

## Performance Considerations

- JSON store suitable for personal/small-team use (<10k transactions)
- For larger datasets, plan migration to SQLite
- All in-memory operations for speed
- Consider caching for frequently generated reports

## Security Notes

- Input validation on all user-provided data
- JSON store file permissions should be restricted (mode 0600)
- Sensitive data should be handled carefully
- No encryption at rest currently (consider for future versions)

## Contributing

See [CONTRIBUTORS.md](CONTRIBUTORS.md) and README.md for contribution guidelines.
