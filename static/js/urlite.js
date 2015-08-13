// Extending jQuery for urlite
(function($){
    $.fn.extend({
        bsAlert: function(message, type, title) {
            type = type || 'success'
            if (type == 'error') {
                type = 'danger'
            }
            var html='<div class="alert alert-' + type.toLowerCase() + ' alert-dismissable text-center"><button type="button" class="close" data-dismiss="alert" aria-hidden="true">&times;</button>';
            if (typeof title !== 'undefined' &&  title !== '') {
                html += '<h4>' + title + '</h4>';
            }
            html += '<span>' + message + '</span></div>';
            $(this).html(html);
            return $(this);
        },
        bsHideAlert: function() {
            var bsalert = $(this);
            bsalert.fadeTo(5000, 500).slideUp(500, function(){
                bsalert.alert('close');
            });
        },
        isUrl: function() {
            return isUrl($(this).val());
        }
    });
})(jQuery);

// Page initailize method call
$(function() {
    if (typeof wpr === "function") { wpr(); }     
});

function queryParams(id) {
    var params = {};
    $(id).find('input[name]').each(function () {
        params[$(this).attr('name')] = $(this).val();
    });
    return params;
}

// Bootstrap Table Formatters
function domainFormatter(v) {
    var d = v.join(', ')
    return d == '*' ? 'All Domains' : d;
}

function prettyDate(v, nd) {
    if (v == '0001-01-01T00:00:00Z') {
        return nd || 'Not yet';
    }

    return $.format.prettyDate(v);
}

function dateFormatter(v) {
    return $.format.prettyDate(v);
}

// Reference: http://stackoverflow.com/a/2723190/1343356
function isUrl(url) {
    return /^(https?|s?ftp):\/\/(((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:)*@)?(((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))|((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.?)(:\d*)?)(\/((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)+(\/(([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)*)*)?)?(\?((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|[\uE000-\uF8FF]|\/|\?)*)?(#((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|\/|\?)*)?$/i.test(url);
}
