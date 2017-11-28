"use_strict";
(function() {
	function savePDF() {
		html2pdf($("#content")[0], {
		  margin:       0,
		  filename:     'report.pdf',
		  image:        { type: 'jpeg', quality: 1 },
		  html2canvas:  { dpi: 192, letterRendering: true },
		  jsPDF:        { unit: 'in', format: 'letter', orientation: 'portrait' }
		});
	}
	window.SMB = savePDF
})();