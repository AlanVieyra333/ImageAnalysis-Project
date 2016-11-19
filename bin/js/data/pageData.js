var myApp = angular.module('myApp1', []);
myApp.controller('OptionsController', ['$scope', function($scope) {
	$scope.idd = 333;
	$scope.idStart = function(){
		$scope.idd = 0;
	};
	$scope.idTmp = function(){
		$scope.idd = $scope.idd+1;
		return $scope.idd;
	};
    $scope.option = [
		{name:"Sección 1", items:[
								{name:"Etiquetado por Componentes Conexas", items:""},
								{name:"Binarizar", items:""},
								{name:"N. G.", items:""},
								{name:"Negativo", items:""},
								{name:"Suma", items:""},
								{name:"Resta", items:""},
								{name:"And", items:""},
								{name:"Or", items:""},
								{name:"Binarización por Canal", items:""},
								{name:"N. G. por Canal", items:""}
							]},
		{name:"Sección 2", items:[
								{name:"Mejora del contraste", items:""},
								{name:"Histograma", items:""},
								{name:"Propiedades del Histograma", items:""},
								{name:"Desplazamiento", items:""},
								{name:"Ensanchamiento", items:""},
								{name:"Estiramiento", items:""},
								{name:"Ecualización", items:""}
							]},
		{name:"Filtros", items:[
								{name:"Pasa Altas", items:[
														{name:"Laplaciano", items:""},
														{name:"Robert", items:""},
														{name:"Prewitt", items:""},
														{name:"Sobel", items:""}
														]},
								{name:"Pasa Bajas", items:[
														{name:"Promedio", items:""},
														{name:"Promedio Pesado", items:""}
														]},
								{name:"Mediana", items:""},
								{name:"Moda", items:""},
								{name:"Max", items:""},
								{name:"Min", items:""}
							]},
		{name:"Segmentación", items:[
								{name:"N. G.", items:""},
								{name:"N. G. inverso", items:""},
								{name:"Multinivel", items:""},
								{name:"Multinivel inverso", items:""}
							]},
		{name:"Morfología Matemática", items:[
								{name:"Binaria", items:[
														{name:"Apertura", items:""},
														{name:"Clausura", items:""},
														{name:"Cerco convexo interno", items:""},
														{name:"Cerco convexo externo", items:""},
														{name:"Trans. Hit & Miss", items:""},
														{name:"Adelgazamiento", items:""}
														]},
								{name:"Latices", items:[
														{name:"Filtro suavizador", items:""},
														{name:"Gradiente morfológica (+)", items:""},
														{name:"Gradiente morfológica (-)", items:""},
														{name:"Top Mat", items:""},
														{name:"Bot Mat", items:""}
														]}
							]},
		{name:"Watershed", items:""}
	];
}]);