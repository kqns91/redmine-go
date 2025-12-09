# redmine-go

Redmine API の非公式 Go クライアント

[![Go Version](https://img.shields.io/badge/Go-1.25.2%2B-00ADD8?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/kqns91/redmine-go.svg)](https://pkg.go.dev/github.com/kqns91/redmine-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/kqns91/redmine-go)

[English](README.md) | 日本語

## 概要

`redmine-go` は Go で書かれた Redmine REST API クライアントです。Redmine との連携方法を3つの形態で提供しています：

- **SDK** - Redmine と連携するアプリケーションを構築するための Go パッケージ
- **CLI** - ターミナルから Redmine を管理するためのコマンドラインツール
- **MCP サーバー** - Model Context Protocol を使用して AI アシスタントが Redmine と連携するためのサーバー実装

3つの形態はすべて同じ SDK 基盤の上に構築されており、22 の Redmine REST API と 76 のメソッドをサポートしています。

---

## SDK

Redmine REST API と連携するための Go クライアントパッケージです。

### インストール

```bash
go get github.com/kqns91/redmine-go
```

### 基本的な使い方

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kqns91/redmine-go/pkg/redmine"
)

func main() {
    client := redmine.New("https://your-redmine.com", "your-api-key")
    ctx := context.Background()

    // プロジェクト一覧の取得
    projects, err := client.ListProjects(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, project := range projects.Projects {
        fmt.Printf("%s (ID: %d)\n", project.Name, project.ID)
    }

    // 課題の作成
    issue := redmine.Issue{
        ProjectID:   1,
        Subject:     "サンプル課題",
        Description: "課題の説明",
    }

    created, err := client.CreateIssue(ctx, issue)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created issue #%d\n", created.Issue.ID)
}
```

### サポートしている API

SDK は以下の Redmine REST API をサポートしています：

**コアリソース**
- Projects（CRUD、アーカイブ/アンアーカイブ）
- Issues（CRUD、ウォッチャー）
- Users（CRUD）
- Time Entries（CRUD）

**プロジェクト管理**
- Versions（CRUD）
- Issue Relations（CRUD）
- Memberships（CRUD）
- Issue Categories（CRUD）

**コンテンツ**
- Wiki Pages（CRUD）
- News（読み取り）
- Files（読み取り、アップロード）
- Attachments（読み取り、更新、削除）

**管理機能**
- Groups（CRUD、ユーザー管理）
- Roles（読み取り）
- Trackers（読み取り）
- Issue Statuses（読み取り）
- Enumerations（優先度、活動、カテゴリ）

**その他**
- Custom Fields（読み取り）
- Queries（読み取り）
- Journals（読み取り）
- My Account（読み取り、更新）
- Search

詳細な API ドキュメントは [pkg/redmine](pkg/redmine) ディレクトリを参照してください。

---

## CLI

ターミナルから Redmine を管理するためのコマンドラインツールです。

### インストール

```bash
go install github.com/kqns91/redmine-go/cmd/redmine@latest
```

### 設定

config コマンドで対話的に設定：

```bash
redmine config set url https://your-redmine.com
redmine config set api_key your-api-key
```

現在の設定を確認：

```bash
redmine config list
```

設定は `~/.redmine/config.yaml` に保存されます。必要に応じて直接編集することもできます。

環境変数やコマンドラインフラグでも設定できます：

```bash
# 環境変数
export REDMINE_URL="https://your-redmine.com"
export REDMINE_API_KEY="your-api-key"

# コマンドラインフラグ
redmine --url https://your-redmine.com --api-key your-api-key <command>
```

### API キーの取得方法

1. Redmine インスタンスにログイン
2. 右上の「個人設定」をクリック
3. 右サイドバーの「API アクセスキー」を探す
4. 「表示」をクリックしてキーをコピー

### 基本的なコマンド

```bash
# プロジェクト
redmine projects list
redmine projects show <project-id>

# 課題
redmine issues list --project <project-id>
redmine issues show <issue-id>
redmine issues create --project <project-id> --subject "タイトル" --description "説明"
redmine issues update <issue-id> --status <status-id> --assigned-to <user-id>

# ユーザー
redmine users list
redmine users show <user-id>
redmine users current
```

### 出力フォーマット

CLI は3つの出力フォーマットをサポートしています：

**テーブルフォーマット**（デフォルト）
```bash
redmine projects list --format table
```
列を持つ構造化されたテーブルで、ターミナルでの閲覧に適しています。

**JSON フォーマット**
```bash
redmine projects list --format json
```
機械可読な JSON 出力で、スクリプトや統合に便利です。

**テキストフォーマット**
```bash
redmine projects list --format text
```
最小限の書式設定を行ったプレーンテキスト出力です。

### ヘルプ

すべてのコマンドで詳細なヘルプを表示できます：

```bash
redmine --help
redmine projects --help
redmine issues create --help
```

---

## MCP サーバー

MCP（Model Context Protocol）サーバーは、AI アシスタントが Redmine と連携できるようにします。

### インストール

```bash
go install github.com/kqns91/redmine-go/cmd/mcp@latest
```

### 設定

MCP クライアントの設定ファイルに追加します。

例えば、Claude Desktop の場合：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

基本的な設定（すべてのツールを有効化）：

```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp",
      "env": {
        "REDMINE_URL": "https://your-redmine.com",
        "REDMINE_API_KEY": "your-api-key"
      }
    }
  }
}
```

### 利用可能なツール

サーバーは 23 カテゴリにわたる 80 のツールを提供します：

**コアリソース**
- Projects（7 ツール）
- Issues（7 ツール）
- Users（6 ツール）
- Issue Categories（5 ツール）
- Time Entries（5 ツール）
- Versions（5 ツール）

**高度な操作**
- Batch Operations（1 ツール）- 複数の関連タスクを一括作成
- Progress Monitoring（3 ツール）- プロジェクト健全性分析、見積調整、再スケジュール提案

**プロジェクト管理**
- Memberships（5 ツール）
- Issue Relations（4 ツール）
- Groups（7 ツール）

**コンテンツ・ドキュメント**
- Wiki Pages（4 ツール）
- Attachments（3 ツール）
- News（2 ツール）
- Files（2 ツール）

**メタデータ・設定**
- Enumerations（3 ツール）
- Roles（2 ツール）
- Metadata（2 ツール）
- Custom Fields（1 ツール）
- Queries（1 ツール）

**ユーザーアカウント**
- My Account（2 ツール）
- Search（1 ツール）
- Journals（1 ツール）

### バッチ操作

`create_task_tree` ツールは、依存関係を持つ複数の関連タスクを効率的に作成できます：

```json
{
  "project_id": 1,
  "tasks": [
    {
      "ref": "backend",
      "subject": "バックエンド開発",
      "tracker_id": 1,
      "status_id": 1,
      "priority_id": 1,
      "assigned_to_id": 2,
      "estimated_hours": 24,
      "start_date": "2025-11-10",
      "due_date": "2025-11-13"
    },
    {
      "ref": "frontend",
      "subject": "フロントエンド開発",
      "parent_ref": "backend",
      "assigned_to_id": 3,
      "estimated_hours": 20,
      "start_date": "2025-11-14",
      "due_date": "2025-11-17",
      "blocks_refs": ["backend"]
    }
  ]
}
```

機能：
- `parent_ref` による親子タスク関係の設定
- `blocks_refs` と `precedes_refs` によるタスク依存関係
- `assigned_to_id` による担当者の自動割り当て
- 見積時間、開始/期日、カスタムフィールドのサポート
- 一つの機能に対して 10〜30 個の関連チケットを作成するのに最適

### 進捗監視

3つのツールでプロジェクトの進捗を監視・管理：

**`analyze_project_health`** - 包括的なプロジェクト健全性分析：
- 予定通り、リスクあり、遅延中の課題をリスト化
- クリティカルパスのタスクを特定
- 遅延日数と影響度を計算
- 実行可能な推奨事項を提供

**`adjust_estimates`** - スマートな見積調整：
- 実際の作業時間と見積を分析
- 現在の進捗に基づいて完了日を予測
- 子課題を含めた計算
- 現実的なスケジュールの維持を支援

**`suggest_reschedule`** - 自動再スケジューリング：
- 遅延タスクと依存関係の競合を検出
- 設定可能なバッファ日数で新しい日付を提案
- 変更を自動適用またはプレビューのみ
- クリティカルパスのみモードをサポート

### ツール制御

環境変数を使用して、どのツールを有効にするかを制御できます。

#### 特定のツールグループを有効にする

`REDMINE_ENABLED_TOOLS` を使用して有効にするツールグループを指定：

```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp",
      "env": {
        "REDMINE_URL": "https://your-redmine.com",
        "REDMINE_API_KEY": "your-api-key",
        "REDMINE_ENABLED_TOOLS": "projects,issues,search"
      }
    }
  }
}
```

利用可能なツールグループ：
`projects`、`issues`、`users`、`categories`、`time_entries`、`versions`、`memberships`、`issue_relations`、`wiki`、`attachments`、`enumerations`、`groups`、`news`、`files`、`roles`、`metadata`、`my_account`、`search`、`queries`、`custom_fields`、`journals`、`batch_operations`、`progress_monitoring`、`all`

#### 特定のツールを無効にする

`REDMINE_DISABLED_TOOLS` を使用して個別のツールを無効化：

```json
{
  "env": {
    "REDMINE_DISABLED_TOOLS": "delete_project,delete_issue,delete_user"
  }
}
```

この設定は `REDMINE_ENABLED_TOOLS` よりも優先されます。

#### 設定例

**読み取り専用モード**

すべての書き込み操作を無効化：

```json
{
  "env": {
    "REDMINE_ENABLED_TOOLS": "all",
    "REDMINE_DISABLED_TOOLS": "create_project,update_project,delete_project,archive_project,unarchive_project,create_issue,update_issue,delete_issue,add_watcher,remove_watcher,create_user,update_user,delete_user,create_issue_category,update_issue_category,delete_issue_category,create_time_entry,update_time_entry,delete_time_entry,create_version,update_version,delete_version,create_membership,update_membership,delete_membership,create_issue_relation,delete_issue_relation,create_or_update_wiki_page,delete_wiki_page,update_attachment,delete_attachment,upload_file,create_group,update_group,delete_group,add_group_user,remove_group_user,update_my_account"
  }
}
```

**プロジェクトと課題の管理のみ**

```json
{
  "env": {
    "REDMINE_ENABLED_TOOLS": "projects,issues,search",
    "REDMINE_DISABLED_TOOLS": "delete_project,delete_issue"
  }
}
```

---

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 関連リソース

- [Redmine REST API ドキュメント](https://www.redmine.org/projects/redmine/wiki/Rest_api)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
