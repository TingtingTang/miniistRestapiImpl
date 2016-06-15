/**
 * Created by ibmBSM on 4/25/2016.
 */

$(function(){

    //var timeT = 1; //add the hook to prevent the repeat notify.
    //set portrait
    function getFileName(o){
        var pos=o.lastIndexOf("\\");
        return o.substring(pos+1);
    }
    //hidden user poritrait
    var portraithide = function($x, $y, $z){
        $y.hide();
        $x.click(function(){
            $y.trigger("click");

            function getObjectURL(file) {
                var url = null ;
                if (window.createObjectURL!=undefined) { // basic
                    url = window.createObjectURL(file) ;
                } else if (window.URL!=undefined) { // mozilla(firefox)
                    url = window.URL.createObjectURL(file) ;
                } else if (window.webkitURL!=undefined) { // webkit or chrome
                    url = window.webkitURL.createObjectURL(file) ;
                }
                return url ;
            }
            $y.change(function() {
                var eImg = "url("+getObjectURL($(this)[0].files[0])+")";
                $z.css("background-image",eImg); // �� this.files[0] this->input
                $x.trigger("blur");
                localStorage.setItem("portraitpicture", eImg);
            });
        });
    };
    var $portrait = $("#portrait");
    var $portraits = $("#portraitInputFileS");
    var $protraith = $("#portraitInputFileH");
    //if(!localStorage.getItem("portraitpicture")){
    //    $portrait.css("background-image","url("+"../img/port2.png"+")");
    //}else{
    //    var portraitstorage = localStorage.getItem("portraitpicture");
    //    $portrait.css("background-image",portraitstorage);
    //}
    portraithide($portraits, $protraith, $portrait);

    //testSection:login to save token in localstorage
    var token;
    var $update = $("#update");
    var userToken = localStorage.getItem("ls.token");
    var $profilename, $profiletsouser, $profiletsopassword, $profileteam, $profileemail, $profilecreatetime, $profileid, $profiletype;
    $profilename = $("#profilename");
    $profileid = $("#profileid");
    $profiletype = $("#profiletype");
    $profiletsouser = $("#profiletsouser");
    $profiletsopassword = $("#profiletsopassword");
    $profileteam = $("#profileteam");
    $profileemail = $("#profileemail");
    $profilecreatetime = $("#profilecreatetime");
    $profilename.val(JSON.parse(localStorage.getItem("ls.username")));
    //var profileName = localStorage.getItem("ls.username");
    var $getProfile = $(".getProfile");
    $($getProfile).hide();
    function get_user(){
        $($getProfile).show();
        var $paramsGetuser;
        $paramsGetuser = {
            "action" : "get_user",
            "name": $profilename.val()
        };
        $.ajax({
            url:"/user",
            type:"get",
            data: $paramsGetuser,
            dataType:"json",
            beforeSend:function(xhr){
                xhr.setRequestHeader("Authorization", JSON.parse(userToken));
                //xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
            },
            success: function(data){
                $profileid.val(data.id);
                $profileemail.val(data.email);
                $profiletype.val(data.type);
                $profileteam.val(data.tso_user);
                $profiletsopassword.val(data.tso_password);
                $profilecreatetime.val(data.create_time);
                $profiletsouser.val(data.tso_user);
                sessionStorage.setItem("ls.id", data.id);
                sessionStorage.setItem("ls.email", data.email);
                sessionStorage.setItem("ls.type", data.type);
                sessionStorage.setItem("ls.tso_user", data.tso_user);
                sessionStorage.setItem("ls.tso_password", data.tso_password);
                sessionStorage.setItem("ls.create_time", data.create_time);
                sessionStorage.setItem("ls.tso_user", data.tso_user);
                $($getProfile).hide();
            //    user $.notify
            //    $.notify("Get information successful"
            //        , {
            //            status: 'success'
            //            //pos: 'bottom-right'
            //        });
            },
            error: function(){
                $.notify("Get information failed!"
                    , {
                        status: 'danger'
                        //pos: 'bottom-right'
                    });

            }
        })

    }
    //get_user();
    if(sessionStorage.getItem("ls.id")){
        $profileid.val(sessionStorage.getItem("ls.id"));
        $profileemail.val(sessionStorage.getItem("ls.email"));
        $profiletype.val(sessionStorage.getItem("ls.type"));
        $profileteam.val(sessionStorage.getItem("ls.tso_user"));
        $profiletsopassword.val(sessionStorage.getItem("ls.tso_password"));
        $profilecreatetime.val(sessionStorage.getItem("ls.create_time"));
        $profiletsouser.val(sessionStorage.getItem("ls.tso_user"));
    }else{
        get_user();
    }
    //profile uptate
    var ajaxProfile = function(){
        var $params;
        $params = {
            "action": "update_user",
            //"name": JSON.parse($profilename.val()),
            "name": $profilename.val(),
            "tso_user": $profiletsouser.val(),
            "tso_password": $profiletsopassword.val(),
            "team": $profileteam.val()
        };
        $.ajax({
            url:"/user",
            type:"put",
            //processData: false,
            data: JSON.stringify($params),
            dataType:"json",
            crossDomain: true,
            beforeSend:function(xhr){
                xhr.setRequestHeader("Authorization", JSON.parse(userToken));
                //xhr.setRequestHeader("Content-Type", "application/json");
            },
            success: function(data){
                if(data.info == "OK"){
                    $profileemail.val(data.email);
                    $profilecreatetime.val(data.create_time);
                    $profileid.val(data.id);
                    $profiletype.val(data.type);
                    $profileteam.val(data.team);
                    $profiletsouser.val(data.tso_user);
                    $profiletsopassword.val(data.tso_password);
                    // user $.notify
                    //timeT = 1;
                    $.notify("Get information successful."
                        , {
                            status: 'success'
                            //pos: 'bottom-right'
                        });
                }
            },
            error: function(){
                $.toaster({
                    priority : 'danger',
                    title : 'Error',
                    message : 'Update failed!'
                });
            }
        });
    };
    var timeCount; //In order to record the changed input, adding the hook timeCount
    function formChange(){
        $("form input").on('input',function(e){
            timeCount = 1;
            $($update).attr("disabled", false);
        });
    }
    formChange();
    if(timeCount = 1){
        $update.click(function(){
            ajaxProfile();
            $($update).attr("disabled", true);
            timeCount = 0;
        });
    }

});