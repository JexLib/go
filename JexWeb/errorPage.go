package jexweb

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
)

var (
	errorPage_Layout = `
	<!DOCTYPE html>
	<html>
	<style>
		@import url(https://fonts.googleapis.com/css?family=Open+Sans:700);
		body,
		div,
		h1,
		h2 {
			margin: 0;
			padding: 0;
			border: 0;
			outline: 0;
			font-size: 100%;
			vertical-align: baseline;
			background: 0 0;
		}
		
		body {
			line-height: 1;
		}
		
		* {
			box-sizing: border-box;
		}
		
		html {
			background-color: #f06060;
			color: #fff;
			font-size: 16px;
			font-family: 'Open Sans', sans-serif;
		}
		
		@media screen and (max-width:600px) {
			html {
				font-size: 10px;
			}
		}
		
		#container {
			width: calc(100% - 100px);
			max-width: 700px;
			min-width: 300px;
			margin: 0 auto;
			padding: 50px;
		}
		
		h1 {
			font-size: 16rem;
		}
		
		h2 {
			font-size: 6rem;
			color: #f48f8f;
		}


	</style>
	
	<head>
		<meta charset="UTF-8">
		<title>title | 404</title>
	</head>
	
	<body>
		<div id="container">
			 {{template  "yeld" .}}
 
		</div>
		{{if .Devel }}
		<table width="80%" border="1" bgcolor="#000000" cellpadding="10px" style="margin-left:10%">
		  <tr >
			<th>{{.Title}}</th>
		  </tr>
		  <tr>
			<td>{{.Message}}</td>
		  </tr>
		</table>
		
	  {{end}}
		
	   
	</body>
	
	</html>
	`

	errorPages = map[int]string{
		404: `
		<h1>404</h1>
		<h2>Not found <span>:(</span></h1>
			<p>Sorry, but the page you were trying to view does not exist.</p>
			<p>It looks like this was the result of either:</p>
		`,
		500: `
		<h1>500</h1>
		<!-- <h2>Internal Server Error</h2> -->
		<h2>Sorry,服务器内部错误，<span>:(</span></h1>
			<p>Sorry, 服务器内部发生错误,不能执行该请求.</p>
			<p>请稍后刷新重试，或者与网站管理员联系.:</p>
		`,
	}
)

func JexHTTPErrorHandler(err error, c echo.Context) {
	data := map[string]interface{}{
		"Devel": jwconfig.IsDevelopment,
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		data["Title"] = he.Error()
		data["Message"] = he.Message
	} else {
		data["Title"] = err.Error()
		data["Message"] = ""
		fmt.Println("errr:", err.Error())
	}

	tmpl, _ := template.New("Layout").Parse(errorPage_Layout)
	tmpl.New("yeld").Parse(errorPages[code])
	out := new(bytes.Buffer)
	err = tmpl.Execute(out, data)
	if err != nil {
		//panic(err)
		return
	}
	//return c.HTML(200, out.String())

	//errorPage := strings.Replace(errorPage_Layout, "{{.yeld}}", errorPages[code], 1)
	if err := c.HTML(code, out.String()); err != nil {
		c.Logger().Error(err)
	}
	fmt.Println()
	c.Logger().Error(err)

}
