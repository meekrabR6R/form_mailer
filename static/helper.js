var workSnippet = function(index) {
    return "<div id=\"outterWorkWrapper" +index+ "\"  class=\"added\">\
        <div id=\"workWrapper" +index+ "\">\
	     <button id=\"removeWork" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
         <div class=\"form-group\">\
           <label for=\"nameOfWork" +index+ "\">Name of Work</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" name=\"nameOfWork" +index+"\" id=\"nameOfWork" +index+ "\" placeholder=\"Name of Work\">\
         </div>\
         <div class=\"form-group\">\
           <label for=\"descOfWork\">Description of Work</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" name=\"descOfWork" +index+"\" id=\"descOfWork" +index+ "\" placeholder=\"Description of Work\">\
         </div>\
         <button id=\"add-photo-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add Photo</button>\
       </div>\
     </div>";
}

var photoSnippet = function(workIndex, photoIndex) {
	var index = workIndex.toString()+""+photoIndex.toString();
	return "<div id=\"outterPhotoWrapper" +index+ "\" class=\"added\">\
	   <div id=\"photoWrapper" +index+ "\" class=\"wrapper form-inline\">\
	   	 <p class=\"text-warning\">If no model is in the photo, please leave 'Model Name' cell blank.</p>\
	     <div class=\"form-group\">\
	       <label for=\"nameOfPhoto" +index+ "\">Name of Photo</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" name=\"nameOfPhoto" +index+ "\"id=\"nameOfPhoto" +index+ "\" placeholder=\"Name of Photo\">\
	     </div>\
         <button id=\"add-model-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add Model</button>\
	     <button id=\"removePhoto" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
	   </div>";
}

var modelSnippet = function(workIndex, photoIndex, modelIndex) {
	var index = workIndex.toString()+""+photoIndex.toString()+""+modelIndex.toString();
	return "<div id=\"modelWrapper" +index+ "\" class=\"added wrapper form-inline\">\
	          <div class=\"form-group\">\
	            <label for=\"firstNameOfModel" +index+ "\">First Name</label>\
	            <input type=\"text\" class=\"form-control\" name=\"firstNameOfModel" +index+ "\" id=\"firstNameOfModel" +index+ "\" placeholder=\"First Name\">\
	          </div>\
	          <div class=\"form-group\">\
	            <label for=\"lastNameOfModel" +index+ "\">Last Name</label>\
	            <input type=\"text\" class=\"form-control\" name=\"lastNameOfModel" +index+ "\" id=\"lastNameOfModel" +index+ "\" placeholder=\"Last Name\">\
	          </div>\
	          <div class=\"form-group\">\
	            <label for=\"emailOfModel" +index+ "\">Model Email</label>\
	            <input type=\"email\" class=\"form-control\" name=\"emailOfModel" +index+ "\" id=\"emailOfModel" +index+ "\" placeholder=\"Model Email\">\
	          </div>\
	          <button id=\"removeModel" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
	        </div>";	
}

function setFormClickListeners() {
	var workIndex = 0;
	$("#add-work-btn").click(function() {
		var photoIndex = 0;
		var tempWorkIndex = workIndex;

		$("#add-work-btn-wrapper").after(workSnippet(workIndex));
		$("#removeWork" + tempWorkIndex.toString()).click(function() {
			$("#outterWorkWrapper" + tempWorkIndex.toString()).remove();
		});

		$("#add-photo-btn" + tempWorkIndex.toString()).click(function() {
			var modelIndex = 0;
			var tempPhotoIndex = photoIndex;

			$("#workWrapper" + tempWorkIndex.toString())
			                   .after(photoSnippet(tempWorkIndex, 
			                   	                   tempPhotoIndex));
			$("#removePhoto" + tempWorkIndex.toString() + "" +
			                   tempPhotoIndex.toString()).click(function() {
				$("#outterPhotoWrapper" + tempWorkIndex.toString() + "" +
				                          tempPhotoIndex.toString()).remove();
			});

			$("#add-model-btn" + tempWorkIndex.toString() + "" +
			                     tempPhotoIndex.toString()).click(function() {
				var tempModelIndex = modelIndex;
				$("#photoWrapper" + tempWorkIndex.toString() + "" +
				                    tempPhotoIndex.toString())
				                        .after(modelSnippet(tempWorkIndex, 
				                    	                    tempPhotoIndex, 
				                    	                    tempModelIndex))
				$("#removeModel" + tempWorkIndex.toString() + "" +
				                   tempPhotoIndex.toString() + "" +
				                   tempModelIndex.toString()).click(function() {
					$("#modelWrapper" + tempWorkIndex.toString() + "" +
					                    tempPhotoIndex.toString() + "" + 
					                    tempModelIndex.toString()).remove();
				});
				modelIndex++;
			});

			photoIndex++;
		});		
		workIndex++;			
	});
}

$(document).ready(function() {
	setFormClickListeners();
	$(".sigPad").signaturePad({drawOnly : true});
});
