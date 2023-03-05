package controller

import (
	"bmacharia/aws_sam_blog_api/model"
	"bmacharia/aws_sam_blog_api/util"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create Article
func CreateArticle(c *gin.Context) {
	var input model.Article

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article := model.Article{
		Title:   input.Title,
		Content: input.Content,
		UserID:  util.CurrentUser(c).ID,
	}
	savedArticle, err := article.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": savedArticle})
}

// get Articles
func GetArticles(c *gin.Context) {
	var Article []model.Article
	err := model.GetArticles(&Article)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Article)
}

// get Article by id
func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Article model.Article
	err := model.GetArticle(&Article, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Article)
}

// update Article
func UpdateArticle(c *gin.Context) {
	var Article model.Article
	id, _ := strconv.Atoi(c.Param("id"))
	err := model.GetArticle(&Article, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&Article)
	err = model.UpdateArticle(&Article)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Article)
}
