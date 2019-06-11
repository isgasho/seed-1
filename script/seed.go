package script

import (
	"fmt"
	"strings"
)

import qlova "github.com/qlova/script"
import "github.com/qlova/script/language"
import "github.com/qlova/script/language/javascript"

import "github.com/qlova/seed/style/css"

type String = qlova.String
type Object string

type Expression struct {
	seed       Seed
	expression string
}

func (p Promise) Raw() string {
	return p.expression
}

func dashes2camels(s string) string {
	var camel string
	var parts = strings.Split(s, "-")
	for i, part := range parts {
		if i == 0 {
			camel += part
		} else {
			camel += strings.Title(part)
		}
	}
	return camel
}

type Seed struct {
	css.Style

	ID, Native string
	Q          Script
}

func (seed Seed) Set(property, value string) {
	property = dashes2camels(property)

	seed.Javascript(`set(` + seed.Element() + `, "` + property + `", "` + value + `");`)
}

func (seed Seed) Get(property string) string {

	property = dashes2camels(property)

	return string(`getComputedStyle(` + seed.Element() + `).` + property)
}

//TODO
func (seed Seed) Bytes() []byte {
	return nil
}

func (seed Seed) Element() string {
	if seed.Native != "" {
		return seed.Native
	}
	return `get("` + seed.ID + `")`
}

func (seed Seed) Javascript(js string) {
	seed.Q.Raw("Javascript", language.Statement(js))
}

type File Expression

func (f File) Type() String {
	return f.seed.wrap(f.expression + `.type`)
}

func (f File) Name() String {
	return f.seed.wrap(f.expression + `.name`)
}

func (seed Seed) wrap(s string) String {
	return seed.Q.StringFromLanguageType(Javascript.String{
		Expression: language.Statement(s),
	})
}

func (seed Seed) SetText(s String) {
	seed.Javascript(seed.Element() + `.textContent = ` + raw(s) + `;`)
}

func (seed Seed) SetPath(s String) {
	seed.Javascript(seed.Element() + `.src = ` + raw(s) + `;`)
}
func (seed Seed) SetSource(s String) {
	seed.Javascript(seed.Element() + `.src = ` + raw(s) + `;`)
}

func (seed Seed) SetHTML(s String) {
	seed.Javascript(seed.Element() + `.innerHTML = ` + raw(s) + `;`)
}

func (seed Seed) SetLeft(s String) {
	seed.Javascript(`set(` + seed.Element() + `, "left", ` + raw(s) + `);`)
}

func (seed Seed) SetDisplay(s String) {
	seed.Javascript(`set(` + seed.Element() + `, "display", ` + raw(s) + `);`)
}

func (seed Seed) SetVisible() {
	seed.Javascript(`set(` + seed.Element() + `, "display", "inline-flex");`)
}

func (seed Seed) SetHidden() {
	seed.Javascript(`set(` + seed.Element() + `, "display", "none");`)
}

func (seed Seed) Click() {
	seed.Javascript(seed.Element() + `.click();`)
}

var unique int

func Unique() string {
	unique++
	return fmt.Sprint("unique_", unique)
}

func (seed Seed) Play() Promise {
	var variable = Unique()
	seed.Javascript(`let ` + variable + ` = ` + seed.Element() + `.play();`)
	return Promise{q: seed.Q, expression: variable}
}

func (seed Seed) Pause() {
	seed.Javascript(seed.Element() + `.pause();`)
}

func (seed Seed) Focus() {
	seed.Javascript(seed.Element() + `.focus();`)
}

func (seed Seed) Restart() {
	seed.Javascript(seed.Element() + `.load();`)
}

func (seed Seed) Left() String {
	return seed.wrap(seed.Element() + `.style.left`)
}

func (seed Seed) Width() String {
	return seed.wrap(`getComputedStyle(get("` + seed.ID + `")).width`)
}

func (seed Seed) SetValue(value String) {
	seed.Javascript(seed.Element() + `.value = ` + raw(value) + `;`)
}

func (seed Seed) SetPlaceholder(value String) {
	seed.Javascript(seed.Element() + `.placeholder = ` + raw(value) + `;`)
}

func (seed Seed) SetClass(value String) {
	seed.Javascript(seed.Element() + `.className = ` + raw(value) + `;`)
}

func (seed Seed) Value() String {
	return seed.wrap(seed.Element() + `.value`)
}

func (seed Seed) Text() String {
	return seed.wrap(seed.Element() + `.innerText`)
}

func (seed Seed) Location() String {
	return seed.wrap(seed.Element() + `.href`)
}

func (seed Seed) Data(key string) String {
	return seed.wrap(seed.Element() + `.data["` + key + `"]`)
}

//Return the index of this feeditem.
func (seed Seed) Index() String {
	return seed.wrap(seed.Element() + `.index`)
}

func (seed Seed) HTML() String {
	return seed.wrap(seed.Element() + `.innerHTML`)
}

func (seed Seed) File() File {
	return File{seed: seed, expression: seed.Element() + `.files[0]`}
}

func (seed Seed) Display() String {
	return seed.wrap(seed.Element() + `.style.display`)
}

//Temporary method DEPRECIATED
func (f File) Raw() string {
	return f.expression
}

func (seed Seed) Load(f File) {
	seed.Javascript(seed.Element() + `.src = window.URL.createObjectURL(` + f.expression + `);`)
}

//Add a child seed to this seed.
func (seed Seed) Add(child Seed) {
	seed.Javascript(seed.Element() + `.appendChild(` + child.Element() + `);`)
}

func (seed Seed) OnClick(f func()) {
	seed.Javascript(seed.Element() + `.onclick = function() {`)
	f()
	seed.Javascript(`};`)
}

//Animations
func (seed Seed) FadeOut() {
	seed.Set("animation", "fadeOut .5s ease both")

	seed.Javascript(`set(` + seed.Element() + `, "display", "inline-flex");`)
	seed.Javascript(`set(` + seed.Element() + `, "z-index", "50");`)
	seed.Javascript(`animating = true;`)
	seed.Javascript(`setTimeout(function() { set(` + seed.Element() + `, "display", "none"); set(` + seed.Element() + `, "animation", ""); set(` + seed.Element() + `, "z-index", "");animation_complete();  }, 500);`)
}

func (seed Seed) FadeIn() {

	seed.Javascript(`set(` + seed.Element() + `, "display", "none");`)

	seed.Javascript(`setTimeout(function() { `)
	seed.Set("animation", "fadeIn .5s ease both")
	seed.Javascript(`set(` + seed.Element() + `, "display", "inline-flex");`)

	seed.Javascript(`if (!last_page) return;`)
	seed.Javascript(`set(get(last_page), "display", "none");`)
	seed.Javascript(`animating = true;`)
	seed.Javascript(`setTimeout(function() { set(` + seed.Element() + `, "z-index", ""); animation_complete(); }, 500);`)
	seed.Javascript(`}, 500);`)
}

//Animations
func (seed Seed) SlideInFrom(direction complex128) {
	if direction == -1i {
		seed.Set("animation", "slideInFromTop .5s ease both")
	}

	if direction == 1i {
		seed.Set("animation", "slideInFromBottom .5s ease both")
	}

	if direction == 1 {
		seed.Set("animation", "slideInFromRight .5s ease both")
	}

	if direction == -1 {
		seed.Set("animation", "slideInFromLeft .5s ease both")
	}

	seed.Javascript(`if (!last_page) return;`)
	seed.Javascript(`set(get(last_page), "display", "inline-flex");`)
	seed.Javascript(`set(` + seed.Element() + `, "z-index", "50");`)
	seed.Javascript(`animating = true;`)
	seed.Javascript(`setTimeout(function() { set(get(last_page), "display", "none"); set(` + seed.Element() + `, "z-index", ""); animation_complete(); }, 500);`)
}

//Animations
func (seed Seed) SlideOutTo(direction complex128) {
	if direction == -1i {
		seed.Set("animation", "slideOutToTop .5s ease both")
	}

	if direction == 1i {
		seed.Set("animation", "slideOutToBottom .5s ease both")
	}

	if direction == 1 {
		seed.Set("animation", "slideOutToRight .5s ease both")
	}

	if direction == -1 {
		seed.Set("animation", "slideOutToLeft .5s ease both")
	}

	seed.Javascript(`set(` + seed.Element() + `, "display", "inline-flex");`)
	seed.Javascript(`set(` + seed.Element() + `, "z-index", "50");`)
	seed.Javascript(`animating = true;`)
	seed.Javascript(`setTimeout(function() { set(` + seed.Element() + `, "display", "none"); set(` + seed.Element() + `, "animation", ""); set(` + seed.Element() + `, "z-index", "");animation_complete(); }, 500);`)
}

func (seed Seed) Translate(x, y Unit) {
	seed.Javascript(seed.Element() + `.style.transform = "translate(` + x.Raw() + "," + y.Raw() + `)";`)
}
