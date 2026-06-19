# tmux-pane-status — Go 再設計

Python 実装（`tmux_pane_status/`）を Go へ移植するにあたっての設計方針。

## 決定事項サマリ

| 項目 | 決定 |
|---|---|
| 互換性 | 破壊的刷新。旧 `PY_TMUX_PANE_*` は読まない |
| 設定の入力源 | 環境変数のみ |
| テンプレート構文 | `{field}` 自前置換 |
| 装飾 | `STYLE_*` で各フィールドに後付け（Python の OPTIONS 分離を踏襲） |
| コマンド表示 | ベースを上書きせず、末尾に控えめなサフィックスを連結 |

## 状態モデル

排他3状態（default / command / git）をやめ、「ベース1択 + 任意サフィックス」に変更する。

```
ベース   = git リポジトリ ? FORMAT_GIT : FORMAT_DEFAULT      （cwd で決定）
コマンド = 非zshの前景プロセスあり ? FORMAT_COMMAND を末尾連結 : ""
出力     = ベース + コマンド
```

- コマンド実行中でも cwd / git 情報は消えない。
- `FORMAT_COMMAND` を空にするとコマンド機能は実質無効化され、pgrep/ps 呼び出しごとスキップされる。

## パッケージ構成

```
cmd/tmux-pane-status/main.go     # CLI: 引数 cwd pid、設定ロード、描画、出力
internal/config/config.go        # defaults + env を読む。FORMAT_* / STYLE_* / ICON_*
internal/collect/collect.go      # Collector: cwd/pid 保持、生ソースを sync.Once でメモ化＋並行
internal/render/render.go        # 状態判定 → template 選択 → 参照抽出 → 描画
internal/render/fields.go        # 各フィールドの抽出ロジック（現 formatters.py 相当）
```

## データフロー

```
1. config ロード（defaults を env で上書き）
2. 状態判定:
     ベース = git リポジトリか（cwd を walk-up して .git を探す）で FORMAT_GIT / FORMAT_DEFAULT
     コマンド = FORMAT_COMMAND が非空 かつ 非zshの前景子プロセスあり なら連結対象
3. 選択した「ベース + コマンドセグメント」を /{(\w+)}/ で舐めて参照フィールド集合を確定
4. 参照フィールド → 必要な生ソース集合（git remote / status / head / ps / pgrep）
5. 生ソースを goroutine で並行取得（sync.Once で重複排除）
6. 各フィールドを抽出 → STYLE で #[...] を巻く → ICON 適用 → template に差し込み
7. stdout へ
```

### 遅延評価と並行化

- 参照フィールドは `{(\w+)}` の正規表現で静的に抽出できるため、必要な生ソースを先に確定 → goroutine で並行プリウォーム → 描画、という流れが素直に書ける。
- 生ソースは Collector のメソッドとして `sync.Once` でメモ化し、重複呼び出しを排除する。

## Python 実装からの主な変更点（捨てるもの）

- `abc.py` のキャッシュ足場（`cache_enabled` / `_load_cache` / `_store_cache`）— 全て未使用のため削除。
- `shell.py` の NameSpace デコレータ — ただの関数に。
- 位置引数 `_extract_data(*ss)` — Collector のメソッド戻り値を型付きで渡す。
- `git_remote_server` と `git_repository_name` が各々 `git remote -v` を呼ぶ二重実行 — 生ソース共有で1回に。
- 排他3状態の分岐ロジック（`main.py:36-53`）— ベース + サフィックスのモデルに置換。

## 生ソースの最適化判断

- `git status -s` / `git remote -v`: subprocess 維持（移植コスト低・確実）。
- `git rev-parse --abbrev-ref HEAD`: 当面 subprocess 維持。`.git/HEAD` 直読みは将来の最適化として保留（detached HEAD 等のエッジで挙動差が出るため、まず同値移植を優先）。

## 環境変数 命名

旧 `PY_TMUX_PANE_*`（二重アンダースコア）を `PANE_STATUS_*`（フラット）に刷新する。

| 旧 | 新 |
|---|---|
| `PY_TMUX_PANE_FORMAT__DEFAULT` | `PANE_STATUS_FORMAT_DEFAULT` |
| `PY_TMUX_PANE_FORMAT__COMMAND` | `PANE_STATUS_FORMAT_COMMAND` |
| `PY_TMUX_PANE_FORMAT__GIT` | `PANE_STATUS_FORMAT_GIT` |
| `PY_TMUX_PANE_OPTIONS__CWD` | `PANE_STATUS_STYLE_CWD` |
| `PY_TMUX_PANE_OPTIONS__*` | `PANE_STATUS_STYLE_*` |
| `PY_TMUX_PANE_ICON__PYTHON` | `PANE_STATUS_ICON_PYTHON` |
| `PY_TMUX_PANE_ICON__*` | `PANE_STATUS_ICON_*` |
| `PY_TMUX_PANE_OVERRIDE_DEFAULTS` | `PANE_STATUS_NO_DEFAULTS` |

`OPTIONS` は色装飾の指定なので意味を明確にするため `STYLE` に改名。

## 既定値

```
PANE_STATUS_FORMAT_DEFAULT = ' {cwd} '
PANE_STATUS_FORMAT_GIT     = ' {git_remote_server} {git_repository_name}{git_cwd} {git_current_branch} {git_status_icons} {project_python} '
PANE_STATUS_FORMAT_COMMAND = ' ⟩ {current_command} {current_command_elapsed}'   # ベース末尾に連結する控えめサフィックス

PANE_STATUS_STYLE_CWD                    = 'fg=blue'
PANE_STATUS_STYLE_GIT_REMOTE_SERVER      = 'bold,fg=blue'
PANE_STATUS_STYLE_GIT_REPOSITORY_NAME    = 'bold,fg=blue'
PANE_STATUS_STYLE_GIT_CURRENT_BRANCH     = 'bold,fg=magenta'
PANE_STATUS_STYLE_CURRENT_COMMAND_ELAPSED = 'bold,bg=green,fg=black'

PANE_STATUS_ICON_PYTHON    = '🐍'
PANE_STATUS_ICON_GITHUB    = '🐱'
PANE_STATUS_ICON_BITBUCKET = '🥛'
PANE_STATUS_ICON_BRANCH    = ''
```

## 実装手順（テストファースト）

1. `test_formatters.py` / `test_directory.py` を Go のテーブルテストへ写経し、振る舞いの同値を固定する。
2. Collector → fields → render の順に、テストを緑にしながら実装する。
3. `hyperfine` で Python 版 vs Go 版を実測し、起動・並行化の効果を確認する。
