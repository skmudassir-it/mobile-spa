package handler

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var (
	db     *sql.DB
	dbOnce sync.Once
	tmpl   *template.Template
)

type Service struct {
	Slug, Title, Category, Description, Image string
	Features                                  []string
}

type TeamMember struct {
	Name, Role, Bio, Initials string
}

type PricingPlan struct {
	Name, Price, Description string
	Features                 []string
	Popular                  bool
}

type FAQ struct{ Q, A string }

type Testimonial struct {
	Name, Text string
	Rating     int
}

type Value struct{ Title, Text, Icon string }

type PageData struct {
	Title, Active    string
	Services         []Service
	Testimonials     []Testimonial
	WhyUs            []Value
	Team             []TeamMember
	Plans            []PricingPlan
	FAQs             []FAQ
	Submitted        bool
	Message          string
}

func initDB() {
	dbOnce.Do(func() {
		var err error
		db, err = sql.Open("sqlite", ":memory:?_journal_mode=WAL")
		if err != nil {
			log.Fatal(err)
		}
		db.Exec(`CREATE TABLE bookings(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT,email TEXT,phone TEXT,service TEXT,message TEXT,created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	})
}

func loadTemplates() {
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"join": func(sep string, items []string) string { return strings.Join(items, sep) },
		"mul":  func(a, b int) int { return a * b },
		"len":  func(v interface{}) int {
			switch val := v.(type) {
			case []Service:
				return len(val)
			case []Testimonial:
				return len(val)
			case []Value:
				return len(val)
			case []TeamMember:
				return len(val)
			case []PricingPlan:
				return len(val)
			case []FAQ:
				return len(val)
			default:
				return 0
			}
		},
	}).ParseFS(templateFS, "templates/*.html"))
}

func render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func navbarHTML(active string) template.HTML {
	links := []struct{ Href, Label string }{{"/", "Home"}, {"/services", "Services"}, {"/about", "About"}, {"/pricing", "Pricing"}, {"/contact", "Contact"}}
	var buf strings.Builder
	buf.WriteString(`<nav class="bg-white/80 backdrop-blur border-b sticky top-0 z-10"><div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between"><a href="/" class="text-xl font-bold text-purple-700">✨ GlowMobile Spa</a><div class="flex gap-1 text-sm">`)
	for _, l := range links {
		class := "px-3 py-1.5 rounded hover:bg-purple-50 transition-colors"
		if l.Href == active || (active == "home" && l.Href == "/") {
			class += " bg-purple-100 text-purple-700 font-medium"
		}
		fmt.Fprintf(&buf, `<a href="%s" class="%s" hx-get="%s" hx-target="body" hx-push-url="true">%s</a>`, l.Href, class, l.Href, l.Label)
	}
	buf.WriteString(`</div></div></nav>`)
	return template.HTML(buf.String())
}

// ─── Handlers ───

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "GlowMobile Spa | Luxury at Your Door",
		Services: []Service{
			{"swedish-massage", "Swedish Massage", "Massage", "Gentle, flowing strokes to melt away tension.", "/static/images/services/swedish-massage.jpg", []string{"60/90 min", "Hot stone option", "Aromatherapy"}},
			{"deep-tissue-massage", "Deep Tissue Massage", "Massage", "Targeted pressure for chronic tension relief.", "/static/images/services/deep-tissue.jpg", []string{"Trigger point therapy", "Sports recovery", "Pain relief"}},
			{"hydrafacial", "HydraFacial", "Facial", "Advanced 3-step facial for radiant skin.", "/static/images/services/hydrafacial.jpg", []string{"Cleanse & exfoliate", "Hydrating serums", "LED therapy"}},
			{"gel-manicure", "Gel Manicure", "Nails", "Long-lasting gel nails with hand massage.", "/static/images/services/manicure.jpg", []string{"Gel or classic", "Cuticle care", "Hand massage"}},
			{"spa-pedicure", "Spa Pedicure", "Nails", "Herbal soak, exfoliation, and leg massage.", "/static/images/services/pedicure.jpg", []string{"Herbal soak", "Sugar scrub", "Paraffin upgrade"}},
			{"couples-massage", "Couples Massage", "Packages", "Side-by-side massages for two.", "/static/images/services/couples-massage.jpg", []string{"Two therapists", "Side-by-side", "Champagne add-on"}},
		},
		Testimonials: []Testimonial{
			{"Sarah M.", "GlowMobile transformed my living room into a five-star spa. The massage was incredible!", 5},
			{"Jennifer K.", "Booked a couples massage for our anniversary. Best date night ever.", 5},
			{"Michael T.", "Got my wife a surprise spa day at home. She loved every minute.", 5},
			{"Ashley R.", "The HydraFacial left my skin glowing for days. So relaxing at home.", 5},
		},
		WhyUs: []Value{
			{"We Come to You", "No traffic, no waiting rooms. Full spa at your doorstep.", "🏠"},
			{"Licensed Pros", "Every therapist fully licensed, insured, and background-checked.", "✅"},
			{"Premium Products", "Top-tier organic and hypoallergenic products only.", "🌿"},
			{"Flexible Scheduling", "Book 24/7. Same-day appointments often available.", "📅"},
		},
	}
	data.Active = "home"
	render(w, "home", data)
}

func handleServices(w http.ResponseWriter, r *http.Request) {
	allServices := []Service{
		{"swedish-massage", "Swedish Massage", "Massage Therapy", "Gentle, flowing strokes to melt away tension and promote deep relaxation.", "/static/images/services/swedish-massage.jpg", []string{"60/90 min sessions", "Hot stone upgrade", "Aromatherapy oils", "Full-body relaxation"}},
		{"deep-tissue-massage", "Deep Tissue Massage", "Massage Therapy", "Targeted pressure to release chronic muscle tension and knots.", "/static/images/services/deep-tissue.jpg", []string{"Trigger point therapy", "Chronic pain relief", "Sports recovery", "Improved mobility"}},
		{"hot-stone-massage", "Hot Stone Massage", "Massage Therapy", "Heated basalt stones combined with therapeutic massage.", "/static/images/services/hot-stone.jpg", []string{"Heated basalt stones", "Deep muscle relief", "Improved circulation", "Full body"}},
		{"prenatal-massage", "Prenatal Massage", "Massage Therapy", "Gentle, supportive massage for expecting mothers.", "/static/images/services/prenatal-massage.jpg", []string{"Side-lying position", "Pregnancy pillow", "Reduced swelling", "Back pain relief"}},
		{"hydrafacial", "HydraFacial", "Facials & Skincare", "Advanced 3-step facial: cleanse, extract, hydrate.", "/static/images/services/hydrafacial.jpg", []string{"Deep cleanse", "Hydrating serums", "LED light therapy", "No downtime"}},
		{"anti-aging-facial", "Anti-Aging Facial", "Facials & Skincare", "Collagen-boosting facial with peptides and antioxidants.", "/static/images/services/anti-aging-facial.jpg", []string{"Collagen peptides", "Vitamin C infusion", "Microcurrent option", "Firming & lifting"}},
		{"gel-manicure", "Gel Manicure", "Nail Care", "Long-lasting gel nails with cuticle care and hand massage.", "/static/images/services/manicure.jpg", []string{"Gel or classic polish", "Cuticle treatment", "Hand & arm massage", "Nail art available"}},
		{"spa-pedicure", "Spa Pedicure", "Nail Care", "Herbal foot soak, exfoliation, nail care, leg massage.", "/static/images/services/pedicure.jpg", []string{"Herbal foot soak", "Sugar scrub", "Callus treatment", "Paraffin upgrade"}},
		{"couples-massage", "Couples Massage", "Packages", "Side-by-side massages for two — perfect for any occasion.", "/static/images/services/couples-massage.jpg", []string{"Two therapists", "Side-by-side tables", "Champagne add-on", "Customizable"}},
		{"spa-party-package", "Spa Party Package", "Packages", "Host a spa party! Multiple therapists for any celebration.", "/static/images/services/spa-party.jpg", []string{"3+ therapists", "Custom menu", "Group pricing", "Refreshments"}},
	}
	categories := []string{"Massage Therapy", "Facials & Skincare", "Nail Care", "Packages"}
	data := PageData{Title: "Services | GlowMobile Spa", Services: allServices, Active: "services"}
	// We need categories for the template
	render(w, "services", map[string]interface{}{"Data": data, "Categories": categories, "AllServices": allServices})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "About | GlowMobile Spa", Active: "about",
		Team: []TeamMember{
			{"Amara Chen", "Founder & Lead Therapist", "15+ years in luxury spa. Certified in 12+ modalities.", "AC"},
			{"Marcus Rivera", "Senior Massage Therapist", "Deep tissue specialist. Former NFL team therapist.", "MR"},
			{"Priya Patel", "Lead Esthetician", "Medical esthetician. HydraFacial & anti-aging expert.", "PP"},
			{"Jasmine Okonkwo", "Nail Art Specialist", "Award-winning nail artist. Gel, acrylic, nail art.", "JO"},
		},
		WhyUs: []Value{
			{"Excellence", "Every treatment with the highest standards of quality and care.", "✨"},
			{"Convenience", "Five-star treatments to your space — no travel, no hassle.", "🏠"},
			{"Safety", "Licensed, insured pros with hospital-grade sanitization.", "🛡️"},
			{"Sustainability", "Organic, cruelty-free products and eco-friendly practices.", "🌿"},
		},
	}
	render(w, "about", data)
}

func handlePricing(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Pricing | GlowMobile Spa", Active: "pricing",
		Plans: []PricingPlan{
			{"Massage Therapy", "From $99", "Swedish, Deep Tissue, Hot Stone, Prenatal", []string{"60 or 90 min sessions", "Licensed therapist", "Aromatherapy included", "Hot stone upgrade", "Couples add-on"}, true},
			{"Facials & Skincare", "From $129", "HydraFacial, Anti-Aging, Custom Facials", []string{"45 or 60 min sessions", "Medical-grade products", "LED light therapy", "Serum infusions", "Zero downtime"}, false},
			{"Nail Care", "From $55", "Gel Manicure, Spa Pedicure, Nail Art", []string{"Gel or classic polish", "Cuticle treatment", "Hand/foot massage", "Nail art available", "Paraffin upgrade"}, false},
			{"Spa Packages", "From $199", "Couples, Parties, Corporate Wellness", []string{"Multiple therapists", "Custom menus", "Group pricing", "Bridal/birthday packages", "Corporate rates"}, false},
		},
		FAQs: []FAQ{
			{"How do you set up a spa in my home?", "We bring everything — massage tables, equipment, products, linens, music. Just need 6x8 ft clear space and a sink. Setup: 15 min."},
			{"Do I need to provide anything?", "Just yourself! We bring all equipment and products."},
			{"How far in advance should I book?", "24-48 hours recommended. Same-day often available."},
			{"Cancellation policy?", "Free cancellation up to 4 hours before. Late: $25 fee."},
			{"Are therapists licensed?", "Yes! Fully licensed, insured, background-checked professionals."},
			{"Do you do group events?", "Absolutely! Bridal showers, birthdays, corporate wellness — contact us for group pricing."},
		},
	}
	render(w, "pricing", data)
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Contact | GlowMobile Spa", Active: "contact"}
	render(w, "contact", data)
}

func handleContactSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	service := r.FormValue("service")
	message := r.FormValue("message")

	db.Exec(`INSERT INTO bookings(name,email,phone,service,message) VALUES(?,?,?,?,?)`, name, email, phone, service, message)

	data := PageData{
		Title: "Thank You | GlowMobile Spa", Active: "contact",
		Submitted: true,
		Message:  "We've received your booking request! We'll get back to you within 2 hours during business hours.",
	}
	render(w, "contact-success", data)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initDB()
	if tmpl == nil {
		loadTemplates()
	}

	path := r.URL.Path

	if strings.HasPrefix(path, "/static/") {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		subFS, _ := fs.Sub(staticFS, "static")
		http.StripPrefix("/static/", http.FileServer(http.FS(subFS))).ServeHTTP(w, r)
		return
	}

	switch path {
	case "/":
		handleHome(w, r)
	case "/services":
		handleServices(w, r)
	case "/about":
		handleAbout(w, r)
	case "/pricing":
		handlePricing(w, r)
	case "/contact":
		handleContact(w, r)
	case "/api/contact":
		handleContactSubmit(w, r)
	default:
		http.NotFound(w, r)
	}
}

