// article
package controllers

import (
	"bytes"
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BindArticleApi(m *martini.ClassicMartini) {
	m.Post("/1/article/new",
		binding.Json(newArticleForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		newArticleHandler)
	m.Post("/1/article/delete",
		binding.Json(deleteArticleForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		deleteArticleHandler)
	m.Post("/1/article/thumb",
		binding.Json(articleThumbForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		articleThumbHandler)
	m.Get("/1/article/is_thumbed",
		binding.Form(articleIsThumbedForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		articleIsThumbedHandler)
	m.Get("/1/article/timelines",
		binding.Form(articleListForm{}),
		ErrorHandler,
		articleListHandler)
	m.Get("/1/article/get",
		binding.Form(articleInfoForm{}),
		ErrorHandler,
		articleInfoHandler)
	m.Post("/1/article/comments",
		binding.Json(articleCommentsForm{}),
		ErrorHandler,
		articleCommentsHandler)
	m.Get("/1/aritcle/thumbList",
		binding.Form(thumbersForm{}),
		thumbersHandler)
}

type articleJsonStruct struct {
	Id         string           `json:"article_id"`
	Parent     string           `json:"parent_article_id"`
	Author     string           `json:"author"`
	Title      string           `json:"cover_text"`
	Image      string           `json:"cover_image"`
	PubTime    int64            `json:"time"`
	Thumbs     int              `json:"thumb_count"`
	ThumbUsers []string         `json:"thumb_users"`
	NewThumbs  int              `json:"new_thumb_count"`
	Reviews    int              `json:"sub_article_count"`
	NewReviews int              `json:"new_sub_article_count"`
	Contents   []models.Segment `json:"article_segments"`
	Content    string           `json:"content"`
	Rewards    int64            `json:"reward_total"`
	Relation   string           `json:"relation"`
}

var (
	header = `<!DOCTYPE HTML>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=310, initial-scale=1, maximum-scale=1, user-scalable=no">
		<style>
			body{
				font-size:16px;
				line-height:30px;
				background-color:#f6f6f6;
				text-align:center;
				margin: 0;
			}
			p{
				text-align:left;
				padding-left: 5px;
				padding-right: 5px;
				word-wrap:break-word;
			}
			img{
				max-width:97%;
				height:auto;
				margin:auto;
			}
			div.divimg {
				text-align:center;
			}
		</style>
	</head>
	<body>`

	footer = `
	</body>
</html>`
)

func convertArticle(article *models.Article) *articleJsonStruct {
	jsonStruct := &articleJsonStruct{}
	jsonStruct.Id = article.Id.Hex()
	jsonStruct.Parent = article.Parent
	jsonStruct.Author = article.Author
	//jsonStruct.Contents = article.Contents
	jsonStruct.PubTime = article.PubTime.Unix()
	jsonStruct.Thumbs = len(article.Thumbs)
	jsonStruct.Reviews = len(article.Reviews)
	jsonStruct.Rewards = article.TotalReward

	jsonStruct.Title = article.Title
	jsonStruct.Image = article.Image

	if len(article.Content) > 0 {
		jsonStruct.Content = header + article.Content + footer
	}
	if len(article.Contents) > 0 {
		jsonStruct.Content = header + content2Html(article.Contents) + footer
	}

	return jsonStruct
}

func content2Html(contents []models.Segment) string {
	buf := &bytes.Buffer{}
	for _, content := range contents {
		switch strings.ToUpper(content.ContentType) {
		case "TEXT":
			s := strings.Split(content.ContentText, "\n")
			for _, a := range s {
				if a = strings.Trim(a, "\n"); len(a) > 0 {
					buf.WriteString("<p>" + a + "</p>")
				}
			}
		case "IMAGE":
			buf.WriteString("<div class=\"divimg\"><img src=\"" + content.ContentText + "\" /></div>")
		}
	}

	return buf.String()
}

type newArticleForm struct {
	Parent   string           `json:"parent_article_id"`
	Contents []models.Segment `json:"article_segments" binding:"required"`
	Tags     []string         `json:"article_tag"`
	parameter
}

func articleCover(contents []models.Segment) (text string, images []string) {
	for _, seg := range contents {
		if len(text) == 0 && strings.ToUpper(seg.ContentType) == "TEXT" {
			text = seg.ContentText
		}
		if strings.ToUpper(seg.ContentType) == "IMAGE" {
			images = append(images, seg.ContentText)
		}
	}
	return
}
func newArticleHandler(request *http.Request, resp http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(newArticleForm)

	article := &models.Article{
		Author:  user.Id,
		Content: content2Html(form.Contents),
		PubTime: time.Now(),
		Parent:  form.Parent,
		Tags:    form.Tags,
	}
	article.Title, article.Images = articleCover(form.Contents)
	if len(article.Images) > 0 {
		article.Image = article.Images[0]
	}

	if len(article.Tags) == 0 {
		article.Tags = []string{"SPORT_LOG"}
	}

	awards := Awards{}
	parent := &models.Article{}
	if len(form.Parent) > 0 {
		if find, err := parent.FindById(form.Parent); !find {
			e := errors.NewError(errors.NotExistsError, "文章不存在!")
			if err != nil {
				e = errors.NewError(errors.DbError)
			}
			writeResponse(request.RequestURI, resp, nil, e)
			return
		}

		if redis.Relationship(parent.Author, user.Id) == models.RelBlacklist {
			writeResponse(request.RequestURI, resp, nil,
				errors.NewError(errors.AccessError, "对方屏蔽了你!"))
			return
		}

		awards = Awards{Literal: 2 + user.Level(), Score: 2 + user.Level()}
	} else {
		awards = Awards{Literal: 10 + user.Level(), Wealth: 10 * models.Satoshi, Score: 10 + user.Level()}
	}

	if err := article.Save(); err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if err := GiveAwards(user, awards, redis); err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError))
		return
	}

	// comment
	if len(form.Parent) > 0 {
		//u := &models.User{Id: parent.Author}
		author := &models.Account{}
		author.FindByUserid(parent.Author)

		// ws push
		event := &models.Event{
			Type: models.EventArticle,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventComment,
				Id:   parent.Id.Hex(),
				From: user.Id,
				To:   parent.Author,
				Body: []models.MsgBody{
					{Type: "total_count", Content: strconv.Itoa(parent.CommentCount())},
					{Type: "image", Content: parent.Image},
				},
			},
		}
		/*
			if err := event.Save(); err == nil {
				redis.IncrEventCount(parent.Author, event.Data.Type, 1)
			}
		*/
		event.Save()

		event.Data.Body = append(event.Data.Body,
			models.MsgBody{Type: "new_count",
				Content: strconv.Itoa(models.EventCount(event.Data.Type, event.Data.Id, event.Data.To))})
		redis.PubMsg(models.EventArticle, parent.Author, event.Bytes())
		// apple push
		if author.Push {
			go sendApn(client, user.Nickname+"评论了你的主题!", author.Devs...)
		}
	}

	respData := map[string]interface{}{
		//"articles_without_content": convertArticle(article),
		"ExpEffect": awards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type deleteArticleForm struct {
	Id string `json:"article_id" binding:"required"`
	parameter
}

func deleteArticleHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(deleteArticleForm)

	article := &models.Article{}
	article.Author = user.Id
	article.Id = bson.ObjectIdHex(form.Id)

	err := article.Remove()
	writeResponse(request.RequestURI, resp, nil, err)
}

type articleThumbForm struct {
	Id     string `json:"article_id" binding:"required"`
	Status bool   `json:"thumb_status"`
	parameter
}

func articleThumbHandler(request *http.Request, resp http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(articleThumbForm)
	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		e := errors.NewError(errors.NotExistsError, "文章不存在!")
		if err != nil {
			e = errors.NewError(errors.DbError)
		}
		writeResponse(request.RequestURI, resp, nil, e)
		return
	}

	if redis.Relationship(article.Author, user.Id) == models.RelBlacklist {
		writeResponse(request.RequestURI, resp, nil,
			errors.NewError(errors.AccessError, "对方屏蔽了你!"))
		return
	}

	if err := article.SetThumb(user.Id, form.Status); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	awards := Awards{}
	if form.Status {
		awards = Awards{Score: 1, Wealth: 1 * models.Satoshi}
		GiveAwards(user, awards, redis)
	}
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, nil)

	author := &models.Account{Id: article.Author}

	// ws push
	event := &models.Event{
		Type: models.EventArticle,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventThumb,
			Id:   article.Id.Hex(),
			From: user.Id,
			To:   author.Id,
			Body: []models.MsgBody{
				{Type: "total_count", Content: strconv.Itoa(len(article.Thumbs) + 1)},
				{Type: "image", Content: article.Image},
			},
		},
	}

	if form.Status {
		author.FindByUserid(article.Author)
		/*
			if err := event.Save(); err == nil {
				redis.IncrEventCount(article.Author, event.Data.Type, 1)
			}
		*/
		event.Upsert()

		event.Data.Body = append(event.Data.Body,
			models.MsgBody{Type: "new_count",
				Content: strconv.Itoa(models.EventCount(event.Data.Type, event.Data.Id, event.Data.To))})
		redis.PubMsg(models.EventArticle, article.Author, event.Bytes())

		// apple push
		if author.Push {
			go sendApn(client, user.Nickname+"赞了你的主题!", author.Devs...)
		}
	} else {
		//count := author.DelEvent(models.EventThumb, article.Id.Hex(), user.Id, author.Id)
		//redis.IncrEventCount(author.Id, models.EventThumb, -count)
		event.Delete()
	}
	//user.UpdateAction(ActThumb, nowDate())
}

type articleIsThumbedForm struct {
	Id string `form:"article_id" binding:"required"`
	parameter
}

func articleIsThumbedHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(articleIsThumbedForm)

	article := &models.Article{}
	article.Id = bson.ObjectIdHex(form.Id)
	b, err := article.IsThumbed(user.Id)

	respData := map[string]bool{"is_thumbed": b}
	writeResponse(request.RequestURI, resp, respData, err)
}

type articleListForm struct {
	Circle bool   `form:"IsAttentionCircle"`
	Token  string `form:"access_token"`
	models.Paging
	Tag string `form:"article_tag"`
}

func articleListHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleListForm) {
	_, articles, err := models.GetArticles(form.Tag, &form.Paging, true)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	jsonStructs := make([]*articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i] = convertArticle(&articles[i])
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["articles_without_content"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}

type articleInfoForm struct {
	Id    string `form:"article_id" binding:"required"`
	Token string `form:"access_token"`
}

func articleInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form articleInfoForm) {
	article := &models.Article{}
	if find, err := article.FindById(form.Id); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if redis.OnlineUser(form.Token) == article.Author {
		event := &models.Event{}
		event.Data.Type = models.EventThumb
		event.Data.Id = article.Id.Hex()
		event.Data.To = article.Author
		event.Clear()

		event.Data.Type = models.EventComment
		event.Clear()

		event.Data.Type = models.EventReward
		event.Clear()
	}

	jsonStruct := convertArticle(article)
	jsonStruct.Relation = redis.Relationship(redis.OnlineUser(form.Token), article.Author)
	switch jsonStruct.Relation {
	case models.RelFriend:
		jsonStruct.Relation = "FRIENDS"
	case models.RelFollowing:
		jsonStruct.Relation = "ATTENTION"
	case models.RelFollower:
		jsonStruct.Relation = "FANS"
	case models.RelBlacklist:
		jsonStruct.Relation = "DEFRIEND"
	}
	thumbers := article.Thumbs
	if len(article.Thumbs) > 6 {
		thumbers = article.Thumbs[len(article.Thumbs)-6:]
	}
	users, _ := models.FindUsersByIds(0, thumbers...)
	for i := len(thumbers); i > 0; i-- { // reverse
		for j, _ := range users {
			if users[j].Id == thumbers[i-1] {
				jsonStruct.ThumbUsers = append(jsonStruct.ThumbUsers, users[j].Profile)
				break
			}
		}
	}

	writeResponse(request.RequestURI, resp, jsonStruct, nil)
}

type articleCommentsForm struct {
	Id string `json:"article_id"  binding:"required"`
	models.Paging
}

func articleCommentsHandler(request *http.Request, resp http.ResponseWriter,
	form articleCommentsForm) {

	article := &models.Article{Id: bson.ObjectIdHex(form.Id)}
	_, comments, err := article.Comments(&form.Paging, true)

	jsonStructs := make([]*articleJsonStruct, len(comments))
	for i, _ := range comments {
		jsonStructs[i] = convertArticle(&comments[i])
	}

	respData := make(map[string]interface{})
	respData["page_frist_id"] = form.Paging.First
	respData["page_last_id"] = form.Paging.Last
	//respData["page_item_count"] = total
	respData["articles_without_content"] = jsonStructs
	writeResponse(request.RequestURI, resp, respData, err)
}

type thumbersForm struct {
	Id    string `form:"article_id"`
	Index int    `form:"page_index"`
}

func thumbersHandler(r *http.Request, w http.ResponseWriter,
	form thumbersForm) {

	article := &models.Article{}
	article.FindById(form.Id)

	var respData struct {
		Users []*leaderboardResp `json:"members_list"`
	}

	if form.Index < 0 {
		form.Index = 0
	}
	thumbers := article.Thumbs
	end := len(thumbers) - form.Index*20
	start := end - 20

	if end <= 0 {
		respData.Users = []*leaderboardResp{}
		writeResponse(r.RequestURI, w, respData, nil)
		return
	}
	if start < 0 {
		start = 0
	}

	thumbers = article.Thumbs[start:end]
	users, _ := models.FindUsersByIds(1, thumbers...)

	for j := len(thumbers); j > 0; j-- { // reverse
		for i, _ := range users {
			if users[i].Id == thumbers[j-1] {
				respData.Users = append(respData.Users, &leaderboardResp{
					Userid:   users[i].Id,
					Score:    users[i].Props.Score,
					Level:    users[i].Level(),
					Profile:  users[i].Profile,
					Nickname: users[i].Nickname,
					Gender:   users[i].Gender,
					LastLog:  users[i].LastLogin.Unix(),
					Birth:    users[i].Birth,
					Location: users[i].Loc,
					Addr:     users[i].LocAddr,
					Phone:    users[i].Phone,
				})
				break
			}
		}
	}

	writeResponse(r.RequestURI, w, respData, nil)
	return

}
