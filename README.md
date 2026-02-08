# portfolio-go-chat

[![CircleCI](https://circleci.com/gh/AtsuyaOotsuka/portfolio-go-chat.svg?style=svg)](https://circleci.com/gh/AtsuyaOotsuka/portfolio-go-chat)

## 概要

本システムはGoのポートフォリオ用に作成した、チャットAPIサービスになります。

CSRFおよび認証は[こちら](https://github.com/AtsuyaOotsuka/portfolio-go-auth)の認証サービス側のリポジトリに対応しております。

本システムはAPIのみのため、PostmanなどのAPIツールからの動作確認となります。

フロントエンドを実装すると評価軸が分散するため、
本ポートフォリオではバックエンドの設計・テスト・運用に
評価を集中させています。

## 起動方法

- .env.sampleをベースに.envを作成
- secret配下に firebase_credentials.json を設置
- Dockerを起動
- sh run.shを実行

## テスト実行方法

- Dockerを起動
- sh test.sh を実行

これで、自動的に環境構築が行われ、テストが実行されます。
テストカバレッジのファイルは mountディレクトリ配下に設置されます。

## 設計方針

本システムは Clean Architecture の考え方を参考に、
責務の分離と依存方向の一方通行を意識して設計しています。

```
main -> app -> provider -> routing -> middleware -> handler -> service -> usecase -> lib
```

という形で一方通行の構成になっております

## 技術選定

本システムでは echo + MongoDB を使用しています。

現時点では各フレームワークや DB の性能・優位性を十分に比較できていないため、
認証サービスとは異なる構成を採用し、知識構築を目的として選定しました。

## テスト方針

unit test では正常系・異常系を含めた振る舞いの検証を行い、
E2E テストでは実際の MongoDB を使用して API 全体の動作確認を行っています。
CI（CircleCI）上でも同構成でテストが実行されます。

テストカバレッジの状況はinternal 配下に関しては、すべての階層において 100%となっております

```
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/entry_point/api      0.295s  coverage: 46.9% of statements
        github.com/AtsuyaOotsuka/portfolio-go-chat/entry_point/cmd              coverage: 0.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app 0.011s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/cmd 0.003s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/command     1.008s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts      0.002s  coverage: [no statements]
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto 0.003s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler     0.006s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware  0.013s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model       0.002s  coverage: [no statements]
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/provider    0.011s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/routing     0.012s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/routing/room_routing        0.018s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service     0.012s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/cmd_svc     0.008s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc   0.004s  coverage: 100.0% of statements
ok      github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase     0.003s  coverage: 100.0% of statements
```

## おわりに

本システムの構成は、かなり冗長なっております。
ポートフォリオという特性上、できる限り多くの領域を表現する必要があるため、
様々な案件で対応可能とするため、Clean Architecture 風の実装、カバレッジ100％を意識したテストで実装を行いました。

実際に参画させていただいた際には、
コーディング規約、チームの慣習とうに則った対応をさせていただく所存でございます。
