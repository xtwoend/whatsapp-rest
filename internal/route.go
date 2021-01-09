package internal

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/xtwoend/whatsapp-rest/pkg/auth"
	"github.com/xtwoend/whatsapp-rest/pkg/router"

	"github.com/xtwoend/whatsapp-rest/internal/index"
	"github.com/xtwoend/whatsapp-rest/internal/whatsapp"
	"github.com/xtwoend/whatsapp-rest/pkg/server"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {
	// Set Endpoint for Root Functions
	router.Router.Get(router.RouterBasePath, index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	// Set Endpoint for WhatsApp Functions
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/login", whatsapp.WhatsAppLogin)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/logout", whatsapp.WhatsAppLogout)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/ping", whatsapp.WhatsAppPing)

	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/text", whatsapp.WhatsAppSendText)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/location", whatsapp.WhatsAppSendLocation)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/document", whatsapp.WhatsAppSendDocument)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/audio", whatsapp.WhatsAppSendAudio)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/image", whatsapp.WhatsAppSendImage)
	router.Router.With(auth.Basic).Post(router.RouterBasePath+"/send/video", whatsapp.WhatsAppSendVideo)
	

	// serv downloaded media 
	mediaDir := http.Dir(server.Config.GetString("SERVER_UPLOAD_PATH"))
	FileServer("/media", mediaDir)
}

func FileServer(path string, root http.FileSystem) {
	
	r := router.Router
	
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}