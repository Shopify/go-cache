# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Fix docs for the client `Get` method

## [v2.2.0] - 2021-09-28
### Added
- Add a cache client mock

## [v2.1.0] - 2021-09-09
### Added
- Add a `prefixClient` wrapper
- Add a `cachelock` package, to acquire and release locks on keys using a `cache.Client`

## [v2.0.1] - 2021-05-26
### Changed
- Fix `go.mod` path

## [v2.0.0] - 2021-05-26
### Added
- `go mod tidy` CI step
- Add dependabot

### Changed
- **BREAKING**: Add `context.Context` as a parameter to all methods
- Bump requirement to `go1.15`
- Use [`go-encoding`](https://github.com/Shopify/go-encoding) instead of custom encoding
- Mark `TtlForExpiration` as deprecated, it does not need to be exported
- Switch to golangci for linting
- Switch to Github action for CI
- Split integration tests
  - Run unit tests on all go versions
  - Run memcached tests on latest go version
  - Run redis tests on latest go version and redis 4, 5, 6
- Many dependency updates

### Removed
- **BREAKING:** Remove deprecated encodings
- **BREAKING:** Remove deprecated `TtlForExpiration`, use `time.Until` instead.

## [v1.0.0] - 2020-09-23
### Added
- Add a CHANGELOG.md
- Upgrade to `go-redis` v8
- Remove the `pkg` folder from the path

## [v0.1.0] - 2020-09-22
### Added
- `Get`/`Set`/`Add`/`Delete`/`Increment`/`Decrement` methods
- `Memory` and `Redis` clients

### Changed
- Require `go1.13+`
