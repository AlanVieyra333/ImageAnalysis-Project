function uploadFileAjax(file, func) {
	var data = new FormData();
    data.append('file', file);
    
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

function uploadFiles(files, func1, func2) {
    // Loop through the FileList and render image files as thumbnails.
    for (var i = 0, f; f = files[i]; i++) {
        // Only process image files.
        if (!f.type.match('image.*')) {
            alert("Solo se permiten imagenes.");
            continue;
        }
        func1(f,func2);
    }
}

function showNewImage(res){
	if (typeof res == "string")
        res = $.parseJSON(res);

    if (!res) {
        alert("Sin respuesta del servidor.");
        return;
    }

    if(res.Code == 1){
        document.getElementById('anImage').innerHTML = ['<img class="thumb" id="imageEdit" src="uploaded/',
        res.FileNameEdit, '" title="',res.FileName,'"/>'].join('');

        document.getElementById('fileSelect').innerHTML = ['<input type="file" id="files" name="files[]"/>'].join('');
        document.getElementById('files').addEventListener('change', handleFileSelect, false);

		$('#anImage').removeAttr('class');
		$('#anImage').show();
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
		
        var infoJSON = {
			Operation:  ID,
			FileName:   fileName,
			FileNameEdit: fileNameEdit,
			Args: "" +  res.FileName + ";"+ res.FileNameEdit + ";"
		};
				
		var data = JSON.stringify(infoJSON);
		var ok = function(res){
			$('#imageEdit').removeAttr('src');
			$('#imageEdit').attr('src','uploaded/' + res.FileNameEdit);
			$('#imageEdit').show();
		};

		sendAjaxJSON(data, ok);
		
		document.getElementById('imageMaskSelect').innerHTML = ['<input type="file" id="imageMask" name="files[]"/>'].join('');
        document.getElementById('imageMask').addEventListener('change', handleFileSelect, false);
    }else{
        alert(res.Message);
    }
}

// --------------------
function handleFileSelect(evt) {
    var files = evt.target.files; // FileList object
	if(this.id == "files")
    	uploadFiles(files, uploadFileAjax, showNewImage);
	else if(this.id == "imageMask"){
		uploadFiles(files, uploadFileAjax, operationImages);		
	}
}

function handleFileSelectDrop(evt) {
    evt.stopPropagation();
    evt.preventDefault();

    var files = evt.dataTransfer.files; // FileList object.
    uploadFiles(files);
}

function handleDragOver(evt) {
    evt.stopPropagation();
    evt.preventDefault();
    evt.dataTransfer.dropEffect = 'copy'; // Explicitly show this is a copy.
}

// --------------------
function open(){
    document.getElementById('files').click();
}

function save() {
    alert("save");
}

function undo() {
    alert("ok: " + this.id);
}

function redo() {
    alert("ok: " + this.id);
}

function showPopup(id) {
	switch (id){
		case 1:
			document.getElementById('label-args').innerHTML = "Conectividad: {4,8}";
			document.getElementById('param').innerHTML = "<select id='args'>"+
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
		default:
			document.getElementById('label-args').innerHTML = "ID invalido";
	}
	$('.my-overlay-container').fadeIn(function() {		
		window.setTimeout(function(){
			//document.getElementById('args').value = "";			
			$('.my-window-container.my-zoomin').addClass('my-window-container-visible');
			document.getElementById('args').focus();
			$('#op').removeAttr('value');
			$('#op').attr('value',id);
			$('#op').show();
		}, 100);		
	});
}

function closePopup() {
	$('.my-overlay-container').fadeOut().end().find('.my-window-container').removeClass('my-window-container-visible');	
}