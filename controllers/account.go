// account
package controllers

import (
	//"encoding/json"
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	//"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	//"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func BindAccountApi(m *martini.ClassicMartini) {
	m.Post("/1/account/registerV2",
		binding.Json(regFormV2{}),
		ErrorHandler,
		regHandlerV2)
	m.Post("/1/account/loginV2",
		binding.Json(loginFormV2{}),
		ErrorHandler,
		loginHandlerV2)
	m.Get("/1/account/check",
		binding.Form(accountCheckForm{}),
		ErrorHandler,
		accountCheckHandler)

	m.Post("/1/account/register",
		binding.Json(userRegForm{}),
		ErrorHandler,
		registerHandler)
	m.Post("/1/account/login",
		binding.Json(loginForm{}),
		ErrorHandler,
		loginHandler)
	m.Post("/1/account/changePassword",
		binding.Json(setPasswordForm{}, (*Parameter)(nil)),
		checkTokenHandler,
		loadUserHandler,
		setPasswordHandler)
	m.Get("/1/user/getDailyLoginRewardInfo",
		binding.Form(loginAwardsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		loginAwardsHandler)
	m.Post("/1/user/logout",
		binding.Json(logoutForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		logoutHandler)
	m.Get("/1/user/recommend",
		binding.Form(recommendForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		recommendHandler)
	m.Get("/1/user/getInfo",
		binding.Form(getInfoForm{}, (*Parameter)(nil)),
		ErrorHandler,
		userInfoHandler)
	m.Get("/1/user/getRelatedMembersCount",
		binding.Form(friendCountForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		friendCountHandler)
	m.Post("/1/user/setInfo",
		binding.Json(setInfoForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		setInfoHandler)
	m.Post("/1/user/setProfileImage",
		binding.Json(setProfileForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		setProfileHandler)
	/*
		m.Post("/1/account/importFriends",
			binding.Json(importFriendsForm{}, (*Parameter)(nil)),
			ErrorHandler,
			checkTokenHandler,
			importFriendsHandler)
	*/
	m.Post("/1/user/setLifePhotos",
		binding.Json(setPhotosForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		setPhotosHandler)
	m.Post("/1/user/deleteLifePhoto",
		binding.Json(delPhotoForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		delPhotoHandler)
	//m.Get("/1/user/news", binding.Form(userNewsForm{}), ErrorHandler, userNewsHandler)
	//m.Get("/1/users", binding.Form(userListForm{}), ErrorHandler, userListHandler)

	m.Get("/1/user/getPKPropertiesInfo",
		binding.Form(scoreDiffForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		scoreDiffHandler)
	m.Get("/1/user/getPropertiesValue",
		binding.Form(getPropsForm{}),
		getPropsHandler)
	m.Post("/1/user/updateEquipment",
		binding.Json(setEquipForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		setEquipHandler)
	m.Get("/1/user/search",
		binding.Form(searchForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		searchHandler)
	m.Get("/1/user/articles",
		binding.Form(userArticlesForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		userArticlesHandler)
	m.Post("/1/user/importContacts",
		binding.Json(importContactsForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		importContactsHandler)
	m.Post("/1/user/resetPassword",
		binding.Json(resetPasswdForm{}),
		resetPasswdHandler)
	m.Post("/1/user/shareToFriends",
		binding.Json(pkShareForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		pkShareHandler)
	m.Get("/1/user/isNikeNameUsed",
		binding.Form(nicknameForm{}),
		checkNicknameHandler)
	m.Get("/1/user/isPhoneUsed",
		binding.Form(checkPhoneForm{}),
		checkPhoneHandler)
	m.Get("/1/user/gameResults",
		binding.Form(gameResultForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		gameResultHandler)

	m.Get("/1/user/leaderboard",
		binding.Form(userlbForm{}),
		userLeaderBoardHandler)

	m.Post("/1/user/auth/request",
		binding.Json(authRequestForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		userAuthRequestHandler)
	m.Get("/1/user/auth/status",
		binding.Form(authStatusForm{}),
		userAuthStatusHandler)
	m.Get("/1/user/auth/info",
		binding.Form(authInfoForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		userAuthInfoHandler)

	m.Post("/1/user/purchaseSuccess",
		binding.Json(purchaseForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		purchaseHandler)

	m.Get("/1/user/getPayHistory",
		binding.Form(purchaseListForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		purchaseListHandler)

	m.Get("/1/user/isPreSportForm", testHandler)

	m.Post("/1/user/sendheart",
		binding.Json(sendHeartForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		sendHeartHandler)
	m.Post("/1/user/recvheart",
		binding.Json(recvHeartForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		recvHeartHandler)
	m.Get("/1/user/ranks",
		binding.Form(userRanksForm{}),
		ErrorHandler,
		checkTokenHandler,
		loadUserHandler,
		userRanksHandler)
}

type regFormV2 struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     string `json:"type"`
	Nickname string `json:"nickname" binding:"required"`
	Profile  string `json:"profile" binding:"required"`
	Gender   string `json:"gender"`
	Birthday int64  `json:"birthday" binding:"required"`
}

func regHandlerV2(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form regFormV2) {

	exists, err := models.CheckUserExists(form.Id, form.Type)
	if err != nil {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.DbError))
		return
	}
	if exists {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.UserExistError))
		return
	}

	dbw, err := getNewWallet()
	if err != nil {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.DbError, "创建钱包失败"))
		return
	}

	user := &models.Account{
		Role:     form.Type,
		Password: Md5(form.Password),
		Nickname: form.Nickname,
		Profile:  form.Profile,
		Gender:   form.Gender,
		Birth:    form.Birthday,
		Wallet:   *dbw,
		RegTime:  time.Now(),
	}

	switch form.Type {
	case models.AccountPhone:
		user.Phone = form.Id
	case models.AccountWeibo:
		user.Weibo = form.Id
	case models.AccountEmail:
		fallthrough
	default:
		user.Email = form.Id
	}

	if err := user.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	token := Uuid() + "-" + strconv.FormatInt(time.Now().AddDate(0, 0, 30).Unix(), 10)
	data := map[string]string{"access_token": token, "userid": user.Id}
	writeResponse(request.RequestURI, resp, data, nil)

	redis.LogRegister(user.Id, form.Type)
	redis.SetOnlineUser(token, user.Id)

	// ws push
	//regNotice(user.Id, redis)
}

type loginFormV2 struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     string `json:"type"`
}

func loginHandlerV2(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form loginFormV2) {
	user := &models.Account{}
	if exists, _ := models.CheckUserExists(form.Id, form.Type); !exists {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.AuthError, "该用户还未注册"))
		return
	}
	user.FindPass(form.Id, form.Type, Md5(form.Password))
	if len(user.Id) == 0 {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.AuthError))
		return
	}

	if user.TimeLimit < 0 {
		writeResponse(request.RequestURI, resp,
			nil, errors.NewError(errors.AuthError, "账户已禁用"))
		return
	}

	token := Uuid() + "-" + strconv.FormatInt(time.Now().AddDate(0, 0, 30).Unix(), 10)
	redis.SetOnlineUser(token, user.Id)

	data := map[string]interface{}{
		"access_token": token,
		"userid":       user.Id,
		//"last_login_time": user.LastLogin.Unix(),
	}
	writeResponse(request.RequestURI, resp, data, nil)
}

type accountCheckForm struct {
	Id   string `form:"id" binding:"required"`
	Type string `form:"type"`
}

func accountCheckHandler(request *http.Request, resp http.ResponseWriter,
	form accountCheckForm) {
	user := &models.Account{}
	user.Find(form.Id, form.Type)

	data := map[string]string{
		"userid": user.Id,
	}
	writeResponse(request.RequestURI, resp, data, nil)
}

// user register parameter
type userRegForm struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nikename"`
	Password string `json:"password" binding:"required"`
	//Role     string `json:"role"`
}

func regNotice(uid string, redis *models.RedisLogger) {
	notice := &models.Event{
		Type: models.EventWallet,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventTx,
			Id:   uid,
			From: uid,
			Body: []models.MsgBody{
				{Type: "rule", Content: "1"},
			},
		},
	}
	redis.Notice(notice.Bytes())
}

func registerHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form userRegForm) {
	user := &models.Account{}
	t := ""
	if phone, _ := strconv.ParseUint(form.Email, 10, 64); phone > 0 {
		user.Phone = form.Email
		t = models.AccountPhone
	} else {
		user.Email = strings.ToLower(form.Email)
		t = models.AccountEmail
	}

	if exists, _ := user.Exists(t); exists {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.UserExistError))
		return
	}
	//user.Nickname = form.Nickname
	user.Password = Md5(form.Password)
	user.Role = t
	user.RegTime = time.Now()
	dbw, err := getNewWallet()
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.DbError, "wallet: "+err.Error()))
		return
	}
	user.Wallet = *dbw

	if err := user.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
	} else {
		token := Uuid() + "-" + strconv.FormatInt(time.Now().AddDate(0, 0, 30).Unix(), 10)
		data := map[string]string{"access_token": token, "userid": user.Id}
		writeResponse(request.RequestURI, resp, data, nil)

		redis.LogRegister(user.Id, t)
		redis.SetOnlineUser(token, user.Id)

		// ws push
		regNotice(user.Id, redis)
	}
}

func getNewWallet() (*models.DbWallet, error) {
	w := NewWallet()
	id, err := saveWallet("", w)
	var addrs []string
	for _, key := range w.Keys {
		addrs = append(addrs, key.PubKey)
	}
	return &models.DbWallet{Id: id, Key: w.SharedKey, Addr: w.Keys[0].PubKey, Addrs: addrs}, err
}

// user login parameter
type loginForm struct {
	Userid   string `json:"userid"`
	Password string `json:"verfiycode"`
	Type     string `json:"account_type"`
}

/*
func (this loginForm) getUserId() string {
	return this.Userid
}
*/
func weiboLogin(uid, password string, redis *models.RedisLogger) (bool, *models.Account, error) {
	user := &models.Account{}
	exists, err := user.FindByWeibo(strings.ToLower(uid))
	if err != nil {
		return false, nil, err
	}

	p := Md5(password)
	registered := user.RegTime.Unix() > 0

	if registered {
		if user.Password != p {
			user.SetPassword(p)
		}
		return false, user, nil
	}
	weiboUser, err := GetWeiboUserInfo(uid, password)
	if err != nil {
		return false, nil, err
	}

	user.Nickname = weiboUser.ScreenName
	if find, _ := user.Exists("nickname"); find {
		user.Nickname = "wb_" + user.Nickname
	}
	user.Password = p
	if !strings.HasPrefix(weiboUser.Gender, "f") {
		user.Gender = "male"
	} else {
		user.Gender = "female"
	}

	user.Weibo = uid
	user.Url = weiboUser.Url
	user.Profile = weiboUser.Avatar
	//user.Addr = &models.Address{Desc: weiboUser.Location}
	user.About = weiboUser.Description
	user.Role = models.AccountWeibo
	user.RegTime = time.Now()

	dbw, err := getNewWallet()
	if err != nil {
		return true, nil, err
	}
	user.Wallet = *dbw

	if !exists {
		if err := user.Save(); err != nil {
			return true, nil, err
		}
	}
	redis.LogRegister(user.Id, user.Role)

	// ws push
	regNotice(user.Id, redis)

	return true, user, nil
}
func guestLogin(redis *models.RedisLogger) (*models.Account, error) {
	user := &models.Account{}
	user.Id = models.GuestUserPrefix + strconv.Itoa(time.Now().Nanosecond()) + ":" + strconv.Itoa(random.Intn(65536))

	return user, nil
}

var ran = rand.New(rand.NewSource(time.Now().UnixNano()))

func loginHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form loginForm) {
	user := &models.Account{}
	var err error
	var reg bool
	token := Uuid() + "-" + strconv.FormatInt(time.Now().AddDate(0, 0, 30).Unix(), 10)

	switch form.Type {
	case "weibo":
		reg, user, err = weiboLogin(form.Userid, form.Password, redis)
	case "weixin":
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.UnimplementedError))
		return
	case "usrpass":
		fallthrough
	default:
		var find bool
		if find, err = user.FindByUserPass(strings.ToLower(form.Userid), Md5(form.Password)); !find {
			if err == nil {
				err = errors.NewError(errors.AuthError)
			}
		}
	}

	if err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if user.TimeLimit < 0 {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AuthError, "账户已禁用"))
		return
	}

	redis.SetOnlineUser(token, user.Id)

	data := map[string]interface{}{
		"access_token":    token,
		"userid":          user.Id,
		"register":        reg,
		"last_login_time": user.LastLogin.Unix(),
		"ExpEffect":       Awards{},
	}
	writeResponse(request.RequestURI, resp, data, nil)
}

type logoutForm struct {
	parameter
}

func logoutHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	redis.DelOnlineUser(p.TokenId())
	if len(user.Devs) > 0 {
		user.RmDevice(user.Devs[0])
	}
	writeResponse(request.RequestURI, resp, nil, nil)

}

type recommendForm struct {
	models.Paging
	parameter
}

func recommendHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(recommendForm)

	excludes := redis.Friends(models.RelFollowing, user.Id)
	excludes = append(excludes, redis.Friends(models.RelBlacklist, user.Id)...)

	users, _ := user.Recommend(excludes)
	var list []*leaderboardResp
	for i, _ := range users {
		if users[i].Id == user.Id {
			continue
		}
		rel := redis.Relationship(user.Id, users[i].Id)
		if rel == models.RelFollowing || rel == models.RelFriend {
			continue
		}

		lb := &leaderboardResp{
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
			Actor:    users[i].Actor,
		}
		lb.Distance, _ = redis.RecStats(users[i].Id)
		lb.Status = users[i].LatestArticle().Title
		list = append(list, lb)
	}

	respData := map[string]interface{}{
		"members_list":  list,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(r.RequestURI, w, respData, nil)
}

type userJsonStruct struct {
	Userid   string `json:"userid"`
	Nickname string `json:"nikename"`
	Email    string `json:"email"`
	Phone    string `json:"phone_number"`
	Type     string `json:"account_type"`
	About    string `json:"about"`
	Profile  string `json:"profile_image"`
	RegTime  int64  `json:"register_time"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Birth    int64  `json:"birthday"`

	Sign        string `json:"sign"`
	Emotion     string `json:"emotion"`
	Profession  string `json:"profession"`
	Hobby       string `json:"fond"`
	Hometown    string `json:"hometown"`
	OftenAppear string `json:"oftenAppear"`

	Actor string `json:"actor"`
	//Rank  string `json:"rankName"`
	//Followed bool   `json:"beFriend"`
	Online bool `json:"beOnline"`

	Props models.Props `json:"proper_info"`

	Addr string `json:"location_desc"`
	models.Location

	Gender    string `json:"sex_type"`
	Follows   int    `json:"attention_count"`
	Followers int    `json:"fans_count"`
	Posts     int    `json:"post_count"`

	Photos     []string     `json:"user_images"`
	CoverImage string       `json:"cover_image"`
	Equips     models.Equip `json:"user_equipInfo"`

	Wallet   string           `json:"wallet"`
	Relation string           `json:"relation"`
	LastLog  int64            `json:"last_login_time"`
	Setinfo  bool             `json:"setinfo"`
	Ban      int64            `json:"ban_time"`
	Auth     *models.UserAuth `json:"auth_info"`

	Pet       string `json:"pet"`
	ScoreRank int    `json:"score_rank"`
}

func convertUser(user *models.Account, redis *models.RedisLogger) *userJsonStruct {
	info := &userJsonStruct{
		Userid:   user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Type:     user.Role,
		About:    user.About,
		Profile:  user.Profile,
		RegTime:  user.RegTime.Unix(),
		Height:   user.Height,
		Weight:   user.Weight,
		Birth:    user.Birth,
		Actor:    user.Actor,
		Location: user.Loc,
		Addr:     user.LocAddr,

		Sign:        user.Sign,
		Emotion:     user.Emotion,
		Profession:  user.Profession,
		Hometown:    user.Hometown,
		OftenAppear: user.Oftenappear,
		Hobby:       user.Hobby,

		//Rank:   userRank(user.Level),
		Online: redis.IsOnline(user.Id),
		Gender: user.Gender,
		//Follows:   len(redis.Follows(user.Id)),
		//Followers: len(redis.Followers(user.Id)),
		Posts: user.ArticleCount(),

		//Props: redis.UserProps(user.Id),
		Props: models.Props{
			Physical: user.Props.Physical,
			Literal:  user.Props.Literal,
			Mental:   user.Props.Mental,
			Wealth:   redis.GetCoins(user.Id),
			Score:    user.Props.Score,
			Level:    user.Level(),
		},

		Photos:     user.Photos,
		CoverImage: user.CoverImage,

		Wallet:  user.Wallet.Addr,
		LastLog: user.LastLogin.Unix(),
		Setinfo: user.Setinfo,
		Ban:     user.TimeLimit,
		Auth:    user.Auth,
	}

	balance, _ := getBalance(user.Wallet.Addrs)
	var wealth int64
	if balance != nil {
		for _, b := range balance.Addrs {
			wealth += (b.Confirmed + b.Unconfirmed)
		}
	}
	info.Props.Wealth = wealth

	info.Follows, info.Followers, _, _ = redis.FriendCount(user.Id)

	/*
		if user.Privilege == 5 {
			info.Actor = "coach"
		} else if user.Privilege == 10 {
			info.Actor = "admin"
		}

			if user.Addr != nil {
				info.Addr = user.Addr.String()
			}
	*/

	if user.Equips != nil {
		info.Equips = *user.Equips
	}

	if info.Auth != nil {
		if info.Auth.IdCard == nil {
			info.Auth.IdCard = info.Auth.IdCardTmp
		}
		if info.Auth.Cert == nil {
			info.Auth.Cert = info.Auth.CertTmp
		}
		if info.Auth.Record == nil {
			info.Auth.Record = info.Auth.RecordTmp
		}
	}

	return info
}

type getInfoForm struct {
	Userid string `form:"userid" binding:"required"`
	parameter
}

func userInfoHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, p Parameter) {
	user := &models.Account{}
	form := p.(getInfoForm)
	if find, err := user.FindByUserid(form.Userid); !find {
		if err == nil {
			err = errors.NewError(errors.NotExistsError)
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	info := convertUser(user, redis)
	//info.Pet = userPet(info.Props.Level)
	info.ScoreRank, _ = user.PropIndex("score", user.Props.Score)

	if uid := redis.OnlineUser(p.TokenId()); len(uid) > 0 {
		relation := redis.Relationship(uid, user.Id)
		switch relation {
		case models.RelFriend:
			info.Relation = "FRIENDS"
		case models.RelFollowing:
			info.Relation = "ATTENTION"
		case models.RelFollower:
			info.Relation = "FANS"
		case models.RelBlacklist:
			info.Relation = "DEFRIEND"
		}
	}

	writeResponse(request.RequestURI, resp, info, nil)
}

func userPet(level int64) string {
	config := &models.Config{}
	config.Find()

	if len(config.Pets) == 0 {
		return ""
	}

	if level < 5 {
		return config.Pets[0]
	}

	if level < 10 && len(config.Pets) > 1 {
		return config.Pets[1]
	}

	if len(config.Pets) > 2 {
		return config.Pets[2]
	}

	return ""
}

type friendCountForm struct {
	parameter
}

func friendCountHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {
	follows, followers, friends, blacklist := redis.FriendCount(user.Id)
	respData := map[string]int{
		"friend_count":    friends,
		"attention_count": follows,
		"fans_count":      followers,
		"defriend_count":  blacklist,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type setInfoForm struct {
	models.UserInfo
	parameter
}

func calcInfo(user *models.Account, setinfo *models.SetInfo) int {
	n := 4.0

	if len(user.Photos) > 0 {
		n++
	}
	if setinfo.Height > 0 || user.Height > 0 {
		n++
	}
	if setinfo.Weight > 0 || user.Weight > 0 {
		n++
	}
	if setinfo.Sign != "" || user.Sign != "" {
		n++
	}
	if setinfo.Emotion != "" || user.Emotion != "" {
		n++
	}
	if setinfo.Profession != "" || user.Profession != "" {
		n++
	}
	if setinfo.Hobby != "" || user.Hobby != "" {
		n++
	}
	if setinfo.Hometown != "" || user.Hometown != "" {
		n++
	}
	if setinfo.OftenAppear != "" || user.Oftenappear != "" {
		n++
	}
	return int((n / 13.0) * 100.0)
}

func setInfoHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger,
	user *models.Account, p Parameter) {

	form := p.(setInfoForm)
	setinfo := &models.SetInfo{
		Phone:    form.UserInfo.Phone,
		Nickname: form.UserInfo.Nickname,
		Height:   form.UserInfo.Height,
		Weight:   form.UserInfo.Weight,
		Birth:    form.UserInfo.Birth,
		Gender:   form.UserInfo.Gender,

		Sign:        form.UserInfo.Sign,
		Emotion:     form.UserInfo.Emotion,
		Profession:  form.UserInfo.Profession,
		Hobby:       form.UserInfo.Hobby,
		Hometown:    form.UserInfo.Hometown,
		OftenAppear: form.UserInfo.OftenAppear,
		CoverImage:  form.UserInfo.CoverImage,

		Setinfo:    user.Setinfo,
		SetinfoAll: user.SetinfoAll,
	}

	addr := &models.Address{
		Country:  form.UserInfo.Country,
		Province: form.UserInfo.Province,
		City:     form.UserInfo.City,
		Area:     form.UserInfo.Area,
		Desc:     form.UserInfo.LocDesc,
	}
	if addr.String() != "" {
		setinfo.Address = addr
	}

	if len(setinfo.Phone) > 0 && setinfo.Phone != user.Phone {
		user.Phone = setinfo.Phone
		//setinfo.Setinfo = false
		if b, _ := user.Exists("phone"); b {
			writeResponse(request.RequestURI, resp, nil,
				errors.NewError(errors.UserExistError, "绑定失败，当前手机号已绑定"))
			return
		}
	}

	if len(setinfo.Nickname) > 0 && setinfo.Nickname != user.Nickname {
		u := &models.Account{}
		if find, _ := u.FindByNickname(setinfo.Nickname); find {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.UserExistError, "昵称已被占用"))
			return
		}
	}

	awards := Awards{}
	if !user.Setinfo || !user.SetinfoAll {
		ratio := calcInfo(user, setinfo)

		if ratio >= 80 && !user.Setinfo {
			setinfo.Setinfo = true
			awards.Wealth = 30 * models.Satoshi
			awards.Score = 30
		}
		if ratio == 100 && !user.SetinfoAll {
			setinfo.SetinfoAll = true
			awards.Wealth = 50 * models.Satoshi
			awards.Score = 50
		}
	}

	if err := user.SetInfo(setinfo); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	GiveAwards(user, awards, redis)

	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, nil)

	user.UpdateAction(ActInfo, nowDate())
	//redis.SetOnlineUser(form.Token, user, false)
}

type setProfileForm struct {
	ImageId string `json:"image_id" binding:"required"`
	parameter
}

func setProfileHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {
	form := p.(setProfileForm)
	err := user.ChangeProfile(form.ImageId)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}

type setPhotosForm struct {
	Pics []string `json:"pic_ids"`
	parameter
}

func setPhotosHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger,
	user *models.Account, p Parameter) {
	form := p.(setPhotosForm)

	awards := Awards{}
	if !user.PhotoSet {
		awards.Wealth = 10 * int64(len(form.Pics)) * models.Satoshi
		awards.Score = 10 * int64(len(form.Pics))
	}
	err := user.AddPhotos(form.Pics)
	GiveAwards(user, awards, redis)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": awards}, err)
}

type delPhotoForm struct {
	Photo string `json:"pic_id"`
	parameter
}

func delPhotoHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {
	err := user.DelPhoto(p.(delPhotoForm).Photo)
	writeResponse(request.RequestURI, resp, nil, err)
}

type loginAwardsForm struct {
	parameter
}

func loginAwards(days, level int) Awards {
	awards := Awards{}

	// calc wealth
	scale := 1.0
	factor := 0.5
	r := ran.Intn(level) + 1
	if days > 7 {
		scale = 1.5
		factor = 1.0
		r = ran.Intn(level*2) + 1
	}
	awards.Wealth = int64(float64(days)*scale+float64(level)*factor+float64(r)) * models.Satoshi

	// calc score
	scale = 5.0
	factor = 1.0
	r = ran.Intn(level) + 1
	if days > 7 {
		scale = 10.0
		factor = 1.5
		r = ran.Intn(level*5) + 1
	}
	awards.Score = int64(float64(days)*scale + float64(level)*factor + float64(r))

	return awards
}

func loginAwardsHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

	a := user.LoginAwards
	if (user.LoginDays-1)%7 == 0 || len(a) == 0 {
		a = []int64{}
		startDay := ((user.LoginDays - 1) / 7) * 7
		level := user.Level()
		score := user.Props.Score

		for i := 0; i < 7; i++ {
			awards := loginAwards(startDay+i+1, int(level))
			a = append(a, awards.Wealth, awards.Score)
			score = score + awards.Score
			level = models.Score2Level(score)
		}

		user.SetLoginAwards(a)
	}

	index := (user.LoginDays - 1) % 7
	awards := Awards{Wealth: a[index*2], Score: a[index*2+1]}
	awards.Level = models.Score2Level(user.Props.Score+awards.Score) - user.Level()
	GiveAwards(user, awards, redis)

	loginAwards := []int64{}
	for i := 0; i < 7; i++ {
		loginAwards = append(loginAwards, a[i*2])
	}
	respData := map[string]interface{}{
		"continuous_logined_days": user.LoginDays,
		"login_reward_list":       loginAwards,
	}
	writeResponse(request.RequestURI, resp, respData, nil)
}

type scoreDiffForm struct {
	Uid string `form:"userid" binding:"required"`
	parameter
}

func scoreDiffHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {

	other := &models.Account{}
	other.FindByUserid(p.(scoreDiffForm).Uid)
	respData := map[string]int64{
		"physique_times":   other.Props.Physical - user.Props.Physical,
		"literature_times": other.Props.Literal - user.Props.Literal,
		"magic_times":      other.Props.Mental - user.Props.Mental,
	}

	writeResponse(request.RequestURI, resp, respData, nil)
}

type getPropsForm struct {
	Uid string `form:"userid" binding:"required"`
}

func getPropsHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, form getPropsForm) {
	user := &models.Account{}
	user.FindByUserid(form.Uid)
	//user.Props.Wealth = redis.GetCoins(form.Uid)
	balance, _ := getBalance(user.Wallet.Addrs)
	var wealth int64
	if balance != nil {
		for _, b := range balance.Addrs {
			wealth += (b.Confirmed + b.Unconfirmed)
		}
	}
	user.Props.Wealth = wealth
	user.Props.Level = user.Level()

	writeResponse(r.RequestURI, w, user.Props, nil)
}

type setEquipForm struct {
	Equips models.Equip `json:"user_equipInfo"`
	parameter
}

func setEquipHandler(request *http.Request, resp http.ResponseWriter,
	user *models.Account, p Parameter) {

	form := p.(setEquipForm)
	err := user.SetEquip(form.Equips)
	writeResponse(request.RequestURI, resp, map[string]interface{}{"ExpEffect": Awards{}}, err)
}

type searchForm struct {
	Nearby   int    `form:"search_nearby"`
	Nickname string `form:"search_nickname"`
	models.Paging
	parameter
}

func searchHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account, form searchForm) {
	users := []models.Account{}
	var err error

	if form.Nearby > 0 {
		form.Paging.Count = 50
		users, err = user.SearchNear(&form.Paging, 50000)
	} else {
		users, err = models.SearchUsers(form.Nickname, &form.Paging)
	}

	var list []*leaderboardResp
	for i, _ := range users {
		if users[i].Id == user.Id {
			continue
		}
		lb := &leaderboardResp{
			Userid:   users[i].Id,
			Score:    users[i].Props.Score,
			Level:    users[i].Level(),
			Profile:  users[i].Profile,
			Nickname: users[i].Nickname,
			Gender:   users[i].Gender,
			LastLog:  users[i].LastLogin.Unix(),
			Birth:    users[i].Birth,
			Location: users[i].Loc,
			Phone:    users[i].Phone,
			Actor:    users[i].Actor,
		}
		list = append(list, lb)
	}

	respData := map[string]interface{}{
		"members_list":  list,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(r.RequestURI, w, respData, err)
}

/*
type importFriendsForm struct {
	Type     string `json:"account_type"`
	Uid      string `json:"userid" binding:"required"`
	AppKey   string `json:"appkey"`
	AppToken string `json:"verfiycode" binding:"required"`
	parameter
}


func (this importFriendsForm) getTokenId() string {
	return this.Token
}

func importFriendsHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(importFriendsForm)

	switch form.Type {
	case "weibo":
		log.Println("get weibo friends")
		friends, err := GetWeiboFriends(form.AppKey, form.Uid, form.AppToken)
		if err != nil {
			writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError, err.Error()))
			return
		}
		log.Println("import weibo friends", len(friends))
		for _, friend := range friends {
			//log.Println(friend.Id, friend.ScreenName)
			u := &models.Account{
				Id:       strconv.FormatInt(int64(friend.Id), 10),
				Nickname: friend.ScreenName,
				Profile:  friend.Avatar,
				Role:     "weibo",
				Gender:   friend.Gender,
				Addr:     &models.Address{Desc: friend.Location},
			}
			if find, _ := u.Exists(""); find {
				if u.RegTime.Unix() > 0 { // registered users only
					redis.ImportFriend(user.Id, u.Id)
				}
				continue
			}
			if err := u.Save(); err == nil {
				redis.SetWBImport(user.Id, u.Id)
			}
		}
	default:
	}
	writeResponse(r.RequestURI, w, map[string]interface{}{"ExpEffect": Awards{}}, nil)
}
*/
type userArticlesForm struct {
	Id   string `form:"userid" binding:"required"`
	Type string `form:"article_type"`
	models.Paging
	parameter
}

func userArticlesHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(userArticlesForm)

	u := &models.Account{}
	u.FindByUserid(form.Id)
	author := convertUser(u, redis)
	_, articles, err := u.Articles(form.Type, &form.Paging)

	jsonStructs := make([]*articleJsonStruct, len(articles))
	for i, _ := range articles {
		jsonStructs[i] = convertArticle(user, &articles[i], author)
	}

	respData := make(map[string]interface{})
	if len(articles) > 0 {
		respData["page_frist_id"] = form.Paging.First
		respData["page_last_id"] = form.Paging.Last
		//respData["page_item_count"] = total
	}
	respData["articles_without_content"] = jsonStructs

	writeResponse(request.RequestURI, resp, respData, err)
}

type importContactsForm struct {
	Contacts []string `json:"contacts"`
	parameter
}

func importContactsHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(importContactsForm)
	var result []*userJsonStruct

	users, _ := models.FindUsersByPhones(form.Contacts)
	ids := redis.Friends(models.RelFollowing, user.Id)
	for j, _ := range users {
		i := 0
		for ; i < len(ids); i++ {
			if users[j].Id == ids[i] || users[j].Id == user.Id {
				break
			}
		}
		if i >= len(ids) {
			result = append(result, convertUser(&users[j], redis))
		}
	}

	writeResponse(r.RequestURI, w, map[string]interface{}{"users": result}, nil)
}

type resetPasswdForm struct {
	Phone    string `json:"phone_number"`
	Password string `json:"password"`
	parameter
}

func resetPasswdHandler(r *http.Request, w http.ResponseWriter,
	form resetPasswdForm) {
	user := &models.Account{}
	if b, err := user.FindByPhone(form.Phone); !b {
		e := errors.NewError(errors.NotExistsError, "未绑定手机,无法重置密码")
		if err != nil {
			e = errors.NewError(errors.DbError)
		}
		writeResponse(r.RequestURI, w, nil, e)
		return
	}

	err := user.SetPassword(Md5(form.Password))
	writeResponse(r.RequestURI, w, nil, err)
}

type setPasswordForm struct {
	Old string `json:"password"`
	New string `json:"passwordNew"`
	parameter
}

func setPasswordHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account, p Parameter) {
	form := p.(setPasswordForm)

	if user.Password != Md5(form.Old) {
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.PasswordError, "原密码错误"))
		return
	}

	err := user.SetPassword(Md5(form.New))
	writeResponse(r.RequestURI, w, nil, err)
}

type pkShareForm struct {
	parameter
}

func pkShareHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {
	awards := Awards{
		Wealth: 30 * models.Satoshi,
		Score:  30,
	}
	if err := GiveAwards(user, awards, redis); err != nil {
		log.Println(err)
		writeResponse(r.RequestURI, w, nil, errors.NewError(errors.DbError))
		return
	}
	writeResponse(r.RequestURI, w,
		map[string]interface{}{"ExpEffect": awards}, nil)
}

type nicknameForm struct {
	Nickname string `form:"nikename"`
}

func checkNicknameHandler(r *http.Request, w http.ResponseWriter, form nicknameForm) {
	user := &models.Account{Nickname: form.Nickname}
	find, err := user.Exists("nickname")

	writeResponse(r.RequestURI, w, map[string]bool{"is_used": find}, err)
}

type checkPhoneForm struct {
	Phone string `form:"phone"`
}

func checkPhoneHandler(r *http.Request, w http.ResponseWriter, form checkPhoneForm) {
	user := &models.Account{Phone: form.Phone}
	find, err := user.Exists("phone")

	writeResponse(r.RequestURI, w, map[string]bool{"is_used": find}, err)
}

type gameResultForm struct {
	Type  string `form:"game_type"`
	Score int    `form:"game_score"`
	Count int    `form:"page_item_count"`
	parameter
}

func gameResultHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(gameResultForm)
	var respData struct {
		Total     []*leaderboardResp `json:"total_list"`
		Friends   []*leaderboardResp `json:"friends_list"`
		Score     int                `json:"total_score"`
		Percent   int                `json:"percent"`
		PerFriend int                `json:"percentFri"`
	}

	if form.Count == 0 {
		form.Count = 3
	}

	gt := gameType(form.Type)
	if scores := redis.UserGameScores(gt, user.Id); len(scores) == 1 {
		respData.Score = int(scores[0])
	}
	redis.SetGameScore(gt, user.Id, form.Score) // current score

	kvs := redis.GameScores(gt, 0, form.Count)
	var ids []string
	for _, kv := range kvs {
		ids = append(ids, kv.K)
	}
	ranks := redis.UserGameRanks(gt, user.Id)

	redis.SetGameScore(gt, user.Id, respData.Score) // recover max score

	n := redis.GameUserCount(gt) - 1
	//log.Println(ranks, n)
	if len(ranks) == 1 && n > 0 && form.Score > 0 {
		respData.Percent = int(float64(n-ranks[0]) / float64(n) * 100.0)
	}

	//log.Println(ids)
	users, _ := models.FindUsersByIds(1, ids...)
	index := 0
	for _, kv := range kvs {
		for i, _ := range users {
			if users[i].Id == kv.K {
				respData.Total = append(respData.Total, &leaderboardResp{
					Userid:   users[i].Id,
					Score:    kv.V,
					Rank:     index + 1,
					Level:    users[i].Level(),
					Profile:  users[i].Profile,
					Nickname: users[i].Nickname,
					Gender:   users[i].Gender,
					LastLog:  users[i].LastGameTime(gt).Unix(),
					Birth:    users[i].Birth,
					Location: users[i].Loc,
					Phone:    users[i].Phone,
					Actor:    users[i].Actor,
				})
				index++

				break
			}
		}
	}

	ids = redis.Friends(models.RelFriend, user.Id)
	if len(ids) > 0 {
		total := len(ids)
		ids = append(ids, user.Id)
		scores := redis.UserGameScores(gt, ids...)

		if len(scores) != len(ids) {
			scores = make([]int64, total)
		}
		kvs = make([]models.KV, total+1)
		for i, _ := range ids {
			kvs[i].K = ids[i]
			kvs[i].V = scores[i]
			if ids[i] == user.Id {
				kvs[i].V = int64(form.Score)
			}
		}

		sort.Sort(sort.Reverse(models.KVSlice(kvs)))
		lb := kvs
		if len(kvs) > 3 {
			kvs = kvs[0:3]
		}
		ids = []string{}
		for _, kv := range kvs {
			ids = append(ids, kv.K)
		}
		users, _ = models.FindUsersByIds(1, ids...)
		index := 0
		rank := 0

		for i, _ := range lb {
			if lb[i].K == user.Id {
				rank = i
				break
			}
		}

		for _, kv := range kvs {
			for i, _ := range users {
				if users[i].Id == kv.K {
					respData.Friends = append(respData.Friends, &leaderboardResp{
						Userid:   users[i].Id,
						Score:    kv.V,
						Rank:     index + 1,
						Level:    users[i].Level(),
						Profile:  users[i].Profile,
						Nickname: users[i].Nickname,
						Gender:   users[i].Gender,
						LastLog:  users[i].LastGameTime(gt).Unix(),
						Birth:    users[i].Birth,
						Location: users[i].Loc,
						Phone:    users[i].Phone,
						Actor:    users[i].Actor,
					})
					index++

					break
				}
			}
		}

		if total > 0 {
			respData.PerFriend = int(float64(total-rank) / float64(total) * 100.0)
		}
	}

	writeResponse(r.RequestURI, w, respData, nil)
}

type purchaseForm struct {
	Coins int64 `json:"coin_value"`
	Value int   `json:"value"`
	Time  int64 `json:"time"`
	parameter
}

func purchaseHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(purchaseForm)

	awards := Awards{Wealth: form.Coins}
	if err := GiveAwards(user, awards, redis); err != nil {
		writeResponse(r.RequestURI, w, nil, err)
		return
	}

	tx := &models.Tx{
		Uid:   user.Id,
		Coins: form.Coins,
		Value: form.Value,
		Time:  time.Unix(form.Time, 0),
	}

	err := tx.Save()

	respData := map[string]interface{}{"ExpEffect": awards}
	writeResponse(r.RequestURI, w, respData, err)
}

type purchaseListForm struct {
	models.Paging
	parameter
}

type purchaseStruct struct {
	Coins int64 `json:"coin_value"`
	Value int   `json:"value"`
	Time  int64 `json:"time"`
}

func purchaseListHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {
	form := p.(purchaseListForm)

	_, txs, _ := user.Txs(&form.Paging)

	list := []*purchaseStruct{}
	for _, tx := range txs {
		list = append(list, &purchaseStruct{
			Coins: tx.Coins,
			Value: tx.Value,
			Time:  tx.Time.Unix(),
		})
	}

	respData := map[string]interface{}{
		"payCoinList":   list,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(r.RequestURI, w, respData, nil)
}

func testHandler(r *http.Request, w http.ResponseWriter) {
	writeResponse(r.RequestURI, w, map[string]interface{}{"is_preSportForm": true}, nil)
}

type userlbForm struct {
	Type string `form:"type"`
	models.Paging
}

func userLeaderBoardHandler(r *http.Request, w http.ResponseWriter, form userlbForm) {
	users, _ := models.UserLeaderBoard(form.Type, &form.Paging)

	start, _ := strconv.Atoi(form.Paging.First)
	lb := make([]*leaderboardResp, len(users))
	for i, _ := range users {
		lb[i] = &leaderboardResp{
			Userid: users[i].Id,
			//Score:    users[i].Props,
			Rank:     start + i + 1,
			Level:    users[i].Level(),
			Profile:  users[i].Profile,
			Nickname: users[i].Nickname,
			Gender:   users[i].Gender,
			LastLog:  users[i].LastLogin.Unix(),
			Birth:    users[i].Birth,
			Location: users[i].Loc,
			Phone:    users[i].Phone,
			Actor:    users[i].Actor,
		}
		switch form.Type {
		case "physique":
			lb[i].Score = users[i].Props.Physical
		case "literature":
			lb[i].Score = users[i].Props.Literal
		case "magic":
			lb[i].Score = users[i].Props.Mental
		}
	}

	respData := map[string]interface{}{
		"members_list":  lb,
		"page_frist_id": form.Paging.First,
		"page_last_id":  form.Paging.Last,
	}
	writeResponse(r.RequestURI, w, respData, nil)
}

type authRequestForm struct {
	Type   string   `json:"auth_type"`
	Images []string `json:"auth_images"`
	Desc   string   `json:"auth_desc"`
	parameter
}

func userAuthRequestHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(authRequestForm)
	err := user.SetAuthInfo(form.Type, form.Images, form.Desc)
	redis.AddAuthCoach(user.Id)

	writeResponse(r.RequestURI, w, nil, err)
}

type authStatusForm struct {
	Userid string `form:"userid"`
}

func userAuthStatusHandler(r *http.Request, w http.ResponseWriter,
	form authStatusForm) {

	user := &models.Account{}
	user.FindByUserid(form.Userid)

	idcard := models.AuthUnverified
	cert := models.AuthUnverified
	record := models.AuthUnverified

	if user.Auth != nil {
		if user.Auth.IdCardTmp != nil {
			idcard = user.Auth.IdCardTmp.Status
		}
		if user.Auth.IdCard != nil {
			idcard = user.Auth.IdCard.Status
		}

		if user.Auth.CertTmp != nil {
			cert = user.Auth.CertTmp.Status
		}
		if user.Auth.Cert != nil {
			cert = user.Auth.Cert.Status
		}

		if user.Auth.RecordTmp != nil {
			record = user.Auth.RecordTmp.Status
		}
		if user.Auth.Record != nil {
			record = user.Auth.Record.Status
		}
	}

	respData := map[string]string{
		"id_card": idcard,
		"cert":    cert,
		"record":  record,
	}
	writeResponse(r.RequestURI, w, respData, nil)
}

type authInfoForm struct {
	Type string `form:"auth_type"`
	parameter
}

func userAuthInfoHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account, p Parameter) {

	info := &models.AuthInfo{
		Status: models.AuthUnverified,
	}

	if user.Auth == nil {
		writeResponse(r.RequestURI, w, info, nil)
		return
	}

	form := p.(authInfoForm)
	switch form.Type {
	case models.AuthIdCard:
		if user.Auth.IdCardTmp != nil {
			info = user.Auth.IdCardTmp
		}
		if user.Auth.IdCard != nil {
			info = user.Auth.IdCard
		}
	case models.AuthCert:
		if user.Auth.CertTmp != nil {
			info = user.Auth.CertTmp
		}
		if user.Auth.Cert != nil {
			info = user.Auth.Cert
		}
	case models.AuthRecord:
		if user.Auth.RecordTmp != nil {
			info = user.Auth.RecordTmp
		}
		if user.Auth.Record != nil {
			info = user.Auth.Record
		}
	}

	writeResponse(r.RequestURI, w, info, nil)
}

type sendHeartForm struct {
	Record string `json:"record_id"`
	parameter
}

func sendHeartHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(sendHeartForm)

	receivers := redis.HeartReceivers(user.Id)

	if len(receivers) == 0 {
		writeResponse(r.RequestURI, w, nil, nil)
	}

	var target string
	for _, recv := range receivers {
		if recv != user.Id {
			target = recv
			break
		}
	}
	// ws push
	event := &models.Event{
		Type: models.EventSystem,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventSendHeart,
			Id:   user.Id,
			From: user.Id,
			To:   target,
			Body: []models.MsgBody{
				{Type: "record_id", Content: form.Record},
				{Type: "userid", Content: user.Id},
			},
		},
	}
	if len(target) > 0 {
		event.Save()
		redis.PubMsg(models.EventSystem, target, event.Bytes())
		redis.LogHeartSend(user.Id)
		redis.SetHeartRecv(target, true)
	}

	writeResponse(r.RequestURI, w, nil, nil)
}

type recvHeartForm struct {
	Sender string `json:"sender"`
	Accept bool   `json:"accept"`
	parameter
}

func recvHeartHandler(r *http.Request, w http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account, p Parameter) {

	form := p.(recvHeartForm)

	// ws push
	event := &models.Event{
		Type: models.EventSystem,
		Time: time.Now().Unix(),
		Data: models.EventData{
			Type: models.EventRecvHeart,
			Id:   user.Id,
			From: user.Id,
			To:   form.Sender,
			Body: []models.MsgBody{
				{Type: "userid", Content: user.Id},
			},
		},
	}

	awards := Awards{}
	if form.Accept {
		redis.SetRelationship(user.Id, []string{form.Sender}, models.RelFriend, true)
		event.Save()
		redis.PubMsg(models.EventSystem, form.Sender, event.Bytes())
		awards.Wealth = 1 * models.Satoshi
		GiveAwards(user, awards, redis)
	}
	redis.SetHeartRecv(user.Id, false)
	writeResponse(r.RequestURI, w, map[string]interface{}{"ExpEffect": awards}, nil)
}

type userRanksForm struct {
	parameter
}

func userRanksHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account) {
	var ranks struct {
		Userid         string `json:"userid"`
		ScoreRank      int    `json:"score_rank"`
		PhysiqueRank   int    `json:"physique_rank"`
		LiteratureRank int    `json:"literature_rank"`
		MagicRank      int    `json:"magic_rank"`
	}

	//log.Println(user.Props)
	ranks.Userid = user.Id
	ranks.ScoreRank, _ = user.PropIndex("score", user.Props.Score)
	ranks.PhysiqueRank, _ = user.PropIndex("physique", user.Props.Physical)
	ranks.LiteratureRank, _ = user.PropIndex("literature", user.Props.Literal)
	ranks.MagicRank, _ = user.PropIndex("magic", user.Props.Mental)

	writeResponse(r.RequestURI, w, ranks, nil)
}
