# Changelog

## [v0.5.1](https://github.com/yteraoka/yabumi/compare/v0.5.0...v0.5.1) - 2026-08-30

- feat: renovate で minimumReleaseAge を 7 days に設定 by @yteraoka in https://github.com/yteraoka/yabumi/pull/48
- Update goreleaser/goreleaser-action action to v7.1.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/49
- Update goreleaser/goreleaser-action action to v7.2.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/50
- Update dependency go to v1.26.3 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/51
- Update actions/create-github-app-token action to v3.2.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/52
- Update goreleaser/goreleaser-action action to v7.2.2 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/53
- Update golangci/golangci-lint-action action to v9.2.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/54
- Update dependency go to v1.26.4 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/55
- Update actions/checkout action to v6.0.3 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/56
- Update actions/checkout action to v7 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/57
- Update actions/setup-go action to v6.5.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/58
- Update goreleaser/goreleaser-action action to v7.2.3 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/59
- Update golangci/golangci-lint-action action to v9.3.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/60
- Rename homebrew-tap repository to homebrew-cask by @yteraoka in https://github.com/yteraoka/yabumi/pull/63
- Update dependency go by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/61
- Update actions/setup-go action to v7 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/62
- Update actions/checkout action to v7.0.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/64
- Update dependency go to v1.27.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/65
- tagpr ワークフローの追加と GitHub App シークレット名の変更 by @yteraoka in https://github.com/yteraoka/yabumi/pull/66
- tagpr 専用の GitHub App を使う by @yteraoka in https://github.com/yteraoka/yabumi/pull/68
- Renovate: minor 以下の更新を auto merge する by @yteraoka in https://github.com/yteraoka/yabumi/pull/67

## [v0.5.0](https://github.com/yteraoka/yabumi/compare/v0.4.0...v0.5.0) - 2026-04-11

- Update actions/create-github-app-token action to v3.1.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/46
- fix: create-github-app-token の app-id を client-id に変更 by @yteraoka in https://github.com/yteraoka/yabumi/pull/47
- Update dependency go to v1.26.2 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/45
- Update actions/setup-go action to v6.4.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/44

## [v0.4.0](https://github.com/yteraoka/yabumi/compare/v0.3.1...v0.4.0) - 2026-03-24

- fix: parseField() で区切り文字が不足する入力時に panic が発生する問題を修正 by @yteraoka in https://github.com/yteraoka/yabumi/pull/38
- fix: エラーハンドリングを統一し一貫性を持たせる by @yteraoka in https://github.com/yteraoka/yabumi/pull/39
- test: テストカバレッジを拡充する by @yteraoka in https://github.com/yteraoka/yabumi/pull/40
- feat: リトライに exponential backoff を追加し 4xx はリトライしないようにする by @yteraoka in https://github.com/yteraoka/yabumi/pull/41
- ci: PR 時に Go テストを実行する GitHub Actions ワークフローを追加 by @yteraoka in https://github.com/yteraoka/yabumi/pull/42
- refactor: 定数化・HTTP クライアント再利用・parseBool 簡略化 by @yteraoka in https://github.com/yteraoka/yabumi/pull/43

## [v0.3.1](https://github.com/yteraoka/yabumi/compare/v0.3.0...v0.3.1) - 2026-03-17

- homebrew-tools を homebrew-tap に rename by @yteraoka in https://github.com/yteraoka/yabumi/pull/30
- Update dependency go to v1.26.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/31
- Update actions/setup-go action to v6.3.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/32
- Update actions/create-github-app-token action to v3 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/34
- golangci-lint の指摘を修正 by @yteraoka in https://github.com/yteraoka/yabumi/pull/35
- GitHub Actions に golangci-lint ワークフローを追加 by @yteraoka in https://github.com/yteraoka/yabumi/pull/36
- Update golangci/golangci-lint-action action to v9 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/37

## [v0.3.0](https://github.com/yteraoka/yabumi/compare/v0.2.1...v0.3.0) - 2026-03-02

- Update actions/checkout action to v6 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/22
- Update actions/setup-go action to v6.1.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/21
- Update dependency go to v1.25.5 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/19
- Update actions/checkout action to v6.0.2 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/25
- Update actions/setup-go action to v6.2.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/23
- Update dependency go to v1.25.6 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/24
- Homebrew casks 対応 by @yteraoka in https://github.com/yteraoka/yabumi/pull/28
- Update dependency go to v1.26.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/26
- Update goreleaser/goreleaser-action action to v7 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/27
- README 更新 by @yteraoka in https://github.com/yteraoka/yabumi/pull/29

## [v0.2.1](https://github.com/yteraoka/yabumi/compare/v0.2.0...v0.2.1) - 2025-10-25

- Update actions/checkout action to v4.3.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/14
- Update actions/checkout action to v5 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/15
- Update goreleaser/goreleaser-action action to v6.4.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/17
- Update actions/setup-go action to v6 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/18
- Update dependency go to v1.25.3 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/16

## [v0.2.0](https://github.com/yteraoka/yabumi/compare/v0.1.6...v0.2.0) - 2025-08-09

- delete arm binary in make clean by @yteraoka in https://github.com/yteraoka/yabumi/pull/5
- go mod tidy by @yteraoka in https://github.com/yteraoka/yabumi/pull/6
- Configure Renovate by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/7
- Update module github.com/bitly/go-simplejson to v0.5.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/8
- Update module github.com/jessevdk/go-flags to v1.5.0 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/9
- Update module github.com/jessevdk/go-flags to v1.6.1 by @renovate[bot] in https://github.com/yteraoka/yabumi/pull/11
- use goreleaser by @yteraoka in https://github.com/yteraoka/yabumi/pull/12
- update go.mod by @yteraoka in https://github.com/yteraoka/yabumi/pull/13

## [v0.1.6](https://github.com/yteraoka/yabumi/compare/v0.1.5...v0.1.6) - 2019-04-03

- build arm binary by @yteraoka in https://github.com/yteraoka/yabumi/pull/4

## [v0.1.5](https://github.com/yteraoka/yabumi/compare/v0.1.4...v0.1.5) - 2019-04-03

- go module by @yteraoka in https://github.com/yteraoka/yabumi/pull/3

## [v0.1.4](https://github.com/yteraoka/yabumi/compare/v0.1.3...v0.1.4) - 2018-12-06

- Add `-m` / `--message` option by @yteraoka in https://github.com/yteraoka/yabumi/pull/1
- Bump version to 0.1.4 by @yteraoka in https://github.com/yteraoka/yabumi/pull/2

## [v0.1.3](https://github.com/yteraoka/yabumi/compare/v0.1.2...v0.1.3) - 2018-09-24

## [v0.1.2](https://github.com/yteraoka/yabumi/compare/v0.1.1...v0.1.2) - 2018-09-24

## [v0.1.1](https://github.com/yteraoka/yabumi/commits/v0.1.1) - 2018-09-24
