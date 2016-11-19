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

	if(  this.id != "send" && ((ID >= 1 && ID<= 2) || (ID >= 9 && ID<= 10) || (ID >= 18 && ID<= 27) ) ){ // Get value.
		showPopup(ID);
	}else if(ID >= 5 && ID <= 8 ){ // Get other image.
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
			}else if(ID == 23){
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
					$('#histInfo').html("<h1>Propiedades</h1>"+
					"<br>Media: "+"<br>"+res.Data.Data2[0]+
					"<br>Varianza: "+"<br>"+res.Data.Data2[1]+
					"<br>Asimetria: "+"<br>"+res.Data.Data2[2]+
					"<br>Energia: "+"<br>"+res.Data.Data2[3]+
					"<br>Entropia:"+"<br>"+res.Data.Data2[4]);
					$('.info').css("display", "block");
					
					P = res.Data.Data1;
					$('#chart_div').css("display", "block");
					google.charts.load('current', {'packages':['corechart']});
					google.charts.setOnLoadCallback(drawChart);	
					break;
				}
			}		
		};
		
		sendAjaxJSON(data, ok, true);		
	}
}

function drawChart() {
	var data = google.visualization.arrayToDataTable([
	  ['Nivel de gris', 'Imagen'],
	  [0,  0.0]
	]);	
	//alert(P);
	for(var k=0; k<P.length; k++){
		data.addRow( [k, P[k]*100] );
	}	
	
	var options = {
		title: 'Histograma',
		hAxis: {title: 'Nivel de gris',  titleTextStyle: {color: '#333'}},
		vAxis: {title: 'Frecuencia (%)', titleTextStyle: {color: '#333'}, minValue: 0}
		//legend: {position: 'none'}
	};
	
	var chart = new google.visualization.AreaChart(document.getElementById('chart_div'));
	chart.draw(data, options);
}