/* =============================================================================================
 * >> mini-IST <<
 *
 * Customized AngularJS for mini-IST web app
 * 
 * 1. Sign (login/register)
 * 2. Top Navbar (logout)
 * 3. Machine List
 * 4. Testcase Suite
 * 5. Dashboard
 * ============================================================================================= */

/* ===============
 * 1. Sign Controller
 * 	- login
 * 	- register
 * ===============*/

(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('MistSignController', MistSignController);
    
    MistSignController.$inject = ['$http', '$state', '$scope', 'Notify', 'toaster', 'localStorageService', '$rootScope'];
    function MistSignController($http, $state, $scope, Notify, toaster, localStorageService, $rootScope) {
    	var vm = this; //ViewModel = this controller?
    	
    	activate();
    	
    	////////////////
    	
    	function activate() {
    		// bind here all data from the form
          	vm.account = {};
          	// place the message if something goes wrong
          	//vm.authMsg = '';
          	
          	//--@ toggle reg/login--
          	$scope.visible = true;
          	vm.toggleSign = function() {
          		$scope.visible = !$scope.visible;
          	};
          	//--& toggle reg/login--
          	
          	//--@ check existing token
          	//>> when already there, then refresh page without login
          	vm.checkToken = function() {
          		//try to get token in local storage
                vm.username = localStorageService.get("username");
                vm.token = localStorageService.get("token");
                vm.moment = localStorageService.get("moment");
                
                // now is when??
                var now = moment();
                
                if (vm.token != "" && vm.token != null && vm.token != undefined
                 && vm.username != "" && vm.username != null && vm.username != undefined ) {
                	//user already logined in
                	//save user name on page
            		//$rootScope.user.name = vm.username;		//can't get it here...?!
        			
        			// if 5 days go, then delete previous localstorage!
        			//--need re-login
        			var diffDays = parseInt(now.diff(vm.moment,'days'));
	                if ( diffDays > 5 ) {
	                	//call 'logout' api to delete previous token
		                vm.paramsLogout = {
							"action": "logout",
							"name": vm.username
						};

            			$http.defaults.headers.post["Authorization"] = vm.token;
		                $http
		          			.post('/logout', JSON.stringify(vm.paramsLogout))
			            	.success(function(data, status, config, headers, statusText) {
			              		// assumes if ok, response is an object with some data, if not, a string with error
			              		// customize according to your api
			              		if ( data.ret != 0 ) {
			                		Notify.alert( 
				                		"[SIGN007] Failed to delete previous user token! \n"+'ret:  '+
										'data.ret:   ' +data.ret+'\n'+
										'data.info:  ' +data.info+'\n'+
										'status:     ' +status+'\n',
										//'config:     ' +config+'\n'+
										//'headers:    ' +headers+'\n'+
										//'statusText: ' +statusText+'\n', 
				                		{status: 'danger', timeout: 2000}
				            		);
			              		} else {
			                		//clear all in local storage
			                		localStorageService.clearAll();
			                		toaster.pop({
						  				type: 'error',
						  				title: 'Session Expired!', 
						  				body: "[SIGN010] Local Storage Has Expired and Cleaned! <"+vm.moment+"> Please re-login!"
						  			});
			              		}
			            	})
			            	.error(function(data, status, config, headers, statusText) {
			          			//vm.authMsg = 'Server Request Error';
			                  	//notify
			                    Notify.alert( 
					                "[SIGN008] $http (AJAX) Failed for requesting removel of user token! \n"+
					                'data.ret:   ' +data.ret+'\n'+
									'data.info:  ' +data.info+'\n'+
									'status:     ' +status+'\n',
									//'config:     ' +config+'\n'+
									//'headers:    ' +headers+'\n'+
									//'statusText: ' +statusText+'\n',
					                {status: 'danger', timeout: 2000}
					            );
			        		});
		                
		                return false;
	                }
        			
        			//*--if goes here and continue, resume local storage...
        			
        			//roll up the 'login page' curtain
        			//toaster - good message
		  			toaster.pop({
		  				type: 'success',
		  				title: 'Welcome back!', 
		  				body: 'Nice to see you again, '+vm.username+'! <br/>Resume <'+vm.moment+'> <br/>'+
		  					  (5-diffDays) +' day(s) to expire.',
		  				onShowCallback: function () {
			                //>>jQuery here
			                //>>jQuery here
	            			$(".curtain.mist-sign").animate({
					    		opacity: 0,
					    		top: '-=300',
					    		zIndex: "-50"
					  		}, 500);
						}
		  			});	
                }
          	};
          	//--& check existing token
          	//immediately run when start/refresh page
          	vm.checkToken();
          	
          	//--@ login--
          	vm.login = function() {
            	//vm.authMsg = '';
            	var loginput = $scope.sign.account.loginName;

				var $userinput = $scope.sign.account.loginName;
				if(((/.+@.+\.[a-zA-Z]{2,4}$/).test($userinput))){
					vm.paramsLogin = {
						"action": "login",
						"email": $userinput,
						"password": $scope.sign.account.loginPassword
					};
				}else{
					vm.paramsLogin = {
						"action": "login",
						"name": $userinput,
						"password": $scope.sign.account.loginPassword
					};
				}

            	if(vm.loginForm.$valid) {
              		$http
                		.post('/login', JSON.stringify(vm.paramsLogin))
                		.success(function(data, status, config, headers, statusText) {
                  			// assumes if ok, response is an object with some data, if not, a string with error
                  			// customize according to your api
                  			if ( data.ret != 0 ) {
                    			//vm.authMsg = 'Incorrect credentials.';
                    			//notify
                    			Notify.alert( 
		                			"[SIGN003] Login Failed! Sorry! \n"+'ret:  '+
									'data.ret:   ' +data.ret+'\n'+
									'data.info:  ' +data.info+'\n'+
									'status:     ' +status+'\n',
									//'config:     ' +config+'\n'+
									//'headers:    ' +headers+'\n'+
									//'statusText: ' +statusText+'\n', 
		                			{status: 'danger', timeout: 2000}
		            			);
                  			} else {
                    			//$state.go('app.dashboard');
                    			//login good
                    			
                    			//save token in local storage
                    			//**save timestamp for expiration decision: YYYY-MM-DD hh:mm:ss
                    			var savetime = moment();
                    			localStorageService.set("username", data.name);
                    			localStorageService.set("email", data.email);
                    			localStorageService.set("token", data.Token);
                    			localStorageService.set("moment", savetime.format("YYYY-MM-DD hh:mm:ss"));
                    			
	                    		//save user name on page
	                    		$rootScope.user.name = data.name; 	//for framework (available now)
                    			
                    			//roll up the 'login page' curtain
                    			//toaster - good message
					  			toaster.pop({
					  				type: 'success',
					  				title: 'Login completed!', 
					  				body: 'Welcome, '+data.name+'! <br/>'+'Created <'+localStorageService.get("moment")+'>',
					  				onShowCallback: function () {
						                //>>jQuery here
		                    			$(".curtain.mist-sign").animate({
								    		opacity: 0,
								    		top: '-=300',
								    		zIndex: "-50"
								  		}, 500);
        							}
					  			});	
                  			}
                		})
                		.error(function(data, status, config, headers, statusText) {
                  			//vm.authMsg = 'Server Request Error';
		                  	//notify
		                    Notify.alert( 
				                "[SIGN004] Login $http (AJAX) Failed! Sorry! \n"+
				                'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n',
								//'config:     ' +config+'\n'+
								//'headers:    ' +headers+'\n'+
								//'statusText: ' +statusText+'\n',
				                {status: 'danger', timeout: 2000}
				            );
                		});
            	} else {
              		// set as dirty if the user click directly to login so we show the validation messages
              		/*jshint -W106*/
              		vm.loginForm.login_name.$dirty = true;
              		vm.loginForm.login_password.$dirty = true;
            	}
          	};
          	//--& login--
          	
          	//--@ register--
          	vm.register = function() {
          		//vm.authMsg = '';
          		vm.paramsReg = {
					"action": "join",
					"name": $scope.sign.account.regName,
					"email": $scope.sign.account.regEmail,
					"password": $scope.sign.account.regPassword1
				};
          		
          		if (vm.registerForm.$valid) {
              		$http
                		.post('/join', JSON.stringify(vm.paramsReg))
                		.success(function(data, status, config, headers, statusText) {
                  			// assumes if ok, response is an object with some data, if not, a string with error
                  			// customize according to your api
                  			if ( data.ret != 0 ) {
                    			//vm.authMsg = response;
                    			//notify
                    			Notify.alert( 
		                			"[SIGN002] Register $http (AJAX) Failed! Sorry! \n"+
									'data.ret:   ' +data.ret+'\n'+
									'data.info:  ' +data.info+'\n'+
									'status:     ' +status+'\n',
									//'config:     ' +config+'\n'+
									//'headers:    ' +headers+'\n'+
									//'statusText: ' +statusText+'\n', 
		                			{status: 'danger', timeout: 2000}
		            			);
                  			} else {
                    			//register good
                    			vm.toggleSign(); //to login block
                    			
                    			//toaster - good message
					  			toaster.pop({
					  				type: 'success',
					  				title: 'Registration Complete!', 
					  				body: "Hello, "+data.name+"! Welcome to join, please sign in to continue."
					  			});
                  			}
                		}) 
                		.error(function(data, status, config, headers, statusText) {
                  			//vm.authMsg = 'Server Request Error';
                  			//notify
		                    Notify.alert( 
				                "[SIGN005] Login $http (AJAX) Failed! Sorry! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n',
								//'config:     ' +config+'\n'+
								//'headers:    ' +headers+'\n'+
								//'statusText: ' +statusText+'\n', 
				                {status: 'danger', timeout: 2000}
				            );
                		});
        		} else {
          			// set as dirty if the user click directly to login so we show the validation messages
		              /*jshint -W106*/
		              vm.registerForm.reg_name.$dirty = true;
		              vm.registerForm.reg_email.$dirty = true;
		              vm.registerForm.reg_password1.$dirty = true;
		              vm.registerForm.reg_password2.$dirty = true;
		              vm.registerForm.reg_agreed.$dirty = true;
          
        		}
          	};
          	//--& register--
          	
    	}
    	
    	////////////////
    	
    }
    
})();

/* ========================
 * 2. Top Navbar controller
 * 	- in top-navbar.html
 * 	- log out button here
 * ========================*/
(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('MistLogoutController', MistLogoutController);

	MistLogoutController.$inject = ['$http', '$scope', 'Notify', 'toaster', 'localStorageService', '$rootScope', '$state'];
	function MistLogoutController($http, $scope, Notify, toaster, localStorageService, $rootScope, $state) {
		var vm = this; //ViewModel = this controller?
    	
    	activate();
    	
    	////////////////
    	
    	function activate() {
    		//--@ logout--
          	vm.logout = function() {
            	vm.paramsLogout = {
					"action": "logout",
					"name": $rootScope.user.name
				};

            	$http.defaults.headers.post["Authorization"] = localStorageService.get("token");
            	
        		//--post('/logout', JSON.stringify(vm.paramsLogout))
          		$http
          			.post('/logout', JSON.stringify(vm.paramsLogout))
	            	.success(function(data, status, config, headers, statusText) {
	              		// assumes if ok, response is an object with some data, if not, a string with error
	              		// customize according to your api
	              		if ( data.ret != 0 ) {
	                		//vm.authMsg = 'Incorrect credentials.';
	                		//notify
	                		Notify.alert( 
		                		"[LOGOUT001] Logout Failed! Sorry! \n"+'ret:  '+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n',
								//'config:     ' +config+'\n'+
								//'headers:    ' +headers+'\n'+
								//'statusText: ' +statusText+'\n', 
		                		{status: 'danger', timeout: 2000}
		            		);
	              		} else {
	                		//logout good
	                		
	                		//clear all in local storage
	                		var moment = localStorageService.get("moment");
	                		localStorageService.clearAll();
	                	
	                		//clear user name on page
	                    	$rootScope.user.name = "";	
                    			
	                		//roll down the 'login page' curtain
	                		toaster.pop({
				  				type: 'success',
				  				title: 'Logout completed!', 
				 				body: 'Goodbye, '+data.name+', see you next time! <br/>'+'Deleted <'+moment+'>',
				  				onShowCallback: function () {
					                //>>jQuery here
			                		$(".curtain.mist-sign").css({"zIndex": "180"}).animate({
							    		opacity: 1,
							    		top: '+=300',
							  		}, 2000, function() {
							  			$state.go('app.welcome');
							  		});
    							}
				  			});
	              		}
	            	})
	            	.error(function(data, status, config, headers, statusText) {
	          			//vm.authMsg = 'Server Request Error';
	                  	//notify
	                    Notify.alert( 
			                "[LOGOUT002] Logout $http (AJAX) Failed! Sorry! \n"+
			                'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
							//'config:     ' +config+'\n'+
							//'headers:    ' +headers+'\n'+
							//'statusText: ' +statusText+'\n',
			                {status: 'danger', timeout: 2000}
			            );
	        		});
          	};
          	//--& logout--
    	}
	}

})();

/* ===============
 * 3. Machine List
 * ===============*/

//service to get initial machine data
(function() {
    'use strict';

    angular
        .module('app.custom')
        .service('MistMachineLoader', MistMachineLoader);

    MistMachineLoader.$inject = ['$http', 'Notify'];
    function MistMachineLoader($http, Notify) {
        this.getMachineData = getMachineData;

        ////////////////

        function getMachineData(onReady) {
            
            $http.get('data/sample-machine.json')
		        .success(onReady)
		        .error(function(data, status, headers, config) {
		        	Notify.alert( 
		                "[MCHN001] Failure Getting Sample Data for Machine...\n\n"+
						'status:     ' +status+'\n',
						//'config:     ' +config+'\n'+
						//'headers:    ' +headers+'\n'+
						//'statusText: ' +statusText+'\n',
		                {status: 'danger', timeout: 2000}
		            );
		        });
            
          
        }
    }
})();
			
 
(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('MistMachineController', MistMachineController);
    
    MistMachineController.$inject = ['$http', '$scope', 'MistMachineLoader', '$filter', 'Notify',
    					   'ngTableParams', '$resource', '$timeout', 'ngTableDataService'];
    function MistMachineController($http, $scope, MistMachineLoader, $filter, Notify, 
    								ngTableParams, $resource, $timeout, ngTableDataService) {
    	var vm = this;
		
		activate();        
		
        ////////////////
        
        function activate() {
        	$scope.popShow = false;
        	$scope.topShow = false;
        	$scope.leftShow = false;
        	$scope.rightShow = false;
        	
        	$(".init-hide").removeClass("init-hide");
        	
        	//initially load machine data
        	MistMachineLoader.getMachineData(machineReady);
        	
        	//onReady = machineReady
        	function machineReady(data) {
        		//data from $http
        		vm.machines = data.machine_list;
        		
        		//set summary initially
        		$scope.totalNum = vm.machines.length;
        		$scope.myown = 3;
        		$scope.myuse = 5;
        		
        		//ngTable population
        		loadMachineNgTable();
        		
        		//animation show
        		$timeout(animationShow, 500);
        	}
        	
        	//load machine NgTable
        	function loadMachineNgTable() {
        		//build new tableParams everytime the data changes
        		vm.tableParams = new ngTableParams(
	                {
	                    page: 1, // show first page
	                    count: 10, // count per page
	                }, {
	            		total: vm.machines.length, // length of data
	            		getData: function($defer, params) {
	                		var filteredData = params.filter() ? $filter('filter')(vm.machines, params.filter()) : vm.machines;
	                		var orderedData = params.sorting() ? $filter('orderBy')(filteredData, params.orderBy()) : vm.machines;
	
	                		params.total(orderedData.length); // set total for recalc pagination
	                		$defer.resolve(orderedData.slice((params.page() - 1) * params.count(), params.page() * params.count()));
	            		}
	        		}
	        	); 
        	}
        	
        	function animationShow() {
        		$scope.topShow = true;
        		$timeout(function() {
        			$scope.leftShow = true;	
        		}, 500);
        		$timeout(function() {
        			$scope.rightShow = true;
        		}, 500);
        	}
        	
        	//* -----------------------------------------------------------------------------------------------
        	
        	// (onbeforesave) check name
        	vm.checkName = function(data) {
            	if (data == "" || data == null) {
            		Notify.alert( 
		                "[MCHN002] Empty name to save, invalid: \n"+data,
						'status:     ' +status+'\n',
						//'config:     ' +config+'\n'+
						//'headers:    ' +headers+'\n'+
						//'statusText: ' +statusText+'\n',
		                {status: 'danger', timeout: 2000}
		            );
            		
              		return '[ERROR] Empty name!';
            	}
          	};
        	
        	// on after save
        	vm.saveMachine = function(data, name, id) {
            	//vm.machine not updated yet
            	angular.extend(data, {
            		id: id,
            		name: name
            	});
            	
            	//reload ngTables
            	//vm.tableParams.count(10);
            	//vm.tableParams.reload();
            	loadMachineNgTable();
            	
            	Notify.alert( 
	                '[MCHN003] Saved Machine: '+ data.name + ' id:'+data.id,
					//'config:     ' +config+'\n'+
					//'headers:    ' +headers+'\n'+
					//'statusText: ' +statusText+'\n',
	                {status: 'success', timeout: 2000}
	            );
            	
            	// return $http.post('/saveUser', data);
          	}; 
        	
        	// remove machine (ng-click)
          	vm.removeMachine = function($index, name, id) {
          		for (var i = 0; i < vm.machines.length; i++) {
          			if (vm.machines[i].name == name & vm.machines[i].id == id) {
          				vm.machines.splice(i, 1);
          				break; //stop looping
          			}
          		}
            	//vm.machines.splice(index, 1);
            	$scope.totalNum = vm.machines.length;
            	
            	//reload ngTables
            	//vm.tableParams.count(10);
            	//vm.tableParams.reload();
            	loadMachineNgTable();
            	
            	Notify.alert( 
	                '[MCHN004] Removed Machine: '+ name + ' id:'+id,
					//'config:     ' +config+'\n'+
					//'headers:    ' +headers+'\n'+
					//'statusText: ' +statusText+'\n',
	                {status: 'success', timeout: 2000}
	            );
            	// return $http.post('/saveUser', data);
          	};
        	
        	// add machine
        	vm.addMachine = function() {
        		if ($scope.popShow) {
        			//do nothing, already shown
        			Notify.alert( 
		                '[MCHN005] Pop-up window Already shown!',
		                {status: 'danger', timeout: 2000}
		            );
        		} else {
        			$scope.popShow = true; //show pop-up for machine add
        		}
	            
          	};
          	
          	// add machine confirm (within pop-up)
          	vm.addMachineConfirm = function() {
          		var vfyName = $scope.addMachineName,
          			vfyHost = $scope.addMachineHost,
          			vfyIP = $scope.addMachineIP,
          			vfyDesc = $scope.addMachineDesc;
          		
          		// valid machine name?
          		if (vfyName == "" || vfyName == null || vfyName == undefined) {
          			Notify.alert( 
		                '[MCHN006] Need "machine name" for Machine-Add!',
		                {status: 'danger', timeout: 2000}
		            );
          			return false;
          		}
          		
          		// valid machine IP?
          		if (vfyIP == "" || vfyIP == null || vfyName == undefined) {
          			Notify.alert( 
		                '[MCHN007] Need "machine IP" for Machine-Add!',
		                {status: 'danger', timeout: 2000}
		            );
          			return false;
          		}
          		
          		// add new machine
          		vm.inserted = {
	              id: vm.machines.length+1,
	              name: vfyName,
	              hostname: (vfyHost ? vfyHost : "Default hostname"),
	              ip: vfyIP,
	              //id: null,
	              desc: (vfyDesc ? vfyDesc : "Default description"),
	              isNew: true
	            };
            	//vm.machines.push(vm.inserted);
            	vm.machines.unshift(vm.inserted);
          		
          		//hide pop-up
          		$scope.popShow = false;
          		
          		//reload ngTables
            	loadMachineNgTable();
            	
            	Notify.alert( 
	                '[MCHN008] New machine added! \n name: '+vfyName + '\n IP:'+vfyIP,
	                {status: 'success', timeout: 2000}
	            );
          	};
        }	
    }
	
})();

/* =================
 * 4. Testcase Suite
 * =================*/

//service to get initial testsuite data
(function() {
    'use strict';

    angular
        .module('app.custom')
        .service('MistTestsuiteLoader', MistTestsuiteLoader);

    MistTestsuiteLoader.$inject = ['$http', 'Notify', 'localStorageService'];
    function MistTestsuiteLoader($http, Notify, localStorageService) {
        this.getTestsuiteData = getTestsuiteData;

        ////////////////

        function getTestsuiteData(onReady) {
            
            /*$http.get('data/sample-testsuite.json')
		        .success(onReady)
		        .error(function(data, status, headers, config) {
		        	Notify.alert( 
		                "[TCST001] Failure Getting Sample Data for Testsuite...\n\n"+
						'status:     ' +status+'\n',
						//'config:     ' +config+'\n'+
						//'headers:    ' +headers+'\n'+
						//'statusText: ' +statusText+'\n',
		                {status: 'danger', timeout: 2000}
		            );
		        });*/
		    
		    //get initial 1-100 testsuites
		    // e.g. "/testsuites?user=mike32432487293&startpos=1&counter=10 "
    		var username = localStorageService.get("username"),
    			url="/testsuites?user="+username+"&startpos=1&counter=100";
    		
    		//call api
    		$http
        		.get(url, {
					    headers: {"Authorization": localStorageService.get("token")}
					 })
        		.success(onReady) 
        		.error(function(data, status, config, headers) {
                    Notify.alert( 
		                "[TSUITE008] Get Testsuite(s) $http (AJAX) GET Failed! \n"+
						'data.ret:   ' +data.ret+'\n'+
						'data.info:  ' +data.info+'\n'+
						'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
        		});
		    
        }
    }
})();

(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('TestsuiteController', TestsuiteController);
    
    TestsuiteController.$inject = ['$filter', 'ngTableParams', '$resource', '$timeout', 'Notify', 'SweetAlert', 'localStorageService',
    							   'ngTableDataService', '$http', '$scope', '$rootScope', 'MistTestsuiteLoader'];
    function TestsuiteController($filter, ngTableParams, $resource, $timeout, Notify, SweetAlert, localStorageService,
    								ngTableDataService, $http, $scope, $rootScope, MistTestsuiteLoader) {
    	var vm = this;	//should be this controller
    	
        activate();
		
        ////////////////
        
        function activate() {
        	//intialization
        	$scope.isPanelSlideShown = false;
        	$scope.popShow = false;
        	$scope.topShow = false;
        	$scope.leftShow = false;
        	$scope.rightShow = false;
        	
        	$(".init-hide").removeClass("init-hide");
        	
        	$scope.newsuite = {
    			name: "",
    			team: "",
    			lib: "",
    			tag: "",
    			desc: "",
    			sscript: "",
    			cscript: "",
    			escript: "",
    			testcase: ""
    		};
        	
        	//initially load testsuite data
        	MistTestsuiteLoader.getTestsuiteData(testsuiteReady);
        	
        	//onReady = testsuiteReady
        	function testsuiteReady(data) {
        		//data from $http
        		vm.testsuites = data.testsuites;
        		vm.fullcount = data.total;
        		
        		//set summary initially
        		testsuiteSummary();
        		
        		//ngTable population
        		loadTestsuiteNgTable();
        		
        		//animation show
        		$timeout(animationShow, 500);
        	}
        	
        	function testsuiteSummary() {
        		$scope.totalNum = vm.testsuites.length;
        		$scope.myown = 2;	//testing only
        		$scope.myuse = 3;	//testing only
        		
        		$scope.pageset = [];
        		
        		//build up pageset as per fullcount
        		var loop = parseInt(vm.fullcount/100)+1;
        		for (var i = 0; i < loop; i++) {
        			$scope.pageset.push((100*i+1)+"-"+100*(i+1));
        		}
        		$scope.currentSet = $scope.pageset[0];
        		//$scope.pageset = ["101-200","201-300","301-400","401-500","501-600","601-700","701-800","801-900","901-1000"];
        		
        		//resettableForm
        		$scope.resetForm = function(form) {
			        form.$setPristine();
			        $scope.newsuite = {
	        			name: "",
	        			team: "",
	        			lib: "",
	        			tag: "",
	        			desc: "",
	        			sscript: "",
	        			cscript: "",
	        			escript: "",
	        			testcase: ""
	        		};
			    };
        	}
        	
        	function loadTestsuiteNgTable() {
        		//build new tableParams everytime the data changes
        		vm.tableParams = new ngTableParams(
	                {
	                    page: 1, // show first page
	                    count: 10, // count per page
	                }, {
	            		total: vm.testsuites.length, // length of data
	            		getData: function($defer, params) {
	                		var filteredData = params.filter() ? $filter('filter')(vm.testsuites, params.filter()) : vm.testsuites;
	                		var orderedData = params.sorting() ? $filter('orderBy')(filteredData, params.orderBy()) : vm.testsuites;
	
	                		params.total(orderedData.length); // set total for recalc pagination
	                		$defer.resolve(orderedData.slice((params.page() - 1) * params.count(), params.page() * params.count()));
	            		}
	        		}
	        	); 
        	}
        	
        	function animationShow() {
        		$scope.topShow = true;
        		$timeout(function() {
        			$scope.leftShow = true;	
        		}, 500);
        		$timeout(function() {
        			$scope.rightShow = true;
        		}, 500);
        	}
        	
        	//*------------------------------------------------------------------------------------------------
        	
        	//fetch page set (ngTable)
        	vm.fetchSet = function(set, which) {
        		// set in "1-100, 101-200, 201-300, ..."
        		var startPos = parseInt(set.split("-")[0]), bound = parseInt(vm.fullcount);
        		
        		switch (which) {
        			case "prev":
        				if ( (startPos-100) < 0 ) {
        					return false; // already the 1st set
        				} else {
        					startPos = startPos - 100;
        				}
        				break;
        			case "next":
        				if ( (startPos+100) > bound) {
        					return false; // already the last set
        				} else {
        					startPos = startPos + 100;
        				}
        				break;
        			case "this":
        			default:
        		};
        		
        		$scope.currentSet = startPos+"-"+(startPos+100-1);
        		
			    // e.g. "/testsuites?user=mike32432487293&startpos=1&counter=10 "
	    		var username = localStorageService.get("username"),
	    			url="/testsuites?user="+username+"&startpos="+startPos+"&counter=100";
	    		
	    		//call api
	    		$http
	        		.get(url, {
						    headers: {"Authorization": localStorageService.get("token")}
						 })
	        		.success(function(data) {
	        			//data from $http
		        		vm.testsuites = data.testsuites;		        		
		        		
		        		//ngTable reload
		        		loadTestsuiteNgTable();
	        		}) 
	        		.error(function(data, status, config, headers) {
	                    Notify.alert( 
			                "[TSUITE008] Get Testsuite(s) $http (AJAX) GET Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
	        		});
        	};
        	
        	vm.checkAndxEdit = function($index, suitename, suitecreator) {
        		// $http to get detail for this testsuite
        		// e.g. "/testsuite?user=mike32432487293&name=testsuite_test01"
        		var url="/testsuite?user="+$rootScope.user.name+"&name="+suitename;
        		
        		//call api
        		$http
            		.get(url, {
					    headers: {"Authorization": localStorageService.get("token")}
					 })
            		.success(function(data, status, config, headers, statusText) {
              			if ( data.ret != 0 ) {
                			Notify.alert( 
	                			"[TSUITE006] Get Testsuite $http (AJAX) Failed! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n', 
	                			{status: 'danger', timeout: 2000}
	            			);
              			} else {
              				//Good! put into slide panel
              				$scope.thisSuite = {
			        			id: data.id,
								name: data.name,
								creator: data.creator,
								modifier: data.modifier,
								create_time: data.create_time,
								update_time: data.update_time,
								tag: data.tag ? data.tag : [],
								team: data.team,
								desc: data.desc,
								library: data.library,
								setup_script: data.setup_script,
								clean_script: data.clean_script,
								exe_script: data.exe_script,
			        		};
	            			
	            			// open sliding panel for details
        					$scope.isPanelSlideShown = true;
              			}
            		}) 
            		.error(function(data, status, config, headers, statusText) {
	                    Notify.alert( 
			                "[TSUITE008] Get Testsuite $http (AJAX) GET Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
            		});
        	};
        	
        	//tag related (slide panel)
        	vm.newTag = function() {
        		$scope.thisSuite.tag.push("new tag");	//default tag text: "new tag"
        	};
        	
        	vm.removeTag = function(index, thisTag) {
        		$scope.thisSuite.tag.splice(index, 1);
        	};
        	
        	//testsuite related
        	vm.removeTestsuite = function(index, name, id) {
        		// DELETE api
        		var paramDel = {
        			"action": "delete_testsuite",
        			"user": $rootScope.user.name,
        			"name": name	        			
        		};
        		
        		
        		/*
        		 * $http
            		.delete('/testsuite', JSON.stringify(paramDel),
            			{ headers: {"Authorization": localStorageService.get("token")} })
        		 */
        		$http({  
			        method: "DELETE",  
			        url: '/testsuite',  
			        data: JSON.stringify(paramDel),  
			        headers: {"Authorization": localStorageService.get("token")}  
			    })
        		.success(function(data, status, config, headers, statusText) {
          			if ( data.ret != 0 ) {
            			Notify.alert( 
                			"[TSUITE010] Delete Testsuite Exception! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n', 
                			{status: 'danger', timeout: 2000}
            			);
          			} else {
          				//Good! rebuild ngTable
          				for (var i = 0; i < vm.testsuites.length; i++) {
		          			if (vm.testsuites[i].id == id) {
		          				vm.testsuites.splice(i, 1);
		          				break; //stop looping
		          			}
		          		}
		            	//vm.machines.splice(index, 1);
		            	$scope.totalNum = vm.testsuites.length;
		            	
		            	//reload ngTables
		            	loadTestsuiteNgTable();
		            	
		            	Notify.alert( 
			                '[TSUITE011] Removed Testsuite: '+ name,
			                {status: 'success', timeout: 2000}
			            );
          			}
        		}) 
        		.error(function(data, status, config, headers, statusText) {
                    Notify.alert( 
		                "[TSUITE012] Delete Testsuite $http (DELETE) Failed! \n"+
						'data.ret:   ' +data.ret+'\n'+
						'data.info:  ' +data.info+'\n'+
						'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
        		});
        	};
        	
        	vm.saveTestsuite = function() {
        		// PUT api
        		var paramUpdate = {
        			"action": "update_testsuite",
        			"user": $rootScope.user.name, 
        			"testsuite": {
        				"id": $scope.thisSuite.id,
        				"name": $scope.thisSuite.name,
        				"tag": $scope.thisSuite.tag,
        				"team": $scope.thisSuite.team,
        				"library": $scope.thisSuite.library,
	        			"desc": $scope.thisSuite.desc,
	        			"setup_script": $scope.thisSuite.setup_script, 
	        			"clean_script": $scope.thisSuite.clean_script, 
	        			"exe_script": $scope.thisSuite.exe_script
        			}	
        		};
        		
        		/*debug
        		 * 
        		SweetAlert.swal(
        			"Great, you're to create the below:", 
        			'id:     '+$scope.thisSuite.id+'\n'+
        			'name:     '+$scope.thisSuite.name+'\n'+
        			'team:     '+$scope.thisSuite.team+'\n'+
        			'lib:      '+$scope.thisSuite.library+'\n'+
        			'tag:      '+$scope.thisSuite.tag+'\n'+
        			'desc:     '+$scope.thisSuite.desc+'\n'+
        			's_script: '+$scope.thisSuite.setup_script+'\n'+
        			'c_script: '+$scope.thisSuite.clean_script+'\n'+
        			'e_script: '+$scope.thisSuite.exe_script+'\n', 
        			'warning'
        		);*/
        		
        		$http
            		.put('/testsuite', JSON.stringify(paramUpdate), 
            			{ headers: {"Authorization": localStorageService.get("token")} })
            		.success(function(data, status, config, headers, statusText) {
              			if ( data.ret != 0 ) {
                			Notify.alert( 
	                			"[TSUITE009] Update Testsuite $http (PUT) Failed! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n', 
	                			{status: 'danger', timeout: 2000}
	            			);
              			} else {
              				//Good! put into slide panel
              				$scope.thisSuite = {
			        			id: data.id,
								name: data.name,
								creator: data.creator,
								modifier: data.modifier,
								create_time: data.create_time,
								update_time: data.update_time,
								tag: data.tag,
								team: data.team,
								desc: data.desc,
								library: data.library,
								setup_script: data.setup_script,
								clean_script: data.clean_script,
								exe_script: data.exe_script,
			        		};
			        		
			        		//rebuild ngTable
	          				for (var i = 0; i < vm.testsuites.length; i++) {
			          			if (vm.testsuites[i].id == $scope.thisSuite.id) {
			          				//sync update
			          				vm.testsuites[i].name = $scope.thisSuite.name;
			          				vm.testsuites[i].update_time = $scope.thisSuite.update_time;
			          				vm.testsuites[i].team = $scope.thisSuite.team;
			          				vm.testsuites[i].library = $scope.thisSuite.library;
			          				
			          				break; //stop looping
			          			}
			          		}
			          		
			          		//reload ngTables
		            		loadTestsuiteNgTable();
	            			
	            			Notify.alert( 
	                			"[TSUITE001] Testsuite Update Saved! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+ 
								data.name+' - '+data.id,
	                			{status: 'success', timeout: 2000}
	            			);
              			}
            		}) 
            		.error(function(data, status, config, headers, statusText) {
	                    Notify.alert( 
			                "[TSUITE007] Get Testsuite $http (GET) Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
            		});
            	
        	};
        	
        	// add testsuite
        	vm.addTestsuite = function() {
        		if ($scope.popShow) {
        			//do nothing, already shown
        			Notify.alert( 
		                '[TSUITE002] Pop-up window Already shown!',
		                {status: 'danger', timeout: 2000}
		            );
        		} else {
        			$scope.popShow = true; //show pop-up
        		}
          	};
        	
        	// submit add testsuite (form ng-submit)
        	vm.submitAddTestsuite = function(form) {
        		/*debug:
        		SweetAlert.swal(
        			"Great, you're to create the below:", 
        			'name:     '+$scope.newsuite.name+'\n'+
        			'team:     '+$scope.newsuite.team+'\n'+
        			'lib:      '+$scope.newsuite.lib+'\n'+
        			'tag:      '+$scope.newsuite.tag+'\n'+
        			'desc:     '+$scope.newsuite.desc+'\n'+
        			's_script: '+$scope.newsuite.sscript+'\n'+
        			'c_script: '+$scope.newsuite.cscript+'\n'+
        			'e_script: '+$scope.newsuite.escript+'\n'+
        			'testcase: '+$scope.newsuite.testcase'\n'+, 
        			'success'
        		);*/
        		
        		var paramInsert = {
        			"action": "insert_testsuite",
        			"user": $rootScope.user.name,
        			"testsuite": {
        				"name": $scope.newsuite.name,
        				"team": $scope.newsuite.team,
        				"library": $scope.newsuite.lib,
	        			"desc": $scope.newsuite.desc,
	        			"tag": $scope.newsuite.tag.trim() == "" ? [] : $scope.newsuite.tag.trim().split(" "),
	        			"setup_script": $scope.newsuite.sscript, 
	        			"clean_script": $scope.newsuite.cscript, 
	        			"exe_script": $scope.newsuite.escript,
	        			"testcases": $scope.newsuite.testcase.trim() == "" ? [] : $scope.newsuite.testcase.trim().split(" ")
        			}	
        		};
        		
        		//call api
        		$http.defaults.headers.post["Authorization"] = localStorageService.get("token"); 
        		$http
            		.post('/testsuite', JSON.stringify(paramInsert))
            		.success(function(data, status, config, headers, statusText) {
              			if ( data.ret != 0 ) {
                			Notify.alert( 
	                			"[TSUITE003] Insert Testsuite $http (POST) Failed! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n', 
	                			{status: 'danger', timeout: 2000}
	            			);
              			} else {
              				//Good! put into ngTable for display
			          		vm.inserted = {
				              id: data.id ,
				              creator: data.creator,
				              name: data.name,
				              team: data.team,
				              library: data.library,
				              update_time: data.update_time,
				              isNew: true
				            };
			            	vm.testsuites.unshift(vm.inserted);
	            			
	            			//reload ngTable population
        					loadTestsuiteNgTable();
        					
                			Notify.alert( 
	                			"[TSUITE004] New Testsuite Inserted! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+ 
								data.name+'-'+data.id,
	                			{status: 'success', timeout: 2000}
	            			);
	            			
	            			//debug:
	            			$scope.resetForm(form);
	            			$scope.popShow = false; //hide pop-up
              			}
            		}) 
            		.error(function(data, status, config, headers, statusText) {
              			//vm.authMsg = 'Server Request Error';
              			//notify
	                    Notify.alert( 
			                "[TSUITE005] Insert Testsuite $http (POST) Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
            		});
        		
        		//normally do this whether success or fail
        		//--$scope.popShow = false; //hide pop-up
        	};
        	
        	// cancel new testsuite
        	vm.cancelTestsuite = function() {
        		$scope.popShow = false; //hide popup
          	};
        	
        	///////////////////
        	
        	//*---Sliding panel related
        	
        	initSlidePanel();
        	
        	function initSlidePanel() {
        		$scope.thisSuite = {
        			id:"5729ed8f537d9362b7ae8386",
					name:"testsuite_test02",
					creator:"mike32432487293",
					modifier:"mike32432487293",
					create_time:"2016-05-04T20:39:43.774+08:00",
					update_time:"2016-05-04T20:57:04.367+08:00",
					team:"testteam02",
					desc:"desc02",
					library:"lib02",
					setup_script:"setupScript02",
					clean_script:"cleanScript02",
					exe_script:"exeScript02",
					tag:["tagAA","tagBB"]
        		};
        	}
        	
        	// PANEL COLLAPSE EVENTS (use panel id name for the boolean flag)
        	/*$scope.$watch('testsuiteSlidePanel', function(newVal) {
        		Notify.alert('Testing panel collapsed', {status: 'warning', timeout: 2000});
        	});*/
        	
        	// PANEL DISMISS EVENTS - Before remove panel
        	$scope.$on('panel-remove', function(event, id, deferred){
            
            // Here is obligatory to call the resolve() if we pretend to remove the panel finally
            // Not calling resolve() will NOT remove the panel
            // It's up to your app to decide if panel should be removed or not
            //deferred.resolve();
          		//Notify.alert('Testing closing panel', {status: 'warning', timeout: 2000});
          		$scope.isPanelSlideShown = false;
          	});
          	
          	// PANEL REFRESH EVENTS
          	$scope.$on('panel-refresh', function(event, id) {
            	Notify.alert('Testing refreshing panel'+id, {status: 'warning', timeout: 2000});
            	
            	$timeout(function(){
              		// directive listen for to remove the spinner 
              		// after we end up to perform own operations
              		$scope.$broadcast('removeSpinner', id);
              	}, 2000);

          	});
        }
        
    }
    
})();

/* ============
 * 5. Dashboard
 * ============*/
(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('MistDashboardController', MistDashboardController);
        
    MistDashboardController.$inject = ['$timeout', '$http', '$scope', '$rootScope', 'Notify'];
    function MistDashboardController($timeout, $http, $scope, $rootScope, Notify) {
    	var vm = this;	//this controller (vm = view model??)
    	
    	$scope.popShow = false;
    	
    	activate();
    	
    	////////////////
        
        function activate() {
        	vm.togglePopShow = function() {
        		$scope.popShow = !$scope.popShow;
        		
        		Notify.alert(
	                'Toggled popShow, now:'+$scope.popShow,
					//'config:     ' +config+'\n'+
					//'headers:    ' +headers+'\n'+
					//'statusText: ' +statusText+'\n',
	                {status: 'success', timeout: 1000}
	            );
        		
        	};
        }
    	
    }

})();    
 
/* =====================
 * 6. WorkloadController
 * =====================*/
(function() {
    'use strict';

    angular
        .module('app.custom')
        .service('WorkloadLoader', WorkloadLoader);

    WorkloadLoader.$inject = ['$http', 'Notify', 'localStorageService'];
    function WorkloadLoader($http, Notify, localStorageService) {
        this.getWorkloadGroupData = getWorkloadGroupData;

        ////////////////

        function getWorkloadGroupData(onReady) {
            
            $http.get('data/sample-workload-groups.json')
		        .success(onReady)
		        .error(function(data, status, headers, config) {
		        	Notify.alert( 
		                "[WKLD000] Failure Getting Sample Data for Workload Groups...\n\n"+
						'status:     ' +status+'\n',
						//'config:     ' +config+'\n'+
						//'headers:    ' +headers+'\n'+
						//'statusText: ' +statusText+'\n',
		                {status: 'danger', timeout: 2000}
		            );
		        });
		    
		    /*
    		var username = localStorageService.get("username"),
    			url="/testsuites?user="+username+"&startpos=1&counter=100";
    		
    		//call api
    		$http
        		.get(url, {
					    headers: {"Authorization": localStorageService.get("token")}
					 })
        		.success(onReady) 
        		.error(function(data, status, config, headers) {
                    Notify.alert( 
		                "[TSUITE008] Get Testsuite(s) $http (AJAX) GET Failed! \n"+
						'data.ret:   ' +data.ret+'\n'+
						'data.info:  ' +data.info+'\n'+
						'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
        		});*/
		    
        }
    }
})();

(function() {
	'use strict';
	
	angular
        .module('app.custom')
        .controller('WorkloadController', WorkloadController);
        
    WorkloadController.$inject = ['$timeout', '$http', '$scope', '$rootScope', 'Notify', 'WorkloadLoader', 'localStorageService',
    							  'editableOptions', 'editableThemes', '$resource'];
    function WorkloadController($timeout, $http, $scope, $rootScope, Notify, WorkloadLoader, localStorageService,
    							 editableOptions, editableThemes, $resource) {
    	var vm = this;	//this controller (vm = view model??)
    	
    	activate();
    	
    	////////////////
        
        function activate() {
        	editableOptions.theme = 'bs3';
        	
        	editableThemes.bs3.inputClass = 'input-sm';
          	editableThemes.bs3.buttonsClass = 'btn-sm';
          	editableThemes.bs3.submitTpl = '<button type="submit" class="btn btn-success"><span class="fa fa-check"></span></button>';
          	editableThemes.bs3.cancelTpl = '<button type="button" class="btn btn-default" ng-click="$form.$cancel()">'+
                                           '<span class="fa fa-times text-muted"></span>'+
            	                           '</button>';
        	
        	$scope.isSrchforShown = false;
        	$scope.isSrchforCollapsed = false;
        	$scope.rightShow = false;
        	$scope.popShow1 = false;
        	$scope.popShow2 = false;
        	$scope.isLevel1Shown = false;
        	$scope.isLevel2Shown = false;
        	$scope.isLevel3ProfileShown = false;
        	$scope.isLevel3RunShown = false;
        	$scope.workloadProfile = {
        		id: "",
        		name: "",
        		creator: "",
        		modifier: "",
        		create_time: "",
        		update_time: "",
        		team: "",
        		runuser: "",
        		desc: "",
        		tags: [],
        		setup_script: "",
        		clean_script: "",
        		exe_script: "",
        		env_script: "",
        		groups: [],
        		machines: [],
        		testsuites: []	
        	};
        	
        	$(".init-hide").removeClass("init-hide");
        	
        	vm.groupSrchfor = {
        		name: "n/a",
        		team: "n/a",
        		workload: "n/a"
        	};
        	
        	
        	// PANEL DISMISS EVENTS - Before remove panel
        	$scope.$on('panel-remove', function(event, id, deferred){
          		$scope.isSrchforShown = false;
          	});
        	
        	//initially load testsuite data
        	WorkloadLoader.getWorkloadGroupData(workloadGrounpReady);
        	
        	//onReady = testsuiteReady
        	function workloadGrounpReady(data) {
        		//data from $http
        		vm.groups = data.groups;
        		vm.fullcount = data.total_count;
        		vm.groupSets = [];	// 2-dimension
        		
        		//build up groups array for HEX display
        		buildGroupArray();
        		
        		$scope.currentSet = '1-100';
        		$scope.pageSets = ['1-100','101-200','201-300','301-400','401-500'];
        		
        		$timeout(function(){
              		$scope.isLevel1Shown = true;
              		$scope.rightShow = true;
              	}, 500);
        		
        	}
        	
        	function buildGroupArray() {
        		// 10 groups in one set
        		var tempset = [];
        		for (var i = 0; i < vm.groups.length; i++) {
        			if ( i > 0 & i % 10 == 0) {
        				vm.groupSets.push(tempset);
        				tempset = [];
        			}
        			tempset.push(vm.groups[i]);
        		}
        		vm.groupSets.push(tempset);
        		
        		//init 1st set
        		vm.groupSet = vm.groupSets[0];
        		vm.pagTotalNum = vm.groups.length;
        		vm.pagCurrentPage = 1;
        		vm.pagMaxSize = 5;
        		vm.pagNumPages = vm.groupSets.length;
        		
        		vm.groupSet.forEach(function(thisOne) {
        			thisOne.bgColorClass = vm.whatBgClass();
        			thisOne.workloadCount = thisOne.workloads ? thisOne.workloads.length : 0;
        		});
        	}
        	
        	/////////////////////////////////////
        	
        	vm.isHexGap = function(index) {
        		return (index % 7 == 0) ? true : false;
        	};
        	
        	vm.pagGroupChanged = function(currentPage) {
        		vm.groupSet = vm.groupSets[currentPage-1];
        		vm.groupSet.forEach(function(thisOne) {
        			thisOne.bgColorClass = vm.whatBgClass();
        			thisOne.wokloadCount = thisOne.workloads.length;
        		});
        	};
        	
        	vm.whatBgClass = function() {
        		var n = Math.floor(Math.random() * 6) + 1, m = Math.floor(Math.random() * 6) + 1;
        		return "bg-linear-color"+n+"-angle"+m;
        	};
        	
        	vm.animate = function() {
        		var n = Math.floor(Math.random() * 5) + 1, m = Math.floor(Math.random() * 5) + 1, usethis;
        		switch (n) {
        			case 1: usethis = "animate-fadeUp"; break;
        			case 2: usethis = "animate-flipX"; break;
        			case 3: usethis = "animate-roll"; break;
        			case 4: usethis = "animate-zoom"; break;
        			case 5: usethis = "animate-rotate"; break;
        			default: usethis = "animate-fadeUp";
        		}
        		return usethis;
        	};
        	
        	vm.groupSlideShow = function() {
        		//emulate $http for the group
        		$http.get('data/sample-group-1.json')
		        .success(function(data) {
		        	$scope.thisGroup = {
	        			id: data.id,
						name: data.name,
						creator: data.creator,
						modifier: data.modifier,
						create_time: data.create_time,
						update_time: data.update_time,
						team: data.team,
						desc: data.desc,
						workload_cnt: data.workloads.length,
						workloads: data.workloads
	        		};
	        		$scope.isGroupSlideShown = true;
		        })
		        .error(function(data, status, headers, config) {
		        	Notify.alert( 
		                "[WKLD004] Failure Getting Sample Data for One Group...\n\n"+'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
		        });
        	};
        	
        	vm.groupEditWorkloadRomove = function(index, workload) {
        		$scope.thisGroup.workloads.splice(index,1);
        		$scope.thisGroup.workload_cnt = $scope.thisGroup.workloads.length;
        	};
        	
        	vm.groupEditSave = function() {
        		//rest API for removing workload in the group
        		Notify.alert( 
	                "[Debug] Here calls the rest API for removing workload in the group: \n\n"
	                	+'gorup id:     ' +$scope.thisGroup.id+'\n'
	                	+'gorup name:   ' +$scope.thisGroup.name+'\n'
	                	+'gorup team:   ' +$scope.thisGroup.team+'\n'
	                	+'gorup desc:   ' +$scope.thisGroup.desc+'\n'
	                	+'workload:     ' +$scope.thisGroup.workloads.length+'\n',
	                {status: 'warning', timeout: 2000}
	            );
        	};
        	
        	vm.groupCheckIn = function(gid, gname) {
        		//go into level2 to show workload(s) in the group
        		
        		//emulate $http for the group
        		$http.get('data/sample-workload.json')
		        .success(function(data) {
		        	vm.workloadSet = data.workloads;
		        	vm.workloadSet.forEach(function(thisOne) {
	        			thisOne.bgColorClass = vm.whatBgClass();
	        		});
		        	
		        	$scope.isLevel1Shown = false;
		        	$scope.isLevel2Shown = true;
		        })
		        .error(function(data, status, headers, config) {
		        	Notify.alert( 
		                "[WKLD001] Failure Getting Sample Data for Workloads...\n\n"+'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
		        });
        	};
        	
        	vm.backToGroup = function() {
        		$scope.isLevel2Shown = false;
		        $scope.isLevel1Shown = true;
        	};
        	
        	vm.checkWorkloadProfile = function(workload, index) {
        		vm.checkWorkloadIndex = index;
        		
        		//go into level3 to show specific workload profile
        		$scope.workloadProfile = {
        			id: workload.id,
	        		name: workload.name,
	        		creator: workload.creator,
	        		modifier: workload.modifier,
	        		create_time: workload.create_time,
	        		update_time: workload.update_time,
	        		team: workload.team,
	        		runuser: workload.run_as_user,
	        		desc: workload.desc,
	        		tags: workload.tag,
	        		setup_script: workload.setup_script,
	        		clean_script: workload.clean_script,
	        		exe_script: workload.exe_script,
	        		env_script: workload.env_script,
	        		machines: workload.machine,
	        		groups: workload.group,
	        		testsuites: workload.testsuites	
	        	};
        		
        		$scope.isLevel3ProfileShown = true;
        		$scope.isLevel2Shown = false;
        	};
        	
        	vm.backFromWorkloadProfile = function() {
        		$scope.isLevel3ProfileShown = false;
		        $scope.isLevel2Shown = true;
        	};
        	
        	vm.runWorkload = function(workload) {
        		$scope.workloadProfile = {
        			id: workload.id,
	        		name: workload.name,
	        		creator: workload.creator,
	        		modifier: workload.modifier,
	        		create_time: workload.create_time,
	        		update_time: workload.update_time,
	        		team: workload.team,
	        		runuser: workload.run_as_user,
	        		desc: workload.desc,
	        		tags: workload.tag,
	        		setup_script: workload.setup_script,
	        		clean_script: workload.clean_script,
	        		exe_script: workload.exe_script,
	        		env_script: workload.env_script,
	        		machines: workload.machine,
	        		groups: workload.group,
	        		testsuites: workload.testsuites	
	        	};
        		
        		$scope.isLevel3RunShown = true;
        		$scope.isLevel2Shown = false;
        	};
        	
        	vm.backFromWorkloadRun = function() {
        		$scope.isLevel3RunShown = false;
		        $scope.isLevel2Shown = true;
        	};
        	
        	vm.removeWorkload = function(workload, index) {
        		var paramDel = {
        			"action": "delete_workload",
        			"user": $rootScope.user.name,
        			"name": workload.name	        			
        		};
        		
        		$http({  
			        method: "DELETE",  
			        url: '/workload',  
			        data: JSON.stringify(paramDel),  
			        headers: {"Authorization": localStorageService.get("token")}  
			    })
        		.success(function(data, status, config, headers, statusText) {
          			if ( data.ret != 0 ) {
            			Notify.alert( 
                			"[DEBUG] Delete Worklaod Exception! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n', 
                			{status: 'danger', timeout: 2000}
            			);
          			} else {
          				//Good!
          				vm.workloadSet.splice(index,1);
          				
		            	Notify.alert( 
			                '[WKLD006] Removed Workload: '+ data.name,
			                {status: 'success', timeout: 2000}
			            );
          			}
        		}) 
        		.error(function(data, status, config, headers, statusText) {
                    Notify.alert( 
		                "[TSUITE012] Delete Testsuite $http (DELETE) Failed! \n"+
						'data.ret:   ' +data.ret+'\n'+
						'data.info:  ' +data.info+'\n'+
						'status:     ' +status+'\n',
		                {status: 'danger', timeout: 2000}
		            );
        		});
        	};
        	
        	///////////////////
        	
        	//*Sliding Group panel (Edit)
        	
        	/*initSlidePanel();
        	function initSlidePanel() {
        		//group slide
        		$scope.thisGroup = {
        			id:"testingGroupId001",
					name:"group_test01",
					creator:"mars487293",
					modifier:"mars487293",
					create_time:"2016-05-04T20:39:43.774+08:00",
					update_time:"2016-05-04T20:57:04.367+08:00",
					team:"Mars",
					desc:"This is a sample workload group for testing only",
					workload_cnt: 7,
					total_run: 100,
					good_run: 98,
					bad_run: 2,
					tag:["tagAA","tagBB"]
        		};
        	}*/
        	
        	// PANEL COLLAPSE EVENTS (use panel id name for the boolean flag)
        	/*$scope.$watch('testsuiteSlidePanel', function(newVal) {
        		Notify.alert('Testing panel collapsed', {status: 'warning', timeout: 2000});
        	});*/
        	
        	// PANEL DISMISS EVENTS - Before remove panel
        	$scope.$on('panel-remove', function(event, id, deferred){
          		Notify.alert('Testing closing panel, id='+id, {status: 'warning', timeout: 2000});
          		$scope.isGroupSlideShown = false;
          	});
          	
          	// PANEL REFRESH EVENTS
          	$scope.$on('panel-refresh', function(event, id) {
            	Notify.alert('Testing refreshing panel'+id, {status: 'warning', timeout: 2000});
            	
            	$timeout(function(){
              		// directive listen for to remove the spinner 
              		// after we end up to perform own operations
              		$scope.$broadcast('removeSpinner', id);
              	}, 2000);

          	});
        	
        	///////////////////
        	
        	//*---Creating new group
        	vm.groupadd = {
        		name: "",
        		team: "",
        		desc: "",
        		workloads: []
        	};
        	
        	//resettableForm
    		vm.resetGroupForm = function(form) {
		        form.$setPristine();
		        vm.groupadd = {
	        		name: "",
	        		team: "",
	        		desc: "",
	        		workloads: []
	        	};
		    };
        	
        	//emulate $resource for the workload selection list
        	var Wlist = $resource('data/sample-workload-candidate.array', {}, {'query': {method:'GET', isArray:true} });
        		
        	vm.wlist = Wlist.query(); 
        	
        	vm.groupWorkloadAdd = function(workload) {
        		vm.groupadd.workloads.push(workload);
        	};
        	
        	vm.groupWorkloadRomove = function(index, workload) {
        		vm.groupadd.workloads.splice(index,1);
        	};
        	
        	vm.groupAddCancel = function(form) {
	        	vm.resetGroupForm(form);
	        	$scope.popShow1 = false;
        	};
        	
        	vm.groupAddConfirm = function(form) {
        		//**rest API: create group
        		Notify.alert('Testing: create group API calling here..\n'
        			+'\n name: '+vm.groupadd.name
        			+'\n team: '+vm.groupadd.team
        			+'\n desc: '+vm.groupadd.desc
        			+'\n workloads: '+vm.groupadd.workloads.length, 
        			{status: 'warning', timeout: 2000}
        		);
        		
        		vm.resetGroupForm(form);
	        	$scope.popShow1 = false;
        	};
        	
        	///////////////////////
        	
        	//*---Creating new workload
        	vm.wkldadd = {
        		name: "",
        		team: "",
        		desc: "",
        		runuser: "",
        		tags: [],
        		sscript: "",
        		cscript: "",
        		xscript: "",
        		escript: "",
        		groups: [],
        		machines: [],
        		testsuites: []
        	};
        	
        	//resettableForm
    		vm.resetWorkloadForm = function(form) {
		        form.$setPristine();
		        vm.wkldadd = {
	        		name: "",
	        		team: "",
	        		desc: "",
	        		runuser: "",
	        		tags: [],
	        		sscript: "",
	        		cscript: "",
	        		xscript: "",
	        		escript: "",
	        		groups: [],
	        		machines: [],
	        		testsuites: []
	        	};
	        	
	        	$scope.tagInput = "";
	        	$scope.gSelect = "";
	        	$scope.tsSelect = "";
	        	$scope.mchSelect = "";
		    };
        	
        	//emulate $resource for the workload selection list
        	var TSlist = $resource('data/sample-testsuite-candidate.array', {}, {'query': {method:'GET', isArray:true} }),
        		Glist = $resource('data/sample-group-candidate.array', {}, {'query': {method:'GET', isArray:true} }),
        		Mlist = $resource('data/sample-machine-candidate.array', {}, {'query': {method:'GET', isArray:true} });
        		
        	vm.tslist = TSlist.query();
        	vm.glist = Glist.query();
        	vm.mchlist = Mlist.query(); 
        	
        	vm.workloadNewItemAdd = function(thisname, category) {
        		switch (category) {
        			case 'tag': 		vm.wkldadd.tags.push(thisname); $scope.tagInput = ""; break;
        			case 'group':		vm.wkldadd.groups.push(thisname); break;
        			case 'testsuite':	vm.wkldadd.testsuites.push(thisname); break;
        			case 'machine':		vm.wkldadd.machines.push(thisname); break;
        			default:
        		};
        	};
        	
        	vm.workloadNewItemRemove = function(index, thisname, category) {
        		Notify.alert('[Debug] remove category:'+category+' index:'+index, {status: 'warning', timeout: 2000});
        		switch (category) {
        			case 'tag': 		vm.wkldadd.tags.splice(index,1); break;
        			case 'group':		vm.wkldadd.groups.splice(index,1); break;
        			case 'testsuite':	vm.wkldadd.testsuites.splice(index,1); break;
        			case 'machine':		vm.wkldadd.machines.splice(index,1); break;
        			default:
        		};
        	};
        	
        	vm.workloadAddCancel = function(form) {
        		vm.resetWorkloadForm(form);
	        	$scope.popShow2 = false;
        	};
        	
        	vm.workloadAddConfirm = function(form) {
        		//**rest API: create workload
        		var paramInsert = {
        			"action": "insert_workload",
        			"user": $rootScope.user.name,
        			"workload": {
        				"name": vm.wkldadd.name,
        				"team": vm.wkldadd.team,
        				"run_as_user": vm.wkldadd.runuser,
	        			"desc": vm.wkldadd.desc,
	        			"tag": vm.wkldadd.tags,
	        			"setup_script": vm.wkldadd.sscript, 
	        			"clean_script": vm.wkldadd.cscript, 
	        			"exe_script": vm.wkldadd.xscript,
	        			"env_script": vm.wkldadd.escript,
	        			"machine": vm.wkldadd.machines,
	        			"testsuites": vm.wkldadd.testsuites,
	        			"groups": vm.wkldadd.groups
        			}	
        		};
        		
        		//call api
        		$http.defaults.headers.post["Authorization"] = localStorageService.get("token");
        		$http.defaults.useXDomain = true; 
        		$http
            		.post('/workload', JSON.stringify(paramInsert))
            		.success(function(data, status, config, headers, statusText) {
              			if ( data.ret != 0 ) {
                			//in case of any problem
              			} else {
              				Notify.alert('[WKLD002] New Workload Created!\n'
			        			+'\n name: '+data.name
			        			+'\n team: '+data.team
			        			+'\n desc: '+data.desc
			        			+'\n run_as_user: '+data.run_as_user
			        			+'\n tags: '+data.tag.length
			        			//+'\n groups: '+vm.wkldadd.groups.length
			        			+'\n testsuites: '+data.testsuites.length
			        			+'\n machines: '+data.machine.length,
			        			{status: 'success', timeout: 2000}
			        		);
	            			
	            			//debug:
	            			vm.resetWorkloadForm(form);
	        				$scope.popShow2 = false;
              			}
            		}) 
            		.error(function(data, status, config, headers, statusText) {
	                    Notify.alert( 
			                "[WKLD003] Insert Workload $http (POST) Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
            		});
            	
        	};
        	
        	////////////////////////
        	
        	//workload profile edit (level3)
        	vm.workloadEditTagRemove = function(index, thisTag) {
        		$scope.workloadProfile.tags(index, 1);
        	};
        	
        	vm.workloadEditSave = function() {
        		$scope.editShow1 = false;
        		$scope.editShow2 = false;
        		$scope.editShow3 = false;
        		
        		// PUT to update
        		var paramUpdate = {
        			"action": "update_workload",
        			"user": $rootScope.user.name, 
        			"workload": {
        				"id": $scope.workloadProfile.id,
        				"name": $scope.workloadProfile.name,
        				"team": $scope.workloadProfile.team,
        				"run_as_user": $scope.workloadProfile.runuser,
        				"desc": $scope.workloadProfile.desc,
        				"tag": $scope.workloadProfile.tags,
	        			"setup_script": $scope.workloadProfile.setup_script,
                        "clean_script": $scope.workloadProfile.clean_script,
                        "exe_script": $scope.workloadProfile.exe_script,
                        "env_script": $scope.workloadProfile.env_script,
                        "machine": $scope.workloadProfile.machines,
                        "group": $scope.workloadProfile.groups,
                        "testsuites": $scope.workloadProfile.testsuites
        			}	
        		};
        		
        		$http
            		.put('/workload', JSON.stringify(paramUpdate), 
            			{ headers: {"Authorization": localStorageService.get("token")} })
            		.success(function(data, status, config, headers, statusText) {
              			if ( data.ret != 0 ) {
                			Notify.alert( 
	                			"[DEBUG] Update Workload $http (PUT) Failed! \n"+
								'data.ret:   ' +data.ret+'\n'+
								'data.info:  ' +data.info+'\n'+
								'status:     ' +status+'\n', 
	                			{status: 'danger', timeout: 2000}
	            			);
              			} else {
              				//Good! update current workloadSet for display
              				var i = vm.checkWorkloadIndex;
              				
              				if (data.id == vm.workloadSet[i].id) {
              					vm.workloadSet[i].name = data.name;
	              				vm.workloadSet[i].team = data.team;
	              				vm.workloadSet[i].desc = data.desc;
	              				vm.workloadSet[i].tag = data.tag;
	              				
	              				Notify.alert( 
		                			"[TSUITE004] Workload Updated! \n"+
									'data.ret:   ' +data.ret+'\n'+
									'data.info:  ' +data.info+'\n'+ 
									data.name+' - '+data.id,
		                			{status: 'success', timeout: 2000}
		            			);
              				} else {
              					Notify.alert( 
		                			"[DEBUG] Update Workload - Wrong index Recorded! \n"+
									'data.id:   ' +data.ret+'\n'+
									'data.name:  ' +data.info+'\n'+
									'check index: '+i+
									'status:     ' +status+'\n', 
		                			{status: 'danger', timeout: 2000}
		            			);
              				}
              			}
            		}) 
            		.error(function(data, status, config, headers, statusText) {
	                    Notify.alert( 
			                "[WKLD005] Update Workload $http (PUT) Failed! \n"+
							'data.ret:   ' +data.ret+'\n'+
							'data.info:  ' +data.info+'\n'+
							'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
            		});
        		
        	};
        	
        	vm.workloadEditCancel = function() {
        		$scope.editShow1 = false;
        		$scope.editShow2 = false;
        		$scope.editShow3 = false;
        	};
        	
        	vm.workloadEditItemAdd = function(thisname, category) {
        		switch (category) {
        			case 'tag': 		
        				$scope.workloadProfile.tags.push(thisname); 
        				$scope.tagInput = "";
        				break;
        			case 'group':		
        				$scope.workloadProfile.groups.push(thisname);
        				$scope.editShow3 = false;
        				break;
        			case 'testsuite':	
        				$scope.workloadProfile.testsuites.push({name: thisname.name, weight: thisname.weight});
        				$scope.editShow1 = false; 
        				break;
        			case 'machine':	
        				$scope.workloadProfile.machines.push(thisname);
        				$scope.editShow2 = false; 
        				break;
        			default:
        		};
        	};
        	
        	vm.workloadEditItemRemove = function(index, thisname, category) {
        		Notify.alert('[Debug] remove category:'+category+' index:'+index, {status: 'warning', timeout: 2000});
        		switch (category) {
        			case 'tag':
        				$scope.workloadProfile.tags.splice(index,1); break;
        			case 'group':
        				$scope.workloadProfile.groups.splice(index,1); break;
        			case 'testsuite':
        				$scope.workloadProfile.testsuites.splice(index,1); break;
        			case 'machine':
        				$scope.workloadProfile.machines.splice(index,1); break;
        			default:
        		};
        	};
        	
        	///////////////////
        }
    	
    }

})();

//*--workload coder (code editor)
(function() {
    'use strict';

    angular
        .module('app.custom')
        .controller('WorkloadCoder', WorkloadCoder);
	
	WorkloadCoder.$inject = ['$rootScope', '$scope', '$http', '$ocLazyLoad', 'SweetAlert'];
    function WorkloadCoder($rootScope, $scope, $http, $ocLazyLoad, SweetAlert) {
        var vm = this;

        layout();
        activate();

        ////////////////
        
        function layout() {
        	// Setup the layout mode 
        	$rootScope.app.useFullLayout = false;
          	$rootScope.app.hiddenFooter = false;
          	$rootScope.app.layout.isCollapsed = false;
          	
          	// Restore layout for demo
          	$scope.$on('$destroy', function(){
              	$rootScope.app.useFullLayout = false;
              	$rootScope.app.hiddenFooter = false;
              	$rootScope.app.layout.isCollapsed = false;
          	});	
        }
        
        function activate() {
        	// Available themes
          	vm.editorThemes = ['3024-day','3024-night','ambiance-mobile','ambiance','base16-dark','base16-light','blackboard',
          					   'cobalt','eclipse','elegant','erlang-dark','lesser-dark','mbo','mdn-like','midnight','monokai',
          					   'neat','neo','night','paraiso-dark','paraiso-light','pastel-on-dark','rubyblue','solarized',
          					   'the-matrix','tomorrow-night-eighties','twilight','vibrant-ink','xq-dark','xq-light'];
        	
        	vm.editorOpts = {
            	mode: 'javascript',
            	lineNumbers: true,
            	matchBrackets: true,
            	theme: 'mbo',
            	viewportMargin: Infinity
          	};
          	
          	vm.refreshEditor = 0;
          	
          	// lazy load Theme
          	vm.loadTheme = function(thisTheme) {
            	var BASE = 'vendor/codemirror/theme/';
            	$ocLazyLoad.load(BASE + thisTheme + '.css');
            	vm.refreshEditor = !vm.refreshEditor;
          	};
          	
          	// load default theme
          	vm.loadTheme(vm.editorOpts.theme);
          	
          	// load JCL Template for submition
          	vm.loadTemplate = function() {
          		$http.get('data/sample-coder.jcl')
			    	.success(function(data) {
			     		var codes = data,
			     			runuser = $scope.workloadProfile.runuser,
			     			lpar = $scope.workloadProfile.machines[0],
			     			sscript = $scope.workloadProfile.setup_script,
			     			cscript = $scope.workloadProfile.clean_script,
			     			xscript = $scope.workloadProfile.exe_script,
			     			escript = $scope.workloadProfile.env_script;
			     		
			     		codes = codes.replace(/{{theRunUser}}/g, runuser);
			     		codes = codes.replace(/{{theMachine}}/g, lpar);
			     		codes = codes.replace(/{{theEnvScript}}/g, escript);
			     		codes = codes.replace(/{{theSetupScript}}/g, sscript);
			     		codes = codes.replace(/{{theExeScript}}/g, xscript);
			     		codes = codes.replace(/{{theCleanScript}}/g, cscript);
			     		
			     		vm.code = codes;
			     	})
			        .error(function(data, status, headers, config) {
			        	Notify.alert( 
			                "[WKLD002] Failure Getting Sample Data for Workload Template...\n\n"+'status:     ' +status+'\n',
			                {status: 'danger', timeout: 2000}
			            );
			        });
          	};
          	
          	// watch for ng-show for workloadRun (level3)
          	$scope.$watch('isLevel3RunShown', function(newValue, oldValue) {
          		if (newValue == true) {
          			vm.loadTemplate();
          		}
          	});
          	
          	// submit code
          	vm.submitCode = function() {
          		//codes in vm.code
          		SweetAlert.swal('Good job!', vm.code, 'success');
          	};
          	
        }
        
    }
})();


//*--The End--