/**
 * Created by Fenix on 03/10/2016.
 */

function sendAjaxJSON(data, ok){
	var url = "/json/";
	var type = "POST";
	var contentType = "application/json; charset=utf-8";
	var dataType = "json";
	var async = false;
	var cache = false;
//	alert("AJAX");
	$.ajax({
		url: url,
		data: data,
		type: type,
		contentType: contentType,
		dataType: dataType,
		async: async,
		cache: cache
	}).done(
		function (res) {
//			alert("RES");
			if (typeof res == "string")
				res = $.parseJSON(res);

			if (!res) {
				alert("Sin respuesta del servidor.");
				return;
			}

			if(res.Code == 1){
				ok(res);
			}else{
				alert(res.Message);
			}
		}
	).fail(
		function (jqXHR, textStatus, errorThrown) {
			alert(textStatus + ": Error al conectar con el servidor.");
		}
	);
}
  
function editImage() {	
    var ID;
	if( this.id == "send" ){
		ID = parseInt(document.getElementById('op').value);		
	}else{
		ID = parseInt(""+this.id);
	}
	
	closePopup();
	
	if(  this.id != "send" && (ID == 1 || ID == 2 || ID == 9 || ID == 10) ){ // Get value.
		showPopup(ID);
	}else if(ID >= 5 && ID <= 8 ){ // Get other image.
		$('#imageMask').attr('valueOp',ID);
		$('#imageMask').removeAttr('value');
		$('#imageMask').attr('value','');
		$('#imageMask').show();
		document.getElementById('imageMask').click();
	}else{		
		/*	----------------	Send data	-----------------	*/
		var fileName = document.getElementById('imageEdit').title;
		var fileNameEdit = document.getElementById('imageEdit').src;
		var j=0;
		for(var i=0; i<fileNameEdit.length; i++){
			if(fileNameEdit.charAt(i) == '/')
				j=i;
		}
		fileNameEdit = fileNameEdit.substring(j+1);
		
		var infoJSON = {
			Operation:  ID,
			FileName:   fileName,
			FileNameEdit: fileNameEdit,
			Args: ""
		};
		
		if(this.id == "send"){
			if(ID == 9){
				infoJSON.Args = "" + $( "#args" ).val() + ";" + $('input:radio[name=args2]:checked').val() + ";";
			}else if(ID == 10){
				infoJSON.Args = "" + $('input:radio[name=args]:checked').val() + ";";
			}else{
				infoJSON.Args = "" + $( "#args" ).val() + ";";
			}
		}
		
		var data = JSON.stringify(infoJSON);
		var ok = function(res){
			$('#imageEdit').removeAttr('src');
			$('#imageEdit').attr('src','uploaded/' + res.FileNameEdit);
			$('#imageEdit').show();
		};

		sendAjaxJSON(data, ok);
	}
}