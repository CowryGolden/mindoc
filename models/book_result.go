package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"strings"
)

type BookResult struct {
	BookId int              `json:"book_id"`
	BookName string         `json:"book_name"`
	Identify string         `json:"identify"`
	OrderIndex int          `json:"order_index"`
	Description string      `json:"description"`
	PrivatelyOwned int      `json:"privately_owned"`
	PrivateToken string     `json:"private_token"`
	DocCount int            `json:"doc_count"`
	CommentStatus string    `json:"comment_status"`
	CommentCount int        `json:"comment_count"`
	CreateTime time.Time    `json:"create_time"`
	CreateName string 	`json:"create_name"`
	ModifyTime time.Time	`json:"modify_time"`
	Cover string            `json:"cover"`
	Label string		`json:"label"`
	MemberId int            `json:"member_id"`
	RoleId int        	`json:"role_id"`
	RoleName string 	`json:"role_name"`

	LastModifyText string 	`json:"last_modify_text"`
}

func NewBookResult() *BookResult {
	return &BookResult{}
}

func (m *BookResult) FindByIdentify(identify string,member_id int) (*BookResult,error) {
	o := orm.NewOrm()

	book := NewBook()

	err := o.QueryTable(book.TableNameWithPrefix()).Filter("identify", identify).One(book)

	if err != nil {
		return m, err
	}

	relationship := NewRelationship()

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id",book.BookId).Filter("member_id",member_id).One(relationship)

	if err != nil {
		return m,err
	}
	var relationship2 Relationship

	err = o.QueryTable(relationship.TableNameWithPrefix()).Filter("book_id",book.BookId).Filter("role_id",0).One(&relationship2)

	if err != nil {
		return m,ErrPermissionDenied
	}

	member := NewMember()

	err = member.Find(relationship2.MemberId)
	if err != nil {
		return m,err
	}

	m.BookId = book.BookId
	m.BookName = book.BookName
	m.Identify = book.Identify
	m.OrderIndex = book.OrderIndex
	m.Description = strings.Replace(book.Description,"\r\n","<br/>",-1)
	m.PrivatelyOwned = book.PrivatelyOwned
	m.PrivateToken = book.PrivateToken
	m.DocCount = book.DocCount
	m.CommentStatus = book.CommentStatus
	m.CommentCount = book.CommentCount
	m.CreateTime = book.CreateTime
	m.CreateName = member.Account
	m.ModifyTime = book.ModifyTime
	m.Cover = book.Cover
	m.Label = book.Label

	m.MemberId = relationship.MemberId
	m.RoleId = relationship.RoleId

	if m.RoleId == 0{
		m.RoleName = "创始人"
	}else if m.RoleId == 1 {
		m.RoleName = "管理员"
	}else if m.RoleId == 2 {
		m.RoleName = "编辑者"
	}else if m.RoleId == 2 {
		m.RoleName = "观察者"
	}


	doc := NewDocument()

	err = o.QueryTable(doc.TableNameWithPrefix()).Filter("book_id",book.BookId).OrderBy("modify_time").One(doc)

	if err == nil {
		member2 := NewMember()
		member2.Find(doc.ModifyAt)

		m.LastModifyText = member2.Account + " 于 " + doc.ModifyTime.Format("2006-01-02 15:04:05")
	}

	return m,nil

}