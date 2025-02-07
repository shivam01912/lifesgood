package llm

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"log"
	"net/http"
)

const generativeModelName = "gemini-1.5-flash"

type RagServer struct {
	Ctx context.Context
	//WvStore      weaviate.Store
	GeminiClient *googleai.GoogleAI
}

/*
func (rs *RagServer) AddDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type document struct {
		Text string
	}
	type addRequest struct {
		Documents []document
	}
	ar := &addRequest{}

	err := readRequestJSON(req, ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store documents and their embeddings in weaviate
	var wvDocs []schema.Document
	for _, doc := range ar.Documents {
		wvDocs = append(wvDocs, schema.Document{PageContent: doc.Text})
	}
	_, err = rs.WvStore.AddDocuments(rs.Ctx, wvDocs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rs *RagServer) QueryHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type queryRequest struct {
		Content string
	}
	qr := &queryRequest{}
	err := readRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the most similar documents.
	docs, err := rs.WvStore.SimilaritySearch(rs.Ctx, qr.Content, 3)
	if err != nil {
		http.Error(w, fmt.Errorf("similarity search: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	var docsContents []string
	for _, doc := range docs {
		docsContents = append(docsContents, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as
	// context.
	ragQuery := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(docsContents, "\n"))
	respText, err := llms.GenerateFromSinglePrompt(rs.Ctx, rs.GeminiClient, ragQuery, llms.WithModel(generativeModelName))
	if err != nil {
		log.Printf("calling generative model: %v", err.Error())
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	renderJSON(w, respText)
}

const ragTemplateStr = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.

Question:
%s

Context:
%s
`
*/

func (rs *RagServer) ReviewHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type queryRequest struct {
		Storefront string
		Content    []string
	}
	qr := &queryRequest{}
	err := readRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a RAG query for the LLM with the most relevant documents as
	// context.
	ragQuery := fmt.Sprintf(reviewTemplateStr, qr.Storefront, constructBulletPoints(qr.Content))
	fmt.Println(ragQuery)
	llmResponse, err := llms.GenerateFromSinglePrompt(rs.Ctx, rs.GeminiClient, ragQuery, llms.WithModel(generativeModelName))
	if err != nil {
		log.Printf("error calling generative model: %v", err.Error())
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	renderJSON(w, llmResponse)
}

const reviewTemplateStr = `
I am the owner of %s. 
Here are some of the reviews of my restaurant on Google Maps - 
%s
Please provide me the sentiment of each review and give me some actionable insights from all the reviews combined in the below format.
Sentiment : Review number : 1 to 5 (1 - highly negative, 5 - highly positive)
Action : Review number : Reason
`
