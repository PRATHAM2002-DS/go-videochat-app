package server

import(
	"flag"
	"os"
	"time"
	"videochat-project/internal/handlers"
	w "videochat-project/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/gofiber/fiber/middleware/logger"	
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
)

var (
	addr = flag.String("addr",":"+ os.Getenv("PORT"),"")
	cert = flag.String("cert","","")
	key = flag.String("key","","")
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	engine := html.New("./views",".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid",handlers.Room)
	app.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/room/:uuid/chat",handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket",websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket",websocket.New(handlers.RoomViewerWebsocket))
	app.Get("/stream/:ssuid",handlers.Stream)
	app.Get("/stream/:ssuid/websocket", websocket.New(handlers.StreamWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}) )
	app.Get("/stream/:ssuid/chat/websocket", websocket.New(handlers.StreamChatWebsocket))
	app.Get("/stream/:ssuid/viewer/websocket", websocket.New(handlers.StreamViewerWebsocket))
	app.Static("/", "./assets")
    
	w.Rooms = make(map[string]*w.Room)
	w.Streams = make(map[string]*w.Stream)
	go dispatchKeyFrames()
	if *cert != "" {
		return app.ListenTLS(*addr,*cert,*key)
	}
	return app.Listen(*addr)


}  

func dispatchKeyFrames() {
	for range time.NewTicker(3 * time.Second).C {
		for _, room := range w.Rooms{
			room.Peers.DispatchKeyFrames()
		}
	}
	
}