## Build, Lint, and Test

- **Build:** `go build ./...`
- **Lint:** `go vet ./...` and `staticcheck ./...`
- **Test:** `go test ./...`
- **Test a single file:** `go test [path-to-file]`
- **Test a single test:** `go test -run [test-name]`

## Code Style

- **Imports:** Standard Go import grouping.
- **Formatting:** Consistent use of tabs for indentation. Adhere to `gofmt` standards.
- **Types:** Static typing is preferred. Interfaces are used to decouple components (e.g., `core.Sequenceable`).
- **Naming Conventions:**
    - **Variables:** camelCase (e.g., `lastName`).
    - **Functions:** PascalCase for exported functions (e.g., `NewNote`) and camelCase for internal functions.
    - **Structs:** PascalCase (e.g., `Note`).
    - **Interfaces:** Follow Go conventions (e.g., `Storable`, `Playable`).
- **Error Handling:** Errors are returned as the second value from functions. Use `fmt.Errorf` for creating new error messages.
- **Comments:** Document exported functions and structs, explaining their purpose and usage.
- **Logging:** Use the `notify` package for logging and debugging. `notify.Debugf` for debug messages, and `notify.Errorf` for errors.
