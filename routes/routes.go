package routes

import (
	"github.com/bmizerany/pat"
	"github.com/kobeld/duoerl/handlers/brands"
	"github.com/kobeld/duoerl/handlers/feeds"
	"github.com/kobeld/duoerl/handlers/followbrands"
	"github.com/kobeld/duoerl/handlers/products"
	"github.com/kobeld/duoerl/handlers/reviews"
	"github.com/kobeld/duoerl/handlers/sessions"
	"github.com/kobeld/duoerl/handlers/users"
	"github.com/kobeld/duoerl/handlers/wishitems"
	"github.com/kobeld/duoerl/middlewares"
	"github.com/kobeld/mangogzip"
	"github.com/paulbellamy/mango"
	"github.com/shaoshing/train"
	"github.com/sunfmin/mangolog"
	"net/http"
)

func Mux() (mux *http.ServeMux) {
	p := pat.New()
	sessionMW := mango.Sessions("f908b1c425062e95d30b8d30de7123458", "duoerl", &mango.CookieOptions{Path: "/", MaxAge: 3600 * 24 * 7})

	rendererMW := middlewares.ProduceRenderer()
	authenMW := middlewares.AuthenticateUser()
	hardAuthenMW := middlewares.HardAuthenUser()
	rHtml, rJson := middlewares.RespondHtml(), middlewares.RespondJson()

	mainLayoutMW := middlewares.ProduceLayout(middlewares.MAIN_LAYOUT)
	mainStack := new(mango.Stack)
	mainStack.Middleware(mangogzip.Zipper, mangolog.Logger, sessionMW, authenMW, mainLayoutMW, rendererMW, rHtml)

	mainAjaxStack := new(mango.Stack)
	mainAjaxStack.Middleware(mangogzip.Zipper, mangolog.Logger, sessionMW, authenMW, rJson)

	hardAuthenStack := new(mango.Stack)
	hardAuthenStack.Middleware(mangogzip.Zipper, mangolog.Logger, sessionMW, hardAuthenMW, mainLayoutMW, rendererMW, rHtml)

	// User related
	p.Get("/login", mainStack.HandlerFunc(sessions.LoginPage))
	p.Post("/login", mainStack.HandlerFunc(sessions.LoginAction))
	p.Get("/signup", mainStack.HandlerFunc(sessions.SignupPage))
	p.Post("/signup", mainStack.HandlerFunc(sessions.SignupAction))
	p.Get("/logout", mainStack.HandlerFunc(sessions.Logout))

	p.Post("/user/edit", hardAuthenStack.HandlerFunc(users.Update))
	p.Get("/user/edit", hardAuthenStack.HandlerFunc(users.Edit))
	p.Get("/user/:id", mainStack.HandlerFunc(users.Show))

	// Brand related
	p.Get("/brands", mainStack.HandlerFunc(brands.Index))
	p.Get("/brand/new", mainStack.HandlerFunc(brands.New))
	p.Post("/brand/create", mainStack.HandlerFunc(brands.Create))
	p.Get("/brand/:id", mainStack.HandlerFunc(brands.Show))
	p.Get("/brand/:id/edit", mainStack.HandlerFunc(brands.Edit))
	p.Post("/brand/:id/edit", mainStack.HandlerFunc(brands.Update))
	// Follow brand
	p.Post("/brand/follow", mainStack.HandlerFunc(followbrands.Create))
	p.Post("/brand/unfollow", mainStack.HandlerFunc(followbrands.Delete))

	// Product related
	p.Get("/products", mainStack.HandlerFunc(products.Index))
	p.Get("/product/new", mainStack.HandlerFunc(products.New))
	p.Post("/product/create", mainStack.HandlerFunc(products.Create))
	p.Get("/product/:id", mainStack.HandlerFunc(products.Show))
	// p.Get("/product/:id/edit", mainStack.HandlerFunc(products.Edit))
	// p.Post("/product/:id/edit", mainStack.HandlerFunc(products.Update))

	// Review related
	p.Post("/review/create", mainStack.HandlerFunc(reviews.Create))

	// Wish Item related
	p.Post("/wish_item/add", mainAjaxStack.HandlerFunc(wishitems.Create))
	p.Post("/wish_item/remove", mainAjaxStack.HandlerFunc(wishitems.Delete))

	p.Get("/", mainStack.HandlerFunc(feeds.Index))
	mux = http.NewServeMux()
	mux.HandleFunc("/favicon.ico", filterUrl)
	mux.Handle("/", p)
	mux.Handle("/public/", http.FileServer(http.Dir(".")))

	train.ConfigureHttpHandler(mux)
	return
}

func filterUrl(w http.ResponseWriter, r *http.Request) {
	return
}
