var workSnippet = function(index) {
    return "<div class=\"added\">\
       <div id=\"workWrapper" +index+ "\">\
	     <button type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
         <div class=\"form-group\">\
           <label for=\"nameOfWork" +index+ "\">Name of Work</label>\
	       <input type=\"text\" class=\"form-control\" id=\"nameOfWork" +index+ "\" placeholder=\"Name of Work\">\
         </div>\
         <div class=\"form-group\">\
           <label for=\"descOfWork\">Description of Work</label>\
	       <input type=\"text\" class=\"form-control\" id=\"nameOfWork" +index+ "\" placeholder=\"Description of Work\">\
         </div>\
         <button id=\"add-model-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add Model</button>\
       </div>\
     </div>"
}

var modelSnippet = function(workIndex, modelIndex) {
	var index = workIndex.toString()+""+modelIndex.toString();
	console.log(index);
	return "<div class=\"added\">\
	   <div class=\"wrapper form-inline\">\
	   	 <p class=\"text-warning\">If no model is in the photo, please leave 'Model Name' cell blank.</p>\
	     <div class=\"form-group\">\
	       <label for=\"nameOfPhoto" +index+ "\">Name of Photo</label>\
	       <input type=\"text\" class=\"form-control\" id=\"nameOfPhoto" +index+ "\" placeholder=\"Name of Photo\">\
	     </div>\
	     <div class=\"form-group\">\
	       <label for=\"nameOfModel" +index+ "\">Name of Model</label>\
	       <input type=\"text\" class=\"form-control\" id=\"nameOfModel" +index+ "\" placeholder=\"Name of Model\">\
	     </div>\
	     <button type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
	   </div>\
	 </div>"
}


function setFormClickListeners() {
	var workIndex = 0;
	$("#add-work-btn").click(function() {
		var modelIndex = 0;
		var tempWorkIndex = workIndex;
		console.log(workIndex.toString());
		$("#add-work-btn-wrapper").after(workSnippet(workIndex));
		
		$("#add-model-btn" + tempWorkIndex.toString()).click(function() {
			var tempModelIndex = modelIndex;
			$("#workWrapper" + tempWorkIndex.toString()).after(modelSnippet(tempWorkIndex, tempModelIndex));
			modelIndex++;
		});
		
		workIndex++;			
	});
}

$(document).ready(function() {
	setFormClickListeners();
	$(".sigPad").signaturePad();
});
