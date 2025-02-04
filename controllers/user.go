// user
package controllers

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	//"log"
	//"math/rand"
	"net/http"
	//"net/url"
	//"strconv"
	//"strings"
	"time"
)

func BindUserApi(m *martini.ClassicMartini) {
	m.Post("/1/user/send_device_token",
		binding.Json(sendDevForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		sendDevHandler)
	m.Post("/1/user/set_push_enable",
		binding.Json(setPushForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		setPushHandler)
	m.Get("/1/user/is_push_enabled",
		binding.Form(pushStatusForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		pushStatusHandler)
	m.Post("/1/user/enableAttention",
		binding.Json(relationshipForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		followHandler)
	m.Post("/1/user/enableDefriend",
		binding.Json(relationshipForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		checkLimitHandler,
		blacklistHandler)
	m.Get("/1/user/getAttentionFriendsList",
		binding.Form(getFollowsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		getFollowsHandler)
	m.Get("/1/user/getAttentedMembersList",
		binding.Form(getFollowsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		getFollowersHandler)
	m.Get("/1/user/getJoinedGroupsList",
		binding.Form(getFollowsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		getGroupsHandler)
	m.Get("/1/user/getRelatedMembersList",
		binding.Form(socialListForm{}),
		ErrorHandler,
		socialListHandler)
}

type sendDevForm struct {
	Dev string `json:"device_token" binding:"required"`
	parameter
}

func sendDevHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(sendDevForm)
	err := user.AddDevice(form.Token)
	writeResponse(request.RequestURI, resp, nil, err)
}

type setPushForm struct {
	Enabled bool `json:"is_enabled"`
	parameter
}

func setPushHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(setPushForm)
	err := user.SetPush(form.Enabled)
	writeResponse(request.RequestURI, resp, nil, err)
}

type pushStatusForm struct {
	parameter
}

func pushStatusHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	enabled, err := user.PushEnabled()
	writeResponse(request.RequestURI, resp, map[string]bool{"is_enabled": enabled}, err)
}

type relationshipForm struct {
	Userids   []string `json:"userids"`
	Follow    bool     `json:"bAttention"`
	Blacklist bool     `json:"bDefriend"`
	parameter
}

func followHandler(request *http.Request, resp http.ResponseWriter,
	client *ApnClient, redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(relationshipForm)

	for _, peer := range form.Userids {
		if redis.Relationship(peer, user.Id) == models.RelBlacklist {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, "对方已屏蔽了你!"))
			return
		}
	}

	redis.SetRelationship(user.Id, form.Userids, models.RelFollowing, form.Follow)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, nil)

	for _, userid := range form.Userids {
		u := &models.Account{}
		u.FindByUserid(userid)

		event := &models.Event{
			Type: models.EventMsg,
			Time: time.Now().Unix(),
			Data: models.EventData{
				Type: models.EventSub,
				Id:   user.Id + "-" + u.Id,
				From: user.Id,
				To:   u.Id,
				Body: []models.MsgBody{
					{Type: "nikename", Content: user.Nickname},
					{Type: "image", Content: user.Profile},
				},
			},
		}

		if form.Follow {
			/*
				if err := event.Save(); err == nil {
					redis.IncrEventCount(u.Id, event.Data.Type, 1)
				}
			*/
			event.Upsert()

			event.Data.Body = append(event.Data.Body,
				models.MsgBody{Type: "new_count", Content: "1"})
			redis.PubMsg(models.EventMsg, u.Id, event.Bytes())

			// apple push
			if u.Push {
				go sendApn(client, user.Nickname+"关注了你!", u.EventCount("", ""), u.Devs...)
			}
		} else {
			//count := u.DelEvent(models.EventSub, user.Id, user.Id, u.Id)
			//redis.IncrEventCount(u.Id, models.EventSub, -count)
			event.Delete()
		}
	}
}

func blacklistHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(relationshipForm)

	redis.SetRelationship(user.Id, form.Userids, models.RelBlacklist, form.Blacklist)

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, nil)
}

type getFollowsForm struct {
	parameter
}

func getFollowsHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	writeResponse(request.RequestURI, resp, redis.Friends(models.RelFollowing, user.Id), nil)
}

func getFollowersHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	writeResponse(request.RequestURI, resp, redis.Friends(models.RelFollower, user.Id), nil)
}

func getGroupsHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	writeResponse(request.RequestURI, resp, redis.Groups(user.Id), nil)
}

type socialListForm struct {
	Userid    string `form:"userid" binding:"required"`
	ArticleId string `form:"article_id"`
	Type      string `form:"member_type" binding:"required"`
	models.Paging
}

func socialListHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form socialListForm) {

	user := &models.Account{Id: form.Userid}

	var ids []string
	switch form.Type {
	case "FRIENDS":
		ids = redis.Friends(models.RelFriend, user.Id)
	case "ATTENTION":
		ids = redis.Friends(models.RelFollowing, user.Id)
	case "FANS":
		ids = redis.Friends(models.RelFollower, user.Id)
	case "DEFRIEND":
		ids = redis.Friends(models.RelBlacklist, user.Id)
	case "WEIBO":
		ids = redis.Friends("weibo", user.Id)
	}
	users, err := models.Users(ids, &form.Paging)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	lb := make([]leaderboardResp, len(users))
	for i, _ := range users {
		lb[i].Userid = users[i].Id
		lb[i].Score = users[i].Props.Score
		lb[i].Level = users[i].Level()
		lb[i].Profile = users[i].Profile
		lb[i].Nickname = users[i].Nickname
		lb[i].Gender = users[i].Gender
		lb[i].LastLog = users[i].LastLogin.Unix()
		lb[i].Birth = users[i].Birth
		lb[i].Location = users[i].Loc
		lb[i].Phone = users[i].Phone
		lb[i].Actor = users[i].Actor
	}

	respData := map[string]interface{}{
		"members_list":  lb,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(request.RequestURI, resp, respData, nil)

	return
}
