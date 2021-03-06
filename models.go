/*
和MongoDB对应的struct
*/

package main

import (
	"html/template"
	"labix.org/v2/mgo/bson"
	"time"
)

// 用户
type User struct {
	Id_          bson.ObjectId `bson:"_id"`
	Username     string
	Password     string
	Email        string
	Website      string
	Location     string
	Tagline      string
	Bio          string
	Twitter      string
	Weibo        string
	JoinedAt     time.Time
	Follow       []string
	Fans         []string
	IsSuperuser  bool
	IsActive     bool
	ValidateCode string
	ResetCode    string
	Index        int
}

// 用户发表的最近10个主题
func (u *User) LatestTopics() *[]Topic {
	c := db.C("topics")
	var topics []Topic

	c.Find(bson.M{"userid": u.Id_}).Sort("-createdat").Limit(10).All(&topics)

	return &topics
}

// 用户的最近10个回复
func (u *User) LatestReplies() *[]Reply {
	c := db.C("replies")
	var replies []Reply

	c.Find(bson.M{"userid": u.Id_}).Sort("-createdat").Limit(10).All(&replies)

	return &replies
}

// 是否被某人关注
func (u *User) IsFollowedBy(who string) bool {
	for _, username := range u.Fans {
		if username == who {
			return true
		}
	}

	return false
}

// 是否关注某人
func (u *User) IsFans(who string) bool {
	for _, username := range u.Follow {
		if username == who {
			return true
		}
	}

	return false
}

// 节点
type Node struct {
	Id_         bson.ObjectId `bson:"_id"`
	Id          string
	Name        string
	Description string
	TopicCount  int
}

// 回复
type Reply struct {
	Id_       bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId
	TopicId   bson.ObjectId
	Markdown  string
	Html      template.HTML
	CreatedAt time.Time
	topic     Topic
}

// 该回复所属于的用户
func (r *Reply) User() *User {
	c := db.C("users")
	user := User{}
	c.Find(bson.M{"_id": r.UserId}).One(&user)

	return &user
}

// 该回复的主题
func (r *Reply) Topic() *Topic {
	if r.topic.Title == "" {
		c := db.C("topics")
		r.topic = Topic{}
		c.Find(bson.M{"_id": r.TopicId}).One(&r.topic)
	}

	return &r.topic
}

// 主题
type Topic struct {
	Id_             bson.ObjectId `bson:"_id"`
	NodeId          bson.ObjectId
	UserId          bson.ObjectId
	Title           string
	Markdown        string
	Html            template.HTML
	CreatedAt       time.Time
	ReplyCount      int32
	LatestReplyId   string
	LatestRepliedAt time.Time
	Hits            int32
}

// 主题所属节点
func (t *Topic) Node() *Node {
	c := db.C("nodes")
	node := Node{}
	c.Find(bson.M{"_id": t.NodeId}).One(&node)

	return &node
}

// 主题的最近的一个回复
func (t *Topic) LatestReply() *Reply {
	if t.LatestReplyId == "" {
		return nil
	}

	c := db.C("replies")
	reply := Reply{}

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(t.LatestReplyId)}).One(&reply)

	if err != nil {
		return nil
	}

	return &reply
}

// 主题的作者
func (t *Topic) User() *User {
	c := db.C("users")
	user := User{}
	c.Find(bson.M{"_id": t.UserId}).One(&user)

	return &user
}

// 主题下的所有回复
func (t *Topic) Replies() *[]Reply {
	c := db.C("replies")
	var replies []Reply

	c.Find(bson.M{"topicid": t.Id_}).All(&replies)

	return &replies
}

// 状态,MongoDB中只存储一个状态
type Status struct {
	Id_        bson.ObjectId `bson:"_id"`
	UserCount  int
	TopicCount int
	ReplyCount int
	UserIndex  int
}

// 站点分类
type SiteCategory struct {
	Id_  bson.ObjectId `bson:"_id"`
	Name string
}

// 分类下的所有站点
func (sc *SiteCategory) Sites() *[]Site {
	var sites []Site
	c := db.C("sites")
	c.Find(bson.M{"categoryid": sc.Id_}).All(&sites)

	return &sites
}

// 站点
type Site struct {
	Id_         bson.ObjectId `bson:"_id"`
	Name        string
	Url         string
	Description string
	CategoryId  bson.ObjectId
	UserId      bson.ObjectId
}
