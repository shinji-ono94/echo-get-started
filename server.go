package main

import (
    "io"
    "github.com/labstack/echo"
    "html/template"
    "net/http"
    "html"
    "github.com/ipfans/echo-session"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// インスタンスを作成
	var e = echo.New()

	//テンプレートの設定
    t := &Template{
        templates: template.Must(template.ParseGlob("static/template/*.html")),
    }
	
	e.Renderer = t

	//セッションを設定
    store := session.NewCookieStore([]byte("secret-key"))
        //セッション保持時間
    store.MaxAge(86400)
    e.Use(session.Sessions("ESESSION", store))
 
    e.GET("/login", ShowLoginHtml)
    e.POST("/login", Login)
 
    // ポート9000でサーバーを起動
    e.Logger.Fatal(e.Start(":9000"))
 
    // Let's Encrypt から証明書を自動取得してhttpsサーバーを起動
    //e.Logger.Fatal(e.StartAutoTLS(":443"))
}

type LoginForm struct {
    UserId string
    Password string
    ErrorMessage string
}

type CompleteJson struct {
    Success bool `json:"success"`
}
 
func ShowLoginHtml(c echo.Context) error {
    session := session.Default(c)
 
    loginId := session.Get("loginCompleted")
    if loginId != nil && loginId == "completed" {
        completeJson := CompleteJson{
            Success: true,
        }
 
        return c.JSON(http.StatusOK, completeJson)
    }
 
    return c.Render(http.StatusOK, "login", LoginForm{})
}
 
func Login(c echo.Context) error {
    loginForm := LoginForm{
        UserId: c.FormValue("userId"),
        Password: c.FormValue("password"),
    }
 
    userId := html.EscapeString(loginForm.UserId)
    password := html.EscapeString(loginForm.Password)
 
    if userId != "userId" && password != "password" {
        loginForm.ErrorMessage = "ユーザーID または パスワードが間違っています"
        return c.Render(http.StatusOK, "login", loginForm)
    }
 
    //セッションにデータを保存する
    session := session.Default(c)
    session.Set("loginCompleted", "completed")
    session.Save()
 
    completeJson := CompleteJson{
        Success: true,
    }
 
    return c.JSON(http.StatusOK, completeJson)
}