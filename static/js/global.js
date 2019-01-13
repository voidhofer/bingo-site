$(function() {
	// Hide any messages after a few seconds
    hideFlash();
});

function hideFlash(rnum)
{
    if (!rnum) rnum = '0';
    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(1000, function() {
            $(this).css({"visibility":"hidden",display:'block'}).slideUp();
            var that = this;
            _.delay(function() {
				that.remove(); 
			}, 400);
        });
    }, 5000);
}

function showFlash(obj)
{
    $('#flash-container').html();
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
		var message = '<div id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert alert-dismissible '+v.Class+'">'
		+ '<button type="button" id="cg_popup_close" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
		+ v.Message + '</div>';
		$('#flash-container').prepend(message);
		hideFlash(rnum);
    });
}

// flashError display error popup
function flashError(message) {
	var flash = [{Class: "alert-danger", Message: message}];
	showFlash(flash);
}
// flashSuccess display success popup
function flashSuccess(message) {
	var flash = [{Class: "alert-success", Message: message}];
	showFlash(flash);
}
// flashNotice display notice popup
function flashNotice(message) {
	var flash = [{Class: "alert-info", Message: message}];
	showFlash(flash);
}
// flashWarning display erning popup
function flashWarning(message) {
	var flash = [{Class: "alert-warning", Message: message}];
	showFlash(flash);
}