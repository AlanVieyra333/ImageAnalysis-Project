function uploadAjax(file) {
    //var inputFileImage = document.getElementById(“archivoImage”);
    //var file = inputFileImage.files[0];
    var data = new FormData();
    data.append('file', file);
    var url = "./upload/";
    var type = "POST";

    $.ajax({
        url: url,
        type: type,
        contentType: false,
        data: data,
        processData: false,
        cache: false
    }).done(
        function (res) {
            if (typeof res == "string")
                res = $.parseJSON(res);

            if (!res) {
                // no se encontraron registros :(
                alert("Sin respuesta del servidor.");
                return;
            }

            //alert(res.Message);
            if(res.Code == 1){
                //document.getElementById('anImage').innerHTML = ['<img class="thumb" id="imageEdit" src="uploaded/',
                //    res.FileName, '" title="', res.FileName ,'"/>'].join('');
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
    ).fail(
        function (jqXHR, textStatus, errorThrown) {
            alert(textStatus + ": Error al conectar con el servidor.");
        }
    );
}

function uploadFiles(files) {
    // Loop through the FileList and render image files as thumbnails.
    for (var i = 0, f; f = files[i]; i++) {

        // Only process 1 file.
        //var f = files[0];

        // Only process image files.
        if (!f.type.match('image.*')) {
            alert("Solo se permiten imagenes.");
            continue;
        }

        uploadAjax(f);

        /*var reader = new FileReader();
        // Closure to capture the file information.
        reader.onload = (function (theFile) {
            return function (e) {
                // Render thumbnail.
                document.getElementById('anImage').innerHTML = ['<img class="thumb" src="',
                    e.target.result, '" title="', escape(theFile.name), '"/>'].join('');
            };
        })(f);

        // Read in the image file as a data URL.
        reader.readAsDataURL(f);*/
    }
}

// --------------------
function handleFileSelect(evt) {
    var files = evt.target.files; // FileList object
    uploadFiles(files);
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
			break;
		case 2:
			document.getElementById('label-args').innerHTML = "Umbral: [0,1]";
			break;
		default:
			document.getElementById('label-args').innerHTML = "ID invalido";
	}
	$('.my-overlay-container').fadeIn(function() {		
		window.setTimeout(function(){
			document.getElementById('args').value = "";
			document.getElementById('args').focus();
			$('.my-window-container.my-zoomin').addClass('my-window-container-visible');
			$('#op').removeAttr('value');
			$('#op').attr('value',id);
			$('#op').show();
		}, 100);		
	});
}

function closePopup() {
	$('.my-overlay-container').fadeOut().end().find('.my-window-container').removeClass('my-window-container-visible');	
}