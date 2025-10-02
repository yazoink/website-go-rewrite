package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"os"
	"io/ioutil"
	"golang.org/x/exp/slices"
	"github.com/gorilla/mux"
	// "strings"
)

var assetDir = "static/"
var jsonDir = assetDir + "json/"
var templateDir = assetDir + "templates/"
var imageDir = assetDir + "images/"
var graphicsDir = imageDir + "graphics/"
var cssDir = assetDir + "css/"
var baseTemplate = templateDir + "base.html"

type Link struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type NavLinks struct {
	Nav []Link `json:"nav"`
	Socials []Link `json:"socials"`
}

type BlogPost struct {
	Title string `json:"title"`
	Subheading string `json:"subheading"`
	Date string `json:"date"`
	Categories []string `json:"categories"`
}

type BasePage struct {
	PageTitle string
	CssFile string
	HeaderImage string
	NavImage string
	AsideImage string
	NavLinks NavLinks
	BlogPosts []BlogPost
	BlogCategories []string
	SelectedBlogPost BlogPost
	SelectedBlogCategory string
	BlogPostsInSelectedCategory []BlogPost
}

func main() {
	r := mux.NewRouter()
	data := BasePage {
		CssFile: "/" + cssDir + "style.css",
		HeaderImage: "/" + graphicsDir + "header.webp",
		NavImage: "/" + graphicsDir + "nav-image.webp",
		AsideImage: "/" + graphicsDir + "aside-image.webp",
		NavLinks: parseLinksJson(jsonDir + "links.json"),
		BlogPosts: parseBlogJson(jsonDir + "blog.json"),
	}
	data.BlogCategories = getBlogCategories(data.BlogPosts)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Welcome..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "home.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/about", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "About Me..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "about.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/bookmarks", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Bookmarks..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "bookmarks.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/technology", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Technology..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "technology.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/services", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Services..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "services.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/blog", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Blog..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "blog.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/blog/category/{category}", func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data.PageTitle = "Blog..."
		if slices.Contains(data.BlogCategories, vars["category"]) {
			data.SelectedBlogCategory = vars["category"]
			data.BlogPostsInSelectedCategory = getBlogPostsInSelectedCategory(data.BlogPosts, vars["category"])
		}
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "blog.html")
		tmpl.Execute(w, data)
		data.SelectedBlogCategory = ""
		data.BlogPostsInSelectedCategory = nil
	})
	r.HandleFunc("/blog/post/{post}", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = ""
		vars := mux.Vars(r)
		postIndex := -1
		for i := 0; i < len(data.BlogPosts); i++ {
			// if strings.ToLower(strings.ReplaceAll(data.BlogPosts[i].Title, " ", "-")) == vars["post"] {
			if data.BlogPosts[i].Title == vars["post"] {
				postIndex = i
			}
		}
		if postIndex != -1 {
			data.SelectedBlogPost = data.BlogPosts[postIndex]
			fmt.Println("??")
		}
		fmt.Println(postIndex)
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "blog-post.html")
		tmpl.Execute(w, data)
	})
	r.HandleFunc("/gallery", func (w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Art Gallery..."
		tmpl, _ := template.ParseFiles(baseTemplate, templateDir + "gallery.html")
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8080", r)
}

func getBlogPostsInSelectedCategory(blogPosts []BlogPost, category string) []BlogPost {
	var selectedPosts []BlogPost
	for i := 0; i < len(blogPosts); i++ {
		for j := 0; j < len(blogPosts[i].Categories); j++ {
			if blogPosts[i].Categories[j] == category {
				selectedPosts = append(selectedPosts, blogPosts[i])
			}
		}
	}
	return selectedPosts
}

func parseLinksJson(path string) NavLinks {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var navLinks NavLinks
	json.Unmarshal(byteValue, &navLinks)
	return navLinks
}

func parseBlogJson(path string) []BlogPost {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var blogPosts []BlogPost
	json.Unmarshal(byteValue, &blogPosts)
	// fmt.Println(blogPosts)
	return blogPosts
}

func getBlogCategories(blogPosts []BlogPost) []string {
	var categories []string
	for i := 0; i < len(blogPosts); i++ {
		for j :=0; j < len(blogPosts[i].Categories); j++ {
			if len(categories) < 1 || !slices.Contains(categories, blogPosts[i].Categories[j]) {
				categories = append(categories, blogPosts[i].Categories[j])
			}
		}
	}
	// fmt.Println(categories)
	return categories
}
