"use_strict";
(function() {
	function set_events() {
		$(".dropdown-item").click(function(e) {
			e.preventDefault();
			$("#TypeId").val($(this).attr("item-index"));
			$("#dropdownMenuButton").text($(this).text());
			var item_text = $("#" + $(this).attr("item-index")).attr("item-text")
			$("#Description").animate({opacity:0},function(){
		        $(this).text(item_text).animate({opacity:1});  
		    });
		});
	}
	set_events();
})();