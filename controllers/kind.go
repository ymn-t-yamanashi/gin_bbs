package controllers

import (
	"gin_bbs/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KindController struct {
	DB *gorm.DB
}

// 一覧画面: 明示的に index.html を呼び出す
func (ctrl *KindController) Index(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	session.Save()

	var posts []models.Kind
	ctrl.DB.Order("created_at desc").Find(&posts)

	// 第2引数を "index.html" と明示
	c.HTML(http.StatusOK, "kind_index.html", gin.H{
		"posts": posts,
		"flash": flash,
	})
}

// 新規作成画面: 明示的に new.html を呼び出す
func (ctrl *KindController) New(c *gin.Context) {
	c.HTML(http.StatusOK, "kind_new.html", gin.H{
		"post": models.Kind{},
	})
}

// 保存処理
func (ctrl *KindController) Create(c *gin.Context) {
	var post models.Kind
	// フォームデータのバインド
	if err := c.ShouldBind(&post); err != nil {
		c.HTML(http.StatusBadRequest, "kind_new.html", gin.H{
			"errors": []string{"入力に不備があります"},
			"post":   post,
		})
		return
	}

	// データベース保存
	if err := ctrl.DB.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "kind_new.html", gin.H{
			"errors": []string{"保存に失敗しました"},
		})
		return
	}

	session := sessions.Default(c)
	session.AddFlash("投稿を作成しました")
	session.Save()

	// 一覧へリダイレクト
	c.Redirect(http.StatusFound, "/kind/")
}

// 詳細画面: 明示的に show.html を呼び出す
func (ctrl *KindController) Show(c *gin.Context) {
	var post models.Kind
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "kind_show.html", gin.H{
		"post": post,
	})
}

// 編集画面: 明示的に edit.html を呼び出す
func (ctrl *KindController) Edit(c *gin.Context) {
	var post models.Kind
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "kind_edit.html", gin.H{
		"post": post,
	})
}

// 更新処理
func (ctrl *KindController) Update(c *gin.Context) {
	var post models.Kind
	if err := ctrl.DB.First(&post, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.ShouldBind(&post)
	ctrl.DB.Save(&post)
	c.Redirect(http.StatusFound, "/kind/"+c.Param("id"))
}

// 削除処理
func (ctrl *KindController) Delete(c *gin.Context) {
	id := c.Param("id")

	// データベースから削除
	if err := ctrl.DB.Delete(&models.Kind{}, id).Error; err != nil {
		// エラーがあれば一覧に戻してメッセージ（本来はフラッシュ等で通知）
		c.Redirect(http.StatusFound, "/kind/")
		return
	}

	// 削除後は一覧にリダイレクト
	c.Redirect(http.StatusFound, "/kind/")
}
