## v1.6.0 (2024-08-06)

### Feat

- add exception for mime.ErrInvalidMediaParameter

## v1.5.2 (2024-05-30)

### Fix

- Also apply allowed list to value switch statements

## v1.5.1 (2024-05-02)

### Fix

- remove init flagset
- exposed AllowPair fields
- remove init function

## v1.5.0 (2024-05-01)

### Feat

- extends Analyzer with custom allowed errors via functional options
- add few more allowed errors from std

### Refactor

- speed up isAllowedErrAndFunc via map

## v1.4.8 (2024-01-25)

### Fix

- Allow io.EOF from (io.ReadCloser).Read

## v1.4.7 (2023-12-05)

### Fix

- Add context errors as allowed direct comparisons (#63)

## v1.4.6 (2023-12-01)

### Fix

- Allow type assertions of non-error types (fixes #61)

## v1.4.5 (2023-09-08)

### Fix

- Infinite recursion in assignment finder (fixes #57)

## v1.4.4 (2023-08-20)

### Fix

- Is methods should be exempt for type assertions and switches too (#50)
- add missing testdata
- ignore unix errno values

### Refactor

- pass extinfo to type assertions
- Permit matching on full paths of allowed errors

## v1.4.3 (2023-06-30)

### Fix

- ignore io.EOF from some more sources

## v1.4.2 (2023-05-12)

### Fix

- Allowed checker panics in some cases (fixes #42)

## v1.4.1 (2023-05-01)

### Fix

- Tar allowlist should be archive/tar

## v1.4.0 (2023-03-13)

### Feat

- Require Go 1.20 in go.mod

## v1.3.0 (2023-03-07)

### Feat

- Suggest fixes for missing wrap verbs
- Stop linting arguments that are .Error() expressions

## v1.2.0 (2023-02-25)

### Feat

- Add the -errorf-multi flag

## v1.1.0 (2023-02-10)

### Feat

- Look for multiple %w verbs as these are valid starting with Go 1.20

## v1.0.6 (2022-11-24)

## v1.0.5 (2022-09-29)

## v1.0.4 (2022-09-16)

## v1.0.3 (2022-09-09)

## v1.0.2 (2022-08-12)

## v1.0.1 (2022-08-11)

## v1.0.0 (2022-05-10)
