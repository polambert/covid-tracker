
function onload() {
	// add commas to all of the numbers in all of the cells
	$("td").each(function(i, o) {
		var text = $(this).text();
		var numberAttempt = parseInt(text);

		if (!Number.isNaN(numberAttempt)) {
			// Is actually a number
			$(this).text(numberAttempt.toLocaleString());
		}
	});

	$("#main-table").stickyTableHeaders();

	// auto-sort by cases on page load
	var myTH = document.getElementsByTagName("th")[1];
	sorttable.innerSortFunction.apply(myTH, []);
}
