# gin-timecard

gin-timecardはgolang,gin,gormで作成されたタイムカード記録用アプリです。<br>
gin-timecardではユーザー登録（ユーザー名、パスワード）を行い、ユーザーごとにタイムカード情報が記録されます。<br>

# 使い方

##インストール〜起動<br>
`$ git clone https://github.com/takashifkd/gin-timecard.git`<br>
`$ cd gin-timecard/`<br>
`$ docker-compose build`<br>
`$ docker-compose up`<br>
command + [T]で別のコマンドを開く<br>
`$ docker-compose exec gin-test sh`<br>
`$ go run main.go`<br>
http://localhost:8080/を開く<br>

##gin-timecardアプリ操作<br>
登録操作：ユーザー名とパスワードで登録します。ユーザー名は重複登録されません。<br>
ログイン操作：登録したユーザー名とパスワードでログインします。<br>
タイムカード一覧画面：月ごとの開始時間、終了時間、休憩時間を記録します。ログイン時には現在月の一覧を表示します。<br>
  ・初期化ボタン：タイムカード一覧が取得できない場合に表示月の空の一覧を作成します（日付のみ記載）。<br>
  ・表示月変更：タイムカード一覧を表示する月を変更します。<br>
  ・編集ボタン：開始時間、終了時間、休憩時間を編集します。（15分間隔）<br>
  ・削除ボタン：開始時間、終了時間、休憩時間をクリアします。（レコードの削除ではない）<br>
  ・ログアウトボタン：ログアウトします。<br>
