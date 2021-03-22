package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/bobby-wan/streaming-platform/webserver/db"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	pathHTML               = ".\\static\\html\\"
	pathTemplatesHTML      = ".\\static\\html\\templates\\"
	streamingServerAddress = "hls://"
	rtmpAddress            = "rtmp://localhost/live"
	hlsAddress             = "http://127.0.0.1:7002"
)

var userController db.UserControllerInterface
var streamController db.StreamControllerInterface

var regexEmail = regexp.MustCompile(".+@.+\\..+")

// var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var store = sessions.NewCookieStore([]byte("SESSION_KEY"))

// var store, err = sessionmanager.NewSessionStore()

//LoginForm helps with the login process
type LoginForm struct {
	Email    string
	Password string
}

//SignupForm helps with the signup process
type SignupForm struct {
	Email    string
	Password string
	Username string
}

// type PageModel struct {
// 	Username    string
// 	Room        string
// 	IsLoggedIn  bool
// 	IsStreaming bool
// }

type UserInfo struct {
	Username    string
	ID          uint
	IsStreaming bool
	StreamID    uint
	StreamURL   string
	ImagePath   string
}

func generateHLSaddress(appname string, username string) string {
	return fmt.Sprintf("%v/%v/%v.m3u8", hlsAddress, appname, username)
}

// var streamCache map[string]sessionmanager.StreamingSession
func generateStream(username string, title string, category db.ContentCategory) (*db.ActiveStream, error) {
	stream := db.ActiveStream{
		Username: username,
		Title:    title,
		Viewers:  0,
		Category: uint32(category),
		URL:      generateHLSaddress("live", username),
		Active:   false,
	}
	ptrStream, err := streamController.Create(stream)
	return ptrStream, err
}

func getUserInfoFromSession(req *http.Request) *UserInfo {
	var userInfo *UserInfo
	session, err := store.Get(req, "auth")

	if err != nil {
		log.Println(err)
		return nil
	}

	if session != nil && !session.IsNew {
		rawUsername, ok := session.Values["username"]
		if !ok {
			return nil
		}
		rawID, ok := session.Values["userId"]
		if !ok {
			return nil
		}
		// rawStreamingFlag, ok := session.Values["isStreaming"]
		// if !ok {
		// 	return nil
		// }
		// rawStreamURL, ok := session.Values["streamURL"]
		// if !ok {
		// 	return nil
		// }
		// rawStreamImagePath, ok := session.Values["streamImagePath"]
		// if !ok {
		// 	return nil
		// }

		username, ok := rawUsername.(string)
		if !ok {
			return nil
		}

		id, ok := rawID.(uint)
		if !ok {
			return nil
		}

		// isStreaming, ok := rawStreamingFlag.(bool)
		// if !ok {
		// 	return nil
		// }

		// streamURL, ok := rawStreamURL.(string)
		// if !ok {
		// 	return nil
		// }

		// imagePath, ok := rawStreamImagePath.(string)
		// if !ok {
		// 	return nil
		// }

		userInfo = &UserInfo{
			Username: username,
			ID:       id,
			// IsStreaming: isStreaming,
			// StreamURL:   streamURL,
			// ImagePath:   imagePath,
		}
	}
	return userInfo
}

func respondBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	return
}

func respondInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func respondUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	return
}

func streamingServerAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.

		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		if ip != "::1" {
			respondUnauthorized(w)
			return
		}
		h.ServeHTTP(w, req) // call ServeHTTP on the next handler handler
	})
}

func endStreamHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	url, ok := queryParams["url"]
	if !ok {
		respondBadRequest(w)
		return
	}
	splitURL := strings.Split(url[0], "/")
	username := splitURL[len(splitURL)-1]

	// user, err := userController.Get(username)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	_, err := streamController.End(username)
	if err != nil {
		fmt.Println(err)
		respondInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func startStreamHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	username, ok := queryParams["username"]
	if !ok || len(username) != 1 {
		respondBadRequest(w)
		return
	}

	// user, err := userController.Get(username[0])
	// if err != nil || user == nil {
	// 	fmt.Println(err)
	// 	respondBadRequest(w)
	// 	return
	// }
	_, err := streamController.Start(username[0])
	if err != nil {
		fmt.Println(err)
		respondInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.

		session, err := store.Get(req, "auth")
		if err != nil {
			h.ServeHTTP(w, req) // call ServeHTTP on the original handler
		}
		if session.IsNew {
			http.Redirect(w, req, "/login", 401)
		}

		// session.Options.HttpOnly = true
		// session.Options.MaxAge = 1800 //30 minutes

		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(w, req) // call ServeHTTP on the original handler
	})
}

func getHTMLTemplates() (*template.Template, error) {
	templates, err := template.ParseGlob(pathTemplatesHTML + "*")
	if err != nil {
		return nil, err
	}
	return templates.ParseGlob(pathHTML + "*.html")
}

func getTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, err := getHTMLTemplates()
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func internalServerError(w *http.ResponseWriter, message *string) {
	(*w).WriteHeader(http.StatusInternalServerError)
	if message != nil {
		(*w).Write([]byte(*message))
	}
}

func (form LoginForm) validate() *map[string]string {
	errors := make(map[string]string)

	match := regexEmail.Match([]byte(form.Email))
	if match == false {
		errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(form.Password) == "" {
		errors["Password"] = "Please enter password"
	}

	if len(errors) > 0 {
		return &errors
	}

	return nil
}

func getDefaultPageData() map[string]interface{} {
	data := make(map[string]interface{}, 0)
	data["isLoggedIn"] = false
	data["isStreaming"] = false
	data["username"] = nil
	data["room"] = nil
	data["categories"] = getStreamingCategories()

	return data
}

func updatePageData(user *UserInfo, data *map[string]interface{}) {
	if user == nil || data == nil {
		return
	}

	(*data)["username"] = user.Username
	(*data)["isStreaming"] = user.IsStreaming
	(*data)["isLoggedIn"] = true
}
func homeHandler(w http.ResponseWriter, req *http.Request) {
	data := getDefaultPageData()
	userInfo := getUserInfoFromSession(req)
	updatePageData(userInfo, &data)
	getTemplate(w, "index", data)
}

func loginGetHandler(w http.ResponseWriter, req *http.Request) {
	getTemplate(w, "login", nil)
}

func loginPostHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	form := LoginForm{}
	form.Email = req.PostFormValue("email")
	form.Password = req.PostFormValue("password")

	errors := form.validate()
	if errors != nil {
		getTemplate(w, "login", errors)
		return
	}

	user, err := userController.GetByEmail(form.Email)
	if err != nil {
		//add invalid email/password message here
		getTemplate(w, "login", nil)
		return
	}

	if user == nil {
		log.Println("Nil user received from db.")
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password)); err != nil {
		//add invalid email/password message here
		getTemplate(w, "login", nil)
		return
	}

	//maybe email here
	session, err := store.Get(req, "auth")
	//check error here?
	session.Values["username"] = user.Username
	session.Values["userId"] = user.ID
	session.Values["isStreaming"] = false
	// session.Values["streamURL"] = nil
	// session.Values["streamImagePath"] = nil

	session.Options.MaxAge = 1800

	err = session.Save(req, w)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func signupGetHandler(w http.ResponseWriter, req *http.Request) {
	getTemplate(w, "signup", nil)
}

func getStream(streamer string) *db.ActiveStream {
	ptrStream, err := streamController.GetByUsername(streamer)
	if err != nil {
		return nil
	}

	return ptrStream
}

func signupPostHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	form := LoginForm{}
	form.Email = req.PostFormValue("email")
	form.Password = req.PostFormValue("password")
}

func viewStreamHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	streamerUsername := vars["streamer"]

	data := getDefaultPageData()
	data["room"] = streamerUsername

	userInfo := getUserInfoFromSession(req)

	stream := getStream(streamerUsername)
	if stream != nil && stream.Active {
		data["isStreaming"] = true
		data["stream"] = *stream
	}
	if userInfo != nil {
		data["isLoggedIn"] = true
		data["username"] = userInfo.Username
	}
	getTemplate(w, "user", data)
}

func getStreamingCategories() []string {
	categories := make([]string, 0)
	for i := db.Gaming; i <= db.Art; i = i << 1 {
		categories = append(categories, i.String())
	}

	return categories
}

func getContentCategoryByName(name string) (*db.ContentCategory, error) {
	for i := db.Gaming; i <= db.Art; i = i << 1 {
		if i.String() == name {
			return &i, nil
		}
	}

	return nil, errors.New("no such category")
}

func liveHandler(w http.ResponseWriter, req *http.Request) {
	userInfo := getUserInfoFromSession(req)
	data := make(map[string]interface{}, 0)

	if userInfo != nil {
		data["isLoggedIn"] = true
		data["username"] = userInfo.Username
		data["room"] = userInfo.Username
		data["isStreaming"] = userInfo.IsStreaming
		data["categories"] = getStreamingCategories()

		getTemplate(w, "live", data)
	} else {
		http.Redirect(w, req, "/login", 302)
	}
}

func streamHandler(w http.ResponseWriter, req *http.Request) {
	userInfo := getUserInfoFromSession(req)

	if userInfo == nil {
		http.Redirect(w, req, "/login", 401)
		return
	}

	parameters := req.URL.Query()
	title, ok := parameters["title"]

	if !ok || len(title) != 1 {
		respondBadRequest(w)
		return
	}

	categoryName, ok := parameters["category"]
	if !ok || len(categoryName) != 1 {
		respondBadRequest(w)
		return
	}

	cat, err := getContentCategoryByName(categoryName[0])
	if err != nil {
		respondBadRequest(w)
		return
	}

	res, err := http.Get("http://localhost:8090/control/get?room=" + userInfo.Username)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resData := make(map[string]interface{})
	err = json.Unmarshal(body, &resData)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status := resData["status"]
	statusCode, ok := status.(float64)

	if !ok || statusCode != 200 {
		log.Fatal("problem with parsing response status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := resData["data"]
	dataString, ok := data.(string)
	if !ok {
		log.Fatal("problem with parsing response data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	response["server-address"] = rtmpAddress
	response["stream-key"] = dataString

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Fatal("problem with marshalling response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream, err := generateStream(userInfo.Username, title[0], *cat)
	if err != nil {
		respondInternalServerError(w)
		return
	}
	fmt.Println(stream)

	_, err = w.Write(responseBody)
	if err != nil {
		log.Fatal("problem with parsing response data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func getStreamsByCategory(category *string) []db.ActiveStream {
	streams := make([]db.ActiveStream, 0)
	if category == nil {
		dbstreams, err := streamController.GetByViewCount(100)
		if err != nil {
			fmt.Println(err)
			return streams
		}
		return dbstreams
	}

	cat, err := getContentCategoryByName(*category)
	if err != nil {
		fmt.Println(err)
		return streams
	}

	dbstreams, err := streamController.GetByCategories(uint32(*cat))
	if err != nil {
		fmt.Println(err)
		return streams
	}

	return dbstreams
}

func browseHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	category := queryParams["category"]
	data := getDefaultPageData()
	if len(category) == 1 {
		cat, err := getContentCategoryByName(category[0])
		if err == nil {
			streams, err := streamController.GetByCategories(uint32(*cat))
			if err == nil {
				if len(streams) > 20 {
					data["streams"] = streams[0:20]
				} else {
					data["streams"] = streams
				}
			} else {
				http.Redirect(w, req, "/", 302)
			}
		} else {
			http.Redirect(w, req, "/", 302)
		}
	}

	user := getUserInfoFromSession(req)
	updatePageData(user, &data)
	getTemplate(w, "browse", data)
}

func initializeDbControllers() error {
	database, err := db.Initialize()
	if err != nil {
		return err
	}

	ptrUserController, err := db.NewUserControllerGORM(database)
	if err != nil {
		return err
	}

	userController = ptrUserController

	ptrStreamController, err := db.NewStreamControllerGORM(database)
	if err != nil {
		return err
	}

	streamController = ptrStreamController

	return nil
}

func initialize() error {
	initializeDbControllers()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")

	r.HandleFunc("/signup", signupGetHandler).Methods("GET")
	r.HandleFunc("/signup", signupPostHandler).Methods("POST")

	r.HandleFunc("/live", liveHandler).Methods("GET")
	r.HandleFunc("/live/{streamer}", viewStreamHandler).Methods("GET")
	r.HandleFunc("/configure", streamHandler).Methods("GET")

	r.HandleFunc("/browse", browseHandler).Methods("GET")
	r.HandleFunc("/end", streamingServerAuth(endStreamHandler)).Methods("GET")
	r.HandleFunc("/start", streamingServerAuth(startStreamHandler)).Methods("GET")

	http.Handle("/", r)
	return http.ListenAndServe(":8080", nil)
}

func main() {
	log.Fatal(initialize())
}
