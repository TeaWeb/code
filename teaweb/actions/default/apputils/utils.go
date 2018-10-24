package apputils

import (
	"errors"
	"github.com/TeaWeb/code/teaapps"
	"github.com/TeaWeb/code/teaplugins"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
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

func FavorApp(appId string) error {
	_, app := FindApp(appId)
	if app == nil {
		return errors.New("app not exist")
	}

	confFile := files.NewFile(Tea.ConfigFile("apps_favor.conf"))
	m := map[string]interface{}{}
	if confFile.Exists() {
		reader, err := confFile.Reader()
		if err != nil {
			return err
		}
		defer reader.Close()
		err = reader.ReadYAML(&m)
		if err != nil {
			m = map[string]interface{}{}
		}
	}
	appIds, ok := m["apps"]
	if ok {
		appIdStrings, ok := appIds.([]interface{})
		if ok {
			m["apps"] = append(appIdStrings, app.UniqueId())
		} else {
			m["apps"] = []string{app.UniqueId()}
		}
	} else {
		m["apps"] = []string{app.UniqueId()}
	}

	writer, err := confFile.Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(m)

	return err
}

func CancelFavorApp(appId string) error {
	_, app := FindApp(appId)
	if app == nil {
		return errors.New("app not exist")
	}

	confFile := files.NewFile(Tea.ConfigFile("apps_favor.conf"))
	if !confFile.Exists() {
		return nil
	}

	m := map[string]interface{}{}
	reader, err := confFile.Reader()
	if err != nil {
		return err
	}
	defer reader.Close()
	err = reader.ReadYAML(&m)
	if err != nil {
		m = map[string]interface{}{}
	}

	appIds, ok := m["apps"]
	if !ok {
		return nil
	}

	appIdStrings, ok := appIds.([]interface{})
	if !ok {
		return nil
	}
	m["apps"] = lists.Delete(appIdStrings, app.UniqueId())

	writer, err := confFile.Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(m)

	return err
}

func FavorAppContains(uniqueId string) bool {
	confFile := files.NewFile(Tea.ConfigFile("apps_favor.conf"))
	if !confFile.Exists() {
		return false
	}

	m := map[string]interface{}{}
	reader, err := confFile.Reader()
	if err != nil {
		return false
	}
	defer reader.Close()
	err = reader.ReadYAML(&m)
	if err != nil {
		return false
	}

	appIds, ok := m["apps"]
	if !ok {
		return false
	}

	appIdStrings, ok := appIds.([]interface{})
	if !ok {
		return false
	}

	return lists.Contains(appIdStrings, uniqueId)
}
