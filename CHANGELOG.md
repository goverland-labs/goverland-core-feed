# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Changed the path name of the go module
- Updated dependencies versions
- Migrated protocol to this repo as different module

### Added
- Added LICENSE information
- Added info for contributing
- Added github issues templates
- Added linter and unit-tests workflows for github actions
- Added badges with link to the license and passed workflows

### Added
- Feed caching by filters 

## [0.0.25] - 2024-02-09

### Changed
- Changed dao feed sorting from created_at desc to state + created_at

## [0.0.24] - 2024-02-08

### Added
- Filter spam and canceled proposals in feed

## [0.0.23] - 2024-02-06

### Changed
- Up platform events library

## [0.0.22] - 2023-12-06

### Changed
- Actualize source libraries

## [0.0.21] - 2023-11-07

### Changed
- Actualize active feed items condition

## [0.0.20] - 2023-10-07

### Fixed
- Sort order of feeds

## [0.0.19] - 2023-09-18

### Changed
- Extend get feed filters

## [0.0.18] - 2023-08-23

### Added
- Send timeline to the queue

## [0.0.17] - 2023-08-23

### Fixed
- Fixed sending proposals updates to the pipeline if was added non-unique action to the timeline

## [0.0.16] - 2023-08-14

### Added
- Handle proposal ends soon event

## [0.0.15] - 2023-07-15

### Fixed
- Fixed getting subscriber by id in the repo
- Updated platform-events dependency to v0.1.0
- Fixed saving new feed items (on duplicate)
- Fixed payload which is sending to the webhook execution

### Changed
- The feed.id field is exported now

### Added
- Added consumers limits (max ack pendings, rate limit, ack wait)
- Added feed.triggered_at

## [0.0.14] - 2023-07-14

### Fixed
- Fixed feed timeline action to proto conversion

## [0.0.13] - 2023-07-14

### Fixed
- Updated platform-events dependency to v0.0.15

## [0.0.12] - 2023-07-14

### Fixed
- Fixed subscription.id generation

## [0.0.11] - 2023-07-14

### Fixed
- Fixed uuid fields
- Updated core-api dependency to v0.0.11

## [0.0.10] - 2023-07-14

### Fixed
- Updated platform-events dependency to v0.0.13

## [0.0.9] - 2023-07-14

### Added
- Added exporting feed item timeline in gRPC server

### Changed
- Updated core-api protocol version to v0.0.9
- Updated structure of the dao feed
- Updates internal structure and handling feed items

## [0.0.8] - 2023-07-12

### Fixed
- Updated platform-events dependency to v0.0.13

## [0.0.7] - 2023-07-12

### Fixed
- Updated platform-events dependency to v0.0.12

## [0.0.6] - 2023-07-11

### Fixed
- Fixed DeletedAt field in subscriber model
- Fixed Dockerfile
- Updated platform-events dependency to v0.0.11

## [0.0.5] - 2023-07-06

### Added
- Flat feed

## [0.0.4] - 2023-06-16

### Changed
- Rename repository from feed to core-feed

## [0.0.3] - 2023-06-14

### Added
- Inject core-api library

## [0.0.2] - 2023-06-07

### Added
- Manipulating with subscriptions
- Simple auth by subscriber identifier
- Publishing callback events

## [0.0.1] - 2023-05-29

### Added
- Handling daos and proposals feed items
- Registering and updating subscribers
