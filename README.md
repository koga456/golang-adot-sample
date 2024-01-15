# golang-adot-sample

## 概要
このプロジェクトは、Golang の OpenTelemetry SDK を利用して AWS Distro for OpenTelemetry のコレクター経由で　AWS X-Ray に対してトレース情報を送信するサンプルアプリケーションです。  

## 起動方法

BFF(HTTP サーバー兼gRPCクライアント), gRPC サーバー, ADOT コレクターの3つの Docker コンテナを起動します。

### 起動手順

1. `trust-policy.json` に任意のアカウントIDと作成するロールを引き受ける任意のロール名を指定し、以下のコマンドで

`test-adot-role` の IAM Role の作成

```
aws iam create-role \
    --role-name test-adot-role \
    --assume-role-policy-document file://trust-policy.json
```

`test-adot-role` の IAM Role に `AWSXrayWriteOnlyAccess` のマネージドポリシーをアタッチ

```
aws iam attach-role-policy \
    --role-name test-adot-role \
    --policy-arn arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess
```

2. 以下のコマンドで一次認証情報を取得してください。

```
aws sts assume-role --role-arn arn:aws:iam::{アカウントID}:role/test-adot-role --role-session-name "TestSession"
```

3. `template` ファイルをコピーし `{プロジェクトルート}/.env` ファイルを作成し、`2.` で取得した一次認証情報を設定してください。
4. プロジェクトルートで `docker-compose up -d` を実行
5. `curl -X GET "http://localhost:8080/test/hoge/"` コマンドを実行。`{"output":"hoge"}` のレスポンスが返却されたら起動成功です。
6. `test-adot-role` の IAM Role を作成した AWS アカウントのコンソールにログインし、`ap-northeast-1` リージョンの CloudWatch から X-Ray トレースを開くとトレースが記録されているはずです。