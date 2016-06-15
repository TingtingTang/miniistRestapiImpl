// Xeditable Demo
// -----------------------------------

(function(window, document, $, undefined){

    $(function(){

        // Font Awesome support
        $.fn.editableform.buttons =
            '<button type="submit" class="btn btn-primary btn-sm editable-submit">'+
            '<i class="fa fa-fw fa-check"></i>'+
            '</button>'+
            '<button type="button" class="btn btn-default btn-sm editable-cancel">'+
            '<i class="fa fa-fw fa-times"></i>'+
            '</button>';

        //enable / disable
        $('#enable').click(function() {
            $('#user .editable').editable('toggleDisabled');
        });

        //editables
        $('#email_prof').editable({
            validate: function(value) {
                if($.trim(value) === '') return 'This field is required';
            }
        });

        $('#sex').editable({
            prepend: "not selected",
            source: [
                {value: 1, text: 'Male'},
                {value: 2, text: 'Female'}
            ],
            display: function(value, sourceData) {
                var colors = {"": "gray", 1: "green", 2: "blue"},
                    elem = $.grep(sourceData, function(o){return o.value == value;});

                if(elem.length) {
                    $(this).text(elem[0].text).css("color", colors[value]);
                } else {
                    $(this).empty();
                }
            }
        });

        $('#team_prof').editable({
            showbuttons: true
        });
        $('#TSOUser_prof').editable({
            showbuttons: true
        });

        $('#TSOpassword_prof').editable({
            validate: function(value) {
                if($.trim(value) === '') return 'This field is required';
            }
        });

    });

    //headpic
    var $portrait = $("#portrait");
    var $por_hide = $("#por_hide");
    var $head_pic = $("#head_pic");

    function getFileName(o){
        var pos=o.lastIndexOf("\\");
        return o.substring(pos+1);
    }

    //hidden user poritrait
    var headEdit = function(){
        //var fileName;
        //var picName;
        $por_hide.hide();
        $head_pic.find("a").click(function(){
            $por_hide.trigger("click");

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
            $por_hide.change(function() {
                var eImg = $('<img />');
                //console.info($(this)[0].files[0]);
                $portrait.css("background-image","url("+getObjectURL($(this)[0].files[0])+")"); // »ò this.files[0] this->input
                console.info($(this)[0].files);
            });
        })
    };
    headEdit();


    //Team selected
    var $profile_team = $("#profile_team");
    var $profile_teamLI = $("#profile_teamLI");
    var profile_value;
    $("#profile_teamLI li").each(function(){
        $(this).find("a").click(function(){

            profile_value = this.innerText;
            $profile_team.attr("value", profile_value);
        })
    });
    $("#head_pic #motto").focus(function () {
        if(this.value == "Edit your motto......"){
            this.value = "";
        }
        $(this).css("background-color", "white");
    }).blur(function () {
        if(this.value == ""){
            this.value = "Edit your motto......";
        }
        $(this).css("background-color", "transparent");
    });

})(window, document, window.jQuery);

