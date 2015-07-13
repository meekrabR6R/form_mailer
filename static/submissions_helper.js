
function getAllSubmissions() {
	$.getJSON( "/forms_json", function(json) {
  		var tableBody = ""
  		$.each( json, function(i, form) {
  			var artistBody = "<tr><td style=\"word-wrap: break-word\">"+form['updated_at']+"</td>\
	  				                                      <td style=\"word-wrap: break-word\">"+form['first_name']+"</td>\
	  				                                      <td style=\"word-wrap: break-word\">"+form['last_name']+"</td>\
	  				                                      <td style=\"word-wrap: break-word\">"+form['email']+"</td>\
	  				                                      <td style=\"word-wrap: break-word\">"+form['address_one']+ " " + form['address_two'] + ",\
	  				                                                                          "+ form['city']+ " " + form['state'] + "\
	  				                                                                          "+ form['zip']+ " " + form['country']+ "</td>\
	  				                                      <td style=\"word-wrap: break-word\">"+form['link']+"</td>";
  			var worksBody = ""
  			$.each(form['works'], function(j, work) {
  				worksBody = artistBody + "<td style=\"word-wrap: break-word\">"+work['name']+"</td>\
                                <td style=\"word-wrap: break-word\">"+work['description']+"</td>";
  				var photosBody = ""
  				if (work['photos'].length > 0) {
  					$.each(work['photos'], function(k, photo) {	
                  		photosBody = worksBody + "<td style=\"word-wrap: break-word\">"+photo['name']+"</td>\
                    		          <td style=\"word-wrap: break-word\">"+photo['name']+"</td>\
                        		      <td style=\"word-wrap: break-word\">"+work['extra']+"</td>";
                        var modelsBody = ""
                        if (photo['models'].length > 0) {
                        	$.each(photo['models'], function(l, model) {
                    			tableBody += photosBody + "<td style=\"word-wrap: break-word\">"+model['first_name']+"</td>\
                    			              <td style=\"word-wrap: break-word\">"+model['last_name']+"</td>\
                    			              <td style=\"word-wrap: break-word\">"+model['email']+"</td>\
                        					  <td style=\"word-wrap: break-word\">"+model['address_one']+ " " + model['address_two'] + ",\
	  				                                                              "+model['city']+ " " + model['state'] + "\
	  				                                                              "+model['zip']+ " " + model['country']+ "</td></tr>";
                        	});	
                        } else {
                        	tableBody += photosBody + "<td style=\"word-wrap: break-word\">N/A</td>\
                    	                  <td style=\"word-wrap: break-word\">N/A</td>\
                    			          <td style=\"word-wrap: break-word\">N/A</td>\
                        				  <td style=\"word-wrap: break-word\">N/A</td></tr>";
                        }
  					});
  				} else {
  					tableBody += worksBody + "<td style=\"word-wrap: break-word\">N/A</td>\
                    	          <td style=\"word-wrap: break-word\">N/A</td>\
                        	      <td style=\"word-wrap: break-word\">"+work['extra']+"</td>\
                            	  <td style=\"word-wrap: break-word\">N/A</td>\
                    	          <td style=\"word-wrap: break-word\">N/A</td>\
                    			  <td style=\"word-wrap: break-word\">N/A</td>\
                        		  <td style=\"word-wrap: break-word\">N/A</td></tr>";
  				}
  			});
		});
		$("#submissions_table tbody").append(tableBody)
	});
}

$(document).ready(function() {
	getAllSubmissions();
});