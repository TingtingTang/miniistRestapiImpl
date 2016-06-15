/* =============================================================================================
 * >> mini-IST <<
 * 
 * jQuery for sign-in(login), sign-up(register), forget password, etc 
 * ============================================================================================= */

$(document).ready(function() {
	
	//*--==page switch: login or register==--
	//--to register page
	//alert("found?? "+ $(".curtain.miniist-sign a.btn.to-register").length);
	
	$(".curtain.mist-sign a.btn.to-register").click(function(e) {
		//alert("hey! to-reg..");
		e.preventDefault();
		$(this).closest(".one-page").addClass("hidden").removeClass("active")
			.siblings(".one-page.register").removeClass("hidden").addClass("active");
	});
	//--to login page
	$(".curtain.mist-sign a.btn.to-login").click(function(e) {
		//alert("hey! to login..");
		e.preventDefault();
		$(this).closest(".one-page").addClass("hidden").removeClass("active")
			.siblings(".one-page.login").removeClass("hidden").addClass("active");
	});
	
	//*--==login go==--
	//--scroll up the curtain for displaying the dashboard
	
	//*--@ register click go
	$(".register.click-go").click(function(e) {
		e.preventDefault();
		//--some additional validation, ajax for REST API ...
		//--ajax for user join
		//*--REST API /api/join
		/*
		 * {
	   			action: "join",
	   			name: "mike",
	   			email: "xx@yy.com",
	   			password: "xxxx"  
			}
		 */
		
		var $reg = $(".mist-sign .one-page.register"),
			username = $("#signupInputUsername").val(),
			email = $("#signupInputEmail1").val(),
			pass1 = $("#signupInputPassword1").val(),
			pass2 = $("#signupInputRePassword1").val(),
			url = "/join",
			params = {
				"action": "join",
				"name": username,
				"email": email,
				"password": pass1
			};
		
		if (pass1 != pass2) {
			//alert("Wrong password setting! Two times not the same!...");
			$.notify("[SIGN001] Wrong password setting! Two times not the same! \n", 
				{
					status: 'danger'
				});
			return false;
		}
		
		//debug:
		/*alert("Calling REST API: \n"+
			"url:      "+url+'\n'+
			"username: "+params.name+'\n'+
			"email:    "+params.email+'\n'+
			"pass:     "+params.password
		);*/
		
		//ajax
		$.ajax({
			url: url,
			type: 'POST',
			data: JSON.stringify(params),
			dataType: "json",
			success: function(data) {
				//alert("data.success? "+data.success);
				/*alert('REST API, Good! \n'+
					'ret:  ' +data.ret+'\n'+
					'info: ' +data.info+'\n'+
					'name: ' +data.name+'\n'+
					'email: '+data.email+'\n'+
					'type: '+data.type+'\n'+
					'create_time: '+data.create_time);*/
				
				$.notify("Registration Complete! Sign in to continue! \n"+
					'ret:  ' +data.ret+'\n'+
					'info: ' +data.info+'\n'+
					'name: ' +data.name+'\n'+
					'email: '+data.email+'\n'+
					'type: '+data.type+'\n'+
					'create_time: '+data.create_time
					, {
						status: 'success',
						pos: 'bottom-right'
					});
					
				//to register page
				$(".curtain.mist-sign a.btn.to-login").click();
			},
			error: function(jqXHR, textStatus, errorThrown) {
				/*alert('AJAX failed! \n'+
					'XHR:  ' +jqXHR.readyState+'\n'+
					'XHR:  ' +jqXHR.status+'\n'+
					'XHR:  ' +jqXHR.statusText+'\n'+
					'XHR:  ' +jqXHR.responseText+'\n'+
					'textStatus: ' +textStatus+'\n'+
					'errorThrown: '+errorThrown+'\n');*/
					
				$.notify("[SIGN002] Register AJAX Failed! Sorry! \n"+
					'XHR:  ' +jqXHR.readyState+'\n'+
					'XHR:  ' +jqXHR.status+'\n'+
					'XHR:  ' +jqXHR.statusText+'\n'+
					'XHR:  ' +jqXHR.responseText+'\n'+
					'textStatus: ' +textStatus+'\n'+
					'errorThrown: '+errorThrown+'\n', 
					{
						status: 'danger'
					});
			}
		});
		
	});
	//*--& register click go
	
	//*--@ login click go
	$(".login.click-go").click(function() {
		//--some additional validation, ajax for REST API ...
		//--ajax for user join
		//*--REST API /api/join
		/*
		 * 请求路径	/api/login
		 *
		 * {"action": "login", "name": "tester1", "password": "passw0rd"}
		 * 
		 * 响应数据
		 *
		 * {
		 * 		"id":"56f3ec7f58f5798e3c126315",
		 * 		"name":"tester1",
		 * 		"email":"tester1@gmail.com",
		 * 		"create_time":"2016-03-24T21:32:47.   462+08:00",
		 * 		"ret":0,
		 * 		"info":"OK",
		 * 		"Token":"..."
		 * }
		 * */
		
		var $login = $(".mist-sign .one-page.login"),
			userinput = $("#loginUserInput").val(),
			pass1 = $("#loginInputPassword1").val(),
			url = "/login",params;
			//params = {
			//	"action": "login",
			//	"name": userinput,
			//	"password": pass1
			//};

		if(((/.+@.+\.[a-zA-Z]{2,4}$/).test(userinput))){
			params = {
				"action": "login",
				"email": userinput,
				"password": pass1
			};
		}else{
			params = {
				"action": "login",
				"name": userinput,
				"password": pass1
			}
		}
		
		//check user input: Email or Username
		/*if (userinput.indexOf('@') > 0) {
			alert("Calling REST API: \n"+ "-----Email Input----- \n"+
				"url:  "+url+'\n'+
				"name: "+params.name+'\n'+
				"pass: "+params.password
			);
		} else {
			alert("Calling REST API: \n"+ "-----UserName Input----- \n"+
				"url:  "+url+'\n'+
				"name: "+params.name+'\n'+
				"pass: "+params.password
			);
		}*/	
		
		//ajax
		$.ajax({
			url: url,
			type: 'POST',
			data: JSON.stringify(params),
			dataType: "json",
			success: function(data) {
				//alert("ajax success...?!");
				if (data.ret == 0) {
					/*alert('REST API, Good! \n'+
						'ret:   ' +data.ret+'\n'+
						'info:  ' +data.info+'\n'+
						'email: ' +data.email+'\n'+
						'name:  ' +data.name+'\n'+
						//'token: '+data.Token+'\n'+
						'create_time: '+data.create_time);*/
					
					//try angular notify
					//Message("This is anglar notify! Login Good!");
					$.notify("Welcome! "+data.name, {
						status: 'success',
						pos: 'bottom-right'
					});
					
					//clear previous input
					$("#loginUserInput").val("");
					$("#loginInputPassword1").val("");
					
					//save token on HTML page
					$("#userToken").empty().text(data.Token);
					$("#userLoginName").empty().text(data.name);
					//alert("User Token Saved: \n\n\n"+$("#userToken").text()+'\n\n\n');
					
					//set username and avatar (if any) on topnav bar
					//var myapp = angular.module('angle');
					var $body = angular.element(document.body);   // 1
					var $rootScope = $body.scope().$root;         // 2
					$rootScope.$apply(function () {               // 3
					    $rootScope.user.name = data.name;
					});
					
					//testing only
					$(".curtain.mist-sign").animate({
			    		opacity: 0,
			    		top: '-=300',
			    		zIndex: "-50"
			  		}, 500);
				} else {
					/*alert('REST API, ERROR! \n'+
					'ret:  ' +data.ret+'\n'+
					'info: ' +data.info+'\n'+
					'name: ' +data.name+'\n'+
					'token: '+data.token+'\n'+
					'type: '+data.type+'\n'+
					'create_time: '+data.create_time);*/
					
					$.notify("[SIGN003] Login Failed! Sorry! \n"+ 
					'ret:  ' +data.ret+'\n'+
					'info: ' +data.info+'\n',
					{
						status: 'danger'
					});
				}
					
			},
			error: function(jqXHR, textStatus, errorThrown) {
				/*alert('AJAX failed! \n'+
					'XHR:  ' +jqXHR.readyState+'\n'+
					'XHR:  ' +jqXHR.status+'\n'+
					'XHR:  ' +jqXHR.statusText+'\n'+
					'XHR:  ' +jqXHR.responseText+'\n'+
					'textStatus: ' +textStatus+'\n'+
					'errorThrown: '+errorThrown+'\n');
					//'name: ' +data.name+'\n'+
					//'token: '+data.token+'\n'+
					//'type: '+data.type+'\n'+
					//'create_time: '+data.create_time);*/
				
				$.notify("[SIGN004] Login AJAX Failed! Sorry! \n"+
					'XHR:  ' +jqXHR.readyState+'\n'+
					'XHR:  ' +jqXHR.status+'\n'+
					'XHR:  ' +jqXHR.statusText+'\n'+
					'XHR:  ' +jqXHR.responseText+'\n'+
					'textStatus: ' +textStatus+'\n'+
					'errorThrown: '+errorThrown+'\n', 
					{
						status: 'danger'
					});
			}
		});	
	});
	//--& login click go
	
	//*--@ logout click go
	$(".logout.click-go").click(function() {
		/*
		 * request head Authorization: <user token>
		 *	data: {"action": "logout", "name": "tester1"} 
		 */
		 
		//alert("clicking logout..");
		var username = $("#userLoginName").text(),
			usertoken = $("#userToken").text(),
			url = "/logout",
			params = {
				"action": "logout",
				"name": username
			};
		
		//if no username to get
		if (username == "" || username == null || username == undefined) {
			$.notify("[SIGN005] Logout Issue! \n"+
				"Can't get the right user name for logout!", 
				{
					status: 'danger'
				});
			
			return false;
		}
		
		//if no user TOKEN to get
		if (usertoken == "" || usertoken == null || usertoken == undefined) {
			$.notify("[SIGN006] Logout Issue! \n"+
				"Can't get the right user TOKEN for logout!", 
				{
					status: 'danger'
				});
			
			return false;
		}
		
		//AJAX for logout rest API with token
		$.ajax({
			url: url,
			type: 'POST',
			data: JSON.stringify(params),
			dataType: "json",
			beforeSend:function(xhr){
				xhr.setRequestHeader('Authorization', usertoken);
			},
			success: function(data) {
				if (data.ret == 0) {
					$(".curtain.mist-sign").css({"zIndex": "180"}).animate({
			    		opacity: 1,
			    		top: '+=300',
			  		}, 1000, function() {
			  			$.notify("Savely Logged Out! See you next time!", {
							status: 'success',
							pos: 'bottom-right'
						});
			  		});
				} else {
					/*alert('REST API, ERROR! \n'+
					'ret:  ' +data.ret+'\n'+
					'info: ' +data.info+'\n'+
					'name: ' +data.name+'\n'+
					'token: '+data.token+'\n'+
					'type: '+data.type+'\n'+
					'create_time: '+data.create_time);*/
					
					$.notify("[SIGN007] Logout Failed! Sorry! \n"+ 
						'ret:  ' +data.ret+'\n'+
						'info: ' +data.info+'\n',
						{
							status: 'danger'
						});
				}
					
			},
			error: function(jqXHR, textStatus, errorThrown) {
				$.notify("[SIGN008] Logout AJAX Failed! Sorry! \n"+
					'XHR:  ' +jqXHR.readyState+'\n'+
					'XHR:  ' +jqXHR.status+'\n'+
					'XHR:  ' +jqXHR.statusText+'\n'+
					'XHR:  ' +jqXHR.responseText+'\n'+
					'textStatus: ' +textStatus+'\n'+
					'errorThrown: '+errorThrown+'\n', 
					{
						status: 'danger'
					});
			}
		});	
		
	});
	//*--& logout click go
	
});

