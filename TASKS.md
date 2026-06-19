# 実装タスク — Go 移植

[DESIGN.md](./DESIGN.md) に基づく実装タスク一覧。テストファースト（先に失敗するテストを書いてから実装）で進める。各フェーズはおおむね上から順に依存する。

## Phase 0: プロジェクト初期化

- [ ] `go mod init github.com/yhiraki/tmux-pane-status`
- [ ] パッケージ用ディレクトリを作成
  - `cmd/tmux-pane-status/`
  - `internal/config/`
  - `internal/collect/`
  - `internal/render/`
- [ ] `.gitignore` に Go 用エントリ（ビルド成果物バイナリ等）を追加
- [ ] CI（`.github/workflows/`）を Python マトリクスから Go（`go test ./...` / `go vet` / `go build`）へ更新

## Phase 1: ディレクトリ判定（test_directory.py 写経 → 実装）

`internal/collect`（または `internal/dir`）に配置。

- [ ] テスト: `git_root` のテーブルテスト（`test_directory.py:28-90` の7ケースを移植）
  - `.git/` が同階層 → root はそこ
  - 子から walk-up して親の `.git/` を見つける
  - 多段ネスト、最も近い `.git/` を返す
  - `.git` が**ファイル**（ディレクトリでない）→ `nil`
  - `.git` なし → `nil`
  - 空ディレクトリ → `nil`
- [ ] 実装: `walkUp(path)` — resolve してディレクトリ単位で親方向へ
- [ ] 実装: `gitRoot(path)` — `.git` が**ディレクトリ**の階層を返す
- [ ] 実装: `isGit(path)` — `gitRoot != nil`
- [ ] 実装: `isPython(path)` — walk-up で `__init__.py` / `setup.py` / `setup.cfg` / `pytest.ini` / `manage.py` のいずれか検出（`directory.py:26-36`）

## Phase 2: git/ps 出力パース（test_formatters.py 写経 → 実装）

`internal/render/fields.go`（パース関数）。

- [ ] テスト: `gitParseRemote` SSH 形式（`test_formatters.py:15-28`）
  - 入力 `origin\tgit@github.com:yhiraki/tmux-pane-status.git (fetch)` → `(origin, github.com, yhiraki, tmux-pane-status.git, (fetch))`
- [ ] テスト: `gitParseRemote` HTTPS 形式（`test_formatters.py:30-43`）
  - `origin\thttps://github.com/yhiraki/tmux-pane-status.git (fetch)` → 同上タプル
- [ ] テスト: `gitParseStatus`（`test_formatters.py:45-61`）
  - ` A format.py` / ` M git.py` / `?? .envrc` を `[status, path]` へ分割
- [ ] 実装: `gitParseRemote` — `@` 区切り（SSH）と `https://` の2系統を正規表現で分解
- [ ] 実装: `gitGetRemote` — `origin` 優先、なければ最初のものを返す（`formatters.py:26-33`）
- [ ] 実装: `gitParseStatus` — 行ごとに空白分割

## Phase 3: config（既定値 + env）

`internal/config/config.go`。

- [ ] テスト: 既定値がそのまま読めること
- [ ] テスト: env が既定値を上書きすること
- [ ] テスト: `PANE_STATUS_NO_DEFAULTS` で全既定値が空になること（`main.py:18-22` 相当）
- [ ] 実装: 既定値の定義（[DESIGN.md の既定値](./DESIGN.md#既定値)）
- [ ] 実装: `FORMAT_* / STYLE_* / ICON_*` を env から読み、`map[string]string` に正規化（小文字キー）
- [ ] 実装: `NO_DEFAULTS` 処理

## Phase 4: 生ソース取得（collect）

`internal/collect/collect.go`。

- [ ] 実装: subprocess runner（stdout を文字列で返す。stderr 握りつぶし、`shell.py:8-13` 相当）
- [ ] 実装: Collector 構造体（`cwd`, `pid` 保持）
- [ ] 実装: 生ソースメソッドを `sync.Once` でメモ化
  - `gitRemote()` = `git remote -v`
  - `gitStatus()` = `git status -s`
  - `gitBranch()` = `git rev-parse --abbrev-ref HEAD`
  - `psInfo(pid, headers)` = `ps -p PID -o ...`
  - `childPID()` = `pgrep -P PID` → 自プロセス除外しソート（`main.py:39-47`）
- [ ] テスト: 同じソースを2回呼んでも subprocess は1回（メモ化の検証。runner をモック可能に）

## Phase 5: フィールド抽出（fields）

`internal/render/fields.go`。各フィールドは生ソース文字列 → 表示文字列。

- [ ] テスト + 実装: `git_remote_server` — server 名、`github.com`/`bitbucket.org` をアイコン置換（`formatters.py:44-55`）
- [ ] テスト + 実装: `git_repository_name` — `user/name`、`.git` サフィックス除去（`formatters.py:58-65`）
- [ ] テスト + 実装: `git_current_branch` — trim、`branch` アイコン前置（`formatters.py:68-75`）
- [ ] テスト + 実装: `git_status_icons` — status 先頭文字を集合化しソート、`[MA?]` 形式（`formatters.py:78-84`）
- [ ] テスト + 実装: `git_cwd` — git root からの相対パス（`formatters.py:87-93`）
- [ ] テスト + 実装: `cwd` — `$HOME` を `~`、`/private` 除去（`formatters.py:96-104`）
- [ ] テスト + 実装: `project_python` — `isPython` なら python アイコン（`formatters.py:107-116`）
- [ ] テスト + 実装: `current_command` — ps command をベース名化、`ssh` の特別整形（`-p` スキップ、`SSH ->` 置換、`formatters.py:119-146`）
- [ ] テスト + 実装: `current_command_elapsed` — ps etime の trim（`formatters.py:149-154`）
- [ ] 実装: STYLE 適用 — `#[opt]...#[default]` で巻く（`abc.py:8-15`）
- [ ] 実装: ICON 適用 — 各フィールドの `_set_icons` 相当

## Phase 6: 描画（render）

`internal/render/render.go`。

- [ ] テスト: 状態判定
  - git リポジトリ → ベース = `FORMAT_GIT`
  - 非git → ベース = `FORMAT_DEFAULT`
  - 非zshの前景子プロセスあり かつ `FORMAT_COMMAND` 非空 → コマンドサフィックス連結
  - `FORMAT_COMMAND` 空 → pgrep/ps を呼ばずサフィックスなし
- [ ] テスト: 参照フィールド抽出 — `{(\w+)}` 正規表現で「ベース + コマンドセグメント」から集合を得る
- [ ] テスト: template 置換 — `{cwd}` 等が値に差し替わる。未参照フィールドは評価されない
- [ ] 実装: 状態判定（[DESIGN.md の状態モデル](./DESIGN.md#状態モデル)）
- [ ] 実装: 参照フィールド → 必要な生ソース集合の対応表
- [ ] 実装: 生ソースを goroutine で並行プリウォーム（`sync.WaitGroup`）→ フィールド描画 → 置換
- [ ] 実装: 出力文字列を組み立て

## Phase 7: CLI

`cmd/tmux-pane-status/main.go`。

- [ ] 実装: 引数パース `cwd`（Path）`pid`（int）（`main.py:65-73`）
- [ ] 実装: `cwd` へ chdir、resolve
- [ ] 実装: config ロード → render → `stdout` へ書き出し
- [ ] テスト: 結合テスト（一時 git リポジトリを作り、想定文字列が出ることを確認）

## Phase 8: 検証・配布・撤去

- [ ] `hyperfine 'tmux-pane-status <cwd> <pid>'` で Python 版 vs Go 版を実測（起動・並行化の効果確認）
- [ ] README 更新 — インストール手順（`go install` / バイナリ配布）、新しい環境変数名での設定例
- [ ] リリース手順の整備（クロスコンパイル / goreleaser 等、必要なら）
- [ ] 旧 Python 実装の撤去（`tmux_pane_status/`, `setup.py`）— Go 版が同値であることを確認後
- [ ] バージョン更新（破壊的変更なので v1.0.0 を検討）

## 補足: 移植時の同値チェックの要点

- `git_root` の `.git` は**ディレクトリのみ**有効（ファイルの `.git` は無視）。テストで担保済み。
- remote パースは origin 優先・フォールバック first。
- `git_status_icons` は先頭文字の**集合**（重複排除 + ソート）。
- 子プロセス判定で**自プロセス（pgrep が拾う自身）を除外**する点を忘れない。
