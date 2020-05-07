  angular.module("integral").directive("integralMath", function() {
       return  {
          controller: function ($scope, $element, $attrs, geolocation) {
             $scope.n = 1;
             $scope.geolocation = geolocation;
          },
          controllerAs: "derevo"

      };
  });
