package main

type Router struct {
	currentView view
	sidebarSel  int
	history     viewHistory
}

func newRouter() Router {
	return Router{
		currentView: startupView,
	}
}

func (r *Router) NavigateTo(to view) view {
	r.history.push(r.currentView)
	r.currentView = to
	r.syncSidebar()
	return r.currentView
}

func (r *Router) NavigateBack() view {
	if prev, ok := r.history.pop(); ok {
		r.currentView = prev
	} else {
		switch r.currentView {
		case fileBrowserView, dashboardView, analyzeView, resultsView, localaiView, settingsView, aboutView, exportView, reviewView, validationView, helpView:
			r.currentView = dashboardView
		}
	}
	r.syncSidebar()
	return r.currentView
}

func (r *Router) SetView(to view) {
	r.currentView = to
	r.syncSidebar()
}

func (r *Router) PushHistory() {
	r.history.push(r.currentView)
}

func (r *Router) CycleSidebar(dir int) {
	n := len(sidebarEntries)
	r.sidebarSel = (r.sidebarSel + dir + n) % n
}

func (r *Router) ActivateSidebarTab() int {
	entry := sidebarEntries[r.sidebarSel]
	if entry.vid == resultsView && entry.tab >= 0 {
		return entry.tab
	}
	return -1
}

func (r *Router) ActivateSidebar() view {
	entry := sidebarEntries[r.sidebarSel]
	if entry.vid != r.currentView || (entry.vid == resultsView && r.currentView == resultsView) {
		r.history.push(r.currentView)
		r.currentView = entry.vid
	}
	return r.currentView
}

func (r *Router) ViewForSidebar(i int) view {
	if i >= len(sidebarEntries) {
		return dashboardView
	}
	return sidebarEntries[i].vid
}

func (r *Router) syncSidebar() {
	for i, e := range sidebarEntries {
		if e.vid == r.currentView {
			r.sidebarSel = i
			return
		}
	}
}
