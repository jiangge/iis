package view

import (
	"bytes"
	"encoding/base64"
	"log"
	"strconv"

	"github.com/coyove/iis/cmd/ch/config"
	"github.com/coyove/iis/cmd/ch/ident"
	"github.com/coyove/iis/cmd/ch/manager"
	"github.com/coyove/iis/cmd/ch/mv"
	"github.com/gin-gonic/gin"
)

var m *manager.Manager

type ArticlesTimelineView struct {
	Articles   []ArticleView
	Next       string
	Prev       string
	SearchTerm string
	IsAdmin    bool
	IsTagInbox bool
	Index      bool
	UUID       string
	User       *mv.User
}

type ArticleRepliesView struct {
	Articles        []ArticleView
	ParentArticle   ArticleView
	CanDeleteParent bool
	CurPage         int
	TotalPages      int
	Pages           []int
	ReplyView       struct {
		UUID      string
		Challenge string
		ShowReply bool
		RAuthor   string
		RContent  string
		EError    string
	}
}

func SetManager(mgr *manager.Manager) {
	m = mgr
}

func Index(g *gin.Context) {
	//	var pl = ArticlesTimelineView{
	//		SearchTerm: g.Param("tag"),
	//	}
	//	var opt uint64
	//	var idtag = ident.IDTagGeneral
	//
	//	if strings.HasPrefix(pl.SearchTerm, "@") {
	//		pl.SearchTerm = pl.SearchTerm[1:]
	//		pl.User, _ = m.GetUser(pl.SearchTerm)
	//
	//		opt |= _abstract
	//		idtag = ident.IDTagAuthor
	//	} else if strings.HasPrefix(pl.SearchTerm, "inbox:") {
	//		pl.SearchTerm = pl.SearchTerm[6:]
	//		pl.IsTagInbox = true
	//
	//		if u, ok := g.Get("user"); !ok || (u.(*mv.User).ID != pl.SearchTerm && !u.(*mv.User).IsMod()) {
	//			g.Redirect(302, "/cat")
	//			return
	//		} else if u.(*mv.User).ID == pl.SearchTerm {
	//			m.UpdateUser(pl.SearchTerm, func(u *mv.User) error {
	//				u.Unread = 0
	//				return nil
	//			})
	//		}
	//
	//		opt |= _abstract
	//		idtag = ident.IDTagInbox
	//	} else if strings.HasPrefix(pl.SearchTerm, "search:") {
	//		pl.SearchTerm = mv.SoftTrunc(pl.SearchTerm[7:], 50)
	//		a := []*mv.Article{}
	//		for _, id := range mv.Search(pl.SearchTerm) {
	//			x, _ := m.Get(id)
	//			if x != nil {
	//				a = append(a, x)
	//			}
	//		}
	//		fromMultiple(&pl.Articles, a, _abstract)
	//		pl.Index = true
	//		pl.SearchTerm = "\"" + pl.SearchTerm + "\""
	//		g.HTML(200, "index.html", pl)
	//		return
	//	}
	//
	//	u, ok := g.Get("user")
	//	pl.IsAdmin = ok && u.(*mv.User).IsAdmin()
	//	if pl.IsAdmin {
	//		pl.UUID, _ = ident.MakeToken(g)
	//	}
	//
	//	cursor := ident.ParseID(g.Query("n")).String()
	//	a, next, err := m.Walk(idtag, pl.SearchTerm, cursor, int(config.Cfg.PostsPerPage))
	//	if err != nil {
	//		Error(500, "INTERNAL: "+err.Error(), g)
	//		return
	//	}
	//
	//	fromMultiple(&pl.Articles, a, opt)
	//
	//	pl.Next = next
	//	pl.Index = cursor == ""
	//
	//	if u, _ := url.Parse(g.Request.Referer()); u != nil {
	//		pl.Prev = u.Query().Get("n")
	//		if pl.Prev <= pl.Next || pl.Index {
	//			// If we are at the front page, or the prev page is smaller than the next page
	//			// then we consider the prev page invalid
	//			pl.Prev = ""
	//		}
	//	}
	//
	//	g.HTML(200, "index.html", pl)
}

func Timeline(g *gin.Context) {
	var pl ArticlesTimelineView
	var u *mv.User
	var single bool

	if uid := g.Param("user"); uid != "" {
		u, _ = m.GetUser(uid)
		if u == nil {
			Error(404, "USER NOT FOUND", g)
			return
		}
		single = true
	} else {
		u2, _ := g.Get("user")
		u, _ = u2.(*mv.User)
		if u == nil {
			g.Redirect(302, "/user")
			return
		}
	}

	cursors := []ident.ID{}
	pendingFCursor := ""

	if single {
		cursors = append(cursors, ident.NewID(ident.IDTagAuthor).SetTag(u.ID))
	} else if g.Request.Method == "POST" {
		for cbuf, _ := base64.StdEncoding.DecodeString(g.PostForm("cursors")); len(cbuf) > 0; {
			if cbuf[0] == 0 {
				list, next := m.GetFollowingList(u, string(cbuf[1:]), 1e6)
				for _, id := range list {
					cursors = append(cursors, ident.NewID(ident.IDTagAuthor).SetTag(id.ID))
				}
				pendingFCursor = next
				break
			}

			id := ident.UnmarshalID(cbuf)
			if !id.Valid() {
				break
			}
			cbuf = cbuf[id.Size():]
			cursors = append(cursors, id)
		}
	} else {
		list, next := m.GetFollowingList(u, "", 1e6)
		for _, id := range list {
			cursors = append(cursors, ident.NewID(ident.IDTagAuthor).SetTag(id.ID))
		}
		pendingFCursor = next
	}

	a, next := m.WalkMulti(int(config.Cfg.PostsPerPage), cursors...)
	fromMultiple(&pl.Articles, a, 0)

	nextbuf := bytes.Buffer{}
	nextbuftmp := [32]byte{}
	for _, c := range next {
		p := c.Marshal(nextbuftmp[:])
		nextbuf.Write(p)
	}

	if pendingFCursor != "" {
		nextbuf.WriteByte(0)
		nextbuf.WriteString(pendingFCursor)
	}

	pl.Next = base64.StdEncoding.EncodeToString(nextbuf.Bytes())
	g.HTML(200, "timeline.html", pl)
}

func Replies(g *gin.Context) {
	var pl ArticleRepliesView
	var pid = g.Param("parent")
	var opt = _richtime

	parent, err := m.Get(pid)
	if err != nil || parent.ID == "" {
		Error(404, "NOT FOUND", g)
		log.Println(pid, err)
		return
	}

	pl.ParentArticle.from(parent, opt)
	pl.ParentArticle.Index = 0
	pl.ParentArticle.SubIndex = ""

	if u, ok := g.Get("user"); ok {
		pl.CanDeleteParent = u.(*mv.User).ID == pl.ParentArticle.Author || u.(*mv.User).IsMod()
	}

	{
		pl.ReplyView.RContent = g.Query("content")
		pl.ReplyView.RAuthor = g.Query("author")
		pl.ReplyView.EError = g.Query("error")
		pl.ReplyView.ShowReply = g.Query("refresh") == "1" || pl.ReplyView.EError != ""
		if pl.ReplyView.RAuthor == "" {
			pl.ReplyView.RAuthor, _ = g.Cookie("id")
		}
		pl.ReplyView.UUID, pl.ReplyView.Challenge = ident.MakeToken(g)
	}

	pl.CurPage, _ = strconv.Atoi(g.Query("p"))
	pl.TotalPages = intdivceil(int(pl.ParentArticle.Replies), config.Cfg.PostsPerPage)

	//m.IncrCounter(g, parent.ID)

	if pl.CurPage == 0 {
		pl.CurPage = 1
	}
	pl.CurPage = intmin(pl.CurPage, pl.TotalPages)
	if pl.CurPage <= 0 {
		pl.CurPage = pl.TotalPages
	}
	if pl.CurPage <= 0 {
		pl.CurPage = 1
	}

	if pl.TotalPages > 0 {
		start := intmin(int(pl.ParentArticle.Replies), (pl.CurPage-1)*config.Cfg.PostsPerPage)
		end := intmin(int(pl.ParentArticle.Replies), pl.CurPage*config.Cfg.PostsPerPage)

		fromMultiple(&pl.Articles, m.GetReplies(parent.ID, start+1, end+1), opt)

		// Fill in at most 7 page numbers for display
		pl.Pages = make([]int, 0, 8)
		for i := pl.CurPage - 3; i <= pl.CurPage+3 && i <= pl.TotalPages; i++ {
			if i > 0 {
				pl.Pages = append(pl.Pages, i)
			}
		}
		for last := pl.Pages[len(pl.Pages)-1]; len(pl.Pages) < 7 && last+1 <= pl.TotalPages; last = pl.Pages[len(pl.Pages)-1] {
			pl.Pages = append(pl.Pages, last+1)
		}
		for first := pl.Pages[0]; len(pl.Pages) < 7 && first-1 > 0; first = pl.Pages[0] {
			pl.Pages = append([]int{first - 1}, pl.Pages...)
		}
	}

	g.HTML(200, "post.html", pl)
}