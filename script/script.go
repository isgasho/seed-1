package script

import (
	"github.com/qlova/seed/internal"
	"github.com/qlova/seed/style/css"

	qlova "github.com/qlova/script"
	"github.com/qlova/script/language"
)

//Ctx is a script context. Providing access to script behaviours.
type Ctx struct {
	*ctx
}

type AnyCtx = qlova.AnyCtx

func CtxFromAnyCtx(any AnyCtx) Ctx {
	if q, ok := any.(Ctx); ok {
		return q
	}
	var context = internal.NewContext()
	return Ctx{&ctx{
		Ctx:     any.RootCtx(),
		Context: context,
	}}
}

type ctx struct {
	internal.Context
	qlova.Ctx

	js   js
	Time time
}

//Require inserts the provided dependency string in the head of the document.
func (q Ctx) Require(dependency string) {

	//Subdependencies.
	if dependency == Goto {
		q.Require(Get)
		q.Require(Set)
	}

	if _, ok := q.Dependencies[dependency]; ok {
		return
	}
	q.Dependencies[dependency] = struct{}{}
}

func (q Ctx) wrap(s string) qlova.String {
	return String{language.Expression(q, s)}
}

//ToJavascript returns the given script encoded as Javascript.
func ToJavascript(f func(q Ctx), ctx ...internal.Context) []byte {
	if f == nil {
		return nil
	}

	var context internal.Context
	if len(ctx) > 0 {
		context = ctx[0]
	} else {
		context = internal.NewContext()
	}

	return toJavascript(f, context)
}

func toJavascript(f func(q Ctx), context internal.Context) []byte {
	var program = qlova.Script(func(q qlova.Ctx) {
		var s = Ctx{&ctx{
			Ctx:     q,
			Context: context,
		}}
		s.js.q = s
		s.Time.Ctx = s
		//s.Go.Script = s
		f(s)
	})

	source := language.Javascript(program)

	return source
}

//Run runs a Javascript function with the given arguments.
func (q Ctx) Run(f Function, args ...qlova.Type) {
	q.Javascript(string(f) + "();")
}

//Unit is a display unit, eg. px, em, %
type Unit qlova.String

//Unit returns a script.Unit from the given unit.
func (q Ctx) Unit(unit complex128) Unit {
	return Unit(String{language.Expression(q, string(css.Decode(unit)))})
}

//SetClipboard is the JS code requried for Clipboard support.
const SetClipboard = `
	const setClipboard = str => {
		const el = document.createElement('textarea');
		el.value = str;
		el.setAttribute('readonly', '');
		el.style.position = 'absolute';
		el.style.left = '-9999px';
		document.body.appendChild(el);
		const selected =
			document.getSelection().rangeCount > 0 ? document.getSelection().getRangeAt(0) : false;
		el.select();
		document.execCommand('copy');
		document.body.removeChild(el);
		if (selected) {
			document.getSelection().removeAllRanges();
			document.getSelection().addRange(selected);
		}
	};
`

//SetClipboard sets the clipboard to the provided string.
func (q Ctx) SetClipboard(text String) {
	q.Require(SetClipboard)
	q.js.Run(`setClipboard`, text)
}
