package beego

import (
	"encoding/json"
	"mime"
	"path/filepath"

	"github.com/qasico/beego/session"
)

//
func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

func registerSession() error {
	if BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		if sessionConfig == "" {
			conf := map[string]interface{}{
				"cookieName":      BConfig.WebConfig.Session.SessionName,
				"gclifetime":      BConfig.WebConfig.Session.SessionGCMaxLifetime,
				"providerConfig":  filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig),
				"secure":          BConfig.Listen.EnableHTTPS,
				"enableSetCookie": BConfig.WebConfig.Session.SessionAutoSetCookie,
				"domain":          BConfig.WebConfig.Session.SessionDomain,
				"cookieLifeTime":  BConfig.WebConfig.Session.SessionCookieLifeTime,
			}
			confBytes, err := json.Marshal(conf)
			if err != nil {
				return err
			}
			sessionConfig = string(confBytes)
		}
		if GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, sessionConfig); err != nil {
			return err
		}
		go GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	if err := BuildTemplate(BConfig.WebConfig.ViewsPath); err != nil {
		if BConfig.RunMode == DEV {
			Warn(err)
		}
		return err
	}
	return nil
}

func registerDocs() error {
	if BConfig.WebConfig.EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
	return nil
}

func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {
		go beeAdminApp.Run()
	}
	return nil
}
