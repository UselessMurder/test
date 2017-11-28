"use_strict";
(function() {
	function set_hight() {
		if(navigator.userAgent.toLowerCase().indexOf('firefox') > -1){
    		$('.tab-pane').height($(window).height() - $("#top-navbar").height() * 3.8);
		} else {
   			$('.tab-pane').height($(window).height() - $("#top-navbar").height() * 2.7);
   		}
	}
	function submitForm(val) {
		$("#QuestionStatus").val(val);
		if (val === 2) {
			if(($('#KTextArea').val().trim().length > 5)) {
				$('form#questionform').submit();
			} else {
				if(!$('#KTextArea').hasClass("is-invalid")) {
					$('#KTextArea').removeClass("is-valid").addClass("is-invalid");
					$('<label class="invalid-feedback">Поле должно быть заполнено!</label>').insertAfter("#KTextArea")
				}
			}
		} else {
			$('form#questionform').submit();
		}
	}
	window.SMB = submitForm
	set_hight();
})();