function uploadFileAjax(file, func, myType) {
	var data = new FormData();
    data.append('file', file);
	data.append('type', myType);
    
    var url = "./upload/";
    var type = "POST";
	var contentType = false;	
	var processData = false;
	var cache = false;

    $.ajax({
        url: url,
        type: type,
        contentType: contentType,
        data: data,
        processData: processData,
		async: true,
        cache: cache
    }).done(
        function (res) {
            func(res);
        }
    ).fail(
        function (jqXHR, textStatus, errorThrown) {
            alert(textStatus + ": Error al conectar con el servidor.");
        }
    );
}

function uploadFiles(files, func1, func2, myType) {
    // Loop through the FileList and render image files as thumbnails.
    //alert(files.length);
	for (var i = 0; i<files.length; i++) {
		var f = files[i];
        // Only process image files.
        if (!f.type.match('image.*')) {
            alert("Solo se permiten imagenes.");
            continue;
        }
        func1(f,func2, myType);
    }
}

function showNewImage(res){
	if (typeof res == "string")
        res = $.parseJSON(res);

    if (!res) {
        alert("Sin respuesta del servidor.");
        return;
    }

	var IDClient = document.cookie.split(';')[0];
	var idx1 = IDClient.indexOf("=")+1;
	var idx2 = IDClient.indexOf(",");
    IDClient = IDClient.substring(idx1,idx2);
    if(res.Code == 1){
        document.getElementById('anImage').innerHTML = ['<img class="thumb" id="imageEdit" src="uploaded/', IDClient,'/',
        res.FileNameEdit, '" title="',res.FileName,'"/>',
		'<!-- Chart Section --><div class="thumb" id="chart_div"></div>',
		'<!-- Histogram propierties Section --><div class="info" id="histInfo"></div>'].join('');

        document.getElementById('fileSelect').innerHTML = ['<input type="file" id="files" name="files[]"/>'].join('');
        document.getElementById('files').addEventListener('change', handleFileSelect, false);

		$('#anImage').removeAttr('class');
		$('#anImage').show();
		
		$('#save').removeAttr('href');
		$('#save').attr('href','uploaded/'+ IDClient+'/' + res.FileNameEdit);
		$('#save').removeAttr('download');
		$('#save').attr('download', res.FileName);
		$('#save').show();
    }else{
        alert(res.Message);
    }
}

function operationImages(res){
	if (typeof res == "string")
        res = $.parseJSON(res);

    if (!res) {
        alert("Sin respuesta del servidor.");
        return;
    }
	var IDClient = document.cookie.split(';')[0];
	var idx1 = IDClient.indexOf("=")+1;
	var idx2 = IDClient.indexOf(",");
    IDClient = IDClient.substring(idx1,idx2);
	// -------------------Send operation between images-------------------------
    if(res.Code == 1){
		var ID = parseInt($('#imageMask').attr('valueOp'));
		var fileName = document.getElementById('imageEdit').title;
		var fileNameEdit = document.getElementById('imageEdit').src;
		var j=0;
		for(var i=0; i<fileNameEdit.length; i++){
			if(fileNameEdit.charAt(i) == '/')
				j=i;
		}
		fileNameEdit = fileNameEdit.substring(j+1);
		//alert(res.FileName + " " + res.FileNameEdit);
        var infoJSON = {
			Operation:  ID,
			FileName:   fileName,
			FileNameEdit: fileNameEdit,
			Args: "" +  res.FileName + ";"+ res.FileNameEdit + ";"
		};
		
		if( (ID >= 31 && ID <= 32) || (ID >= 35 && ID <= 41) ) {
			infoJSON.Args += $( "#args" ).val() + ";";
		}
				
		var data = JSON.stringify(infoJSON);
		var ok = function(res){
			$('#imageEdit').removeAttr('src');
			$('#imageEdit').attr('src','uploaded/'+ IDClient+'/' + res.FileNameEdit);
			$('#imageEdit').show();
			
			$('#save').removeAttr('href');
			$('#save').attr('href','uploaded/'+ IDClient + "/" + res.FileNameEdit);
			$('#save').removeAttr('download');
			$('#save').attr('download', res.FileName);
			$('#save').show();
		};

		sendAjaxJSON(data, ok, false);
		
		document.getElementById('imageMaskSelect').innerHTML = ['<input type="file" id="imageMask" name="files[] valueOp=0"/>'].join('');
		document.getElementById('imageMask').addEventListener('change', handleFileSelect, false);
    }else{
        alert(res.Message);
    }
}

// --------------------
function handleFileSelect(evt) {
	//alert("is:" + this.id);
    var files = evt.target.files; // FileList object
	if(this.id == "files"){
		/*	----------------	Delete all data (0)	-----------------	*/		
		uploadFiles(files, uploadFileAjax, showNewImage, "0");    	
	}else if(this.id == "imageMask"){
		uploadFiles(files, uploadFileAjax, operationImages, "1");		
	}
}

function handleFileSelectDrop(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    var files = evt.dataTransfer.files; // FileList object.
	
	uploadFiles(files, uploadFileAjax, showNewImage, "0");
}

function handleDragOver(evt) {
    evt.stopPropagation();
    evt.preventDefault();
    evt.dataTransfer.dropEffect = 'copy'; // Explicitly show this is a copy.
}

// --------------------

function showPopup(id) {
	switch (id){
		case 1:
			document.getElementById('label-args').innerHTML = "Conectividad: {4,8}";
			document.getElementById('param').innerHTML = "<select name='args' id='args'>"+
				"    <option value='4'>4</option>"+
				"    <option value='8'>8</option>"+
				"</select>";
			break;
		case 2:
			document.getElementById('label-args').innerHTML = "Umbral: [0,1]";
			document.getElementById('param').innerHTML = "<input type='range' name='args' min='0' max='1' step='0.01' value='0.5' id='args'>";
			break;
		case 9:
			document.getElementById('label-args').innerHTML = "Umbral y canal: [0,1],{Rojo, Verde, Azul}";
			$('#param').html("<input type='range' name='args' min='0' max='1' step='0.01' value='0.5' id='args'>");
			$('#param').append("<input type='radio' name='args2' value='1' checked> Rojo<br>"+
	  			"<input type='radio' name='args2' value='2'> Verde<br>"+
				"<input type='radio' name='args2' value='3'> Azul<br>");
			break;
		case 10:
			document.getElementById('label-args').innerHTML = "Canal: {Rojo, Verde, Azul}";
			document.getElementById('param').innerHTML = "<input type='radio' name='args' value='1' checked id='args'> Rojo<br>"+
	  			"<input type='radio' name='args' value='2'> Verde<br>"+
				"<input type='radio' name='args' value='3'> Azul<br>";
			break;
		case 13:
			document.getElementById('label-args').innerHTML = "Valor de desplazamiento: [-255,255]";
			$('#param').html("<input type='number' name='args' min='-255' max='255' value='0' step='1' required id='args'>");
			break;
		case 14:
		case 15:
			document.getElementById('label-args').innerHTML = "Maximo y minimo valor deseado: [0,255]";
			$('#param').html("<input type='number' name='args' min='0' max='255' value='0' step='1' required id='args'>");
			$('#param').append("<input type='number' name='args2' min='0' max='255' value='255' step='1' required id='args2'>");
			break;
		case 17:
		case 18:
		case 19:
		case 20:
			document.getElementById('label-args').innerHTML = "Umbral: [0,1]";
			document.getElementById('param').innerHTML = "<input type='range' name='args' min='0' max='1' step='0.01' value='0.5' id='args'>";
			break;
		case 21:
		case 23:
		case 24:
		case 25:
		case 26:
			document.getElementById('label-args').innerHTML = "Tamaño de ventana (nxn): {n | impar >= 3}";
			document.getElementById('param').innerHTML = "<input type='number' name='args' min='3' value='3' step='2' required id='args'>";
			break;
		case 22:
			document.getElementById('label-args').innerHTML = "Tamaño de ventana(nxn) y peso(N): {n | impar >= 3}, N>1";
			$('#param').html("<input type='number' name='args' min='3' value='3' step='2' required id='args'>");
			$('#param').append("<input type='number' name='args2' min='2' value='2' required id='args2'>");
			break;
		case 31:
		case 32:
		case 35:
		case 37:
		case 38:
		case 39:
		case 40:
		case 41:
			document.getElementById('label-args').innerHTML = "Ingrese las coordenadas del origen del EE: {x,y) | x,y >= 0<br>(A continuacion seleccione la imagen que usara como EE)";
			$('#param').html("(X,Y)=<input type='text' name='args' value='0,0' pattern='^[\\d]+[,][\\d]+$' required id='args'>");
			break;
		default:
			document.getElementById('label-args').innerHTML = "ID invalido";
			$('#param').html("");
			return 0;
	}
	$('.my-overlay-container').fadeIn(function() {		
		window.setTimeout(function(){
			//document.getElementById('args').value = "";
			$('.my-window-container.my-zoomin').addClass('my-window-container-visible');			
			$('#op').removeAttr('value');
			$('#op').attr('value',id);
			$('#op').show();
			document.getElementById('args').focus();
		}, 100);		
	});
}

function closePopup() {
	$('.my-overlay-container').fadeOut().end().find('.my-window-container').removeClass('my-window-container-visible');
	//$('.my-window-container').removeClass('my-window-container-visible');
	//$('.my-overlay-container').hide();
	$(".info").css("display", "none");
	$("#chart_div").css("display", "none");
}

// ----------------------
function open(){
    document.getElementById('files').click();
}

function getInfoJSON(){
	var fileName = document.getElementById('imageEdit').title;
	var fileNameEdit = document.getElementById('imageEdit').src;
	var j=0;
	for(var i=0; i<fileNameEdit.length; i++){
		if(fileNameEdit.charAt(i) == '/')
			j=i;
	}
	fileNameEdit = fileNameEdit.substring(j+1);
	var infoJSON = {
		Operation:  0,
		FileName:   fileName,
		FileNameEdit: fileNameEdit,
		Args: ""
	};
	return infoJSON;
}

function undo() {
	$(".info").css("display", "none");
	$("#chart_div").css("display", "none");
	
	var infoJSON = getInfoJSON();
	infoJSON.Operation = -1;
	infoJSON.Args = "1";
	
	var data = JSON.stringify(infoJSON);
	var ok = function(res){
		var IDClient = document.cookie.split(';')[0];
		var idx1 = IDClient.indexOf("=")+1;
		var idx2 = IDClient.indexOf(",");
	    IDClient = IDClient.substring(idx1,idx2);
		
		$('#imageEdit').removeAttr('src');
		$('#imageEdit').attr('src','uploaded/'+ IDClient + "/" + res.FileNameEdit);
		$('#imageEdit').show();
		
		$('#save').removeAttr('href');
		$('#save').attr('href','uploaded/'+ IDClient + "/" + res.FileNameEdit);
		$('#save').removeAttr('download');
		$('#save').attr('download', res.FileName);
		$('#save').show();
	};
	sendAjaxJSON(data, ok, true);
}

function redo() {
	$(".info").css("display", "none");
	$("#chart_div").css("display", "none");
	
    var infoJSON = getInfoJSON();
	infoJSON.Operation = -1;
	infoJSON.Args = "2";
	
	var data = JSON.stringify(infoJSON);
	var ok = function(res){
		var IDClient = document.cookie.split(';')[0];
		var idx1 = IDClient.indexOf("=")+1;
		var idx2 = IDClient.indexOf(",");
	    IDClient = IDClient.substring(idx1,idx2);
		
		$('#imageEdit').removeAttr('src');
		$('#imageEdit').attr('src','uploaded/'+ IDClient + "/" + res.FileNameEdit);
		$('#imageEdit').show();
		
		$('#save').removeAttr('href');
		$('#save').attr('href','uploaded/'+ IDClient + "/" + res.FileNameEdit);
		$('#save').removeAttr('download');
		$('#save').attr('download', res.FileName);
		$('#save').show();
	};
	sendAjaxJSON(data, ok, true);
}
