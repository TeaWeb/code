package forms

import (
	"reflect"
)

type Group struct {
	Namespace string             `yaml:"namespace" json:"namespace"`
	Elements  []ElementInterface `yaml:"elements" json:"elements"`
	HTML      string             `yaml:"html" json:"html"`
}

func (this *Group) Add(element ElementInterface) {
	s := element.Super()
	s.Namespace = this.Namespace
	s.ClassType = reflect.TypeOf(element).Elem().Name()
	this.Elements = append(this.Elements, element)
}

func (this *Group) Compose() {
	this.HTML = ""
	for _, e := range this.Elements {
		element := e.Super()

		this.HTML += "<tr>\n"
		this.HTML += "	<td class=\"title\">" + element.Title
		if len(element.Subtitle) > 0 {
			this.HTML += " (" + element.Subtitle + ")"
		}
		if element.IsRequired {
			this.HTML += " *"
		}
		this.HTML += "</td>\n"
		this.HTML += "	<td>\n"
		this.HTML += e.Compose()

		if len(element.Comment) > 0 {
			this.HTML += "\n<p class=\"comment\">" + element.Comment + "</p>"
		}

		this.HTML += "\n</td>"
		this.HTML += "</tr>"
	}
}
