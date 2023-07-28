package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grace/api"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type bookOptions struct {
	id          uint
	title       string
	author      string
	publishedAt string
	edition     string
	description string
	genre       string
}

var baseURL string

type Book struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	Edition     string    `json:"edition"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
}

func main() {
	var rootCmd = &cobra.Command{Use: "books"}
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "http://localhost:3000", "Base URL of the REST API")
	rootCmd.AddCommand(NewAddCommand(), NewGetCommand(), NewUpdateCommand(), NewDeleteCommand())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Command execution failed: %v\n", err)
		os.Exit(1)
	}
}

func NewAddCommand() *cobra.Command {
	options := bookOptions{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new Book",
		Run: func(cmd *cobra.Command, args []string) {
			addBook(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.title, "title", "", "Title of the Book")
	flags.StringVar(&options.author, "author", "", "Author of the Book")
	flags.StringVar(&options.publishedAt, "published-at", "", "Publish date of the Book")
	flags.StringVar(&options.edition, "edition", "", "Edition of the Book")
	flags.StringVar(&options.description, "description", "", "Description of the Book")
	flags.StringVar(&options.genre, "genre", "", "Genre of the Book")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("author")
	cmd.MarkFlagRequired("published-at")

	return cmd
}

func NewGetCommand() *cobra.Command {
	options := bookOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Books",
		Run: func(cmd *cobra.Command, args []string) {
			getBooks(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.author, "author", "", "Filter Books by author")
	flags.StringVar(&options.genre, "genre", "", "Filter Books by genre")

	return cmd
}

func NewUpdateCommand() *cobra.Command {
	options := bookOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Book",
		Run: func(cmd *cobra.Command, args []string) {
			updateBook(options)
		},
	}

	flags := cmd.Flags()
	flags.UintVar(&options.id, "id", 0, "ID of the Book to update")
	flags.StringVar(&options.title, "title", "", "New title of the Book")
	flags.StringVar(&options.author, "author", "", "New author of the Book")
	flags.StringVar(&options.publishedAt, "published-at", "", "New publish date of the Book")
	flags.StringVar(&options.edition, "edition", "", "New edition of the Book")
	flags.StringVar(&options.description, "description", "", "New description of the Book")
	flags.StringVar(&options.genre, "genre", "", "New genre of the Book")

	cmd.MarkFlagRequired("id")

	return cmd
}

func NewDeleteCommand() *cobra.Command {
	options := bookOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Book",
		Run: func(cmd *cobra.Command, args []string) {
			deleteBook(options)
		},
	}

	flags := cmd.Flags()
	flags.UintVar(&options.id, "id", 0, "ID of the Book to delete")
	cmd.MarkFlagRequired("id")

	return cmd
}

func addBook(options bookOptions) {
	parsedPublishedAt, err := time.Parse("2006-01-02", options.publishedAt)
	if err != nil {
		fmt.Println("Invalid published date")
		return
	}

	book := Book{
		Title:       options.title,
		Author:      options.author,
		PublishedAt: parsedPublishedAt,
		Edition:     options.edition,
		Description: options.description,
		Genre:       options.genre,
	}

	body, err := json.Marshal(book)
	if err != nil {
		fmt.Printf("Failed to serialize Book object: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/book", baseURL), bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to add Book: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if err := checkHttpResponse(resp); err != nil {
		fmt.Printf("Failed to add Book: %v\n", err)
	} else {
		fmt.Println("Book added successfully")
	}
}

func getBooks(options bookOptions) {
	url := fmt.Sprintf("%s/v1/book", baseURL)
	if options.author != "" {
		url += "?author=" + options.author
	}
	if options.genre != "" {
		url += "?genre=" + options.genre
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch Books: %v\n", err)
		return
	}
	defer resp.Body.Close()

	type BookResponse struct {
		api.Response
		Message []Book `json:"message"`
	}

	var response BookResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Failed to decode response body:", err.Error())
		return
	}

	if response.Type != "sync" || response.StatusCode != 200 {
		fmt.Printf("Failed to get books. Error code: %d. Error: %s\n", response.Code, response.Error)
		return
	}
	printBookTable(response.Message)
}

func updateBook(options bookOptions) {
	var err error
	book := make(map[string]interface{})

	if options.publishedAt != "" {
		_, err = time.Parse("2006-01-02", options.publishedAt)
		if err != nil {
			fmt.Println("Invalid published date")
			return
		}
		book["published_at"] = options.publishedAt
	}

	if options.title != "" {
		book["title"] = options.title
	}
	if options.author != "" {
		book["author"] = options.author
	}
	if options.edition != "" {
		book["edition"] = options.edition
	}
	if options.description != "" {
		book["description"] = options.description
	}
	if options.genre != "" {
		book["genre"] = options.genre
	}

	body, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Failed to serialize Book object")
		return
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/book/%d", baseURL, options.id), bytes.NewReader(body))
	if err != nil {
		fmt.Println("Failed to create request:", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to update Book:", err.Error())
		return
	}
	defer resp.Body.Close()

	if err := checkHttpResponse(resp); err != nil {
		fmt.Println("Failed to update Book:", err)
	} else {
		fmt.Println("Book updated successfully")
	}
}

func deleteBook(options bookOptions) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/book/%d", baseURL, options.id), nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to delete Book: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if response, err := handleHttpResponse(resp); err != nil {
			fmt.Printf("Failed to delete Book and parse responpse body: %v\n", err)
		} else {
			fmt.Printf("Failed to delete Book: %s\n", response.Error)
		}
	} else {
		fmt.Println("Book deleted successfully")
	}
}

func handleHttpResponse(resp *http.Response) (*api.Response, error) {
	var response api.Response
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode HTTP response body into api.Response: %w", err)
	}

	return &response, nil
}

func checkHttpResponse(resp *http.Response) error {
	response, err := handleHttpResponse(resp)
	if err != nil {
		return err
	}

	if response.Type != "sync" || response.StatusCode != 200 {
		return fmt.Errorf("Unexpected response. Expected type 'sync' and status code '200'. Got type '%s' and status code '%d'. Error: %s", response.Type, response.StatusCode, response.Error)
	}

	return nil
}

func printBookTable(books []Book) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Author", "Published At", "Edition", "Description", "Genre"})

	for _, book := range books {
		t.AppendRow([]interface{}{
			book.ID,
			book.Title,
			book.Author,
			book.PublishedAt.Format("2006-01-02"),
			book.Edition,
			book.Description,
			book.Genre,
		})
	}

	t.Render()
}
