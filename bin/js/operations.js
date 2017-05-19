/**
 * Created by Fenix on 03/10/2016.
 */

function sendAjaxJSON(data, ok, async){
	var url = "./json/";
	var type = "POST";
	var contentType = "application/json; charset=utf-8";	
	//var contentType = false;		
	var processData = false;
	var dataType = "json";	
	var cache = false;
	//var async = true;
	
	var beforeSend = function () {};
	var complete = function () {};
	var success = function (res) {
//		alert("RES");
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
	};
	var error = function (jqXHR, textStatus, errorThrown) {
		alert(textStatus + ": Error al conectar con el servidor.");
	};
	
//	alert("AJAX");
	
	$.ajax({
		url: url,
		type: type,
		contentType: contentType,
		data: data,		
		processData: processData,		
		dataType: dataType,
		
		async: async,
		cache: cache,
		beforeSend: beforeSend,
		complete: complete,
		success: success,
		error: error
	});//.done(success).fail(error);	
}

var P;
function editImage() {	
    var ID;
	if( this.id == "send" ){
		ID = parseInt(document.getElementById('op').value);		
	}else{
		ID = parseInt(""+this.id);
	}
	closePopup();

	if(  this.id != "send" && ((ID >= 1 && ID<= 2) || (ID >= 9 && ID<= 10) || (ID >= 13 && ID<= 15) || (ID >= 17 && ID<= 26) || (ID >= 31 && ID<= 32) || (ID >= 35 && ID<= 41)  ) ){ // Get value.
		showPopup(ID);
	}else if( (ID >= 5 && ID <= 8) || (ID >= 31 && ID<= 32) || (ID >= 35 && ID<= 41) ) { // Get other image.		
		$('#imageMask').attr('valueOp',ID);
		$('#imageMask').removeAttr('value');
		$('#imageMask').attr('value','');
		$('#imageMask').show();
		document.getElementById('imageMask').click();
	}else{
		/*	----------------	Send data	-----------------	*/		
		var infoJSON = getInfoJSON();
		infoJSON.Operation = ID;
		
		if(this.id == "send"){
			if(ID == 9){
				infoJSON.Args = "" + $( "#args" ).val() + ";" + $('input:radio[name=args2]:checked').val() + ";";
			}else if(ID == 10){
				infoJSON.Args = "" + $('input:radio[name=args]:checked').val() + ";";
			}else if(ID == 22 || (ID >= 14 && ID <= 15)){
				infoJSON.Args = "" + $( "#args" ).val() + ";" + + $( "#args2" ).val() + ";";
			}else{
				infoJSON.Args = "" + $( "#args" ).val() + ";";
			}
		}
		
		var data = JSON.stringify(infoJSON);
		// hacer al realizar ajax
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
			
			//alert("Data: " + res.Data.Data1[0]);
			for (var k=0; k< 255; k++){
				if(res.Data.Data1[k] != 0.0){
					P = res.Data.Data1;	
					// Crear grafico.
					//$('#chart_div').css("display", "block");
					drawChart();						
					$('#chart_div').css("display", "block");
					
					
					$('#histInfo').html("<h1>Propiedades</h1>"+
					"<br>Media: "+"<br>"+res.Data.Data2[0]+
					"<br>Varianza: "+"<br>"+res.Data.Data2[1]+
					"<br>Asimetria: "+"<br>"+res.Data.Data2[2]+
					"<br>Energia: "+"<br>"+res.Data.Data2[3]+
					"<br>Entropia:"+"<br>"+res.Data.Data2[4]);
					$('.info').css("display", "block");
					
					break;
				}
			}		
		};
		
		sendAjaxJSON(data, ok, true);		
	}
}

function drawChart() {	
	chart = new Highcharts.Chart({
		chart: {
			renderTo: 'chart_div', 	// Asignar div donde se mostrara
			defaultSeriesType: 'line'	// Pongo que tipo de gráfica es
		},
		title: {
			text: 'Histograma'	// Titulo (Opcional)
		},
		subtitle: {
			text: 'Imagen en niveles de gris'				// Subtitulo (Opcional)
		},
		// Pongo los datos en el eje de las 'X'
		xAxis: {
			//categories: ['Feb12','Mar12','Abr12','May12','Jun12','Jul12','Ago12','Sep12','Oct12','Nov12','Dic12','Ene13','Feb13','Mar13','Abr13','May13','Jun13'],
			// Pongo el título para el eje de las 'X'
			title: {
				text: 'Nivel de gris'
			}
		},
		yAxis: {
			// Pongo el título para el eje de las 'Y'
			title: {
				text: 'Frecuencia (%)'
			}
		},
		// Doy formato al la "cajita" que sale al pasar el ratón por encima de la gráfica
		tooltip: {
			enabled: true,
			formatter: function() {
				return '<b>'+ this.series.name +'</b><br/>'+
					this.x +': '+ this.y +' '+this.series.name;
			}
		},
		// Doy opciones a la gráfica
		plotOptions: {
			line: {
				dataLabels: {
					enabled: false
				},
				enableMouseTracking: true
			}
		},
		// Doy los datos de la gráfica para dibujarlas
		series: [{
            name: 'Imagen',
            data: P
        }],
	});
	
}