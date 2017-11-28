"use_strict";
(function() {
	function savePDF() { 
	 	var element = document.getElementById('content');
		html2pdf(element, {
		  margin:       0,
		  filename:     'report.pdf',
		  image:        { type: 'jpeg', quality: 1 },
		  html2canvas:  { dpi: 192, letterRendering: true },
		  jsPDF:        { unit: 'in', format: 'letter', orientation: 'portrait' }
		});
	}
	window.SMB = savePDF
})();