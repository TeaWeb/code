package app

import (
	"github.com/TeaWeb/code/teaapps"
	"github.com/TeaWeb/code/teaplugins"
)

func FindApp(appId string) (*teaplugins.Plugin, *teaapps.App) {
	for _, p := range teaplugins.Plugins() {
		for _, a := range p.Apps {
			if a.Id == appId {
				return p, a
			}
		}
	}
	return nil, nil
}
