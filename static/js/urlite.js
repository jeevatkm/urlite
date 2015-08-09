// Extending jQuery for urlite
(function($){
    $.fn.extend({
        bs_alert: function(message, type, title) {
            type = type || 'success'
            var html='<div class="alert alert-' + type.toLowerCase() + ' alert-dismissable text-center"><button type="button" class="close" data-dismiss="alert" aria-hidden="true">&times;</button>';
            if (typeof title !== 'undefined' &&  title !== '') {
                html += '<h4>' + title + '</h4>';
            }
            html += '<span>' + message + '</span></div>';
            $(this).html(html);
        }
    });
})(jQuery);

// Page initailize method call
$(function() {
    var alert_container = $("#alert_container");
    if (alert_container) {
        alert_container.fadeTo(5000, 500).slideUp(500, function(){
            alert_container.alert('close');
        });
    }

    if (typeof wpr === "function") { wpr(); }     
});