  angular.module("integral").directive("integralMin", function() {
       return  {
          controller: function ($scope, $element, $attrs, geolocation) {
             this.addCities = function() {
                   geolocation.addCity(); 
                   this.message = "New account added";
             };
             this.message = "Press the button to add new accounts";
          },
          controllerAs: "knopka"

      };
  });
