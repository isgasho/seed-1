/*
	Seed is an eco-friendly Go library for creating cross-platform applications that work on the Desktop, Mobile and Web.
	
	Example:
		
		package main

		import "github.com/qlova/seed"

		func main() {
			var App = seed.New()
			App.SetName("Hello World")
			App.SetText("Hello World")
			App.Launch()
		}

*/

package seed

import "github.com/qlova/seed/worker"
import "github.com/qlova/seed/style"

import (
	"net/http"
	"fmt"
	"path"
	"bytes"
	"html"
	"strings"
)

//TODO cleanup
var ServiceWorker = worker.NewServiceWorker()

//TODO cleanup
func RegisterAsset(path string) {
	ServiceWorker.Assets[path] = true
}

//DEPRECIATED
func (seed Seed) ID() string {
	return seed.id
}

/*Set the name of the application generated by this seed.

	**Desktop** 
	this will be shown on the titlebar.
	
	**Mobile**
	This will become the name of the app.
	
	**Web**
	This will become the name of the tab.
*/
func (seed Seed) SetName(name string) {
	seed.manifest.Name = name
	if seed.manifest.ShortName == "" {
		seed.manifest.ShortName = name
	}
}

/*
	Set the short name of the application generated by this seed.
*/
func (seed Seed) SetShortName(name string) {
	seed.manifest.ShortName = name
}


/*Set the description of the application generated by this seed.

	**Desktop** 
	N/A
	
	**Mobile**
	This will become the description of the app.
	
	**Web**
	N/A
*/
func (seed Seed) SetDescription(description string) {
	seed.manifest.Description = description
}

/*Set the icon for the application generated by this seed.

	**Desktop** 
	This will be the icon shown on the titlebar and in the taskbar.
	
	**Mobile**
	This will become the icon for the app.
	
	**Web**
	This will become the icon for the tab.
*/
func (seed Seed) SetIcon(path string) {
	//TODO
	seed.manifest.AddIcon(path)
}

func (seed Seed) SetThemeColor(color string) {
	//TODO
	seed.manifest.ThemeColor = color
}

//TODO Should be Internal.
func (seed Seed) SetClass(class string) {
	seed.class = class
}

//TODO Should be Internal.
func (seed Seed) SetTag(tag string) {
	seed.tag = tag
}
//TODO Should be Internal.
func (seed Seed) SetAttributes(attr string) {
	seed.attr = attr
}
//TODO Should be Internal.
func (seed Seed) Attributes() string {
	return seed.attr
}

//TODO Should be Internal.
func (seed Seed) SetPlaceholder(placeholder string) {
	seed.attr += " placeholder='"+placeholder+"' "
}

//Add a font to the seed.
//TODO merge with style?
/*func (seed Seed) AddFont(name, file, weight string) {
	
	switch weight {
		case "black":
			weight = "900"
		case "semi-bold":
			weight = "600"
		case "regular":
			weight = "400"
		case "light":
			weight = "300"
		case "extra-light":
			weight = "200"
	}
	
	RegisterAsset(file)
	
	seed.fonts.Write([]byte(`@font-face {
	font-family: '`+name+`';
	src: url('`+file+`');
	font-weight: `+weight+`;
}
`))
}*/

//Does this need to be here?
func (seed Seed) GetStyle() *style.Style {
	return &seed.Style
}

func (seed Seed) Page() bool {
	return seed.page
}

func (seed Seed) Require(script string) {
	seed.scripts = append(seed.scripts, script)
}

//Add a child seed to this seed.
func (seed Seed) Add(child Interface) {
	seed.children = append(seed.children, child)
	child.GetSeed().SetParent(seed)
}

//Add a handler to the seed, when this seed is launched as root, the handlers will be executed for each incomming request.
func (seed Seed) AddHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	seed.handlers = append(seed.handlers, handler)
}


func (seed Seed) GetParent() Interface {
	return seed.parent
}


func (seed Seed) SetParent(parent Interface) {
	seed.parent = parent
}

func (seed Seed) GetChildren() []Interface {
	return seed.children
}

//Add text, html or whatever!
func (seed Seed) SetContent(data string) {
	seed.content = []byte(data)
}

//Set the text content of the seed.
func (seed Seed) SetText(data string) {
	data = html.EscapeString(data)
	data = strings.Replace(data, "\n", "<br>", -1)
	data = strings.Replace(data, "  ", "&nbsp;", -1)
	data = strings.Replace(data, "\t", "&emsp;", -1)
	seed.content = []byte(data)
}

type Client struct {
	client
}

func (client Client) WriteString(s string) {
	client.client.ResponseWriter.Write([]byte(s))
}

type client struct {
	http.ResponseWriter
	*http.Request
}

//Set the text content of the seed which will be dynamic at runtime.
func (seed Seed) SetDynamicText(f func(Client)) {
	seed.dynamicText = f
}


func (seed Seed) OnClick(f func(Script)) {
	if seed.onclick == nil {
		seed.onclick = f
	} else {
		var old = seed.onclick
		seed.onclick = func(q Script) {
			old(q)
			f(q)
		}
	}
}

func (seed Seed) OnClickGoto(page Seed) {
	if seed.onclick == nil {
		seed.onclick = func(q Script) {
			q.Goto(page)
		}
	} else {
		var old = seed.onclick
		seed.onclick = func(q Script) {
			old(q)
			func(q Script) {
				q.Goto(page)
			}(q)
		}
	}
}

func (seed Seed) OnReady(f func(Script)) {
	if seed.onready == nil {
		seed.onready = f
	} else {
		var old = seed.onready
		seed.onready = func(q Script) {
			old(q)
			f(q)
		}
	}
}

func (seed Seed) OnPageEnter(f func(Script)) {
	seed.OnReady(func(q Script) {
		q.Javascript(q.Get(seed).Element()+".enterpage = function() {")
		f(q)
		q.Javascript("};")
	})
}

func (seed Seed) OnPageExit(f func(Script)) {
	seed.OnReady(func(q Script) {
		q.Javascript(q.Get(seed).Element()+".exitpage = function() {")
		f(q)
		q.Javascript("};")
	})
}


func (seed Seed) OnChange(f func(Script)) {
	if seed.onchange == nil {
		seed.onchange = f
	} else {
		var old = seed.onchange
		seed.onchange = func(q Script) {
			old(q)
			f(q)
		}
	}
}

func (seed Seed) buildStyleSheet(sheet *style.Sheet) {
	seed.postProduction()
	if data := seed.Style.Bytes(); data != nil {
		seed.styled = true
		sheet.Add("#"+seed.id, seed.Style)
	}
	for _, child := range seed.children {
		child.GetSeed().buildStyleSheet(sheet)
	}
}

func (seed Seed) BuildStyleSheet() style.Sheet {
	var stylesheet = make(style.Sheet)
	seed.buildStyleSheet(&stylesheet)
	return stylesheet
}

func (seed Seed) buildStyleSheetForLandscape(sheet *style.Sheet) {
	seed.postProduction()
	if data := seed.Landscape.Bytes(); data != nil {
		seed.styled = true
		sheet.Add("#"+seed.id, seed.Landscape)
	}
	for _, child := range seed.children {
		child.GetSeed().buildStyleSheetForLandscape(sheet)
	}
}

func (seed Seed) BuildStyleSheetForLandscape() style.Sheet {
	var stylesheet = make(style.Sheet)
	seed.buildStyleSheetForLandscape(&stylesheet)
	return stylesheet
}

func (seed Seed) buildFonts() map[style.Font]struct{} {
	
	var fonts = make(map[style.Font]struct{})
	if seed.font.FontFace.FontFamily != "" {
		fonts[seed.font] = struct{}{}
	}

	for _, child := range seed.children {
		for font := range child.GetSeed().buildFonts() {
			fonts[font] = struct{}{}
		}
	}
	
	return fonts
}

func (seed Seed) BuildFonts() []byte {
	var buffer bytes.Buffer
	
	var fonts = seed.buildFonts()

	for font := range fonts {
		buffer.WriteString("@font-face {")
		buffer.Write(font.Bytes())
		buffer.WriteByte('}')
	}

	return buffer.Bytes()
}

func (seed Seed) buildAnimations(animations *[]Animation, names *[]string) {
	
	if seed.animation != nil {
		*animations = append(*animations, seed.animation)
		*names = append(*names, seed.ID())
	}

	for _, child := range seed.children {
		child.GetSeed().buildAnimations(animations, names)
	}
}

func (seed Seed) BuildAnimations() []byte {
	var buffer bytes.Buffer
	
	var animations = make([]Animation, 0) 
	var names = make([]string, 0) 
	seed.buildAnimations(&animations, &names)

	for i, animation := range animations {
		buffer.WriteString("@keyframes "+names[i]+" {")
		buffer.Write(animation.Bytes())
		buffer.WriteByte('}')
	}

	return buffer.Bytes()
}

type dynamicHandler struct {
	id string
	handler func(Client)
}

func (seed Seed) buildDynamicHandler(handler *[]dynamicHandler) {
	
	if seed.dynamicText != nil {
		(*handler) = append((*handler), dynamicHandler{
			id: seed.id,
			handler: seed.dynamicText,
		})
	}
	
	for _, child := range seed.children {
		child.GetSeed().buildDynamicHandler(handler)
	}
}


func (seed Seed) BuildDynamicHandler() (func(w http.ResponseWriter, r *http.Request)) {
	var handlers = make([]dynamicHandler, 0)
	seed.buildDynamicHandler(&handlers)
	
	if len(handlers) == 0 {
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			w.Write([]byte(`"`))
			w.Write([]byte(handler.id))
			w.Write([]byte(`":"`))
			handler.handler(Client{client{
				Request: r,
				ResponseWriter: w, 
			}})
			w.Write([]byte(`"`))
		}
	}
}

func (seed Seed) HTML() ([]byte) {
	seed.postProduction()

	var html bytes.Buffer
	
	html.WriteByte('<')
	html.WriteString(seed.tag)
	html.WriteByte(' ')
	if seed.attr != "" {
		html.WriteString(seed.attr)
		html.WriteByte(' ')
	}
	html.WriteString("id='")
	html.WriteString(fmt.Sprint(seed.id))
	html.WriteByte('\'')
	
	if seed.attr != "" {
		html.WriteString("class='")
		html.WriteString(seed.class)
		html.WriteByte('\'')
	}
	
	if data := seed.Style.Bytes(); !seed.styled && data != nil {
		html.WriteString(" style='")
		html.Write(data)
		html.WriteByte('\'')
	}
	
	if seed.onclick != nil {
		html.WriteString(" onclick='")
		html.Write(toJavascript(seed.onclick))
		html.WriteByte('\'')
	}
	
	if seed.onchange != nil {
		html.WriteString(" onchange='")
		html.Write(toJavascript(seed.onchange))
		html.WriteByte('\'')
	}
	
	html.WriteByte('>')
	
	if seed.content != nil {
		html.Write(seed.content)
	}
	
	for _, child := range seed.children {
		html.Write(child.GetSeed().Render())
	}
	
	html.WriteString("</")
	html.WriteString(seed.tag)
	html.WriteByte('>')
	
	return html.Bytes()
}

func (seed Seed) Render() []byte {
	return seed.HTML()
}

//Return a fully fully rendered application in HTML for the seed.
func (seed Seed) render(production bool) []byte {
	var style = seed.BuildStyleSheet().Bytes()
	var styleForLandscape = seed.BuildStyleSheetForLandscape().Bytes()
	var html = seed.HTML()
	var fonts = seed.BuildFonts()
	var animations = seed.BuildAnimations()
	var scripts = seed.Scripts()
	var onready = seed.BuildOnReady()

	var buffer bytes.Buffer
	buffer.Write([]byte(`<!DOCTYPE html><html><head>
		<meta name="viewport" content="height=device-height, width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no">

		<meta name="apple-mobile-web-app-capable" content="yes">
		<meta name="apple-mobile-web-app-status-bar-style" content="black">
		<meta name="theme-color" content="`+seed.manifest.ThemeColor+`">
		
		<title>`+seed.manifest.Name+`</title>
		
		

		<link rel="manifest" href="/app.webmanifest">`))

	for _, icon := range seed.manifest.Icons {
		buffer.Write([]byte(`<link rel="apple-touch-icon" sizes="`+icon.Sizes+`" href="`+icon.Source+`">`))
	}

	for script := range scripts {
		if path.Ext(script) == ".js" { 
			buffer.Write([]byte(`<script src="`+script+`"></script>`))
		} else if path.Ext(script) == ".css" {
			buffer.Write([]byte(`<link rel="stylesheet" href="`+script+`" />`))
		}
	}

	buffer.Write([]byte(`<script>
			if ('serviceWorker' in navigator) {
				window.addEventListener('load', function() {
					navigator.serviceWorker.register('/index.js').then(function(registration) {
						console.log('ServiceWorker registration successful with scope: ', registration.scope);
					}, function(err) {
						console.log('ServiceWorker registration failed: ', err);
					});
				});
			}
		</script>
		
		
		<style>
	`))

	buffer.Write(fonts)
	buffer.Write(animations)
	buffer.Write(style)	

	buffer.WriteString(`@media screen and (orientation: landscape) {`)
	buffer.Write(styleForLandscape)
	buffer.WriteString(`}`)

	//Optimise to array
	var PagesArray string
	for i, page := range pages {
		PagesArray += "'"+page.ID()+"'"
		if i < len(pages)-1 {
			PagesArray += ","
		}
	}

		buffer.Write([]byte(`
		</style>
			
		<style>

			* {
				
				flex-shrink: 0;
			}
			
			::-webkit-scrollbar { 
				display: none; 
			}

			a {
				text-decoration: none;
			}
			
			p {
				margin-block-start: 0;
				margin-block-end: 0;
			}
			
			 html, body {
				position: fixed;
				overscroll-behavior: none; 
				-webkit-overscroll-behavior: none; 
				-webkit-overflow-scrolling: none; 
				cursor: pointer; 
				margin: 0; 
				padding: 0;
				height: 100%;
				width: 100%;
				-webkit-touch-callout: none;
				-webkit-user-select: none;
				-khtml-user-select: none;
				-moz-user-select: none;
				-ms-user-select: none;
				user-select: none;
				-webkit-tap-highlight-color: transparent;
				
				/* Some nice defaults for centering content. */
				display: inline-flex;
				align-items: center;
				justify-content: center;
				flex-direction: row;
				overflow: hidden;
			}
		</style>
		
		<script>
			//NO BOUNCE BUGS GODDAMMIT https://github.com/lazd/iNoBounce
			(function(global){var startY=0;var enabled=false;var supportsPassiveOption=false;try{var opts=Object.defineProperty({},"passive",{get:function(){supportsPassiveOption=true}});window.addEventListener("test",null,opts)}catch(e){}var handleTouchmove=function(evt){var el=evt.target;var zoom=window.innerWidth/window.document.documentElement.clientWidth;if(evt.touches.length>1||zoom!==1){return}while(el!==document.body&&el!==document){var style=window.getComputedStyle(el);if(!style){break}if(el.nodeName==="INPUT"&&el.getAttribute("type")==="range"){return}var scrolling=style.getPropertyValue("-webkit-overflow-scrolling");var overflowY=style.getPropertyValue("overflow-y");var height=parseInt(style.getPropertyValue("height"),10);var isScrollable=scrolling==="touch"&&(overflowY==="auto"||overflowY==="scroll");var canScroll=el.scrollHeight>el.offsetHeight;if(isScrollable&&canScroll){var curY=evt.touches?evt.touches[0].screenY:evt.screenY;var isAtTop=startY<=curY&&el.scrollTop===0;var isAtBottom=startY>=curY&&el.scrollHeight-el.scrollTop===height;if(isAtTop||isAtBottom){evt.preventDefault()}return}el=el.parentNode}evt.preventDefault()};var handleTouchstart=function(evt){startY=evt.touches?evt.touches[0].screenY:evt.screenY};var enable=function(){window.addEventListener("touchstart",handleTouchstart,supportsPassiveOption?{passive:false}:false);window.addEventListener("touchmove",handleTouchmove,supportsPassiveOption?{passive:false}:false);enabled=true};var disable=function(){window.removeEventListener("touchstart",handleTouchstart,false);window.removeEventListener("touchmove",handleTouchmove,false);enabled=false};var isEnabled=function(){return enabled};var testDiv=document.createElement("div");document.documentElement.appendChild(testDiv);testDiv.style.WebkitOverflowScrolling="touch";var scrollSupport="getComputedStyle"in window&&window.getComputedStyle(testDiv)["-webkit-overflow-scrolling"]==="touch";document.documentElement.removeChild(testDiv);if(scrollSupport){enable()}var iNoBounce={enable:enable,disable:disable,isEnabled:isEnabled};if(typeof module!=="undefined"&&module.exports){module.exports=iNoBounce}if(typeof global.define==="function"){(function(define){define("iNoBounce",[],function(){return iNoBounce})})(global.define)}else{global.iNoBounce=iNoBounce}})(this);
			
			var get = function(id) {
				return document.getElementById(id)
			};
						
			var pages = [`+PagesArray+`];
			var last_page = null;
			var goto = function(next_page_id) {
				for (let page_id of pages) {
					let page = get(page_id);
					if (getComputedStyle(page).display != "none") {
						if (page.exitpage) page.exitpage();
						set(page, 'display', 'none');
						last_page = page_id;
					}
				}
				let next_page = get(next_page_id);
				if (next_page.enterpage) next_page.enterpage();
				set(next_page, 'display', 'inline-flex');
			};
			var back = function() {
				if (last_page == null) return;
				goto(last_page);
			};

			function setCookie(cname, cvalue, exdays) {
			  var d = new Date();
			  d.setTime(d.getTime() + (exdays*24*60*60*1000));
			  var expires = "expires="+ d.toUTCString();
			  document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
			}

			function getCookie(cname) {
			  var name = cname + "=";
			  var decodedCookie = decodeURIComponent(document.cookie);
			  var ca = decodedCookie.split(';');
			  for(var i = 0; i <ca.length; i++) {
			    var c = ca[i];
			    while (c.charAt(0) == ' ') {
			      c = c.substring(1);
			    }
			    if (c.indexOf(name) == 0) {
			      return c.substring(name.length, c.length);
			    }
			  }
			  return "";
			}
			
		`))

		if !production {
			buffer.Write([]byte(`
			var set = function(element, property, value) {
				if (!(element.id in InternalStyleState)) {
					InternalStyleState[element.id] = {};
				}
				element.style[property] = value;
				InternalStyleState[element.id][property] = element.style[property].trim();
			};
			
			var InternalStyleState = {};
			
			
			if (window.location.hostname.includes("localhost")) {
				let url = new URL('/socket', window.location.href);
				url.protocol = url.protocol.replace('http', 'ws');
				let Socket = new WebSocket(url.href);
				Socket.onclose = function() {
					close();
				}
				Socket.onerror = function() {
					close();
				}
				//Disable refresh on chrome because otherwise the app will close.
				document.onkeydown = function() {    
					switch (event.keyCode) { 
						case 116 : //F5 button
							event.returnValue = false;
							event.keyCode = 0;
							return false; 
						case 82 : //R button
							if (event.ctrlKey) { 
								event.returnValue = false; 
								event.keyCode = 0;  
								return false; 
							} 
					}
				}

				function parseCss(attribute) {
					let css = {};
					if (!attribute) return css;
					
					attribute = attribute.replace(/(\/\*([\s\S]*?)\*\/)|(\/\/(.*)$)/gm, '');
											
					//Gonna have to parse the css.
					let styles = attribute.split(';');
					for (let style of styles) {
						if (style == "") continue;
						let splits = style.split(':');
						let property = splits[0];
						let value = splits[1];
						if (value == undefined) continue;
						
						css[property] = value;
					}

					return css;
				}
				
				var edits = {};
				window.addEventListener('load', function() {
					var observer = new MutationObserver(function(mutations) {
						mutations.forEach(function(mutation) {
							if (mutation.target.id == "") return;
							
							let style = parseCss(mutation.target.getAttribute("style"));
								
							for (let property in style) {
								let value = style[property];
							
								if (mutation.target.id in InternalStyleState && InternalStyleState[mutation.target.id][property] == value.trim()) {
									continue;
								}
								
								if (!(mutation.target.id in edits)) {
									edits[mutation.target.id] = {};
								}
								edits[mutation.target.id][property] = true;
							}
							
							//InternalStyleState[mutation.target][]
						});    
					});
	
					const observerConfig = {
					
						attributes: true, // attribute changes will be observed | on add/remove/change attributes
						attributeOldValue: true, // will show oldValue of attribute | on add/remove/change attributes | default: null
						
						characterData: true, // data changes will be observed | on add/remove/change characterData
						characterDataOldValue: true, // will show OldValue of characterData | on add/remove/change characterData | default: null
						
						childList: true, // target childs will be observed | on add/remove
						subtree: true, // target childs will be observed | on attributes/characterData changes if they observed on target
						
						attributeFilter: ['style'] // filter for attributes | array of attributes that should be observed, in this case only style
					
					};
	
					observer.observe(document, observerConfig);
				});
				window.addEventListener("click", function(event) {
					var an = window.getSelection().anchorNode;
				 	// this is the innermost *element*
				 	var element = an;
				 	if (element == null) return;
				 	while (!( element instanceof Element )) {
				    	element = element.parentElement;
				    	if (element == null) return;
				    }

					if (!(element.id in edits)) {
						edits[element.id] = {};
					}
					edits[element.id].text = true;
				});
				
				window.addEventListener("keypress", function(event) {
					//Edit mode.
					if (event.key == "e" && event.ctrlKey) {

						if (document.designMode == "on") {
							document.designMode = "off";
						} else {
							document.designMode = "on";
						}

						event.preventDefault();
						return true;
					}
					//Save Edits.
					if (event.key == "s" && event.ctrlKey) {
	
						for (let edit in edits) {
							
						
							let style = parseCss(get(edit).getAttribute("style"));
							let change = false;

							let message = "#"+edit+" {";

							if (edits[edit].text) {
								message += "text: `+"`"+`"+get(edit).innerHTML+"`+"`"+`;";
								change = true;
							}
							
							for (let property in style) {
								let value = style[property];
								
								if (edit in InternalStyleState && InternalStyleState[edit][property] == value.trim()) {
									continue;
								}

								message += property.trim()+":"+value.trim()+";";
								change = true;
							}
							message += "}";
							if (change) {
								Socket.send(message)								
							}
						}

						let body = document.querySelector("body");
						body.contentEditable = "false";
						event.preventDefault();
						return true;
					}
				})
			} else {
				history.pushState(null, null, document.URL);
				window.addEventListener('popstate', function () {
					back();
					history.pushState(null, null, document.URL);
				});
			}
			`))
		}

		if production {
			buffer.Write([]byte(`
			var set = function(element, property, value) {
				element.style[property] = value;
			};

			history.pushState(null, null, document.URL);
							window.addEventListener('popstate', function () {
								back();
								history.pushState(null, null, document.URL);
							});`))
		}

		buffer.Write(onready)

		var dynamic = seed.BuildDynamicHandler()

		if dynamic != nil {
			buffer.WriteString(`
			var dynamic = new XMLHttpRequest();
	
			dynamic.onreadystatechange = function() {
				if (this.readyState == 4 && this.status == 200) {
					var updates = JSON.parse(this.responseText);
					for (let id in updates) {
						document.getElementById(id).textContent = updates[id];
					}
				}
			};
	
			dynamic.open("GET", "/dynamic", true);
			dynamic.send();`)
		}

		buffer.Write([]byte(`
				</script>
				
				</head><body>
			`))
	buffer.Write(html)
	buffer.Write([]byte(`</body></html>`))

	
	return buffer.Bytes()
}

//TODO random port, can be set with enviromental variables.
func (seed Seed) Launch() error {
	Launcher{Seed: seed}.Launch()
	return nil
}
