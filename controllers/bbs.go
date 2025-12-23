package controllers

import (
	"gin_bbs/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BBSController struct {
	DB *gorm.DB
}

// 一覧画面: 明示的に index.html を呼び出す
func (ctrl *BBSController) Index(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	session.Save()

	var posts []models.BBS
	ctrl.DB.Order("created_at desc").Find(&posts)

	// 第2引数を "index.html" と明示
	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts": posts,
		"flash": flash,
	})
}

// 新規作成画面: 明示的に new.html を呼び出す
func (ctrl *BBSController) New(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", gin.H{
		"post": models.BBS{},
	})
}

// 保存処理
func (ctrl *BBSController) Create(c *gin.Context) {
	var post models.BBS
	// フォームデータのバインド
	if err := c.ShouldBind(&post); err != nil {
		c.HTML(http.StatusBadRequest, "new.html", gin.H{
			"errors": []string{"入力に不備があります"},
			"post":   post,
		})
		return
	}

	// データベース保存
	if err := ctrl.DB.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "new.html", gin.H{
			"errors": []string{"保存に失敗しました"},
		})
		return
	}

	session := sessions.Default(c)
	session.AddFlash("投稿を作成しました")
	session.Save()

	// 一覧へリダイレクト
	c.Redirect(http.StatusFound, "/bbs/")
}

// 詳細画面: 明示的に show.html を呼び出す
func (ctrl *BBSController) Show(c *gin.Context) {
	var post models.BBS
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "show.html", gin.H{
		"post": post,
	})
}

// 編集画面: 明示的に edit.html を呼び出す
func (ctrl *BBSController) Edit(c *gin.Context) {
	var post models.BBS
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"post": post,
	})
}

// 更新処理
func (ctrl *BBSController) Update(c *gin.Context) {
	var post models.BBS
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.ShouldBind(&post)
	ctrl.DB.Save(&post)
	c.Redirect(http.StatusFound, "/bbs/"+c.Param("id"))
}

// 削除処理
func (ctrl *BBSController) Delete(c *gin.Context) {
	ctrl.DB.Delete(&models.BBS{}, c.Param("id"))
	c.Redirect(http.StatusFound, "/bbs/")
}
