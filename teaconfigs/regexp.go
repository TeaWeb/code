package teaconfigs

import "regexp"

var RegexpExternalURL = regexp.MustCompile("(?i)^(http|https|ftp)://")
var RegexpDigitVariable = regexp.MustCompile("\\${\\d+}")
var RegexpNamedVariable = regexp.MustCompile("\\${[\\w.-]+}")
var RegexpDigitNumber = regexp.MustCompile("^\\d+$")
