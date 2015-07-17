var workSnippet = function(index) {
    return "<div id=\"outterWorkWrapper" +index+ "\"  class=\"added\">\
        <div id=\"workWrapper" +index+ "\">\
	     <button id=\"removeWork" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
         <div class=\"form-group\">\
           <label for=\"nameOfWork" +index+ "\">Name of project being submitted</label>\
	       <input type=\"text\" class=\"form-control work-name\" data-parsley-required=\"true\" name=\"nameOfWork" +index+"\" id=\"nameOfWork" +index+ "\" placeholder=\"Name of Work\">\
         </div>\
         <div class=\"form-group\">\
           <label for=\"descOfWork" +index+ "\">Description (optional)</label>\
	       <input type=\"text\" class=\"form-control\" data-parsley-required=\"false\" name=\"descOfWork" +index+"\" id=\"descOfWork" +index+ "\" placeholder=\"Description of Work\">\
         </div>\
         <div class=\"form-group\">\
         	<label for=\"extraForWork" +index+ "\">Additional notes (optional) ie: Location, Date, HMU, Creative Direction, etc.</label>\
         	<textarea class=\"form-control\" name=\"extraForWork" +index+ "\" id=\"extraForWork" +index+ "\" placeholder=\"Notes...\"></textarea>\
         </div>\
         <button id=\"add-photo-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add File</button>\
         <label for=\"add-photo-btn\">This is to be an itemized list of all files in your download link</label>\
       </div>\
     </div>";
}

var photoSnippet = function(workIndex, photoIndex) {
	var index = workIndex.toString()+""+photoIndex.toString();
	return "<div id=\"outterPhotoWrapper" +index+ "\" class=\"added\">\
	   <div id=\"photoWrapper" +index+ "\" class=\"wrapper form-inline\">\
	     <div class=\"form-group\">\
	       <label for=\"nameOfPhoto" +index+ "\">File Name</label>\
	       <input type=\"text\" class=\"form-control file-name"+workIndex+"\" data-parsley-required=\"true\" name=\"nameOfPhoto" +index+ "\"id=\"nameOfPhoto" +index+ "\" placeholder=\"ie: photo1.jpg\">\
	     </div>\
	     <div class=\"form-group\">\
	       <label for=\"titleOfPhoto" +index+ "\">File Title</label>\
	       <input type=\"text\" class=\"form-control file-title"+workIndex+"\" data-parsley-required=\"true\" name=\"titleOfPhoto" +index+ "\"id=\"titleOfPhoto" +index+ "\" placeholder=\"ie: Photo of the Moon\">\
	     </div>\
	     <div class=\"form-group\">\
         	<label for=\"add-model-btn" +index+ "\" class=\"wrapper form-inline text-warning\">If applicable:</label>\
         	<button id=\"add-model-btn" +index+ "\" type=\"Button\" class=\"btn btn-default\">+ Add Model</button>\
	     </div>\
	     <button id=\"removePhoto" +index+ "\" type=\"button\" class=\"close\" aria-label=\"Close\"><span aria-hidden=\"true\">&times;</span></button>\
	   </div>";
}

var modelSnippet = function(workIndex, photoIndex, modelIndex) {
	var index = workIndex.toString()+""+photoIndex.toString()+"-"+modelIndex.toString();
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

	$("#view-release-btn").click(function() {
		makeAndShowReleaseForm()
	})
	
	$("#add-work-btn").click(function() {
		var photoIndex = 0;
		var tempWorkIndex = workIndex;

		$("#wrk-wrn").after(workSnippet(workIndex));
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
	$("#add-work-btn").click();
	$("#add-photo-btn0").click();
	$("#add-photo-btn0").click();
	$("#add-photo-btn0").click();
}

function makeAndShowReleaseForm() {
	var releaseString = $("#hidden-release-text").val();
	//var firstName     = $("#firstName").val();
	//var lastName      = $("#lastName").val();
	//var fullName      = (firstName + " " + lastName).toUpperCase();
	//var works         = getItemsStringByClass(".work-name")
	//var newString = vsprintf(releaseString, [works]);

	$("#release-text").text(releaseString)//newString)
	$("#release-modal").modal();
}

function getItemsStringByClass(className) {
	var works = "";
	var count = getItemCountByClass(className);
	var index = 0;
	$(className).each(function(){
		console.log(count)
		var str = "["+ $(this).val() + "] ([" + 
			getPhotosAsString(".file-name"+ ((count - 1) - index)) + 
			"]) (\"Images\"),"

  		if (index < count - 1 || count == 1) {
  			works += str;
  		} else if (index == count - 1) {
  			works += "and " + str;
  		}
  		index++;
	});
	return works;
}

function getPhotosAsString(className) {
	var photos = "";
	var count = getItemCountByClass(className);
	var index = 0;
	$(className).each(function() {
		console.log($(this).val());
		if (index < count - 1) {
			photos += $(this).val() + ", ";
		} else {
			photos += $(this).val();
		}
		index++;
	});
	return photos;
}

function getItemCountByClass(className) {
	var count = 0;
	$(className).each(function() {
		count++;
	});
	return count;
}

$(document).ready(function() {
	setFormClickListeners();
	$(".sigPad").signaturePad({drawOnly : true});
});
