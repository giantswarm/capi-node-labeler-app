# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2024-01-29

### Changed

- Migrate from PSP to PSS. Add `global.podSecurityStandards.enforced` to facilitate the change.

## [0.4.0] - 2024-01-22

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

## [0.3.4] - 2021-10-12

### Added

- Added support for `node-role.kubernetes.io/control-plane` label.

## [0.3.3] - 2021-10-11

### Added

- Add `kubernetes.io/role` to applied labels.
- Add more verbose output.

## [0.3.2] - 2021-10-11

### Changed

- Fix Docker file.

## [0.3.1] - 2021-10-11

## [0.3.0] - 2021-10-11

### Changed

- Fix chart name.

## [0.2.0] - 2021-10-11

### Changed

- Use default-catalog.

## [0.1.0] - 2021-10-11

[Unreleased]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.3.4...v0.4.0
[0.3.4]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/capi-node-labeler-app/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/capi-node-labeler-app/releases/tag/v0.1.0
