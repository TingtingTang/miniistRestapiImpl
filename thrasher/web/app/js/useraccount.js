/**
 * Created by ibmBSM on 5/12/2016.
 */
$(function() {
    var $username_account = $("#username_account");
    var $oldpassword_account = $("#oldpassword_account");
    var $newpassword_account = $("#newpassword_account");
    var $confirmpassword_account = $("#confirmpassword_account");
    var $changepassword_account = $("#changepassword_account");
    //var $deleteaccount_account = $("#deleteaccount-account");
    var $okdelete = $("#okdelete");
    var $deleteaccount_name = $("#deleteaccount_name");
    var $deleteaccount_password = $("#deleteaccount_password");
    var $deleteaccount_type = $("#deleteaccount_type");
    //validate the format


    //change password
    $username_account.val(JSON.parse(localStorage.getItem("ls.username")));
    var userToken = localStorage.getItem("ls.token");

    // ajax for changePassword
    function changePassword(){
        var paramsChangepass = {
            "action": "change_password",
            "name": $username_account.val(),
            "password": $oldpassword_account.val(),
            "new_password":$newpassword_account.val()
        };
        $.ajax({
            url: "/user/password",
            type: "post",
            data: JSON.stringify(paramsChangepass),
            dataType: "json",
            beforeSend:function(xhr){
                xhr.setRequestHeader("Authorization", JSON.parse(userToken));
                //xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
            },
            success: function(data){
                if(data.ret == 0 && data.info == "OK"){
                    $("#changepassModal").modal("hide");
                    localStorage.setItem("ls.name", "");
                    localStorage.setItem("ls.token", "");
                    localStorage.setItem("ls.moment", "");
                    localStorage.setItem("ls.username", "");
                    $.notify('Change password succeed!'+"/n"+"But please login again!"
                        , {
                            title: "Success",
                            status: 'success'
                            //pos: 'bottom-right'
                        });
                    //window.location.replace("/");
                    setTimeout(function(){
                        window.location.replace("/");
                    }, 2000);

                }
            },
            error:function(){
                $("#changepassModal").modal("hide");
                var errorMsg = 'The wrong password.';
                $($oldpassword_account).attr("data-original-title", errorMsg);
                $($oldpassword_account).tooltip();
                $($oldpassword_account).toggleClass("onError", true);
                $oldpassword_account.val("");
                $newpassword_account.val("");
                $confirmpassword_account.val("");
                setTimeout(function(){
                    $.notify('The original password is wrong, please input again!'
                        , {
                            title: "Error",
                            status: 'danger'
                            //pos: 'bottom-right'
                        });
                }, 100);
            }
        })
    }
    //validate change password
    $("form#account-form input").blur(function (e) {
        //var okMsg = 'It' + "'" + 's right.';
        //  validate the page of login.html with username
        if ($(this).is('#username_account')) {
            var errorMsg = 'The wrong username.';
            if (this.value != JSON.parse(localStorage.getItem("ls.username"))) {
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                $(this).toggleClass("onError", true);
            } else {
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
        //  validate the page of register.html with username
        if($(this).is('#oldpassword_account')){
            var errorMsg = 'At least SIX please.';
            if(this.value == "" || this.value == null || this.value.length < 6){
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                $(this).toggleClass("onError", true);
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
        //  validate the page of register.html with email
        if($(this).is('#newpassword_account')){
            var errorMsg = "At least SIX please.";
            if(this.value == "" || this.value == null || this.value.length < 6){
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                var thisClass = $(this).attr("class");
                $(this).toggleClass("onError", true);
                //$(this).attr("class", thisClass+" onError");
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
        //      Retry the password
        if($(this).is('#confirmpassword_account')){
            if(this.value == "" || this.value == null || this.value !== $("#newpassword_account").val()){
                var errorMsg= 'Retry password is not right.';
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                var thisClass = $(this).attr("class");
                $(this).toggleClass("onError", true);
                //$(this).attr("class", thisClass+" onError");
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
    }).keyup(function () {
        $(this).triggerHandler("blur");
    }).focus(function () {
        $(this).triggerHandler("blur");
    });//end blur
    $changepassword_account.bind("click", function(){
        $("form#account-form input").trigger('blur');
        var numError = $("form#account-form .onError").length;
        if(numError){
            $.notify('Input error! Please enter again.'
                , {
                    title: "Error",
                    status: 'danger'
                    //pos: 'bottom-right'
                });
            return false;
        }else{
            if($("#oldpassword_account").val() == $("#newpassword_account").val()){
                $.notify('The new password and old password cannot be same .'
                    , {
                        title: "Error",
                        status: 'danger'
                        //pos: 'bottom-right'
                    });
                return false;
            }else{
                $(this).attr("data-target", "#changepassModal");
                $("#okChange").bind("click", function(){
                    changePassword();
                });
            }
        }
    });

    //validate delete account
    $("#deleteaccount_name").val(JSON.parse(localStorage.getItem("ls.username")));
    $("div#deleteaccountModal form input").blur(function(e){
        if($(this).is("#deleteaccount_name")){
            var errorMsg = "The wrong username or email.";
            if(this.value != JSON.parse(localStorage.getItem("ls.username"))){
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                $(this).toggleClass("onError", true);
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
        if($(this).is("#deleteaccount_type")){
            var errorMsg = "Please enter 'delete my account'";
            if(this.value != "delete my account"){
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                $(this).toggleClass("onError", true);
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
        if($(this).is('#deleteaccount_password')){
            var errorMsg = 'At least SIX please.';
            if(this.value == "" || this.value == null || this.value.length < 6){
                $(this).attr("data-original-title", errorMsg);
                $(this).tooltip();
                $(this).toggleClass("onError", true);
            }else{
                $(this).removeAttr("data-original-title");
                $(this).toggleClass("onError", false);
            }
        }
    }).keyup(function(){
        $(this).triggerHandler("blur");
    }).focus(function(){
        $(this).triggerHandler("blur");
    });

    $okdelete.bind("click", function(){
        $("div#deleteaccountModal form input").trigger("blur");
        var numError = $("div#deleteaccountModal form .onError").length;
        if(numError){
            return false;
        }else{
            $(this).attr("data-target", "#deleteaccountModal");
            deleteAccount();
        }
    });
    function deleteAccount(){
        var paramsDeleteaccount = {
            "action": "delete_user",
            "name": $deleteaccount_name.val(),
            "password": $deleteaccount_password.val()
        };
        $.ajax({
            url: "/user",
            type: "delete",
            data: JSON.stringify(paramsDeleteaccount),
            dataType: "json",
            beforeSend:function(xhr){
                xhr.setRequestHeader("Authorization", JSON.parse(userToken));
            },
            success: function(data){
                if(data.ret == 0 && data.info == "OK"){
                    localStorage.setItem("ls.name", "");
                    localStorage.setItem("ls.token", "");
                    localStorage.setItem("ls.moment", "");
                    localStorage.setItem("ls.username", "");
                    window.location.replace("/");
                    $deleteaccount_type.val("");
                }
            },
            error:function(){
                var errorMsg = 'The wrong password.';
                $($deleteaccount_password).attr("data-original-title", errorMsg);
                $($deleteaccount_password).tooltip();
                $($deleteaccount_password).toggleClass("onError", true);
                $deleteaccount_password.val("");
                $($okdelete).attr("data-toggle", "popover");
                $($okdelete).attr("data-content", "Enter the password is wrong, please input again");
                $($okdelete).popover("show");
                setTimeout(function(){
                    $($okdelete).trigger("click");
                    $($okdelete).popover("destroy");
                }, 2000);
            }
        })
    }
});
