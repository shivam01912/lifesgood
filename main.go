package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/googleai"
	"lifesgood/app/admin"
	adminBlogHandler "lifesgood/app/admin/blog"
	"lifesgood/app/requestHandler"
	"lifesgood/llm"
	"log"
	"net/http"
	"os"
)

const embeddingModelName = "text-embedding-004"

func main() {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./data/templates"))
	router.PathPrefix("/css/").Handler(fs)

	//direct REST API handler
	router.HandleFunc("/", requestHandler.HomePageHandler)
	router.HandleFunc("/home", requestHandler.HomePageHandler)
	router.HandleFunc("/blog", requestHandler.BlogHandler)

	//admin login handlers
	router.HandleFunc("/admin/login", admin.LoginHandler)
	router.HandleFunc("/admin/home", admin.AdminHome)
	router.HandleFunc("/admin", admin.ProcessLogin)

	//blog counter update handlers
	router.HandleFunc("/blog/likes", requestHandler.LikesIncrement)
	router.HandleFunc("/blog/views", requestHandler.ViewsIncrement)

	//image handler
	router.HandleFunc("/blog/image", requestHandler.ReadFile)

	//create blog handlers
	router.HandleFunc("/admin/addblog", adminBlogHandler.AddBlogHandler)
	router.HandleFunc("/admin/blog/preview", adminBlogHandler.PreviewBlog)
	router.HandleFunc("/admin/blog/publish", adminBlogHandler.ProcessPublishBlog)

	//common admin handlers
	router.HandleFunc("/admin/modifyblog", adminBlogHandler.DeletePageHandler)

	//update blog handlers
	router.HandleFunc("/admin/updateblog", adminBlogHandler.UpdateBlogPageHandler)
	router.HandleFunc("/blog/update", adminBlogHandler.ProcessUpdateBlog)

	//delete blog handlers
	router.HandleFunc("/blog/delete", adminBlogHandler.DeleteBlogHandler)

	//LLM test application
	ctx := context.Background()
	////apiKey := os.Getenv("GEMINI_API_KEY")
	geminiClient, err := googleai.New(ctx,
		googleai.WithAPIKey(os.Getenv("Gemini_Key")),
		googleai.WithDefaultEmbeddingModel(embeddingModelName))
	if err != nil {
		log.Fatal(err)
	}
	//
	//emb, err := embeddings.NewEmbedder(geminiClient)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//wvStore, err := weaviate.New(
	//	weaviate.WithEmbedder(emb),
	//	weaviate.WithScheme("https"),
	//	weaviate.WithAPIKey(os.Getenv("Weaviate_Key")),
	//	weaviate.WithHost("bn3muyvktqh10vd3vg6xa.c0.asia-southeast1.gcp.weaviate.cloud"),
	//	weaviate.WithIndexName("Document"),
	//)

	server := &llm.RagServer{
		Ctx: ctx,
		//WvStore:      nil,
		GeminiClient: geminiClient,
	}

	//mux := http.NewServeMux()
	//router.HandleFunc("/add", server.AddDocumentsHandler)
	//router.HandleFunc("/query", server.QueryHandler)
	router.HandleFunc("/review", server.ReviewHandler)

	//port := cmp.Or(os.Getenv("SERVERPORT"), "8080")
	port := ":8092"
	//address := "localhost" + port
	//log.Println("listening on", address)
	//log.Fatal(http.ListenAndServe(address, router))

	log.Println("Starting Server.")
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Println("Unknown error : ", err)
		return
	}
	log.Println("Started Server.")
}
