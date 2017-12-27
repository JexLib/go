package permission

var (
	styleStr = `
<style>
.bg {
	min-height: 100%;
	background: rgb(24, 26, 26);
	text-align: center;
	font-size: 14px;
}

.login {
	margin: 0 auto;
	max-width: 400px;
	padding: 50px;
	background: #ffffff;
	border-radius: 10px;
	position: relative;
	margin-top: 10%;
	box-shadow: 0px 0px 1px 1px rgba(161, 159, 159, 0.1);
}

span.logo-4 {
	font-weight: 700;
	font-size: 25px;
	color: #1994b8;
}

span.logo-5 {
	font-weight: 300;
	font-size: 25px;
	color: #dadada;
}

.square {
	width: 100%;
	height: 60px;
	display: inline-block;
	vertical-align: middle;
}

.width-100 {
	width: 100%;
}

.login-form .form-group {
	margin-bottom: 10px;
}

.blue-button {
	background-color: #39ADB4;
	border: none;
	color: white;
	border-radius: 2px;
	padding: 10px 30px;
	text-transform: uppercase;
	transition: all 0.3s ease;
	cursor: pointer;
	overflow: visible;
	margin-bottom: 10px;
}

.alert_info {
	background-color: rgba(0, 137, 204, 0.85);
	color: #fff;
	font-size: 18px;
	height: 0;
	left: 0;
	line-height: 50px;
	opacity: 0;
	overflow: hidden;
	position: fixed;
	text-align: center;
	top: 0;
	transition: all 0.3s ease-in-out 0s;
	width: 100%;
	z-index: 22222;
}

.alert_info.in {
	height: 49px;
	opacity: 1;
}

.alert_info.error {
	background-color: rgba(162, 6, 19, 0.85);
}

.alert_info.success {
	background-color: rgba(6, 162, 45, 0.85);
}

.captcha {
	position: relative;
}

.rucaptcha-image {
	position: absolute;
	top: 2px;
	right: 5px;
	width: 82px;
	cursor: pointer;
}
</style>
`

	html_script = `
<script>
var base_fn = {
	alertInfo: function(info) {
		var obj = $(".alert_info");
		obj.text(info).addClass("in");
		setTimeout(function() {
			obj.removeClass("in error success");
		}, 3500);
	},
	alertError: function(info) {
		$(".alert_info").addClass("error");
		this.alertInfo(info);
	},
	alertSuccess: function(info) {
		$(".alert_info").addClass("success");
		this.alertInfo(info);
	}
};

function ajxformSubmit(id) {
	$(id).ajaxSubmit({
	    // target:'body',
		type: 'POST',
		async: false,
		// dataType: 'json',
		success: function(data) {
			base_fn.alertSuccess(data.message);
			canCloseDlg = true;
			if (data.status == 302) {
				location.href = data.location;
			};
			
		},
		error: function(XmlHttpRequest, textStatus, errorThrown) {
			canCloseDlg = false;
			 if (XmlHttpRequest.responseJSON){
			    switch (XmlHttpRequest.responseJSON.status) {
			       case  "E-CaptchaCode":
					      base_fn.alertError(XmlHttpRequest.responseJSON.message);
					      $("#_rucaptcha").val("");
				          refreshCaptchaImg();
				          break;
			       default:
			         base_fn.alertError(XmlHttpRequest.responseJSON.message);
				}
			 }else{
		 		base_fn.alertError(XmlHttpRequest.responseText);
		 	};
	 
			  
		},
		clearForm: false,
		resetForm: false
	});
	return false;
};

//刷新验证码
function  refreshCaptchaImg(){  
	var img = document.getElementById("captchaImg");  
	img.src = "{{.PathCaptchaCode}}?rnd=" + Math.random();  
};  

</script>`

	html_layout = `
<!DOCTYPE html>
<html lang="en">
{{template  "STYLE"}}
{{template  "SCRIPT" .}}
<head>
	<meta charset="UTF-8">
	<title>loging</title>
	<link href=" https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
	<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
	<script src="http://malsup.github.com/min/jquery.form.min.js"></script>
</head>

<body class="bg">
    <div class="alert_info"></div>
	<div class="login">
		<header class="text-center">
			<div class="square"><span class="logo-4">{{.Title}}</span><span class="logo-5">{{.Subtitle}}</span></div>
		</header>
		 {{template "yeld" .}}
	</div>
</body>

</html>
`
	html_sing_in = `
	<style>
    .spinner {
        margin: 10px auto 0;
        width: 90%;
        display: inline-block;
        text-align: center;
        padding-bottom: 10px;
        margin-bottom: 10px;
        border-bottom: 1px solid #e7e8e8;
    }
    
    .spinner .imgBG {
        width: 50px;
        height: 50px;
        background-color: #39ADB4;
        border-radius: 100%;
        line-height: 0px;
        display: inline-block;
        -webkit-animation: bouncedelay 1.4s infinite ease-in-out;
        animation: bouncedelay 1.4s infinite ease-in-out;
        -webkit-animation-fill-mode: both;
        animation-fill-mode: both;
        padding: 5px;
    }
    
    .spinner .bounce {
        -webkit-animation-delay: -0.16s;
        animation-delay: -0.16s;
        display: inline-block;
        width: 50px;
        margin: 0 5px;
        cursor: pointer;
    }
    
    .spinner .imgBG:hover {
        background-color: rgb(28, 78, 73);
    }
</style>
		 <form class="login-form" method="post">
			 <div class="form-group">
				 <div class="input-group">
					 <div class="input-group-addon"><i class="glyphicon glyphicon-user"></i></div>
					 <input placeholder="用户名" name="user" class="form-control" type="text"  required="required">
				 </div>
			 </div>
			 <div class="form-group">
				 <div class="input-group">
					 <div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>
					 <input placeholder="密码" name="paswd" class="form-control" type="password"  required="required">
				 </div>
			 </div>
			 <div class="form-group captcha">
			     <input class="form-control" placeholder="请输入右边验证码" id="_rucaptcha" name="_rucaptcha" type="text" autocorrect="off" autocapitalize="off" pattern="[0-9a-z]*" maxlength="4" autocomplete="off">
				 <a class="rucaptcha-image-box" href="#">
				   <img id="captchaImg" class="rucaptcha-image" src="{{.PathCaptchaCode}}" title="看不清可点击刷新验证码" onclick="refreshCaptchaImg()" alt="验证码">
				 </a>
		     </div>
			 <br>
			

		 </form>
		 <button type="btn" onclick="ajxformSubmit('.login-form');"  class="blue-button width-100">登陆</button>

		 {{if ne (len .Oauth2Clients) 0}}
		 <div class="spinner">
		 <p  style="color: #979696;">使用第三方登录</p>
		 {{range $k, $v := .Oauth2Clients}}
		   <div class="bounce" onclick="window.open('?oauth={{$k}}','_self')">
		      <div class="imgBG">
			     <img src="{{$v.ImgUrl}}" height="40px">
		      </div>
		      <a>{{$k}}</a>
	       </div>
	     {{end}}
		 <!-- <a class="bounce " href="/auth/github"></a>
		 <a class="bounce github" href="/auth/github"></a>
		 <a class="bounce " href="/auth/github"></a> -->
	    </div>
		{{end}}
		 <p>还没有注册用户吗？ <strong><a href="{{.Sign_up}}">现在注册！</a></strong></p>
`

	html_sing_up = `
		 <form class="login-form" method="post">
		 <div class="form-group">
		 <div class="input-group">
			 <div class="input-group-addon"><i class="glyphicon glyphicon-user"></i></div>	        		
			   <input placeholder="用户名"  name="user" class="form-control" type="text"  required="required">           
		   </div>	
	 </div>
	 <div class="form-group">
		 <div class="input-group">
			 <div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>	        		
			   <input    placeholder="密码" name="paswd" class="form-control" type="password"  required="required">           
		   </div>	
	 </div>
	 <div class="form-group">
		 <div class="input-group">
			 <div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>	        		
			   <input  placeholder="确认密码" name="paswd1" class="form-control" type="password"  required="required">           
		   </div>	
	 </div>
	 <div class="form-group captcha">
	 <input class="form-control" placeholder="请输入右边验证码" id="_rucaptcha" name="_rucaptcha" type="text" autocorrect="off" autocapitalize="off" pattern="[0-9a-z]*" maxlength="4" autocomplete="off">
	 <a class="rucaptcha-image-box" href="#">
	   <img id="captchaImg" class="rucaptcha-image" src="{{.PathCaptchaCode}}" title="看不清可点击刷新验证码" onclick="refreshCaptchaImg()" alt="验证码">
	 </a>
 </div>
	 <br>
	 </form>
	 <button type="btn" onclick="ajxformSubmit('.login-form');" class="blue-button width-100">注册</button>
		 <p>已有账号？ <strong><a href="{{.Sign_in}}">现在登陆！</a></strong></p>
`

	html_changepwd = `
<form class="login-form" method="post">
<div class="form-group">
<div class="input-group">
	<div class="input-group-addon"><i class="glyphicon glyphicon-user"></i></div>	        		
	  <input placeholder="用户名"  name="user" class="form-control" type="text"  required="required">           
  </div>	
</div>
<div class="form-group">
<div class="input-group">
	<div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>	        		
	  <input    placeholder="密码" name="paswd" class="form-control" type="password"  required="required">           
  </div>	
</div>
<div class="form-group">
<div class="input-group">
	<div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>	        		
	  <input    placeholder="新密码" name="newpaswd" class="form-control" type="password"  required="required">           
  </div>	
</div>
<div class="form-group">
<div class="input-group">
	<div class="input-group-addon"><i class="glyphicon glyphicon-lock"></i></div>	        		
	  <input  placeholder="确认新密码" name="newpaswd1" class="form-control" type="password"  required="required">           
  </div>	
</div>
<br>
</form>
<button type="button" onclick="ajxformSubmit('.login-form');" class="blue-button width-100">修改密码</button>
<p>放弃修改密码 ，返回 <strong><a href="/">主页</a></strong></p>
`

	html_logout = `
	<!DOCTYPE html>
	<html>
	<style>
		.bg {
			min-height: 100%;
			background: rgb(24, 26, 26);
			text-align: center;
			font-size: 14px;
		}
		
		.box {
			margin: 0 auto;
			max-width: 500px;
			padding-top: 30px;
			background: #ffffff;
			border-radius: 10px;
			position: relative;
			margin-top: 10%;
			box-shadow: 0px 0px 1px 1px rgba(161, 159, 159, 0.1);
		}

		.alert_info {
			background-color: rgba(0, 137, 204, 0.85);
			color: #fff;
			font-size: 18px;
			height: 0;
			left: 0;
			line-height: 50px;
			opacity: 0;
			overflow: hidden;
			position: fixed;
			text-align: center;
			top: 0;
			transition: all 0.3s ease-in-out 0s;
			width: 100%;
			z-index: 22222;
		}
		
		.alert_info.in {
			height: 49px;
			opacity: 1;
		}
		
		.alert_info.error {
			background-color: rgba(162, 6, 19, 0.85);
		}
		
		.alert_info.success {
			background-color: rgba(6, 162, 45, 0.85);
		}
	</style>
 
	<head>
		<meta charset="UTF-8">
	
	</head>
	<link href=" https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
	<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
	<script src="http://malsup.github.com/min/jquery.form.min.js"></script>
	
	<body class="bg">
	<form class="lform box" method="post">
	     <p>注销后，访问权限将受到限制&hellip;
				<p>
					<strong>确认注销登录？</p> 
					<div class="modal-footer">
					  <button type="button"  onclick="javascript :history.back(-1)" class="btn btn-default">取消</button>
					  <button type="button" onclick="ajxformSubmit('.lform');" class="btn btn-warning">确认</button>
					</div> 
	 
	 <form>
	 </body>
	 
	 </html>
`
)
