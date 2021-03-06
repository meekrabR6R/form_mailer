
/**
* Brutally ugly helper function for generating html table of artists/models
*/
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
          	                  <td style=\"word-wrap: break-word\"><a href=\""+form['link']+"\">"+form['link']+"</td>";
  			var worksBody = ""
  			$.each(form['works'], function(j, work) {
  				worksBody = artistBody + "<td style=\"word-wrap: break-word\">"+work['name']+"</td>\
				            <td style=\"word-wrap: break-word\">"+work['description']+"</td>";
  				var photosBody = ""
  				if (work['photos'].length > 0) {
  					$.each(work['photos'], function(k, photo) {	
                  		photosBody = worksBody + "<td style=\"word-wrap: break-word\">"+photo['name']+"</td>";
                  		if (photo['title'].length > 0) {
							photosBody += "<td style=\"word-wrap: break-word\">"+photo['title']+"</td>"
						} else {
							photosBody += "<td style=\"word-wrap: break-word\">-</td>"
						}
						photosBody += "<td style=\"word-wrap: break-word\">"+work['extra']+"</td>";
                        var modelsBody = ""
                        if (photo['models'].length > 0) {
                        	$.each(photo['models'], function(l, model) {
                    			tableBody += photosBody + "<td style=\"word-wrap: break-word\">"+model['first_name']+"</td>\
                    			              <td style=\"word-wrap: break-word\">"+model['last_name']+"</td>\
                    			              <td style=\"word-wrap: break-word\">"+model['email']+"</td>";
                    			
                    			if (model['address_one'].length > 0) {
                    				console.log(model['address_one'])
    						  		tableBody += "<td style=\"word-wrap: break-word\">"+model['address_one']+ " \
    						  		                                    "+ model['address_two'] + ",\
                                                                        "+model['city']+ " " + model['state'] + " \
                                                                        "+model['zip']+ " " + model['country']+ "</td></tr>";
                    			} else {
                    			          tableBody += "<td style=\"word-wrap: break-word\">-</td></tr>";
                    			}
                        			
                        	});	
                        } else {
                        	tableBody += photosBody + "<td style=\"word-wrap: break-word\">N/A</td>\
                    	                  <td style=\"word-wrap: break-word\">-</td>\
                    			          <td style=\"word-wrap: break-word\">-</td>\
                        				  <td style=\"word-wrap: break-word\">-</td>\
                        				  <td style=\"word-wrap: break-word\">-</td></tr>";
                        }
  					});
  				} else {
  					tableBody += worksBody + "<td style=\"word-wrap: break-word\">N/A</td>\
                    	          <td style=\"word-wrap: break-word\">-</td>\
                        	      <td style=\"word-wrap: break-word\">"+work['extra']+"</td>\
                            	  <td style=\"word-wrap: break-word\">-</td>\
                    	          <td style=\"word-wrap: break-word\">-</td>\
                    			  <td style=\"word-wrap: break-word\">-</td>\
                    			  <td style=\"word-wrap: break-word\">-</td>\
                        		  <td style=\"word-wrap: break-word\">-</td></tr>";
  				}
  			});
		});
		$("#submissions_table tbody").append(tableBody)
	});
}

/**
* Generate CSV from table
*/
$("#make_csv").click(function() {
	var csv = $('#submissions_table').table2CSV({delivery:'value'});
	this.href="data:application/octet-stream," +encodeURIComponent(csv);	
})

$(document).ready(function() {
	getAllSubmissions();
});