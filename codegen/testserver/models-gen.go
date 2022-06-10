// Code generated by github.com/tinhtran24/gqlgen, DO NOT EDIT.

package testserver

type InnerInput struct {
	ID int `json:"id"`
}

type InnerObject struct {
	ID int `json:"id"`
}

type Keywords struct {
	Break       string `json:"break"`
	Default     string `json:"default"`
	Func        string `json:"func"`
	Interface   string `json:"interface"`
	Select      string `json:"select"`
	Case        string `json:"case"`
	Defer       string `json:"defer"`
	Go          string `json:"go"`
	Map         string `json:"map"`
	Struct      string `json:"struct"`
	Chan        string `json:"chan"`
	Else        string `json:"else"`
	Goto        string `json:"goto"`
	Package     string `json:"package"`
	Switch      string `json:"switch"`
	Const       string `json:"const"`
	Fallthrough string `json:"fallthrough"`
	If          string `json:"if"`
	Range       string `json:"range"`
	Type        string `json:"type"`
	Continue    string `json:"continue"`
	For         string `json:"for"`
	Import      string `json:"import"`
	Return      string `json:"return"`
	Var         string `json:"var"`
}

type OuterInput struct {
	Inner InnerInput `json:"inner"`
}

type OuterObject struct {
	Inner InnerObject `json:"inner"`
}

type User struct {
	ID      int    `json:"id"`
	Friends []User `json:"friends"`
}
