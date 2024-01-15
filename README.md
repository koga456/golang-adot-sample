# golang-adot-sample

## 起動方法

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

3. `template` ファイルをコピーし `{project_root}/.env` ファイルを作成し、`2.` で取得した一次認証情報を設定してください。