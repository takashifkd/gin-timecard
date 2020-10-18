package router

import (
	"strconv"
	// "gin-test/controller"
	"gin-test/crypto"
	"gin-test/db"
	"gin-test/models"
	"log"
	"net/http"
	// "os"
	"os/user"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionInfo struct {
	UserId         interface{}
	UserName       interface{}
	IsSessionAlive bool
}

//formで送信されたデータをbind
var month models.Month
var sessionInfo SessionInfo

// Init
func Init() {
	r := gin.Default()
	r.LoadHTMLGlob("./view/*.html")

	// セッションの作成
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	//トップ画面へ
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "top.html", gin.H{})
	})

	//ユーザー登録画面へ
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{})
	})

	// ユーザー登録
	r.POST("/signup", func(c *gin.Context) {
		var form user.User
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
			c.Abort()
		} else {
			username := c.PostForm("username")
			password := c.PostForm("password")
			// 登録ユーザーが重複していた場合にはじく処理
			if err := createUser(username, password); err != nil {
				c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
			}
			c.Redirect(302, "/")
		}
	})

	// ユーザーログイン画面
	r.GET("/login", func(c *gin.Context) {

		c.HTML(200, "login.html", gin.H{})
	})

	// ユーザーログイン
	r.POST("/login", func(c *gin.Context) {

		// DBから取得したユーザーパスワード(Hash)
		dbPassword := getUser(c.PostForm("username")).Password
		log.Println(dbPassword)
		// フォームから取得したユーザーパスワード
		formPassword := c.PostForm("password")

		// ユーザーパスワードの比較
		if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
			log.Println("ログインできませんでした")
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
			c.Abort()
		} else {
			// セッションの作成
			//格納データ取得
			sessionUserID := getUser(c.PostForm("username")).ID
			sessionUserName := getUser(c.PostForm("username")).Username
			//構造体へ格納
			sessionInfo.UserId = sessionUserID
			sessionInfo.UserName = sessionUserName
			sessionInfo.IsSessionAlive = true

			// //セッションへセット（ここら辺よくわからん）
			// session := sessions.Default(c)
			// session.Set("UserId", sessionUserID)
			// session.Set("UserName", sessionUserName)
			// session.Set("alive", true)
			// session.Save()

			//
			log.Println("ログインできました")
			log.Println(sessionInfo)

			//index.htmlへ埋め込むデータを取得
			t := time.Now()
			month.Month = t.Format("2006-01")

			var timecards []models.Timecard

			//ユーザーIDと月を指定してtimecard一覧を取得
			// timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
			timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
			listLen := len(timecards) == 0

			c.HTML(200, "index.html", gin.H{
				"sessioninfo": sessionInfo,
				"timecards":   timecards,
				"month":       month.Month,
				"listLen":     listLen,
			})
		}
	})

	//ログアウト
	r.GET("/logout", func(c *gin.Context) {
		//構造体へ格納したセッション情報を削除
		sessionInfo.UserId = nil
		sessionInfo.UserName = nil
		sessionInfo.IsSessionAlive = false

		month.Month = ""
		log.Println(sessionInfo)
		log.Println(month.Month)

		c.HTML(200, "top.html", gin.H{})

	})

	//表示月を変更
	r.POST("/timecard/select", func(c *gin.Context) {
		//formで送信されたデータをbind
		c.Bind(&month)

		// 表示月データ取得後一覧へ戻る
		var timecards []models.Timecard
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		listLen := len(timecards) == 0
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
			"listLen":     listLen,
		})
	})

	//timecard新規作成画面へ移動
	r.GET("/timecard/new", func(c *gin.Context) {

		c.HTML(http.StatusOK, "new.html", gin.H{
			"sessioninfo": sessionInfo,
			"month":       month.Month,
		})
	})

	// timecardListを初期化
	r.GET("/timecard/init", func(c *gin.Context) {
		//初期化用timecardlist
		var initTimecard models.Timecard
		initTimecard.UserID = sessionInfo.UserId.(uint)
		//月初
		fd, _ := time.Parse("2006-01-02", month.Month+"-01")
		//月末
		ld := fd.AddDate(0, 1, -1)
		//UserIDと日付のみ記述されたレコードを作成
		for i := 0; i < ld.Day(); i++ {
			initTimecard.Day = fd.AddDate(0, 0, i).Format("2006-01-02")
			//timecard新規作成処理
			createTimecard(initTimecard)
		}

		// 初期化処理終了後一覧へ戻る
		var timecards []models.Timecard
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		listLen := len(timecards) == 0
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
			"listLen":     listLen,
		})

	})

	//timecard新規作成
	r.POST("/timecard/new", func(c *gin.Context) {
		//formで送信されたデータをbind
		var t models.Timecard
		c.Bind(&t)
		//セッションよりユーザーIDを取得し格納
		t.UserID = sessionInfo.UserId.(uint)

		//timecard新規作成処理
		createTimecard(t)

		// 新規作成処理終了後一覧へ戻る
		var timecards []models.Timecard
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		listLen := len(timecards) == 0
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
			"listLen":     listLen,
		})
	})

	//timecard一覧表示
	r.GET("/timecard/index", func(c *gin.Context) {
		if err := c.Bind(&month); err != nil {
			t := time.Now()
			month.Month = t.Format("2006-01")
		}

		//timecard一覧格納用
		var timecards []models.Timecard
		// db := db.GetDB()

		//ユーザーIDと月を指定してtimecard一覧を取得
		// timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)

		// db.Find(&timecards) // DBから全てのレコードを取得する
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
		})
	})

	//timecard編集画面へ移動
	r.GET("/timecard/edit/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Println(id)
		var timecard models.Timecard
		// c.Bind(&timecard)
		// id := timecard.ID

		timecard = getTimecard(sessionInfo.UserId.(uint), uint(id))
		log.Println(timecard)

		//編集画面を表示
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"timecard": timecard,
		})
	})

	//timecardを修正
	r.POST("/timecard/edit/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		log.Println(id)
		s := c.PostForm("Start")
		e := c.PostForm("End")
		bt := c.PostForm("BreakTime")
		// timecard := getTimecard(sessionInfo.UserId.(uint), uint(id))
		// timecard.Start = s
		// timecard.End = e
		// timecard.BreakTime = bt

		updateTimecard(uint(id), s, e, bt)
		// c.Redirect(http.StatusFound, "/timecard/index")

		// 編集処理終了後一覧へ戻る
		var timecards []models.Timecard
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		listLen := len(timecards) == 0
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
			"listLen":     listLen,
		})

	})

	//timecardをクリア
	r.GET("/timecard/delete/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		timecard := getTimecard(sessionInfo.UserId.(uint), uint(id))
		timecard.Start = ""
		timecard.End = ""
		timecard.BreakTime = ""

		updateTimecard(uint(id), "", "", "")

		// 削除処理終了後一覧へ戻る
		var timecards []models.Timecard
		timecards = getTimecardList(sessionInfo.UserId.(uint), month.Month)
		listLen := len(timecards) == 0
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sessioninfo": sessionInfo,
			"timecards":   timecards,
			"month":       month.Month,
			"listLen":     listLen,
		})
	})

	r.Run()
	// r.Run(":" + os.Getenv("PORT"))
}

//--------------------------------------------
//以下、controller packageを認識しないためここに記載する
//--------------------------------------------

// ユーザー登録処理
func createUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	// db := gormConnect()
	db := db.GetDB()
	// defer db.Close()
	// Insert処理
	if err := db.Create(&models.User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil

}

// ユーザーを一件取得
func getUser(username string) models.User {
	db := db.GetDB()
	var user models.User
	db.First(&user, "username = ?", username)
	// db.Close()
	return user
}

// 新規作成
func createTimecard(m models.Timecard) {
	db := db.GetDB()
	log.Println("createTimecard：DBに接続しました")
	db.Create(&m)
	log.Println("createTimecard：新規作成しました")
}

// 指定月のリストの初期化
func initTimecardList(l []models.Timecard) {
	db := db.GetDB()
	log.Println("initTimecardList：DBに接続しました")
	db.Create(&l)
	log.Println("initTimecardList：指定月のタイムカードリストを作成しました")

}

//timecard一覧（ユーザー、月指定）取得
// func getTimecardList(UserID string, month string) []models.Timecard {
func getTimecardList(UserID uint, month string) []models.Timecard {
	db := db.GetDB()
	log.Println("getTimecardList：DBに接続しました")
	var timecards []models.Timecard
	selectMonth := month + "%"
	db.Where("user_id = ?", UserID).Where("day LIKE ?", selectMonth).Order("day").Find(&timecards)
	// db.Find(&timecards)
	log.Println("getTimecardList：データを取得しました")
	return timecards

}

//編集用timecard取得
func getTimecard(UserID uint, id uint) models.Timecard {
	db := db.GetDB()
	log.Println("getTimecard：DBに接続しました")
	var timecard models.Timecard
	db.Where("user_id = ? AND id = ?", UserID, id).First(&timecard)
	// db.Close()
	return timecard
}

func updateTimecard(id uint, s string, e string, bt string) {
	db := db.GetDB()
	log.Println("updateTimecard：DBに接続しました")
	// ID := timecard.ID
	log.Println(id)
	db.Model(&models.Timecard{}).Where("id = ?", id).Update("Start", s)
	db.Model(&models.Timecard{}).Where("id = ?", id).Update("End", e)
	db.Model(&models.Timecard{}).Where("id = ?", id).Update("BreakTime", bt)
	log.Println("updateTimecard：データを更新しました")

}

func deleteTimecard(id uint) {
	db := db.GetDB()
	log.Println("deleteTimecard：DBに接続しました")
	db.Delete(&models.Timecard{}, id)
	log.Println("deleteTimecard：データを削除しました")
}
