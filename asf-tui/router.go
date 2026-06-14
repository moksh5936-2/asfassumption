package main

type focusTarget int

const (
	focusContent focusTarget = iota
	focusSidebar
)

type treeNode struct {
	label     string
	vid       view
	tab       int
	expanded  bool
	children  []*treeNode
	isSection bool
}

type Router struct {
	currentView view
	sidebarSel  int
	focus       focusTarget
	history     viewHistory
	sidebarTree []*treeNode
}

var sidebarTreeBase = []*treeNode{
	{label: "CASES", isSection: true},
	{label: "+ New Analysis", vid: analyzeView, tab: -1},
	{label: "AI", isSection: true},
	{label: "🧠 Local AI", vid: localAIView, tab: -1},
	{label: "SYSTEM", isSection: true},
	{label: "⚙ Settings", vid: settingsView, tab: -1},
	{label: "❓ Help", vid: helpView, tab: -1},
	{label: "ℹ About", vid: aboutView, tab: -1},
}

// sidebarVisibleNodes returns all visible nodes (expanded children included).
func (r *Router) sidebarVisibleNodes() []*treeNode {
	var nodes []*treeNode
	for _, n := range r.sidebarTree {
		nodes = append(nodes, n)
		if n.expanded {
			nodes = append(nodes, n.children...)
		}
	}
	return nodes
}

func newRouter() Router {
	tree := make([]*treeNode, len(sidebarTreeBase))
	for i, n := range sidebarTreeBase {
		cp := *n
		tree[i] = &cp
	}
	return Router{
		currentView: analyzeView,
		focus:       focusContent,
		sidebarTree: tree,
	}
}

func (r *Router) rebuildCaseEntries(caseLabels []string) {
	var tree []*treeNode
	for _, n := range sidebarTreeBase {
		cp := *n
		tree = append(tree, &cp)
		if cp.label == "+ New Analysis" {
			for i, label := range caseLabels {
				tree = append(tree, &treeNode{
					label: "📁 " + label,
					vid:   caseView,
					tab:   i,
				})
			}
		}
	}
	r.sidebarTree = tree
	vis := len(r.sidebarVisibleNodes())
	if r.sidebarSel >= vis {
		r.sidebarSel = vis - 1
	}
	if r.sidebarSel < 0 {
		r.sidebarSel = 0
	}
}

func (r *Router) NavigateTo(to view) view {
	r.history.push(r.currentView)
	r.currentView = to
	return r.currentView
}

func (r *Router) NavigateBack() view {
	if prev, ok := r.history.pop(); ok {
		r.currentView = prev
	} else {
		switch r.currentView {
		case analyzeView, caseView, settingsView, aboutView, reportsView, reviewView, validationView, helpView, localAIView:
			r.currentView = analyzeView
		}
	}
	return r.currentView
}

func (r *Router) SetView(to view) {
	r.currentView = to
}

func (r *Router) ToggleFocus() {
	if r.focus == focusContent {
		r.focus = focusSidebar
		n := len(r.sidebarVisibleNodes())
		if r.sidebarSel >= n {
			r.sidebarSel = n - 1
		}
	} else {
		r.focus = focusContent
	}
}

func (r *Router) sidebarSelView() view {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		return nodes[r.sidebarSel].vid
	}
	return analyzeView
}

func (r *Router) sidebarSelTab() int {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		return nodes[r.sidebarSel].tab
	}
	return -1
}

func (r *Router) sidebarSelIsParent() bool {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		n := nodes[r.sidebarSel]
		return len(n.children) > 0
	}
	return false
}

func (r *Router) sidebarSelIsExpanded() bool {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		return nodes[r.sidebarSel].expanded
	}
	return false
}

func (r *Router) sidebarMoveUp() bool {
	nodes := r.sidebarVisibleNodes()
	orig := r.sidebarSel
	for r.sidebarSel > 0 {
		r.sidebarSel--
		if !nodes[r.sidebarSel].isSection {
			return r.sidebarSel != orig
		}
	}
	r.sidebarSel = orig
	return false
}

func (r *Router) sidebarMoveDown() bool {
	nodes := r.sidebarVisibleNodes()
	orig := r.sidebarSel
	for r.sidebarSel < len(nodes)-1 {
		r.sidebarSel++
		if !nodes[r.sidebarSel].isSection {
			return r.sidebarSel != orig
		}
	}
	r.sidebarSel = orig
	return false
}

func (r *Router) sidebarExpand() {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		nodes[r.sidebarSel].expanded = true
	}
}

func (r *Router) sidebarCollapse() {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) {
		nodes[r.sidebarSel].expanded = false
	}
}

func (r *Router) sidebarActivate() (view, int) {
	nodes := r.sidebarVisibleNodes()
	if r.sidebarSel < len(nodes) && nodes[r.sidebarSel].isSection {
		return r.currentView, -1
	}
	n := r.sidebarSelView()
	tab := r.sidebarSelTab()
	if n != r.currentView || (n == caseView && r.currentView == caseView) {
		r.history.push(r.currentView)
		r.currentView = n
	}
	r.focus = focusContent
	return n, tab
}
