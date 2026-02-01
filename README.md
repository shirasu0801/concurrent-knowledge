# Concurrent Knowledge

ジャンル別IT用語に特化した1問1答クイズWebアプリケーション

## 概要

Concurrent Knowledgeは、システムエンジニア向けのIT用語学習アプリです。OS、ネットワーク、データベース、クラウド、セキュリティ、プログラミング、プロジェクト管理の7つのジャンルから、4択クイズ形式で学習できます。

### 主な機能

- ジャンル選択（7ジャンル）
- 難易度選択（初級・中級・上級）
- 4択クイズ（10問、各60秒）
- スコア計算（時間ボーナス付き）
- 結果表示と正解率
- 復習リスト（LocalStorage）
- 問題CRUD管理
- 基本SEO対策

## 技術スタック

- **Backend**: Go 1.25+
- **Database**: SQLite with WAL mode
- **Frontend**: HTML templates + Vanilla JavaScript
- **Router**: Gorilla Mux
- **Storage**: Browser LocalStorage (review list)

## セットアップ

### 前提条件

- Go 1.25以上
- SQLite3 (DB初期化用)

### インストール

1. リポジトリをクローン

```bash
git clone <repository-url>
cd concurrent-knowledge
```

2. 依存関係をインストール

```bash
go mod download
```

3. 環境変数を設定（オプション）

```bash
cp .env.example .env
# .envファイルを編集
```

4. データベースを初期化（2つの方法）

**方法1: 専用ツールを使用（推奨）**
```bash
go run cmd/init-db/main.go
```

**方法2: サーバー起動時に自動初期化**
```bash
# データベースが存在しない場合、サーバー起動時に自動的に初期化されます
go run cmd/server/main.go
```

5. ブラウザでアクセス

```
http://localhost:8080(ポート占有の場合は3000とか)
```

## Makefileコマンド

```bash
make help       # ヘルプを表示
make build      # アプリケーションをビルド
make run        # アプリケーションを実行
make test       # テストを実行
make clean      # ビルド成果物を削除
make init-db    # データベースを初期化
make deps       # 依存関係をインストール
make dev        # DB初期化 + 開発モードで実行
```

## プロジェクト構造

```
concurrent-knowledge/
├── cmd/server/          # アプリケーションエントリーポイント
├── internal/            # 内部パッケージ
│   ├── config/          # 設定管理
│   ├── database/        # DB接続・初期化
│   ├── models/          # データモデル
│   ├── repository/      # データアクセス層
│   ├── service/         # ビジネスロジック
│   ├── handler/         # HTTPハンドラ
│   ├── middleware/      # ミドルウェア
│   └── utils/           # ユーティリティ
├── web/                 # フロントエンド
│   ├── static/          # 静的ファイル
│   │   ├── css/         # スタイルシート
│   │   └── js/          # JavaScript
│   └── templates/       # HTMLテンプレート
├── scripts/             # スクリプト
├── migrations/          # DBマイグレーション
└── data/                # SQLiteファイル（gitignore）
```

## スコア計算

```
Score = Σ(BasePoint × DifficultyMultiplier + (RemainingTime/TotalTime) × TimeBonus)
```

- **BasePoint**: 100点（固定）
- **DifficultyMultiplier**:
  - 初級: 1.0
  - 中級: 1.5
  - 上級: 2.0
- **TimeBonus**: 最大100点（残り時間に比例）

## API エンドポイント

### Quiz API

- `POST /api/quiz/start` - クイズ開始
- `POST /api/quiz/submit` - 回答送信

### Question CRUD API

- `GET /api/questions` - 問題一覧
- `GET /api/questions/{id}` - 問題詳細
- `POST /api/questions` - 問題作成
- `PUT /api/questions/{id}` - 問題更新
- `DELETE /api/questions/{id}` - 問題削除

### Other API

- `GET /api/genres` - ジャンル一覧
- `GET /api/health` - ヘルスチェック

## パフォーマンス

- クイズ開始APIレスポンス: **<200ms**
- 回答判定APIレスポンス: **<200ms**
- ページ初期表示: **<1秒**

### 最適化戦略

- WALモードによる並行読み取り性能向上
- sync.Mapによる問題プールキャッシング
- インメモリクイズセッション管理
- コネクションプール（MaxOpen=25, MaxIdle=5）

## セキュリティ

- XSS防止: html/templateの自動エスケープ
- SQLインジェクション防止: パラメータ化クエリ
- 入力検証: 難易度・ジャンルIDのバリデーション

## デプロイ

### EC2デプロイ

```bash
# ビルド
make build

# EC2にアップロード
scp bin/server.exe user@ec2-instance:/path/to/app/

# 実行
./server.exe
```

### 環境変数

```
PORT=8080
DATABASE_PATH=./data/quiz.db
ENV=production
LOG_LEVEL=info
```

## ライセンス

MIT License

## 作者

Concurrent Knowledge Team

## バージョン

v1.0.0
