package app

import (
	"bytes"
	"github.com/owenzhou/ginrbac/utils"
	"html/template"
	"math"
	"net/url"
	"regexp"
	"strings"
)

func newPagination(total int64, uri string) *pagination {
	page := &pagination{
		Total:          total,
		CurrentPage:    1,
		PageSize:       20,
		PageTotal:      0,
		PrevPageLabel:  "上一页",
		NextPageLabel:  "下一页",
		FirstPageLabel: "首页",
		LastPageLabel:  "尾页",
		MaxButtonCount: 8,
		CurrentURI:     uri,
		PageLabel:      "page",
		PertyUrl:       false,
	}
	page.Apply()
	return page
}

type pagination struct {
	Total          int64  `comment:"数据总数"`
	CurrentPage    int    `comment:"当前页"`
	PageSize       int    `comment:"显示条数"`
	PageTotal      int    `comment:"总页数"`
	PrevPageLabel  string `comment:"上一页标签"`
	NextPageLabel  string `comment:"下一页标签"`
	FirstPageLabel string `comment:"首页标签"`
	LastPageLabel  string `comment:"末页标签"`
	MaxButtonCount int    `comment:"最大链接数量"`
	CurrentURI     string `comment:"当前链接"`
	PageLabel      string `comment:"url中分页的标签"`
	PertyUrl       bool   `comment:"是否seo友好"`
}

//从当前url匹配分页
func (p *pagination) FindCurrentPage() int {
	reg := regexp.MustCompile(`([/]{1}` + p.PageLabel + `/)([^/]*)`)
	if reg.MatchString(p.CurrentURI) {
		p := reg.FindStringSubmatch(p.CurrentURI)
		return utils.Str2Int(p[2])
	}
	return 1
}

//链接生成或追加参数
func (p *pagination) GenerateURL(key string, val interface{}) string {
	//如果参数值为空，直接返回当前链接
	if utils.Int2Str(val) == "" {
		return p.CurrentURI
	}

	var reg *regexp.Regexp
	if p.PertyUrl {
		reg = regexp.MustCompile(`([/]{1}` + key + `/)([^/]*)`)
	} else {
		reg = regexp.MustCompile(`([\?&]{1}` + key + `=)([^&]*)`)
	}

	if !reg.MatchString(p.CurrentURI) {
		if !p.PertyUrl && !strings.Contains(p.CurrentURI, "?") {
			return p.CurrentURI + "?" + key + "=" + utils.Int2Str(val)
		}
		if p.PertyUrl {
			return p.CurrentURI + "/" + key + "/" + utils.Int2Str(val)
		}

		query, _ := url.ParseQuery(p.CurrentURI)
		query.Add(p.PageLabel, utils.Int2Str(val))
		u, _ := url.PathUnescape(query.Encode())
		return u
	}
	return reg.ReplaceAllString(p.CurrentURI, "${1}"+utils.Int2Str(val))
}

//输出翻页链接
func (p *pagination) Links() template.HTML {
	var links []int
	var start int
	var tmpl *template.Template
	preNum := p.MaxButtonCount / 2
	//如果最大链接数超过最大页数，或者当前页数减 preNum 小于0，起始页为 1
	if p.MaxButtonCount > p.PageTotal || p.CurrentPage-preNum < 1 {
		start = 1
	} else {
		//获取起始页
		start = p.CurrentPage - preNum
		//如果当前页加上最大链接数大于总页数
		if p.CurrentPage+p.MaxButtonCount > p.PageTotal {
			start = p.PageTotal - p.MaxButtonCount + 1
		}
	}
	for i, j := start, 0; i <= p.PageTotal; i, j = i+1, j+1 {
		if j >= p.MaxButtonCount {
			break
		}
		links = append(links, i)
	}

	ul := `{{ $page := .page }}
	{{ if ne $page.PageTotal 0 }}
	<ul class="pagination">
		{{ if ne $page.FirstPageLabel "" }}
		<li class="page-item"><a class="page-link" href="{{$page.GenerateURL $page.PageLabel 1}}">{{ $page.FirstPageLabel }}</a></li>
		{{ end }}
		<li class="page-item"><a class="page-link" href="{{$page.GenerateURL $page.PageLabel $page.PrevPage}}">{{ $page.PrevPageLabel }}</a></li>
		{{ range $k, $v := .links }}
		<li class="page-item {{if eq $v $page.CurrentPage}} active {{end}}"><a class="page-link" href="{{ $page.GenerateURL $page.PageLabel $v }}">{{ $v }}</a></li>
		{{ end }}
		<li class="page-item"><a class="page-link" href="{{$page.GenerateURL $page.PageLabel $page.NextPage}}">{{ $page.NextPageLabel }}</a></li>
		{{ if ne $page.LastPageLabel "" }}
		<li class="page-item"><a class="page-link" href="{{$page.GenerateURL $page.PageLabel $page.PageTotal}}">{{ $page.LastPageLabel }}</a></li>
		{{ end }}
	</ul>
	{{ end }}`

	tmpl = template.Must(template.New("page").Parse(ul))
	args := map[string]interface{}{
		"page":  p,
		"links": links,
	}
	b := bytes.NewBuffer([]byte{})
	tmpl.ExecuteTemplate(b, "page", args)
	return template.HTML(b.String())
}

//下一页
func (p *pagination) NextPage() int {
	nextPage := p.CurrentPage + 1
	if nextPage > p.PageTotal {
		nextPage = p.PageTotal
	}
	return nextPage
}

//上一页
func (p *pagination) PrevPage() int {
	prePage := p.CurrentPage - 1
	if prePage < 1 {
		return 1
	}
	return prePage
}

//应用并更新总页数
func (p *pagination) Apply() {
	p.PageTotal = int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
}
