# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## 2.0.0 (2026-01-01)


### ⚠ BREAKING CHANGES

* **logging:** EnqueueFn signature changed to return error, affecting internal interfaces

### Features

* adding api key middleware ([f6c62eb](https://github.com/fardannozami/activity-tracker/commit/f6c62eb5658667e61bdc22e87f567b9ea3174566))
* adding Dockerfile ([032f7e5](https://github.com/fardannozami/activity-tracker/commit/032f7e540b1a71a7ba6b9860cfa0796681cb9a4a))
* adding httpapi ([cfd93e6](https://github.com/fardannozami/activity-tracker/commit/cfd93e6ee3ea9fe7cad4ff9628b5fc3c76aced30))
* **api:** add API hit logging with batching worker ([ef86967](https://github.com/fardannozami/activity-tracker/commit/ef86967085fe0fbbd3bf871f7cf164ad83800af0))
* **app:** add config field to app struct ([8b269db](https://github.com/fardannozami/activity-tracker/commit/8b269db3f87643920cdca7bd5f8900b1d9c8b3cf))
* **auth:** add JWT authentication for API access ([6216671](https://github.com/fardannozami/activity-tracker/commit/6216671cc42118f54fa9aabc24385f62514dcda8))
* **cache:** add Redis caching and usage tracking ([1bb79f1](https://github.com/fardannozami/activity-tracker/commit/1bb79f1206478ce42d024c22b4fc81863c6baeea))
* **db:** add api hit repository and client api key prefix migration ([8cf98e7](https://github.com/fardannozami/activity-tracker/commit/8cf98e7175c4a290436183d7e648ddda0ec93c84))
* **db:** add database migrations for clients and usage tracking ([a0b5597](https://github.com/fardannozami/activity-tracker/commit/a0b5597024e157b3ef1e9aead49ed77d200ca5d2))
* **db:** add error handling for duplicate client email ([12fa6e0](https://github.com/fardannozami/activity-tracker/commit/12fa6e07eaadd570ce0a276367527c95ce1498f8))
* docker compose and config ([02c1965](https://github.com/fardannozami/activity-tracker/commit/02c19656f33ee8a42c36ceedd4e04c9b9e0e43fe))
* **logging:** add error handling to hit recording and batching ([e5e1cd9](https://github.com/fardannozami/activity-tracker/commit/e5e1cd9173f61050f7498da688decf4aeee90cce))
* register ([ddbcbf8](https://github.com/fardannozami/activity-tracker/commit/ddbcbf8d16c80cb34e583351e53a08ba3ffbaf45))
* **usage:** add usage tracking endpoints ([50c1b4a](https://github.com/fardannozami/activity-tracker/commit/50c1b4acba91f5cc0201352cc0cf1cede8208da3))
