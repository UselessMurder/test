"use_strict";
(function() {
	$(document).ready(function(){
		$.validator.addMethod("regx", function(value, element, regexpr) {          
		    return regexpr.test(value);
		}, "Поле заполнено некорректно.");
	    $("#settingsform").validate({
	    	errorClass: "is-invalid",
	       rules:{
	            InOrganizationName:{
	                required: true,
	                regx: /[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$/,
	            },
	            CommercialDesignation:{
	                required: true,
	                regx: /[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$/,
	            },
	            InContactPerson:{
	                required: true,
	                regx: /^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$/,
	            },
	            InContactPersonPost:{
	                required: true,
	                regx: /[А-ЯЁa-яё 1-9]{1,20}/,
	            },
	            InPhone:{
	                required: true,
	                regx: /^\+?\d{1,3}?[- .]?\(?(?:\d{2,3})\)?[- .]?\d\d\d[- .]?\d\d\d\d$/,
	            },
	            InEmail:{
	                required: true,
	                regx: /^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$/,
	            },
	            InAddress:{
	                required: true,
	                regx: /^(.+)\s+(\S+?)(-(\d+))?$/,
	            },
	            InCity:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            InState:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            InCountry:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            InIndex:{
	                required: true,
	                regx: /^\d{6}$/,
	            },
	            InURL:{
	                required: true,
	                regx: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/,
	            },
	            OutOrganizationName:{
	                required: true,
	                regx: /[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$/,
	            },
	            OutContactPerson:{
	                required: true,
	                regx: /^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$/,
	            },
	            OutContactPersonPost:{
	                required: true,
	                regx: /[А-ЯЁa-яё 1-9]{1,20}/,
	            },
	            OutPhone:{
	                required: true,
	                regx: /^\+?\d{1,3}?[- .]?\(?(?:\d{2,3})\)?[- .]?\d\d\d[- .]?\d\d\d\d$/,
	            },
	            OutEmail:{
	                required: true,
	                regx: /^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$/,
	            },
	            OutAddress:{
	                required: true,
	                regx: /^(.+)\s+(\S+?)(-(\d+))?$/,
	            },
	            OutCity:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            OutState:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            OutCountry:{
	                required: true,
	                regx: /[a-zA-Zа-яёА-ЯЁ .]{1,30}/,
	            },
	            OutIndex:{
	                required: true,
	                regx: /^\d{6}$/,
	            },
	            OutURL:{
	                required: true,
	                regx: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/,
	            },
	       },
	       messages:{
	             InOrganizationName:{
	               required: "Это поле обязательно для заполнения",
	            },
	            CommercialDesignation:{
	                required: "Это поле обязательно для заполнения",
	            },
	            InContactPerson:{
	               required: "Это поле обязательно для заполнения",
	            },
	            InContactPersonPost:{
	                required: "Это поле обязательно для заполнения",
	            },
	            InPhone:{
	               required: "Это поле обязательно для заполнения",
	            },
	            InEmail:{
	                required: "Это поле обязательно для заполнения",
	            },
	            InAddress:{
	                required: "Это поле обязательно для заполнения",
	            },
	            InCity:{
	               required: "Это поле обязательно для заполнения",
	            },
	            InState:{
	                required: "Это поле обязательно для заполнения",
	            },
	            InCountry:{
	               required: "Это поле обязательно для заполнения",
	            },
	            InIndex:{
	               required: "Это поле обязательно для заполнения",
	            },
	            InURL:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutOrganizationName:{
	               required: "Это поле обязательно для заполнения",
	            },
	            OutContactPerson:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutContactPersonPost:{
	               required: "Это поле обязательно для заполнения",
	            },
	            OutPhone:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutEmail:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutAddress:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutCity:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutState:{
	               required: "Это поле обязательно для заполнения",
	            },
	            OutCountry:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutIndex:{
	                required: "Это поле обязательно для заполнения",
	            },
	            OutURL:{
	               required: "Это поле обязательно для заполнения",
	            },
	        },
	    });
	});
})();