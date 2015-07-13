
function getAllSubmissions() {
	$.getJSON( "/forms_json", function(json) {
  		
  		$.each( json, function(index, form) {
  			console.log(form)
  			$("#submissions_table tbody").append("<tr><td>"+form['updated_at']+"</td>\
  				                                      <td>"+form['last_name']+", "+form['first_name']+"</td>\
  				                                      <td>Artist</td><td>blahblah</td></tr>")
    		$.each(form, function(key, val) {
    			console.log(key + " :: " + val)
    		})
  		});
	});
}

$(document).ready(function() {
	getAllSubmissions();
});