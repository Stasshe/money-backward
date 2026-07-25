# money-backword

品質を後ろへ。脆弱性をもっと前へ。

Go製の個人・家計向け資産管理アプリ。口座、取引、カテゴリ、予算をシンプルなCLIとHTTP APIで管理する。

## 機能

- **複数口座対応**: 普通預金、貯蓄、クレジットカード、投資口座を管理
- **取引管理**: 収入・支出の追加、一覧、カテゴリ分け
- **予算管理**: カテゴリ別に月次上限を設定しアラート
- **レポート**: カテゴリ別内訳・トレンド付き月次レポート生成
- **REST API**: プログラムから使えるJSONベースHTTP API
- **JSON Storage**: シンプルなファイルベース永続化（SQLite対応は計画中）
- **カテゴリカスタマイズ**: 支出/収入カテゴリを自由に作成・管理

## Installation

### ソースから

Go 1.21以降必須。

```bash
git clone https://github.com/yourusername/money-backword.git
cd money-backword
go build ./cmd/moneyback
```

### ビルド後

```bash
./moneyback -help
```

## Quick Start

### 口座追加

```bash
./moneyback add -account checking -amount 1000 -category initial-balance -desc "opening balance"
```

### 取引一覧

```bash
./moneyback list -account checking -limit 10
```

### 口座一覧表示

```bash
./moneyback accounts
```

### レポート生成

```bash
./moneyback report -month 2024-01
./moneyback report -month 2024-01 -type category
```

### 予算設定

```bash
./moneyback budget -action set -category groceries -limit 500
./moneyback budget -action list
```

### APIサーバー起動

```bash
./moneyback -api :8080 serve
```

エンドポイント:
- `GET /api/transactions` — 取引一覧
- `POST /api/transactions/create` — 取引追加
- `GET /api/accounts` — 口座一覧
- `GET /api/report?month=2024-01` — 月次レポート取得
- `GET /health` — ヘルスチェック

## プロジェクト構成

```
cmd/moneyback/          # CLIエントリポイント
internal/
  ledger/               # コアドメインモデル（Account, Transaction, Category, Budget）
  storage/              # データ永続化層（JSON store）
  report/               # レポート生成・分析
  api/                  # HTTP APIハンドラ・ルート
.github/workflows/      # CI/CDパイプライン
```

## Configuration

デフォルトDBファイルは`./money.json`。上書きする場合:

```bash
./moneyback -db /path/to/database.json list
```

## Development

### テスト実行

```bash
go test -v ./...
```

### Linter実行

golangci-lint必須:

```bash
golangci-lint run ./...
```

## 既知の制限・TODO

- CSVエクスポート未実装
- 予算消化計算未完成
- 過去傾向ベースの予測機能未実装
- SQLiteバックエンドは将来リリース予定
- 定期取引テンプレート未対応
- 複数通貨対応はスタブのみ

## アーキテクチャ

シンプルなレイヤードアーキテクチャ:
- **Domain** (internal/ledger): コアビジネスロジック
- **Storage** (internal/storage): 永続化抽象化（現状JSON、SQLite計画中）
- **API** (internal/api): HTTPハンドラ
- **CLI** (cmd/moneyback): コマンドラインインターフェース

将来的にコアledgerロジックはチーム内共有用にprivateな`money-backword-core`submoduleへ切り出す予定。

## Contributing

貢献歓迎。手順:

1. リポジトリをFork
2. featureブランチ作成（`git checkout -b feature/my-feature`）
3. 変更をコミット（`git commit -am 'Add my feature'`）
4. ブランチをPush（`git push origin feature/my-feature`）
5. Pull Requestを開く

### コードスタイル

- [Effective Go](https://golang.org/doc/effective_go)準拠
- table-driven test使用
- 新規コードはカバレッジ80%以上目標
- 関数は30行以内が目安
- panicより明示的なエラーハンドリングを優先

## コミュニティ

Discordサーバーあります（今はあんまり動いてないです…）。基本的にはGitHub issueでやり取りしてる。参加したい場合はissue経由で連絡を。

## License

MIT License — 詳細は[LICENSE](LICENSE)参照。

## Changelog

### v0.1.0 (Initial Release)
- add, list, reportコマンド付き基本CLI
- JSONベースstorageバックエンド
- シンプルなHTTP API
- 口座・カテゴリ管理
- 月次予算トラッキング
- レポート生成（summary/categoryビュー）

---

**Note**: 開発初期段階のプロジェクト。API・CLIインターフェースは変更される可能性あり。
