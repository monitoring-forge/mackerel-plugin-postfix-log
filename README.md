# mackerel-plugin-postfix-log

Postfixの`maillog`を解析し、メール配信の遅延や転送結果をMackerelに投稿するためのプラグインです。

## 概要

このプラグインはPostfixのSMTP送信ログ（`postfix/smtp`）を解析し、以下の情報をMackerelカスタムメトリックとして出力します。

- 配信遅延時間（total / recving / queuing / connection / transmission）
- 転送結果別の件数（2xx / 4xx / 5xx）
- 転送結果別の比率（2xx / 4xx / 5xx）

## インストール

### mkrコマンドを利用する場合

Mackerelプラグインをインストールします。

```bash
mkr plugin install monitoring-forge/mackerel-plugin-postfix-log
```

インストール先の例: `/opt/mackerel-agent/plugins/bin/mackerel-plugin-postfix-log`

### GitHub Releasesからダウンロードする場合

[GitHub Releases](https://github.com/monitoring-forge/mackerel-plugin-postfix-log/releases/) から、OS・アーキテクチャに応じた最新のバイナリをダウンロードしてください。

```bash
# 例: Linux x86_64 の場合
curl -L -o mackerel-plugin-postfix-log.tar.gz \
  "https://github.com/monitoring-forge/mackerel-plugin-postfix-log/releases/latest/download/mackerel-plugin-postfix-log_linux_amd64.tar.gz"
tar xzf mackerel-plugin-postfix-log.tar.gz
```

## mackerel-agentでの設定

`mackerel-agent.conf`にプラグイン設定を追加してください。

```ini
[plugin.metrics.postfix_log]
command = ["/opt/mackerel-agent/plugins/bin/mackerel-plugin-postfix-log", "--logfile", "/var/log/maillog"]
```

ログファイルのパスが `/var/log/maillog` 以外の場合は、`--logfile` オプションで指定してください。

設定後、mackerel-agentを再起動します。

```bash
sudo systemctl restart mackerel-agent
```

## 使い方

```bash
mackerel-plugin-postfix-log [OPTIONS]
```

### オプション

| オプション | 説明 | デフォルト値 |
| --- | --- | --- |
| `--logfile` | 解析対象のPostfixログファイルパス | `/var/log/maillog` |
| `--posfile-prefix` | 読み込み位置を保存するファイル名の接頭辞 | `maillog` |
| `-v`, `--version` | バージョン情報を表示 | - |
| `-V`, `--verbose` | 詳細なログを表示 | - |
| `-h`, `--help` | ヘルプを表示 | - |

### 実行例

```bash
mackerel-plugin-postfix-log --logfile /var/log/maillog
```

出力例:

```
postfixlog.total_delay.average  0.240476        1555681849
postfixlog.total_delay.99_percentile    0.250000        1555681849
postfixlog.total_delay.95_percentile    0.240000        1555681849
postfixlog.total_delay.90_percentile    0.240000        1555681849
postfixlog.recving_delay.average        0.040000        1555681849
postfixlog.recving_delay.99_percentile  0.040000        1555681849
postfixlog.recving_delay.95_percentile  0.040000        1555681849
postfixlog.recving_delay.90_percentile  0.040000        1555681849
postfixlog.queuing_delay.average        0.000476        1555681849
postfixlog.queuing_delay.99_percentile  0.010000        1555681849
postfixlog.queuing_delay.95_percentile  0.000000        1555681849
postfixlog.queuing_delay.90_percentile  0.000000        1555681849
postfixlog.connection_delay.average     0.090000        1555681849
postfixlog.connection_delay.99_percentile       0.090000        1555681849
postfixlog.connection_delay.95_percentile       0.090000        1555681849
postfixlog.connection_delay.90_percentile       0.090000        1555681849
postfixlog.transmission_delay.average   0.090000        1555681849
postfixlog.transmission_delay.99_percentile     0.090000        1555681849
postfixlog.transmission_delay.95_percentile     0.090000        1555681849
postfixlog.transmission_delay.90_percentile     0.090000        1555681849
postfixlog.transfer_num.2xx_count       1.615385        1555681849
postfixlog.transfer_num.4xx_count       0.000000        1555681849
postfixlog.transfer_num.5xx_count       0.000000        1555681849
postfixlog.transfer_total.count 1.615385        1555681849
postfixlog.transfer_ratio.2xx_percentage        100.000000      1555681849
postfixlog.transfer_ratio.4xx_percentage        0.000000        1555681849
postfixlog.transfer_ratio.5xx_percentage        0.000000        1555681849
```

## 出力データの意味

### 遅延時間メトリック

Postfixログ内の `delay=...` と `delays=...` の値を元にしています。

| メトリック名 | 内容 |
| --- | --- |
| `postfixlog.total_delay.*` | 1通あたりの総配信遅延時間（秒） |
| `postfixlog.recving_delay.*` | 受信処理にかかった遅延時間（秒） |
| `postfixlog.queuing_delay.*` | キュー内で待機していた時間（秒） |
| `postfixlog.connection_delay.*` | 宛先SMTPサーバーへの接確確立にかかった時間（秒） |
| `postfixlog.transmission_delay.*` | メッセージ送信にかかった時間（秒） |

`*` には `average`（平均）、`90_percentile`、`95_percentile`、`99_percentile` が入ります。

### 転送件数メトリック

単位時間あたりの転送件数を出力します。

| メトリック名 | 内容 |
| --- | --- |
| `postfixlog.transfer_num.2xx_count` | 正常終了（DSN 2xx）の転送件数/秒 |
| `postfixlog.transfer_num.4xx_count` | 一時的なエラー（DSN 4xx）の転送件数/秒 |
| `postfixlog.transfer_num.5xx_count` | 恒久的なエラー（DSN 5xx）の転送件数/秒 |
| `postfixlog.transfer_total.count` | 総転送件数/秒 |

### 転送比率メトリック

全転送件数に対する各種ステータスの割合（%）を出力します。

| メトリック名 | 内容 |
| --- | --- |
| `postfixlog.transfer_ratio.2xx_percentage` | 正常終了（DSN 2xx）の割合 |
| `postfixlog.transfer_ratio.4xx_percentage` | 一時的なエラー（DSN 4xx）の割合 |
| `postfixlog.transfer_ratio.5xx_percentage` | 恒久的なエラー（DSN 5xx）の割合 |

## ログの例

以下のようなPostfixのSMTP送信ログを解析対象としています。

```
Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)
```

このログから `delay`、`delays`、`dsn` の値を抽出して集計します。

## 注意事項

- このプラグインは `postfix/smtp` のログ行のみを解析します。他のプロセス（`postfix/smtpd` など）のログは無視されます。
- ログファイルの読み込み位置は `mackerel-agent` のプラグインワークディレクトリ内に保存され、次回実行時に続きから解析します。
- ログファイルへのアクセス権限が必要です。mackerel-agentの実行ユーザーに読み取り権限があるか確認してください。