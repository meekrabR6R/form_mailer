var workSnippet = function(index) {
    return "<div id=\"outterWorkWrapper" +index+ "\"  class=\"added\">\
        <div id=\"workWrapper" +index+ "\">\
	     <button id=\"removeWork" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
         <div class=\"form-group\">\
           <label for=\"nameOfWork" +index+ "\">Name of Work</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" id=\"nameOfWork" +index+ "\" placeholder=\"Name of Work\">\
         </div>\
         <div class=\"form-group\">\
           <label for=\"descOfWork\">Description of Work</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" id=\"nameOfWork" +index+ "\" placeholder=\"Description of Work\">\
         </div>\
         <button id=\"add-photo-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add Photo</button>\
       </div>\
     </div>"
}

var photoSnippet = function(workIndex, photoIndex) {
	var index = workIndex.toString()+""+photoIndex.toString();
	console.log(index);
	return "<div id=\"photoWrapper" +index+ "\" class=\"added wrapper form-inline\">\
	   	 <p class=\"text-warning\">If no model is in the photo, please leave 'Model Name' cell blank.</p>\
	     <div class=\"form-group\">\
	       <label for=\"nameOfPhoto" +index+ "\">Name of Photo</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"true\" id=\"nameOfPhoto" +index+ "\" placeholder=\"Name of Photo\">\
	     </div>\
	     <div class=\"form-group\">\
	       <label for=\"nameOfModel" +index+ "\">Name of Model</label>\
	       <input type=\"text\" class=\"form-control\" id=\"nameOfModel" +index+ "\" placeholder=\"Name of Model\">\
	     </div>\
	     <div class=\"form-group\">\
	       <label for=\"emailOfModel" +index+ "\">Email Address of Model</label>\
	       <input type=\"text\" class=\"form-control\" id=\"emailOfModel" +index+ "\" placeholder=\"Email Address of Model\">\
	     </div>\
	     <button id=\"removePhoto" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
	   </div>"
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
			var tempPhotoIndex = photoIndex;
			$("#workWrapper" + tempWorkIndex.toString()).after(photoSnippet(tempWorkIndex, tempPhotoIndex));
			$("#removePhoto" + tempWorkIndex.toString() + "" + tempPhotoIndex.toString()).click(function() {
				$("#photoWrapper" + tempWorkIndex.toString() + "" + tempPhotoIndex.toString()).remove();
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
