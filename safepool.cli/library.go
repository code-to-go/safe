package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/code-to-go/safepool.lib/api/library"
	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/skratchdot/open-golang/open"
)

func addDocument(l library.Library) {
	for {
		prompt := promptui.Prompt{
			Label: "Local Path",
		}

		localPath, _ := prompt.Run()
		if localPath == "" {
			return
		}

		stat, err := os.Stat(localPath)
		if err != nil {
			color.Red("invalid path '%s'", localPath)
			continue
		}
		if stat.IsDir() {
			color.Red("folders are not supported at the moment")
			continue
		}

		var items []string
		var item string
		parts := strings.Split(localPath, string(os.PathSeparator))
		sort.Slice(parts, func(i, j int) bool { return i > j })
		for _, p := range parts {
			if p != "" {
				item = path.Join(p, item)
				items = append(items, item)
			}
		}

		sel := promptui.Select{
			Label: "Name in the pool",
			Items: items,
		}
		_, name, _ := sel.Run()

		prompt = promptui.Prompt{
			Label:   "Edit Name",
			Default: name,
		}
		name, _ = prompt.Run()
		h, err := l.Send(localPath, name)
		if core.IsErr(err, "cannot upload %s: %v", localPath) {
			color.Red("cannot upload %s", localPath)
		} else {
			color.Green("'%s' uploaded to '%s:%s' with id %d", localPath, l.Pool.Name, name, h.Id)
			return
		}
	}

}

type actionType int

const (
	uploadLocal actionType = iota
	openLocally
	openFolder
	deletelocal
	updateLocal
	downloadTemp
)

type action struct {
	typ      actionType
	document library.Document
}

func actionsOnDocument(l library.Library, ds []library.Document) {
	items := []string{"🔙 Back"}
	var actions []action

	for _, d := range ds {
		if d.LocalPath != "" {
			if d.Mode == library.Newer {
				items = append(items, "send update to the pool")
				actions = append(actions, action{uploadLocal, d})
			}
			items = append(items, "open locally")
			actions = append(actions, action{openLocally, d})
			items = append(items, "open local folder")
			actions = append(actions, action{openFolder, d})
			items = append(items, "delete")
			actions = append(actions, action{deletelocal, d})
		} else {
			i, _ := security.IdentityFromId(d.AuthorId)
			switch d.Mode {
			case library.Newer:
				items = append(items, fmt.Sprintf("receice update from %s", i.Nick))
				actions = append(actions, action{updateLocal, d})
			case library.Conflict:
				items = append(items, fmt.Sprintf("receive replacement from %s", i.Nick))
				actions = append(actions, action{updateLocal, d})
			}
			items = append(items, fmt.Sprintf("download from %s to a temporary location", i.Nick))
			actions = append(actions, action{downloadTemp, d})
		}
	}

	label := fmt.Sprintf("Choose the action on '%s'", ds[0].Name)
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	idx, _, _ := prompt.Run()
	if idx == 0 {
		return
	}
	a := actions[idx-1]
	switch a.typ {
	case openLocally:
		open.Start(a.document.LocalPath)
	case openFolder:
		open.Start(filepath.Dir(a.document.LocalPath))
	case deletelocal:

	case uploadLocal:
		l.Send(a.document.LocalPath, a.document.Name)
	case updateLocal:
		localPath, _ := l.GetLocalPath(a.document.Name)
		l.Receive(a.document.Id, localPath)
	case downloadTemp:
		dest := filepath.Join(os.TempDir(), a.document.Name)
		l.Save(a.document.Id, dest)
	}

}

func documentString(d library.Document) string {
	var mode, author string
	switch d.Mode {
	case library.Same:
		mode = "✓"
	case library.Newer:
		mode = "👶"
	case library.Older:
		mode = "👴"
	case library.Conflict:
		mode = "💣"
	case library.Deleted:
		mode = "🗑"
	}

	if identity, ok, _ := security.GetIdentity(d.AuthorId); ok {
		author = identity.Nick
	} else {
		author = d.AuthorId
	}
	return fmt.Sprintf("%s %s ←%s 🔗%s", mode, d.Name, author, d.LocalPath)
}

func Library(p *pool.Pool) {
	p.Sync()
	l := library.Get(p, "library")

	for {
		documents, err := l.List("")
		if core.IsErr(err, "cannot read document list: %v") {
			color.Red("something wrong")
			return
		}

		m := map[string][]library.Document{}
		for _, d := range documents {
			if d.LocalPath != "" {
				m[d.Name] = append([]library.Document{d}, m[d.Name]...)
			} else if _, ok := m[d.Name]; !ok {
				m[d.Name] = append(m[d.Name], d)
			}
		}

		var d2 [][]library.Document
		for _, ds := range m {
			d2 = append(d2, ds)
		}

		items := []string{"🔙 Back", "⟳ Refresh", "＋ Add"}
		for _, ds := range d2 {
			d := ds[0]
			if d.Mode == library.Folder {
				items = append(items, fmt.Sprintf("📁 %s", d.Name))
			} else {
				items = append(items, documentString(d))
			}
		}

		prompt := promptui.Select{
			Label: "Choose",
			Items: items,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}
		switch idx {
		case 0:
			return
		case 1:
			p.Sync()
		case 2:
			addDocument(l)
		default:
			actionsOnDocument(l, d2[idx-3])
		}
	}
}
