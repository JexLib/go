package permission

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo"

	"github.com/JexLib/golang/JexWeb/session"
	"github.com/JexLib/golang/cache"
	"github.com/JexLib/golang/cache/memory"
	"github.com/JexLib/golang/crypto/base58"
)

var (
	PermUsers *UserState
	_key      = "auth_user"
)

type UserState struct {
	store cache.Cache
}

type User struct {
	Account string
	Email   string
	// RememberMe bool //记住我
	IsAdmin   bool
	IsLogin   bool                   //是否在线
	Confirmed bool                   //是否确认为有效账号
	data      map[string]interface{} // 保存常用用户数据，如角色等
}

type cookieInfo struct {
	Account        string
	IsLogin        bool
	LastLoggedTime int64
}

func (ci *cookieInfo) encode() string {
	bytes, _ := json.Marshal(ci)
	return base58.Encode(bytes)
}

func cookieInfoDecode(str string) (*cookieInfo, error) {
	bytes := base58.Decode(str)
	ci := &cookieInfo{}
	err := json.Unmarshal(bytes, ci)
	return ci, err
}

func GetCookieInfo(c echo.Context) *cookieInfo {
	sess := session.Default(c)
	if cdata, ok := sess.Get(_key).(string); ok {
		if ci, err := cookieInfoDecode(cdata); err == nil {
			return ci
		}
	}
	return nil
}

func (ci *cookieInfo) Save(c echo.Context) {
	sess := session.Default(c)
	sess.Set(_key, ci.encode())
	if err := sess.Save(); err != nil {
		fmt.Println("session save:", err)
	}
}

func (ci *cookieInfo) Delete(c echo.Context) {
	sess := session.Default(c)
	sess.Delete(_key)
	sess.Save()
}

func (u *User) Save(c echo.Context) {
	PermUsers.store.Set(u.Account, u)
}

func UserInfo(c echo.Context) *User {
	cookINFO := GetCookieInfo(c)
	if cookINFO != nil {
		u := PermUsers.User(cookINFO.Account)
		if u != nil {
			u.IsLogin = cookINFO.IsLogin
		}
		return u
	}
	return &User{}
}

func newUsers(mcache cache.Cache) {
	if mcache == nil {
		mcache = memory.NewMemoryCache(time.Hour * 24 * 30)
	}
	PermUsers = &UserState{
		store: mcache,
	}
}

func (state *UserState) User(account string) *User {
	dat := PermUsers.store.Get(account)
	if dat == nil {
		u := User{
			Account: account,
		}
		PermUsers.store.Set(account, &u)
		return &u
	}
	return dat.(*User)
}

func (state *UserState) AddUser(account string) *User {
	return state.User(account)
}

func (state *UserState) RemoveUser(account string) {
	state.store.Delete(account)
}

func (state *UserState) IsAdmin(account string) bool {
	u := state.User(account)
	return u.IsAdmin
}

// SetAdminStatus marks a user as an administrator.
func (state *UserState) SetAdmin(account string) {
	u := state.User(account)
	u.IsAdmin = true
	state.store.Set(account, u)
}

func (state *UserState) RemoveAdmin(account string) {
	u := state.User(account)
	u.IsAdmin = false
	state.store.Set(account, u)
}

func (state *UserState) IsConfirmed(account string) bool {
	u := state.User(account)
	return u.Confirmed
}

func (state *UserState) SetConfirmed(account string, confirmed bool) {
	u := state.User(account)
	u.Confirmed = confirmed
	state.store.Set(account, u)
}

func (state *UserState) Email(account string) string {
	u := state.User(account)
	return u.Email
}

func (state *UserState) SetEmail(account string, email string) {
	u := state.User(account)
	u.Email = email
	state.store.Set(account, u)
}

func (state *UserState) Set(account, key string, v interface{}) {
	u := state.User(account)
	u.data[key] = v
	state.store.Set(account, u)
}

func (state *UserState) Get(account, key string) interface{} {
	u := state.User(account)
	return u.data[key]
}

func (state *UserState) AllUsernames() []string {
	return state.store.Keys()
}

// type users struct {
// 	store cache.Cache
// }

// type User struct {
// 	Account  string //用户唯一主键 primary
// 	Password string
// 	IsAdmin  bool
// 	IsLogin  bool
// 	Data     map[string]interface{} //可以在用户登录完成时，保存常用用户数据，如角色等
// }

// type sessionStruct struct {
// 	UUID       string
// 	Account    string
// 	RememberMe bool //记住我
// }

// func newUsers(mcache cache.Cache) *users {
// 	return &users{
// 		store: mcache,
// 	}
// }

// //新建用户信息
// func (us *users) getUser(sess session.Session, id ...string) *User {
// 	if len(id) == 0 {
// 		//新建用户id
// 		uuid := guid.New()
// 		base58.Encode([]byte(uuid))
// 		id = append(id, guid.New())
// 	}
// 	u := us.store.Get(id[0])
// 	if u == nil {
// 		u = User{
// 			Account: id[0],
// 			Data:    make(map[string]interface{}),
// 		}
// 		sess.Set(_key, id[0])
// 		sess.Save()
// 		us.store.Set(id[0], u)
// 	}
// 	uu := u.(User)
// 	return &uu
// }

// //获取当前用户信息
// func (us *users) User(c echo.Context) *User {
// 	sess := session.Default(c)

// 	if sess == nil {
// 		return nil
// 	}
// 	id := sess.Get(_key)
// 	if id != nil {
// 		return us.getUser(sess, id.(string))
// 	}
// 	return us.getUser(sess)
// }

// func (us *users) Count() int64 {
// 	return us.store.Count()
// }

// func (u *User) Save(c echo.Context) {
// 	sess := session.Default(c)
// 	if _config.Cache == nil {
// 		sess.Set(_key, u.Bytes())
// 		sess.Save()
// 	} else {
// 		_users.store.Set(u.Account, *u)
// 	}
// }

// func (u *User) Delete() {
// 	_users.store.Delete(u.Account)
// }

// func (u *User) Bytes() []byte {
// 	bytes, _ := json.Marshal(u)
// 	return bytes
// }

// func newUser(bytes []byte) (*User, error) {
// 	u := &User{}
// 	if err := json.Unmarshal(bytes, u); err == nil {
// 		return u, nil
// 	} else {
// 		return nil, err
// 	}
// }
