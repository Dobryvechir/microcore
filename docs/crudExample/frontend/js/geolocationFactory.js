  angular.module("integral").factory("geolocation", function($http) {
     var factory = {};
     factory.cities=[];

     $http.get("/api/v1/account").then(function successCallback(response) {
                  if (response && response.data) {
                       factory.cities = response.data;
                  }
          }, function errorCallback(response) {
                 window.console.error("error:",response);
          }
     );

     factory.addCity = function() {
          var newItem = {name:"",amount:""};
          $http.post("/api/v1/account",newItem).then(function(res){
                factory.cities.push(res);
                window.console.log("Created ", res," final",factory.cities);
          },function(res){
                window.console.log("Failed to create ", newItem," final",factory.cities);
          });	
     } 

     factory.updateCity = function(index) {
        var item = factory.cities[index];
        if (item && item.id) {  
           $http.put("/api/v1/account/"+item.id,item).then(function(){
                window.console.log("Updated at "+index, item," final",factory.cities);
           },function(){
                window.console.log("Failed to update at "+index, item," final",factory.cities);
           });
        }	
     }

     factory.deleteCity = function(index) {
        var items = factory.cities.splice(index, 1);
        var id = items && items[0] && items[0].id;
        if (id) {  
           $http.delete("/api/v1/account/"+id).then(function(){
                window.console.log("Removed at "+index, id," final",factory.cities);
           },function(){
                window.console.log("Failed to remove at "+index, id," final",factory.cities);
           });
        }	
     }
     
     factory.pidsumok = function() {
         var sum=0;
         for(var i=0;i<factory.cities.length;i++) {
             var current = factory.cities[i].amount;
             if (typeof current === "string") {
                  current = parseFloat(current);
             }
             if (current) {
                  sum += current;
             }
         }
         return sum;
     }	
     return factory;
  });
