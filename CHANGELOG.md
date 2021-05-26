# Changelog

## v2

- Add context as first argument
- Remove deprecated encoding, use github.com/Shopify/go-encoding directly instead.
- Remove deprecated TtlForExpiration, use `time.Until` instead.
- Bump required Go version to 1.15

## v1

- First stable release
- [Remove the pkg folder](https://github.com/Shopify/go-cache/pull/10), clients will need to update their import paths.
- [Upgrade to go-redis v8](https://github.com/Shopify/go-cache/pull/8)
