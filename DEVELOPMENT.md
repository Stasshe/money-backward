# Development Guide

## アーキテクチャ概要

money-backwordはレイヤードアーキテクチャ:

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

## パッケージ構成

### `internal/ledger`
ビジネスドメイン表すコアモデル:
- `Account`: 口座
- `Transaction`: 1件の取引
- `Category`: 支出/収入分類
- `Budget`: 月次予算上限

バリデーションメソッド付き素のGo struct。

### `internal/storage`
Storage抽象化層。`Store`インターフェースと実装:
- `JSONStore`: ファイルベースJSON永続化

将来実装予定:
- `SQLiteStore`: DBバックエンド永続化

### `internal/api`
HTTP APIハンドラ・ルーティング。RESTfulエンドポイント:
- `/api/transactions` — transaction CRUD
- `/api/accounts` — 口座管理
- `/api/budgets` — 予算管理
- `/api/report` — 財務レポート

### `internal/report`
レポート生成・財務分析:
- 月次サマリー
- カテゴリ別分析
- トレンド計算
- 将来: CSV/PDFエクスポート

### `cmd/moneyback`
CLIアプリ。コマンドラインフラグ解析しビジネスロジックへ委譲。

## 開発フロー

### 新機能追加時

1. `internal/ledger/`にドメインモデル定義
2. `Store`インターフェースと実装にstorageメソッド追加
3. CLIコマンドまたはAPIハンドラ追加
4. テスト作成
5. READMEと本ファイルのドキュメント更新

### テスト

```bash
# 全テスト実行
go test -v ./...

# race detection付き
go test -race ./...

# カバレッジレポート生成
go test -cover ./...

# カバレッジHTMLレポート
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### コードスタイル

[Effective Go](https://golang.org/doc/effective_go)準拠:
- フォーマットは`gofmt`
- コード品質チェックは`go vet`
- エラーハンドリングは明示的に
- table-driven test使う
- 関数は小さく単一責務に保つ

### DBスキーマ（将来のSQLite実装用）

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

## 既知のTODO

- [ ] レポートのCSVエクスポート実装
- [ ] SQLite対応をstorageの選択肢として追加
- [ ] 予算消化トラッキング実装
- [ ] 過去傾向ベースの予測機能追加
- [ ] 定期取引対応
- [ ] 複数通貨対応
- [ ] CSV/OFX形式からのインポート
- [ ] 取引照合機能
- [ ] API認証

## パフォーマンス上の注意

- JSON storeは個人/小規模チーム利用向け（1万件未満想定）
- 大規模データセットはSQLiteへの移行を検討
- 速度優先で全てインメモリ操作
- 頻繁に生成するレポートはキャッシュ検討

## セキュリティ上の注意

- ユーザー入力は全てバリデーション
- JSON storeファイルのパーミッションは制限すべき（mode 0600）
- 機密データの取り扱いは慎重に
- 現状保存時暗号化なし（将来バージョンで検討）

## Contributing

貢献ガイドラインは[CONTRIBUTORS.md](CONTRIBUTORS.md)とREADME.md参照。
