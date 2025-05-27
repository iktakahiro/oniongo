# Go DDD & オニオンアーキテクチャ 実装例とテクニック

[![Go Version](https://img.shields.io/badge/go-1.24.3-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/iktakahiro/oniongo)

[English](README.md) | 日本語

**注意**: このリポジトリは「GoアプリケーションでDDDアーキテクチャを実装する方法」を示すサンプルです。参考にする場合は、本番環境にデプロイする前に認証とセキュリティを実装してください！

* 関連するPython実装: [dddpy](https://github.com/iktakahiro/dddpy)

## 技術スタック

* **gRPC & Connect-Go**: HTTP/2対応のモダンなRPCフレームワーク
* **Ent**: コード生成機能付きの型安全なGo用ORM
* **SQLite**: 開発用の軽量データベース
* **Buf**: Protocol Bufferの管理とコード生成
* **Samber/do**: 依存性注入コンテナ
* **Atlas**: データベースマイグレーションツール

## プロジェクトセットアップ

1. 依存関係をインストール:

```bash
make install
```

2. コード生成（protobuf、ent、mock）:

```bash
make buf-generate
make ent-generate
make mockgen
```

3. データベースマイグレーション実行:

```bash
make migrate-up
```

4. gRPCサーバー起動:

```bash
make server
```

サーバーはデフォルトでポート8080で起動します。`PORT`環境変数を設定することで変更できます。

## コードアーキテクチャ

ディレクトリ構造はオニオンアーキテクチャに基づいています：

```
internal/
├── domain/           # ドメイン層（エンティティ、値オブジェクト、リポジトリインターフェース）
│   └── todo/
├── application/      # アプリケーション層（ユースケース）
│   ├── todoapp/
│   └── uow/         # Unit of Workパターン
├── infrastructure/  # インフラストラクチャ層（リポジトリ実装、外部サービス）
│   ├── ent/         # Ent ORM（スキーマ、生成コード、リポジトリ）
│   ├── grpc/        # gRPC生成コード
│   ├── sqlite/      # データベースマイグレーション
│   └── di/          # 依存性注入設定
└── api/             # プレゼンテーション層（gRPCハンドラー）
    └── grpc/
```

### ドメイン層

ドメイン層はコアビジネスロジックを含み、外部の関心事から独立しています。以下を含みます：

1. **エンティティ**: アイデンティティを持つコアビジネスオブジェクト
2. **値オブジェクト**: 特性を記述する不変オブジェクト
3. **リポジトリインターフェース**: データ永続化の契約
4. **ドメインサービス**: 単一のエンティティに属さないビジネスロジック

#### 1. エンティティ

`Todo`エンティティはコアビジネスオブジェクトを表現します：

```go
// Todo is the entity that represents a todo item.
type Todo struct {
    id          TodoID
    title       string
    body        string
    status      TodoStatus
    createdAt   time.Time
    updatedAt   time.Time
    completedAt *time.Time
}

// NewTodo creates a new Todo.
func NewTodo(title string, body string) (*Todo, error) {
    now := time.Now()
    return &Todo{
        id:          NewTodoID(),
        title:       title,
        body:        body,
        status:      TodoStatusNotStarted,
        createdAt:   now,
        updatedAt:   now,
        completedAt: nil,
    }, nil
}
```

エンティティの主な特徴：

* ビジネスルールと不変条件をカプセル化
* 状態遷移のメソッドを提供（`Start()`、`Complete()`）
* バリデーションによりデータ整合性を維持
* 型安全性のために値オブジェクトを使用（`TodoID`、`TodoStatus`）

#### 2. 値オブジェクト

値オブジェクトは型安全性を確保し、バリデーションロジックをカプセル化します：

```go
// TodoID represents a unique identifier for a Todo.
type TodoID uuid.UUID

// NewTodoID creates a new TodoID.
func NewTodoID() TodoID {
    return TodoID(uuid.New())
}

// TodoStatus represents the status of a Todo.
type TodoStatus string

const (
    TodoStatusNotStarted TodoStatus = "not_started"
    TodoStatusInProgress TodoStatus = "in_progress"
    TodoStatusCompleted  TodoStatus = "completed"
)
```

#### 3. リポジトリインターフェース

リポジトリインターフェースは実装詳細を指定せずにデータ永続化の契約を定義します：

```go
// TodoRepository is the interface that wraps the basic CRUD operations for Todo.
type TodoRepository interface {
    Create(ctx context.Context, todo *Todo) error
    Update(ctx context.Context, todo *Todo) error
    FindAll(ctx context.Context) ([]*Todo, error)
    FindByID(ctx context.Context, id TodoID) (*Todo, error)
    Delete(ctx context.Context, id TodoID) error
}
```

### インフラストラクチャ層

インフラストラクチャ層はドメイン層で定義されたインターフェースの実装を含みます。以下を含みます：

1. **リポジトリ実装**: Ent ORMを使用した具体的な実装
2. **データベーススキーマ**: Entスキーマ定義
3. **外部サービス統合**: gRPC、HTTPクライアントなど

#### 1. リポジトリ実装

`todoRepository`はEnt ORMを使用してドメインリポジトリインターフェースを実装します：

```go
// todoRepository is the implementation of the TodoRepository interface.
type todoRepository struct{}

// Create creates the Todo.
func (r todoRepository) Create(ctx context.Context, todo *todo.Todo) error {
    tx, err := db.GetTx(ctx)
    if err != nil {
        return err
    }

    status := todoschema.Status(todo.Status().String())
    _, err = tx.TodoSchema.Create().
        SetTitle(todo.Title()).
        SetBody(todo.Body()).
        SetStatus(status).
        Save(ctx)
    if err != nil {
        return fmt.Errorf("failed to create todo: %w", err)
    }
    return nil
}
```

リポジトリインターフェースとは異なり、インフラストラクチャ層の実装コードは特定の技術（この例ではEnt ORMとSQLite）に固有の詳細を含むことができます。

#### 2. データマッピング

インフラストラクチャ層はドメインエンティティとデータベースモデル間の変換を処理します：

```go
// convertEntToTodo converts ent.TodoSchema to domain Todo
func convertEntToTodo(v *entgen.TodoSchema) (*todo.Todo, error) {
    status, err := todo.NewTodoStatusFromString(string(v.Status))
    if err != nil {
        return nil, fmt.Errorf("failed to convert status %v: %w", v.Status, err)
    }
    return todo.ReconstructTodo(v.ID, v.Title, *v.Body, status, v.CreatedAt, v.UpdatedAt), nil
}
```

### アプリケーション層

アプリケーション層はアプリケーション固有のビジネスルールを含みます。以下を含みます：

1. **ユースケース実装**: ドメインオブジェクトを調整するアプリケーションサービス
2. **トランザクション管理**: データ整合性のためのUnit of Workパターン
3. **エラーハンドリング**: アプリケーション固有のエラー処理

#### 1. ユースケース実装

各ユースケースは単一の`Execute`メソッドを持つ別々の構造体として実装されます：

```go
// CreateTodoUseCase is the interface that wraps the basic CreateTodo operation.
type CreateTodoUseCase interface {
    Execute(ctx context.Context, req CreateTodoRequest) error
}

// createTodoUseCase is the implementation of the CreateTodoUseCase interface.
type createTodoUseCase struct {
    todoRepository todo.TodoRepository
    txRunner       uow.TransactionRunner
}

// Execute creates a new Todo.
func (u createTodoUseCase) Execute(ctx context.Context, req CreateTodoRequest) error {
    todo, err := todo.NewTodo(req.Title, req.Body)
    if err != nil {
        return fmt.Errorf("failed to create todo: %w", err)
    }
    
    err = u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
        if err := u.todoRepository.Create(ctx, todo); err != nil {
            return fmt.Errorf("failed to save todo: %w", err)
        }
        return nil
    })
    if err != nil {
        return fmt.Errorf("failed to execute transaction: %w", err)
    }
    return nil
}
```

ユースケースの主な特徴：

* 単一責任の原則
* Unit of Workパターンによるトランザクション管理
* 明確なインターフェース定義
* コンストラクタによる依存性注入

#### 2. トランザクション管理

アプリケーション層はデータ整合性を確保するためにUnit of Workパターンを使用します：

```go
// TransactionRunner provides transaction management capabilities.
type TransactionRunner interface {
    RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### プレゼンテーション層

プレゼンテーション層はgRPCリクエストとレスポンスを処理します。以下を含みます：

1. **gRPCハンドラー**: protobufメッセージとドメインオブジェクト間の変換
2. **Protocol Buffer定義**: API契約定義
3. **入力バリデーション**: buf validateを使用したリクエストバリデーション
4. **エラーハンドリング**: ドメインエラーをgRPCステータスコードに変換

#### 1. gRPCハンドラー

ハンドラーは`api/grpc`ディレクトリ下に整理されています：

```go
func (h *todoHandler) CreateTodo(
    ctx context.Context,
    req *connect.Request[v1.CreateTodoRequest],
) (*connect.Response[v1.CreateTodoResponse], error) {
    // Extract body value if present
    body := ""
    if req.Msg.Body != nil {
        body = *req.Msg.Body
    }

    // Create use case request
    useCaseReq := todoapp.CreateTodoRequest{
        Title: req.Msg.Title,
        Body:  body,
    }

    // Execute use case
    if err := h.createTodoUseCase.Execute(ctx, useCaseReq); err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    // Return response
    return connect.NewResponse(&v1.CreateTodoResponse{}), nil
}
```

#### 2. Protocol Buffer定義

API契約はバリデーション付きのProtocol Buffersを使用して定義されます：

```protobuf
// TodoService provides all todo-related operations
service TodoService {
  // CreateTodo creates a new todo item
  rpc CreateTodo(CreateTodoRequest) returns (CreateTodoResponse);
  
  // GetTodo retrieves a todo item by its ID
  rpc GetTodo(GetTodoRequest) returns (GetTodoResponse);
  
  // GetTodos retrieves all todo items
  rpc GetTodos(GetTodosRequest) returns (GetTodosResponse);
  
  // UpdateTodo updates an existing todo item
  rpc UpdateTodo(UpdateTodoRequest) returns (UpdateTodoResponse);
  
  // StartTodo changes the todo status to in progress
  rpc StartTodo(StartTodoRequest) returns (StartTodoResponse);
  
  // CompleteTodo changes the todo status to completed
  rpc CompleteTodo(CompleteTodoRequest) returns (CompleteTodoResponse);
  
  // DeleteTodo deletes a todo item
  rpc DeleteTodo(DeleteTodoRequest) returns (DeleteTodoResponse);
}

message CreateTodoRequest {
  string title = 1 [(buf.validate.field).string.min_len = 1];
  optional string body = 2;
}
```

## 使用方法

1. このリポジトリをクローン
2. 依存関係をインストール: `make install`
3. コード生成: `make buf-generate && make ent-generate`
4. マイグレーション実行: `make migrate-up`
5. サーバー起動: `make server`
6. gRPCサーバーが`localhost:8080`で利用可能になります

### grpcurlを使用したサンプルリクエスト

* 新しいTodoを作成:

```bash
grpcurl -plaintext -d '{
  "title": "DDD アーキテクチャの実装",
  "body": "DDD原則を使用したサンプルアプリケーションの作成"
}' localhost:8080 oniongo.v1.TodoService/CreateTodo
```

* 全てのTodoを取得:

```bash
grpcurl -plaintext -d '{}' localhost:8080 oniongo.v1.TodoService/GetTodos
```

* 特定のTodoを取得:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/GetTodo
```

* Todoを開始:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/StartTodo
```

* Todoを完了:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/CompleteTodo
```

## 開発

### コード生成

このプロジェクトは複数のコード生成ツールを使用します：

```bash
# protobufコード生成
make buf-generate

# Ent ORMコード生成
make ent-generate

# テスト用モック生成
make mockgen
```

### データベースマイグレーション

```bash
# 新しいマイグレーションを作成
make migrate-diff name=add_new_field

# マイグレーションを適用
make migrate-up
```

### テスト実行

```bash
go test ./...
```

### コード品質

このプロジェクトはコード品質を維持するために複数のツールを使用します：

* **golangci-lint**: Go用の包括的なリンター
* **buf**: Protocol bufferのリンティングと破壊的変更検出
* **mockery**: テスト用のモック生成

```bash
# コードフォーマット
make fmt

# リンター実行
make lint
```

## 主要な設計パターン

### 1. 依存性注入

プロジェクトは`samber/do`ライブラリを使用して依存性注入を行い、層間の疎結合を確保します：

```go
// Dependency injection setup
func DependencyInjection() *do.Injector {
    injector := do.New()
    
    // Register dependencies
    do.Provide(injector, todorepo.NewTodoRepository)
    do.Provide(injector, todoapp.NewCreateTodoUseCase)
    do.Provide(injector, todohandler.NewTodoServiceHandler)
    
    return injector
}
```

### 2. Unit of Workパターン

トランザクション管理はUnit of Workパターンで処理され、複数の操作にわたってデータ整合性を確保します：

```go
err = u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
    // 単一トランザクション内での複数のリポジトリ操作
    return nil
})
```

### 3. リポジトリパターン

リポジトリパターンはデータアクセスロジックを抽象化し、ドメイン層を特定のデータベース技術から独立させます：

```go
// ドメイン層がインターフェースを定義
type TodoRepository interface {
    Create(ctx context.Context, todo *Todo) error
    // ... その他のメソッド
}

// インフラストラクチャ層が実装を提供
type todoRepository struct{}
```

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています - 詳細は[LICENSE](LICENSE)ファイルを参照してください。

## このプロジェクトについて

このリポジトリは、Goでドメイン駆動設計（DDD）とオニオンアーキテクチャを実装する方法を示し、スケーラブルなアプリケーションを構築するためのクリーンで保守可能、テスト可能なコードベース構造を提供します。
