# redmine api クライアント sdk の実装

## 役割

あなたは、Model Context Protocol について詳しく、Go らしい実装が得意で Go を効率的に実装できるシニアエンジニアです。

## 実装方針

- https://github.com/modelcontextprotocol/go-sdk (v1.0.0) を使用して Redmine の MCP サーバーを実装してください。
- redmine の api クライアント sdk は `./pkg/redmine` に実装してあります（76メソッド、全22API）。
- 実装するツールは29個で、以下のサービスに絞ります：
  - Projects（7ツール）
  - Issues（7ツール）
  - Users（6ツール）
  - Issue Categories（5ツール）
  - Search（2ツール）
  - Trackers（1ツール）
  - Issue Statuses（1ツール）

## ディレクトリ構造

```
redmine-go/
├── cmd/
│   ├── mcp-server/          # MCPサーバーのエントリーポイント
│   │   └── main.go          # MCP server起動ロジック
│   └── cli/                 # 将来のCLI用（未実装）
│
├── internal/
│   ├── mcp/                 # MCP固有の実装
│   │   ├── server.go        # MCPサーバーの初期化・設定
│   │   ├── handlers/        # MCPツールハンドラー（サービス別）
│   │   │   ├── projects.go     # Projectsツール（7ツール）
│   │   │   ├── issues.go       # Issuesツール（7ツール）
│   │   │   ├── users.go        # Usersツール（6ツール）
│   │   │   ├── categories.go   # Issue Categoriesツール（5ツール）
│   │   │   ├── search.go       # Searchツール（2ツール）
│   │   │   ├── trackers.go     # Trackersツール（1ツール）
│   │   │   └── statuses.go     # Issue Statusesツール（1ツール）
│   │   └── types.go         # MCP用の共通型定義
│   │
│   ├── usecase/             # ビジネスロジック（MCP/CLI共通）
│   │   ├── project.go       # プロジェクト操作のユースケース
│   │   ├── issue.go         # チケット操作のユースケース
│   │   ├── user.go          # ユーザー操作のユースケース
│   │   ├── category.go      # カテゴリ操作のユースケース
│   │   ├── search.go        # 検索操作のユースケース
│   │   └── metadata.go      # メタデータ取得（trackers, statuses）
│   │
│   └── config/              # 設定管理
│       ├── config.go        # 設定構造体
│       └── loader.go        # 環境変数からの設定読み込み
│
├── pkg/
│   └── redmine/             # Redmine APIクライアントSDK（既存）
```

### データフロー

```
MCP Client (Claude)
  ↓
cmd/mcp-server/main.go
  ↓
internal/mcp/server.go
  ↓
internal/mcp/handlers/*.go
  ↓
internal/usecase/*.go
  ↓
pkg/redmine/*.go
  ↓
Redmine API Server
```

## 注意点

- エントリーポイントは `./cmd/mcp-server/main.go` に作成
- 将来のCLI実装を考慮し、`internal/usecase/`にビジネスロジックを集約
- `internal/mcp/handlers/`はMCP固有の入出力型変換のみを担当
- 設定は環境変数から読み込み（`REDMINE_URL`, `REDMINE_API_KEY`）
- 実装はサービスごとに進める
- 実装方針を todo 化してそれに従って丁寧に進める
- コミットは明確に適切な粒度でおこなう
- コミット前には `golangci-lint run` と `go test ./...` が通ることを担保する
